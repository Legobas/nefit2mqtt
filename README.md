# Nefit2MQTT

[![Docker Pulls](https://badgen.net/docker/pulls/legobas/nefit2mqtt?icon=docker&label=pulls)](https://hub.docker.com/r/legobas/nefit2mqtt)
[![Docker Image Size](https://badgen.net/docker/size/legobas/nefit2mqtt?icon=docker&label=image%20size)](https://hub.docker.com/r/legobas/nefit2mqtt)

**Nefit2MQTT** is a bridge between **Nefit Easy** thermostats and MQTT-based home-automation systems. It allows you to monitor and control your thermostat from Home Assistant, Domoticz, Node-RED, or any other MQTT platform.

Under the hood, Nefit2MQTT uses the [Nefit Easy HTTP Server](https://github.com/robertklep/nefit-easy-http-server) created by Robert Klep. A few patches are applied to use it with a recent LTS version of Node.

---

## Key Features

* **Simple Setup:** Single YAML configuration file with sensible defaults.
* **Smart Polling & Minimal API Usage:** Polls every 20 seconds when the thermostat is active and reduces polling to once every 90 minutes during inactivity to protect the Bosch/Nefit backend. Immediate updates are triggered after any MQTT command.
* **Home Assistant Auto-Discovery:** Automatically registers your Nefit Easy thermostat and sensors in Home Assistant via MQTT Discovery.
* **Secure Storage:** Sensitive session data and tokens are stored encrypted in `session.dat`.

---

## Quick Start

### 1. Create a data directory
Create a directory on your host to store the configuration and encrypted session data:

```bash
mkdir -p /opt/nefit2mqtt
```

### 2. Create the configuration file
Create `/opt/nefit2mqtt/nefit2mqtt.yml` with your MQTT broker details and Nefit Easy credentials:

```yaml
mqtt:
  url: tcp://mqttbroker:1883
  username: your_mqtt_username
  password: your_mqtt_password

nefit:
  serial_number: 'your_serial_number'
  access_key: 'your_access_key'
  password: 'your_nefit_password'

homeassistant:
  discovery: true
```

### 3. Start with Docker Compose
Save the following as `compose.yml`:

```yaml
services:
  nefit2mqtt:
    image: legobas/nefit2mqtt:latest
    container_name: nefit2mqtt
    init: true
    environment:
      - LOGLEVEL=info
      - TZ=Europe/Amsterdam
    volumes:
      - /opt/nefit2mqtt:/data:rw
    restart: unless-stopped
```

Start the container:

```bash
docker compose up -d
```

### 4. Verify logs and MQTT messages
Check the startup logs to confirm both the Nefit backend and MQTT broker connected successfully:

```bash
docker compose logs -f nefit2mqtt
```

Subscribe to the base topic (default: `nefit2mqtt`) to verify incoming messages:

```bash
mosquitto_sub -h mqttbroker -t 'nefit2mqtt/#' -v
```

---

## Configuration

All available options for `nefit2mqtt.yml`:

### `mqtt`
* `url` — MQTT broker URL, e.g., `tcp://mqttbroker:1883` (**required**)
* `username` — MQTT username (optional, default: empty)
* `password` — MQTT password (optional, default: empty)
* `qos` — MQTT Quality of Service: `0`, `1`, or `2` (default: `0`)
* `retain` — Retain published MQTT messages (default: `false`)
* `base_topic` — Base topic prefix for all messages (default: `nefit2mqtt`)

### `nefit`
* `serial_number` — Nefit Easy serial number without spaces or hyphens (**required**)
* `access_key` — Nefit Easy access key without spaces or hyphens (**required**)
* `password` — Your personal Nefit Easy app password (**required**)

### `homeassistant`
* `discovery` — Enable Home Assistant MQTT Discovery (`true` or `false`, default: `false`)

---

## Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `LOGLEVEL` | Set logging verbosity (`debug`, `info`, `warn`, `error`) | `info` |
| `TZ` | Container timezone (e.g., `Europe/Amsterdam`) | `UTC` |

---

## Home Assistant Integration

When `homeassistant.discovery` is set to `true`, Nefit2MQTT automatically publishes MQTT discovery payloads. Once the container starts, your Nefit Easy thermostat will appear automatically under **Settings → Devices & Services → MQTT** in Home Assistant, including climate controls and diagnostic sensors (such as room temperature, target temperature, and system pressure).

---

## Security & Session Storage

Session tokens and sensitive data are stored in `/data/session.dat`. This file is encrypted at rest to protect your credentials.

---

## Building from Source

To build the binary manually:

```bash
go build -ldflags "-X main.Version=1.0.0" -o nefit2mqtt .
```

When running outside Docker, the application looks for `nefit2mqtt.yml` in the following order:
1. `/data/nefit2mqtt.yml`
2. `~/.data/nefit2mqtt.yml`
3. `./nefit2mqtt.yml` (current working directory)
4. `./data/nefit2mqtt.yml`

---

## Credits

* [Nefit Easy HTTP Server](https://github.com/robertklep/nefit-easy-http-server) by Robert Klep
* [Eclipse Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang) - MQTT connectivity
* [ZeroLog](https://github.com/rs/zerolog) - Structured JSON logging
