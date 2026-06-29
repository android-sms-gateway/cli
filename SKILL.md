---
name: smsgate-cli
description: >
  Use when you need to send SMS messages, manage webhooks, check delivery status,
  retrieve logs, or issue CA certificates via the SMS Gateway for Android API.
  Triggered on SMS sending, webhook management, message status checks, and log retrieval.
license: Apache-2.0
---

# SMSGate CLI

A CLI tool for interacting with the [SMS Gateway for Android](https://sms-gate.app) API. Ships two binaries:

- **`smsgate`** — send SMS, batch operations, manage webhooks, fetch logs
- **`smsgate-ca`** — issue TLS certificates for private deployments

## Installation

### Option 1: GitHub Releases

Download the latest archive from the [Releases page](https://github.com/android-sms-gateway/cli/releases/latest). Archives are built with GoReleaser for Linux, macOS, and Windows. Each archive contains both `smsgate` and `smsgate-ca`.

Archive naming: `smsgate_Linux_x86_64.tar.gz`, `smsgate_Windows_x86_64.zip`, `smsgate_Darwin_arm64.tar.gz`, etc.

```bash
# Use bash: extract archive and move binaries into PATH
tar xzf smsgate_Linux_x86_64.tar.gz
sudo mv smsgate smsgate-ca /usr/local/bin/
```

### Option 2: Install using Go

```bash
# Use bash: install via Go toolchain
go install github.com/android-sms-gateway/cli/cmd/smsgate@latest
go install github.com/android-sms-gateway/cli/cmd/smsgate-ca@latest
```

> **Note**: Replace `@latest` with a specific version tag (e.g. `@v1.2.3`) for reproducible installs.

### Option 3: Docker

```bash
# Use bash: run via Docker
docker run -it --rm --env-file .env ghcr.io/android-sms-gateway/cli \
  send --phones '+12025550123' 'Hello!'
```

The Docker image includes both `smsgate` (entrypoint) and `smsgate-ca`.

## Updating

- **Go**: Re-run the `go install` commands above with the target version tag.
- **Docker**: `docker pull ghcr.io/android-sms-gateway/cli`

## Configuration

Credentials can be passed via CLI flags, environment variables, or a `.env` file.

| Flag | Env Var | Description | Default |
|------|---------|-------------|---------|
| `--endpoint`, `-e` | `ASG_ENDPOINT` | API endpoint URL | `https://api.sms-gate.app/3rdparty/v1` |
| `--username`, `-u` | `ASG_USERNAME` | Username | required |
| `--password`, `-p` | `ASG_PASSWORD` | Password | required |
| `--format`, `-f` | — | Output format | `text` |

The `.env` file in the working directory is loaded automatically.

## Commands

### `smsgate send`

Send a text or data (binary) SMS message.

```bash
smsgate send [flags] <message>
```

| Flag | Description | Default |
|------|-------------|---------|
| `--phones`, `-p` | Recipient phone number(s), E.164 format (repeatable or comma-separated) | required |
| `--id` | Custom message ID | auto-generated |
| `--device-id`, `--device` | Specific device ID | auto |
| `--sim-number`, `--sim` | SIM slot (one-based) | device default |
| `--delivery-report` | Enable delivery report | `true` |
| `--priority` | Priority (-128 to 127; >= 100 bypasses limits) | `0` |
| `--ttl` | Time-to-live (e.g. `1h30m`) | unlimited |
| `--valid-until` | Expiration time (RFC3339) | unlimited |
| `--schedule-at` | Schedule delivery time (RFC3339) | immediate |
| `--skip-phone-validation` | Skip phone validation | `false` |
| `--device-active-within` | Filter by device activity (hours) | `0` (no filter) |
| `--data` | Send data message (content must be base64) | `false` |
| `--data-port` | Destination port for data message | `53739` |

### `smsgate status`

Check the delivery state of a sent message.

```bash
smsgate status <message-id>
```

### `smsgate batch send`

Send messages in bulk from a CSV or XLSX file.

```bash
smsgate batch send [flags] <filename>
```

| Flag | Description | Default |
|------|-------------|---------|
| `--map` | Column mapping: `phone=Col,text=Col` (comma-separated) | required |
| `--sheet` | Sheet name (XLSX only) | first sheet |
| `--delimiter` | CSV delimiter | `,` |
| `--header` | Treat first row as header | `true` |
| `--dry-run` | Validate and preview without sending | `false` |
| `--validate-only` | Validate input only (no preview) | `false` |
| `--concurrency` | Number of concurrent workers | CPU cores |
| `--continue-on-error` | Continue after per-row failures | `false` |

Also accepts the shared flags from `send` (`--device-id`, `--sim-number`, `--priority`, `--ttl`, `--valid-until`, `--delivery-report`, `--skip-phone-validation`, `--device-active-within`), applied to every message in the batch.

**Column mapping fields:**

| Field | Required | Description |
|-------|----------|-------------|
| `phone` | yes | Phone number |
| `text` | yes | Message text |
| `id` | no | Custom message ID |
| `device_id` | no | Device identifier |
| `sim_number` | no | SIM slot number |
| `priority` | no | Message priority |

**Workflow modes:**

1. **Validate only** — check file and mapping: `--validate-only`
2. **Dry run** — preview all parsed rows: `--dry-run`
3. **Full send** — send with progress: `--concurrency=5`

### `smsgate webhooks`

Manage webhooks for event notifications.

```bash
# Register a webhook
smsgate webhooks register --event <event> [--id ID] [--device-id DEVICE] <url>

# List all webhooks
smsgate webhooks list

# Delete a webhook
smsgate webhooks delete <id>
```

Events: `sms:received`, `sms:sent`, `sms:failed`, `device:connected`, `device:disconnected`

### `smsgate logs`

Retrieve logs within a time range.

```bash
smsgate logs [--from TIME] [--to TIME]
```

Dates use RFC3339 format (e.g. `2024-01-15T10:30:00Z`). Defaults to the last 24 hours.

### `smsgate-ca`

Issue TLS certificates for private SMS Gateway deployments.

```bash
smsgate-ca [--timeout DURATION] webhooks [--out FILE] [--keyout FILE] <ip-address>
smsgate-ca [--timeout DURATION] private [--out FILE] [--keyout FILE] <ip-address>
```

| Flag | Description | Default |
|------|-------------|---------|
| `--timeout`, `-t` | Request timeout | `30s` |
| `--out` | Certificate output file | `server.crt` |
| `--keyout` | Private key output file | `server.key` |

The IP address must be a private (RFC 1918) address.

## Output Formats

All commands support `--format` (or `-f`) with four options:

- **`text`** — human-readable key:value pairs (default)
- **`json`** — pretty-printed JSON
- **`raw`** — compact one-line JSON
- **`table`** — tab-aligned columns

Error messages are always printed to stderr in plain text regardless of format.

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Invalid options or arguments |
| 2 | Server request error |
| 3 | Output formatting error |
| 4 | Internal error |

## Examples

```bash
# Send a text message
smsgate send --phones '+12025550123' 'Hello, Dr. Turk!'

# Send to multiple phones
smsgate send --phones '+12025550123,+12025550124' 'Hello!'

# Send with priority (bypasses limits)
smsgate send --phones '+12025550123' --priority 100 'Urgent'

# Send a data message
echo -n 'hello' | base64
smsgate send --phones '+12025550123' --data --data-port 12345 'aGVsbG8='

# Check message status
smsgate status zXDYfTmTVf3iMd16zzdBj

# Batch send from CSV
smsgate batch send contacts.csv --map phone=Phone,text=Message

# Batch send from XLSX with a specific sheet
smsgate batch send campaign.xlsx --sheet Sheet1 --map phone=Phone,text=Message

# Dry-run batch (validate + preview)
smsgate batch send contacts.csv --map phone=Phone,text=Message --dry-run

# Register a webhook
smsgate webhooks register --event sms:received https://example.com/hook

# List webhooks (table format)
smsgate --format table webhooks list

# Delete a webhook
smsgate webhooks delete wh_abc123

# Get logs for a specific time range
smsgate logs --from '2024-01-15T00:00:00Z' --to '2024-01-15T23:59:59Z'

# Issue a CA certificate for webhooks
smsgate-ca webhooks 192.168.1.100 --out server.crt --keyout server.key

# Issue a CA certificate for a private server
smsgate-ca private 10.0.0.5 --out myserver.crt --keyout myserver.key
```
