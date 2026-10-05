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

- [📖 About](#-about)
- [📚 Table of Contents](#-table-of-contents)
- [⭐ Features](#-features)
- [💻 Installation](#-installation)
- [💻 Usage](#-usage)
- [💻 Configuration](#-configuration)
- [📚 Documentation](#-documentation)
- [🤝 Contributing](#-contributing)
- [⚖️ License](#️-license)
- [📜 Legal Notice](#-legal-notice)

## ⭐ Features

- Send SMS to one or multiple recipients with delivery options: device and SIM selection, priority, TTL, scheduled delivery, delivery reports
- Batch sending from CSV and XLSX files with column mapping, dry-run and validate-only modes
- Data messages with base64 payloads and custom destination ports
- Webhook management (register, list, delete) for event notifications
- Log retrieval for arbitrary time ranges
- Multiple output formats: `text`, `json`, `raw`, `table`
- Certificate issuing for private server deployments via `smsgate-ca`

## 💻 Installation

### Option 1: Download from GitHub Releases

Binaries are published for Linux, macOS, and Windows on the [Releases page](https://github.com/android-sms-gateway/cli/releases/latest). Each archive contains both `smsgate` and `smsgate-ca`.

```bash
curl -LO https://github.com/android-sms-gateway/cli/releases/latest/download/smsgate_Linux_x86_64.tar.gz
tar -xzf smsgate_Linux_x86_64.tar.gz
chmod +x smsgate smsgate-ca
sudo mv smsgate smsgate-ca /usr/local/bin/
```

### Option 2: Install using Go

Requires Go 1.25+. Make sure `$(go env GOPATH)/bin` is in your `PATH`.

```bash
go install github.com/android-sms-gateway/cli/cmd/smsgate@latest
go install github.com/android-sms-gateway/cli/cmd/smsgate-ca@latest
```

### Option 3: Docker

```bash
docker run -it --rm --env-file .env ghcr.io/android-sms-gateway/cli \
  send --phones '+12025550123' 'Hello, Dr. Turk!'
```

## 💻 Usage

```bash
smsgate [global options] command [command options] [arguments...]
```

| Group    | Commands                                                        | Description                                     |
| -------- | --------------------------------------------------------------- | ----------------------------------------------- |
| Messages | `send`, `status` (alias: `state`), `batch send`                 | Send and track messages, bulk from CSV/XLSX     |
| Webhooks | `webhooks register`, `webhooks list`, `webhooks delete`         | Register, audit, and remove webhook endpoints   |
| Logs     | `logs` (alias: `log`)                                           | Retrieve logs for a specific time range         |

Credentials are read from environment variables (see Configuration); never put them on the command line.

```bash
# Send a message
smsgate send --phones '+12025550123' 'Hello, Dr. Turk!'
```

```bash
# Batch send from a CSV file with column mapping
smsgate batch send --map phone=Phone,text=Message contacts.csv
```

```bash
# Check delivery status
smsgate status zXDYfTmTVf3iMd16zzdBj
```

```bash
# Register a webhook for received messages
smsgate webhooks register --event sms:received https://example.com/webhook
```

```bash
# Fetch logs for a specific time range
smsgate logs --from '2024-01-15T00:00:00Z' --to '2024-01-15T23:59:59Z'
```

`batch send` supports `--dry-run` and `--validate-only` to test input files before sending, plus shared message options such as `--device-id`, `--sim-number`, `--priority`, and `--ttl`.

### smsgate-ca

```bash
# Issue a certificate for a private server
smsgate-ca private --out server.crt --keyout server.key 203.0.113.10
```

### Exit Codes

| Code | Meaning                    |
| ---- | -------------------------- |
| `0`  | Success                    |
| `1`  | Invalid options/arguments  |
| `2`  | Server request error       |
| `3`  | Output formatting error    |

Errors are printed to stderr without formatting. Full command reference: [CLI docs](https://docs.sms-gate.app/integration/cli/).

## 💻 Configuration

The CLI is configured via command-line flags, environment variables, or a `.env` file in the working directory (loaded automatically). The `smsgate` and `smsgate-ca` binaries use the variables below.

| Flag                    | Env Var         | Description      | Default value                          |
| ----------------------- | --------------- | ---------------- | -------------------------------------- |
| `--endpoint`, `-e`      | `ASG_ENDPOINT`  | The endpoint URL | `https://api.sms-gate.app/3rdparty/v1` |
| `--username`, `-u`      | `ASG_USERNAME`  | Your username    | **required**                           |
| `--password`, `-p`      | `ASG_PASSWORD`  | Your password    | **required**                           |
| `--format`, `-f`        | n/a             | Output format    | `text`                                 |
| `--timeout`, `-t`       | `ASG_CA_TIMEOUT`| Request timeout  | `30s` (smsgate-ca only)                |

Output formats: `text` (human-readable, default), `json` (pretty printed), `raw` (one-line JSON), `table` (columnar output for lists).

## 📚 Documentation

- [CLI reference](https://docs.sms-gate.app/integration/cli/)
- [Central docs](https://docs.sms-gate.app/)
- [GitHub repository](https://github.com/android-sms-gateway/cli)

## 🤝 Contributing

Contributions are welcome. Open an issue or pull request; follow the repo's `make fmt` / `make lint` / `make test` checks and the [README style guide](https://docs.sms-gate.app/).

## ⚖️ License

Apache-2.0. See [LICENSE](LICENSE).

## 📜 Legal Notice

Android is a trademark of Google LLC.

<!-- Reference-style badge URLs: style=for-the-badge is mandatory -->
[contributors-shield]: https://img.shields.io/github/contributors/android-sms-gateway/cli?style=for-the-badge
[contributors-url]: https://github.com/android-sms-gateway/cli/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/android-sms-gateway/cli?style=for-the-badge
[forks-url]: https://github.com/android-sms-gateway/cli/network/members
[stars-shield]: https://img.shields.io/github/stars/android-sms-gateway/cli?style=for-the-badge
[stars-url]: https://github.com/android-sms-gateway/cli/stargazers
[issues-shield]: https://img.shields.io/github/issues/android-sms-gateway/cli?style=for-the-badge
[issues-url]: https://github.com/android-sms-gateway/cli/issues
[license-shield]: https://img.shields.io/github/license/android-sms-gateway/cli?style=for-the-badge
[license-url]: https://github.com/android-sms-gateway/cli/blob/main/LICENSE
