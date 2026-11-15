package nefit

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	baseUrl        = "http://nefit:5533"
	requestTimeout = 10 * time.Second
)

type NefitSettings struct {
	Username string
	Password string
}

func New(username string, password string) (*NefitSettings, error) {
	var err error

	log.Trace().Msgf("Username: %s", username)

	rc := &NefitSettings{
		Username: username,
		Password: password,
	}

	if rc.NefitOnline() {
		err = nil
	} else {
		err = errors.New("could not connect to Nefit")
	}

	return rc, err
}

func (rc NefitSettings) SendToNefit(command string, jsonMsg string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	log.Debug().Msgf("jsonMsg: %s", jsonMsg)

	url := baseUrl + command
	async := false

	log.Debug().Msgf("Nefit url: %s", url)
	log.Debug().Msgf("Nefit json: %s", jsonMsg)

	payload := bytes.NewReader([]byte(jsonMsg))
	log.Trace().Msgf("Nefit Request:\n%v", jsonMsg)

	req, err := http.NewRequestWithContext(ctx, "POST", url, payload)
	if err != nil {
		log.Error().Err(err).Msg("HTTP Request error")
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	if rc.Username != "" {
		req.SetBasicAuth(rc.Username, rc.Password)
	}

	httpClient := &http.Client{Timeout: requestTimeout}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Fatal().Err(err).Msg("HTTP Client error")
	}
	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal().Err(err)
		return "", err
	}

	responseMessage := string(response)

	log.Trace().Msgf("Nefit Response:\n%s", responseMessage)

	if async {
		log.Debug().Msgf("Job: %s", responseMessage)
	}

	return responseMessage, nil
}

func (rc NefitSettings) NefitOnline() bool {
	online := false
	for i := 1; i < 10; i++ {
		time.Sleep(time.Second * 1)
		response, err := rc.SendToNefit("/core/version", "{}")
		if err == nil {
			online = true
			data := map[string]any{}
			err := json.Unmarshal([]byte(response), &data)
			if err != nil {
				log.Fatal().Err(err)
			}
			log.Info().Msgf("Nefit version: %s", data["version"])
			break
		}
	}
	return online
}
