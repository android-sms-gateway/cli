# ⌨️ SMSGate CLI

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stars][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

A command-line interface for the SMS Gateway for Android API. Ships two binaries: `smsgate` (messages, webhooks, logs) and `smsgate-ca` (certificates for private deployments). Part of the SMSGate ecosystem - see the [CLI docs](https://docs.sms-gate.app/integration/cli/) for the full reference.

## 📖 About

The CLI lets you automate SMS Gateway for Android without writing code: send messages (including batch operations from CSV and Excel files), manage webhooks, retrieve logs, and - via `smsgate-ca` - issue certificates for private server deployments. It can be used interactively or embedded in CI/CD pipelines, shell scripts, and cron jobs.

## 📚 Table of Contents

- [⌨️ SMSGate CLI](#️-smsgate-cli)
  - [📖 About](#-about)
  - [📚 Table of Contents](#-table-of-contents)
  - [⭐ Features](#-features)
  - [💻 Installation](#-installation)
    - [Option 1: Download from GitHub Releases](#option-1-download-from-github-releases)
    - [Option 2: Install using Go](#option-2-install-using-go)
    - [Option 3: Docker](#option-3-docker)
  - [💻 Usage](#-usage)
    - [Sending messages](#sending-messages)
    - [Batch message sending](#batch-message-sending)
    - [Getting message status](#getting-message-status)
    - [Managing webhooks](#managing-webhooks)
    - [Getting logs](#getting-logs)
    - [Output formats](#output-formats)
    - [smsgate-ca](#smsgate-ca)
    - [Exit Codes](#exit-codes)
  - [💻 Configuration](#-configuration)
  - [📚 Documentation](#-documentation)
  - [🤝 Contributing](#-contributing)
  - [⚖️ License](#️-license)
  - [📜 Legal Notice](#-legal-notice)

## ⭐ Features

- Send SMS to one or multiple recipients with delivery options: device and SIM selection, priority, TTL, scheduled delivery, delivery reports
- Send MMS with subject and multiple attachments, and data messages with base64 payloads and custom destination ports
- Batch sending from CSV and XLSX files with column mapping, dry-run and validate-only modes
- Webhook management (register, list, delete) for event notifications
- Log retrieval for arbitrary time ranges
- Multiple output formats: `text`, `json`, `raw`, `table`
- Certificate issuing for private server deployments via `smsgate-ca`

## 💻 Installation

### Option 1: Download from GitHub Releases

1. Go to the [Releases page](https://github.com/android-sms-gateway/cli/releases/latest) of this repository.
2. Download the appropriate binary archive for your operating system and architecture. Each archive contains both `smsgate` and `smsgate-ca`.
3. Extract the archive to a directory of your choice.
4. Move the binaries to a directory in your system's PATH.

```bash
curl -LO https://github.com/android-sms-gateway/cli/releases/latest/download/smsgate_Linux_x86_64.tar.gz
tar -xzf smsgate_Linux_x86_64.tar.gz
chmod +x smsgate smsgate-ca
sudo mv smsgate smsgate-ca /usr/local/bin/
```

### Option 2: Install using Go

Requires Go 1.26+ (for building from source). Make sure `$(go env GOPATH)/bin` is in your `PATH`.

```bash
go install github.com/android-sms-gateway/cli/cmd/smsgate@latest
go install github.com/android-sms-gateway/cli/cmd/smsgate-ca@latest
```

### Option 3: Docker

Docker (optional, for containerized execution).

```bash
docker run -it --rm --env-file .env ghcr.io/android-sms-gateway/smsgate \
  send --phones '+12025550123' 'Hello, Dr. Turk!'
```

## 💻 Usage

```bash
smsgate [global options] command [command options] [arguments...]
```

| Group    | Commands                                                | Description                                   |
| -------- | ------------------------------------------------------- | --------------------------------------------- |
| Messages | `send`, `status` (alias: `state`), `batch send`         | Send and track messages, bulk from CSV/XLSX   |
| Webhooks | `webhooks register`, `webhooks list`, `webhooks delete` | Register, audit, and remove webhook endpoints |
| Logs     | `logs` (alias: `log`)                                   | Retrieve logs for a specific time range       |

For a complete list of available commands, you can:
- Run `smsgate help` or `smsgate --help` in your terminal.
- Visit the official documentation at [docs.sms-gate.app](https://docs.sms-gate.app/integration/cli/#commands).

Credentials are read from environment variables (see Configuration); never put them on the command line.

### Sending messages

The `send` command supports various options to customize message delivery:

```bash
# Send a message
smsgate send --phones '+12025550123' 'Hello, Dr. Turk!'

# Send to multiple numbers
smsgate send --phones '+12025550123' --phones '+12025550124' 'Hello, doctors!'
# or
smsgate send --phones '+12025550123,+12025550124' 'Hello, doctors!'

# Send with explicit device selection
smsgate send --phones '+12025550123' --device-id device123 'Message'

# Send with SIM number selection (1-based)
smsgate send --phones '+12025550123' --sim-number 2 'Message'

# Send with priority (>=100 bypasses limits)
smsgate send --phones '+12025550123' --priority 100 'Urgent message'

# Send with time-to-live (TTL)
smsgate send --phones '+12025550123' --ttl 1h30m 'Expiring message'

# Send with expiration date (RFC3339 format)
smsgate send --phones '+12025550123' --valid-until '2026-12-31T23:59:59Z' 'Message'

# Disable delivery report
smsgate send --phones '+12025550123' --delivery-report=false 'Message'

# Skip phone number validation
smsgate send --phones '+12025550123' --skip-phone-validation 'Message'

# Filter by device activity (devices active within last 12 hours)
smsgate send --phones '+12025550123' --device-active-within 12 'Message'

# Schedule delivery (RFC3339 format)
smsgate send --phones '+12025550123' --schedule-at '2026-12-31T23:59:59Z' 'Message'

# Send data message (base64 encoded)
echo -n 'hello world' | base64
smsgate send --phones '+12025550123' --data --data-port 12345 'aGVsbG8gd29ybGQ='

# Send MMS with subject and multiple attachments
smsgate send --mms --subject 'Project update' --phones '+12025550123' \
  --attachment /path/to/report.pdf \
  --attachment /path/to/chart.png \
  'Here are the latest documents'

# Send MMS with attachments only (no text)
smsgate send --mms --phones '+12025550123' \
  --attachment /path/to/photo.jpg

# Send MMS to multiple recipients
smsgate send --mms --phones '+12025550123' --phones '+12025550124' \
  --subject 'Team update' --attachment /path/to/screenshot.png \
  'Please review'
```

**Notes:**
- `--data` and `--mms` are mutually exclusive.
- At least one of text or an attachment is required for MMS.
- MIME type is detected from the file extension (falls back to `application/octet-stream`).

**Send command options:**

| Option                      | Description                                                                                                                                               | Default Value | Example                 |
| --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------- | ----------------------- |
| `--id`                      | A unique message ID. If not provided, one will be automatically generated.                                                                                | empty         | `zXDYfTmTVf3iMd16zzdBj` |
| `--device-id`, `--device`   | Optional device ID for explicit selection. If not provided, a random device will be selected.                                                             | empty         | `oi2i20J8xVP1ct5neqGZt` |
| `--phones`, `--phone`, `-p` | Specifies the recipient's phone number(s). This option can be used multiple times or accepts comma-separated values. Numbers must be in E.164 format.     | **required**  | `+12025550123`          |
| `--sim-number`, `--sim`     | The one-based SIM card slot number. If not specified, the device's SIM rotation feature will be used.                                                     | empty         | `2`                     |
| `--delivery-report`         | Enables delivery report for the message.                                                                                                                  | `true`        | `true` / `false`        |
| `--priority`                | Sets the priority of the message. Messages with priority >= 100 bypass all limits and delays. Range: -128 to 127.                                         | `0`           | `100`                   |
| **Data Message**            |                                                                                                                                                           |               |                         |
| `--data`                    | Send data message instead of text (content must be base64 encoded).                                                                                       | `false`       | `true`                  |
| `--data-port`               | Destination port for data message (1 to 65535).                                                                                                           | `53739`       | `12345`                 |
| **MMS Message**             |                                                                                                                                                           |               |                         |
| `--mms`                     | Send MMS message instead of text; text is optional when at least one attachment is provided.                                                              | `false`       | `true`                  |
| `--subject`                 | MMS subject line. Requires `--mms`.                                                                                                                       | empty         | `Photo update`          |
| `--attachment`              | Path to an attachment file. Repeatable; MIME type is detected from the file extension. Requires `--mms`.                                                  | empty         | `/path/to/image.jpg`    |
| **Options**                 |                                                                                                                                                           |               |                         |
| `--ttl`                     | Time-to-live (TTL) for the message. Duration format (e.g., `1h30m`). If not provided, the message will not expire.<br>**Conflicts with `--valid-until`.** | empty         | `1h30m`                 |
| `--valid-until`             | The expiration date and time for the message. RFC3339 format (e.g., `2006-01-02T15:04:05Z07:00`).<br>**Conflicts with `--ttl`.**                          | empty         | `2026-12-31T23:59:59Z`  |
| `--skip-phone-validation`   | Skip phone number validation.                                                                                                                             | `false`       | `true`                  |
| `--device-active-within`    | Time window in hours for device activity filtering. `0` means no filtering.                                                                               | `0`           | `12`                    |
| `--schedule-at`             | Schedule delivery time (RFC3339 format). If not provided, the message is sent immediately.<br>**Example:** `2026-12-31T23:59:59Z`.                        | immediate     | `2026-12-31T23:59:59Z`  |

### Batch message sending

The CLI supports sending messages in bulk from CSV and Excel files. The `batch send` command also supports the shared delivery/device options from `send` (such as `--device-id`, `--sim-number`, `--priority`, `--ttl`, `--valid-until`, `--delivery-report`, `--skip-phone-validation`, `--device-active-within`), allowing fine-grained control over each message in the batch.

**Supported file formats:**
- **CSV** (Comma-Separated Values)
- **XLSX** (Excel files)

**Basic usage:**

```bash
# Batch send from a CSV file with column mapping
smsgate batch send --map phone=Phone,text=Message contacts.csv

# Send messages from an Excel file with specific sheet
smsgate batch send campaign.xlsx --sheet Sheet1 --map phone=Phone,text=Message

# Send with custom delimiter and no header
smsgate batch send data.csv --delimiter ';' --header=false --map phone=col_1,text=col_2
```

**Batch-specific options:**

| Option                | Description                                        | Default Value | Example                    |
| --------------------- | -------------------------------------------------- | ------------- | -------------------------- |
| `--sheet`             | Sheet name (defaults to first sheet)               | empty         | `Sheet1`                   |
| `--delimiter`         | CSV delimiter character                            | `,`           | `;`                        |
| `--header`            | Treat first row as header                          | `true`        | `false`                    |
| `--map`               | Column mapping (required)                          | **required**  | `phone=Phone,text=Message` |
| `--dry-run`           | Validate and print normalized rows without sending | `false`       | `true`                     |
| `--validate-only`     | Validate input only (no preview, no sending)       | `false`       | `true`                     |
| `--concurrency`       | Number of concurrent send workers                  | CPU cores     | `5`                        |
| `--continue-on-error` | Continue sending after per-row failures            | `false`       | `true`                     |

**Inherited message options:**
The shared delivery/device options from the `send` command are also available (e.g., `--device-id`, `--sim-number`, `--priority`, `--ttl`, `--valid-until`, `--delivery-report`, `--skip-phone-validation`, `--device-active-within`). These apply to every message sent in the batch.

**Column mapping:**

The `--map` option defines how columns in your file map to message fields:

| Field        | Required | Description                                   |
| ------------ | -------- | --------------------------------------------- |
| `phone`      | ✅        | Phone number column                           |
| `text`       | ✅        | Message text column                           |
| `id`         | ❌        | Message ID column (UUID generated when empty) |
| `device_id`  | ❌        | Device identifier column                      |
| `sim_number` | ❌        | SIM number column (1-255)                     |
| `priority`   | ❌        | Message priority column (-128 to 127)         |

**Mapping examples:**

```bash
# Basic mapping with headers
smsgate batch send --map phone=Phone,text=Message contacts.csv

# Excel file with specific sheet
smsgate batch send --sheet Sheet1 --map phone=Phone,text=Message campaign.xlsx

# No headers, column positions
smsgate batch send --header=false --map phone=col_1,text=col_2 data.csv

# Full mapping with optional fields
smsgate batch send --map phone=Phone,text=Message,device_id=Device,sim_number=SIM,priority=Priority contacts.csv
```

**Workflow modes:**

1. **Validation Only** - Validates file format and column mapping, checks required fields, exits without sending:
   ```bash
   smsgate batch send --map phone=Phone,text=Message --validate-only contacts.csv
   ```

2. **Dry Run** - Validates and processes all rows, shows what would be sent without actually sending:
   ```bash
   smsgate batch send --map phone=Phone,text=Message --dry-run contacts.csv
   ```

3. **Full Send** - Sends all messages with real-time progress:
   ```bash
   smsgate batch send --map phone=Phone,text=Message --concurrency=5 contacts.csv
   ```

**Output and error handling:**

- **Summary**: `Batch send summary: total=100 enqueued=95 failed=3 skipped=2`
- **Real-time progress**: Shows each message's UUID and state during sending
- **Error handling**: By default stops on first error; use `--continue-on-error` to send all rows even if some fail

**Best practices:**

1. Always use `--dry-run` or `--validate-only` first to test your configuration
2. Ensure phone numbers are in E.164 format
3. Start with lower concurrency values and increase as needed
4. Use `--continue-on-error` for non-critical bulk sends
5. Combine with inherited message options (e.g., `--device-id`, `--priority`) for advanced scenarios

### Getting message status

```bash
# Check delivery status
smsgate status zXDYfTmTVf3iMd16zzdBj
```

### Managing webhooks

```bash
# Register a webhook for SMS received events
smsgate webhooks register -e sms:received https://example.com/webhook

# Register a device-specific webhook
smsgate webhooks register -e sms:sent --device-id device123 https://example.com/webhook

# List all registered webhooks
smsgate webhooks list

# Delete a webhook by ID
smsgate webhooks delete <webhook-id>
```

**Supported webhook events:** `sms:received`, `sms:sent`, `sms:failed`, `device:connected`, `device:disconnected`

### Getting logs

The `logs` command retrieves logs for a specific time range. Dates should be in RFC3339 format (e.g., `2024-01-15T10:30:00Z`).

```bash
# Get logs for the last 24 hours (default)
smsgate logs

# Fetch logs for a specific time range
smsgate logs --from '2024-01-15T00:00:00Z' --to '2024-01-15T23:59:59Z'

# Get logs with custom time range and output format
smsgate --format json logs --from '2024-01-15T10:00:00+07:00' --to '2024-01-15T18:00:00+07:00'
```

### Output formats

**Text** (human-readable, default)

```text
ID: zXDYfTmTVf3iMd16zzdBj
State: Pending
IsHashed: false
IsEncrypted: false
Recipients:
        +12025550123    Pending
        +12025550124    Pending
```

**JSON** (pretty-printed)

```json
{
  "id": "zXDYfTmTVf3iMd16zzdBj",
  "state": "Pending",
  "isHashed": false,
  "isEncrypted": false,
  "recipients": [
    {
      "phoneNumber": "+12025550123",
      "state": "Pending"
    },
    {
      "phoneNumber": "+12025550124",
      "state": "Pending"
    }
  ],
  "states": {}
}
```

**Raw** (one-line JSON)

```json
{"id":"zXDYfTmTVf3iMd16zzdBj","state":"Pending","isHashed":false,"isEncrypted":false,"recipients":[{"phoneNumber":"+12025550123","state":"Pending"},{"phoneNumber":"+12025550124","state":"Pending"}],"states":{}}
```

**Table** (tab-aligned columnar output for lists)

```text
ID                                    EVENT          URL                              DEVICE ID
123e4567-e89b-12d3-a456-426614174000  sms:received   https://example.com/webhook      dev-abc
def45678-e89b-12d3-a456-426614174000  sms:sent       https://example.com/other
```

### smsgate-ca

```bash
# Issue a certificate for a private server
smsgate-ca private --out server.crt --keyout server.key 192.168.1.10
```

The IP address must be a private one (RFC 1918, e.g. `10.x.x.x`, `172.16.x.x`, `192.168.x.x`).

### Exit Codes

| Code | Meaning                   |
| ---- | ------------------------- |
| `0`  | Success                   |
| `1`  | Invalid options/arguments |
| `2`  | Server request error      |
| `3`  | Output formatting error   |
| `4`  | Internal error            |

Errors are printed to stderr without formatting. Full command reference: [CLI docs](https://docs.sms-gate.app/integration/cli/).

## 💻 Configuration

The CLI can be configured using environment variables or command-line flags. You can also use a `.env` file in the working directory to set these variables (loaded automatically via [godotenv](https://github.com/joho/godotenv)).

| Flag               | Env Var          | Applies to   | Description                                                   | Default value                          |
| ------------------ | ---------------- | ------------ | ------------------------------------------------------------- | -------------------------------------- |
| `--endpoint`, `-e` | `ASG_ENDPOINT`   | `smsgate`    | The endpoint URL                                              | `https://api.sms-gate.app/3rdparty/v1` |
| `--username`, `-u` | `ASG_USERNAME`   | `smsgate`    | Your username                                                 | **required**                           |
| `--password`, `-p` | `ASG_PASSWORD`   | `smsgate`    | Your password                                                 | **required**                           |
| `--format`, `-f`   | n/a              | `smsgate`    | Output format                                                 | `text`                                 |
| `--timeout`, `-t`  | `ASG_CA_TIMEOUT` | `smsgate-ca` | Certificate polling deadline, checked between status requests | `30s`                                  |

Output formats: `text` (human-readable, default), `json` (pretty-printed), `raw` (one-line JSON), `table` (columnar output for lists).

`smsgate-ca` talks to the certificate authority at `https://ca.sms-gate.app/api/v1` and accepts no credentials; its configurable options are `--timeout` (env `ASG_CA_TIMEOUT`), `--out`, and `--keyout`.

## 📚 Documentation

- [CLI reference](https://docs.sms-gate.app/integration/cli/)
- [Central docs](https://docs.sms-gate.app/)
- [GitHub repository](https://github.com/android-sms-gateway/cli)

## 🤝 Contributing

Contributions are welcome. Open an issue or pull request; follow the repo's `make fmt` / `make lint` / `make test` checks and the [contributing guidelines](https://docs.sms-gate.app/contributing/).

## ⚖️ License

Apache-2.0. See [LICENSE](LICENSE).

## 📜 Legal Notice

Android is a trademark of Google LLC.

[contributors-shield]: https://img.shields.io/github/contributors/android-sms-gateway/cli?style=for-the-badge
[contributors-url]: https://github.com/android-sms-gateway/cli/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/android-sms-gateway/cli?style=for-the-badge
[forks-url]: https://github.com/android-sms-gateway/cli/network/members
[stars-shield]: https://img.shields.io/github/stars/android-sms-gateway/cli?style=for-the-badge
[stars-url]: https://github.com/android-sms-gateway/cli/stargazers
[issues-shield]: https://img.shields.io/github/issues/android-sms-gateway/cli?style=for-the-badge
[issues-url]: https://github.com/android-sms-gateway/cli/issues
[license-shield]: https://img.shields.io/github/license/android-sms-gateway/cli?style=for-the-badge
[license-url]: https://github.com/android-sms-gateway/cli/blob/master/LICENSE
