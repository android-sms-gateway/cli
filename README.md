<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>
<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![Apache-2.0 License][license-shield]][license-url]

<br />

<div align="center">
  <h3 align="center">SMSGate CLI</h3>

  <p align="center">
    A command-line interface for interacting with the SMSGate API
    <br />
    <a href="https://docs.sms-gate.app/integration/cli/"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/android-sms-gateway/cli/issues/new?labels=bug">Report Bug</a>
    ·
    <a href="https://github.com/android-sms-gateway/cli/issues/new?labels=enhancement">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
- [📱 About The Project](#-about-the-project)
  - [⚙️ Built With](#️-built-with)
- [💻 Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
    - [Option 1: Download from GitHub Releases](#option-1-download-from-github-releases)
    - [Option 2: Install using Go](#option-2-install-using-go)
    - [Option 3: Docker](#option-3-docker)
- [💻 Configuration](#-configuration)
  - [Available Options](#available-options)
  - [Output Formats](#output-formats)
- [💻 Usage](#-usage)
  - [Commands](#commands)
  - [Exit codes](#exit-codes)
  - [Examples](#examples)
    - [Sending messages](#sending-messages)
    - [Batch message sending](#batch-message-sending)
    - [Getting message status](#getting-message-status)
    - [Getting logs](#getting-logs)
    - [Output formats](#output-formats-1)
- [👥 Contributing](#-contributing)
- [©️ License](#️-license)
- [⚠️ Legal Notice](#️-legal-notice)


<!-- ABOUT THE PROJECT -->
## 📱 About The Project

There are two CLI tools in this repository: `smsgate` and `smsgate-ca`. The first one is for SMS Gateway for Android itself, and the second one is for the Certificate Authority.

This CLI provides a robust interface for:
- Sending and managing SMS messages (including batch operations from CSV and Excel files)
- Configuring webhook integrations
- Issuing certificates for private deployments

### ⚙️ Built With

- [![Go][Go-shield]][Go-url]
- [![Goreleaser][Goreleaser-shield]][Goreleaser-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## 💻 Getting Started

### Prerequisites

- Go 1.25+ (for building from source)
- Docker (optional, for containerized execution)

### Installation

#### Option 1: Download from GitHub Releases

1. Go to the [Releases page](https://github.com/android-sms-gateway/cli/releases/latest) of this repository.
2. Download the appropriate binary for your operating system and architecture.
3. Extract the archive to a directory of your choice.
4. Move the binary to a directory in your system's PATH.

#### Option 2: Install using Go

```bash
go install github.com/android-sms-gateway/cli/cmd/smsgate@latest
```

This will download, compile, and install the latest version of the CLI tool. Make sure your Go bin directory is in your system's PATH.

#### Option 3: Docker

```bash
docker run -it --rm --env-file .env ghcr.io/android-sms-gateway/cli \
  send --phone '+12025550123' 'Hello, Dr. Turk!'
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 💻 Configuration

The CLI can be configured using environment variables or command-line flags. You can also use a `.env` file in the working directory to set these variables.

### Available Options

| Option             | Env Var        | Description      | Default value                          |
| ------------------ | -------------- | ---------------- | -------------------------------------- |
| `--endpoint`, `-e` | `ASG_ENDPOINT` | The endpoint URL | `https://api.sms-gate.app/3rdparty/v1` |
| `--username`, `-u` | `ASG_USERNAME` | Your username    | **required**                           |
| `--password`, `-p` | `ASG_PASSWORD` | Your password    | **required**                           |
| `--format`, `-f`   | n/a            | Output format    | `text`                                 |

### Output Formats

The CLI supports four output formats:

1. `text`: Human-readable text output (default)
2. `json`: Pretty printed JSON-formatted output
3. `raw`: One-line JSON-formatted output
4. `table`: Tab-aligned columnar output for lists and sub-tables

Please note that when the exit code is not `0`, the error description is printed to stderr without any formatting.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 💻 Usage

```bash
smsgate [global options] command [command options] [arguments...]
```

### Commands

The CLI offers three main groups of commands:

- **Messages**: Commands for sending messages and checking their status, including batch operations from CSV and Excel files.
- **Webhooks**: Commands for managing webhooks, including creating, updating, and deleting them.
- **Logs**: Commands for retrieving logs for a specific time range.

For a complete list of available commands, you can:
- Run `smsgate help` or `smsgate --help` in your terminal.
- Visit the official documentation at [docs.sms-gate.app](https://docs.sms-gate.app/integration/cli/#commands).

### Exit codes

The CLI uses exit codes to indicate the outcome of operations:

- `0`: success
- `1`: invalid options or arguments
- `2`: server request error
- `3`: output formatting error

### Examples

For security reasons, it is recommended to pass credentials using environment variables or a `.env` file.

Credentials can also be passed via CLI options:

```bash
smsgate -u <username> -p <password> send --phones '+12025550123' 'Hello, Dr. Turk!'
```

#### Sending messages

The `send` command supports various options to customize message delivery:

```bash
# Send a simple text message
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
smsgate send --phones '+12025550123' --valid-until '2024-12-31T23:59:59Z' 'Message'

# Disable delivery report
smsgate send --phones '+12025550123' --delivery-report=false 'Message'

# Skip phone number validation
smsgate send --phones '+12025550123' --skip-phone-validation 'Message'

# Filter by device activity (devices active within last 12 hours)
smsgate send --phones '+12025550123' --device-active-within 12 'Message'

# Send data message (base64 encoded)
echo -n 'hello world' | base64
smsgate send --phones '+12025550123' --data --data-port 12345 'aGVsbG8gd29ybGQ='
```

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
| **Options**                 |                                                                                                                                                           |               |                         |
| `--ttl`                     | Time-to-live (TTL) for the message. Duration format (e.g., `1h30m`). If not provided, the message will not expire.<br>**Conflicts with `--valid-until`.** | empty         | `1h30m`                 |
| `--valid-until`             | The expiration date and time for the message. RFC3339 format (e.g., `2006-01-02T15:04:05Z07:00`).<br>**Conflicts with `--ttl`.**                          | empty         | `2024-12-31T23:59:59Z`  |
| `--skip-phone-validation`   | Skip phone number validation.                                                                                                                             | `false`       | `true`                  |
| `--device-active-within`    | Time window in hours for device activity filtering. `0` means no filtering.                                                                               | `0`           | `12`                    |

#### Batch message sending

The CLI supports sending messages in bulk from CSV and Excel files. The `batch send` command also supports the shared delivery/device options from `send` (such as `--device-id`, `--sim-number`, `--priority`, `--ttl`, `--valid-until`, `--delivery-report`, `--skip-phone-validation`, `--device-active-within`), allowing fine-grained control over each message in the batch.

**Supported file formats:**
- **CSV** (Comma-Separated Values)
- **XLSX** (Excel files)

**Basic usage:**

```bash
# Send messages from a CSV file
smsgate batch send contacts.csv --map phone=Phone,text=Message

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

**File format examples:**

**CSV with Headers:**
```csv
Phone,Message,Device,Priority
+12025550123,"Hello Dr. Turk!",device1,1
+12025550124,"Hello Dr. Smith!",device1,2
+12025550125,"Hello Dr. Jones!",device2,1
```

**CSV without Headers:**
```csv
+12025550123,"Hello Dr. Turk!",device1,1
+12025550124,"Hello Dr. Smith!",device1,2
+12025550125,"Hello Dr. Jones!",device2,1
```

**Excel Files:**
- Supports multiple sheets (use `--sheet` to specify)
- First row treated as headers by default
- Column mapping works the same as CSV

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

#### Getting message status

```bash
# Get the status of a sent message
smsgate status zXDYfTmTVf3iMd16zzdBj
```

#### Getting logs

The `logs` command retrieves logs for a specific time range. Dates should be in RFC3339 format (e.g., `2024-01-15T10:30:00Z`).

```bash
# Get logs for the last 24 hours (default)
smsgate logs

# Get logs for a specific time range
smsgate logs --from '2024-01-15T00:00:00Z' --to '2024-01-15T23:59:59Z'

# Get logs with custom time range and output format
smsgate --format json logs --from '2024-01-15T10:00:00+07:00' --to '2024-01-15T18:00:00+07:00'
```

#### Output formats

**Text**

```text
ID: zXDYfTmTVf3iMd16zzdBj
State: Pending
IsHashed: false
IsEncrypted: false
Recipients:
        +12025550123    Pending
        +12025550124    Pending
```

**JSON**

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

**Raw**

```json
{"id":"zXDYfTmTVf3iMd16zzdBj","state":"Pending","isHashed":false,"isEncrypted":false,"recipients":[{"phoneNumber":"+12025550123","state":"Pending"},{"phoneNumber":"+12025550124","state":"Pending"}],"states":{}}
```

**Table**

```text
ID                                    EVENT          URL                              DEVICE ID
123e4567-e89b-12d3-a456-426614174000  sms:received   https://example.com/webhook      dev-abc
def45678-e89b-12d3-a456-426614174000  sms:sent       https://example.com/other        
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 👥 Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## ©️ License

Distributed under the Apache-2.0 license. See [LICENSE](LICENSE) for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## ⚠️ Legal Notice

Android is a trademark of Google LLC.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
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
[Go-shield]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white
[Go-url]: https://go.dev
[Goreleaser-shield]: https://img.shields.io/badge/Goreleaser-FF007A?style=for-the-badge&logo=goreleaser&logoColor=white
[Goreleaser-url]: https://goreleaser.com
