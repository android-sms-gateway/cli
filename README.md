<a name="readme-top"></a>

# SMS Gateway for Android™ CLI

A command-line interface for working with SMS Gateway for Android.

## Table of Contents

- [SMS Gateway for Android™ CLI](#sms-gateway-for-android-cli)
  - [Table of Contents](#table-of-contents)
  - [Installation](#installation)
    - [Option 1: Download from GitHub Releases](#option-1-download-from-github-releases)
    - [Option 2: Install using Go](#option-2-install-using-go)
  - [Configuration](#configuration)
    - [Environment Variables](#environment-variables)
    - [Command-line Flags](#command-line-flags)
      - [Output Formats](#output-formats)
  - [Usage](#usage)
    - [Commands](#commands)
    - [Send Message](#send-message)
      - [Get Message Status](#get-message-status)
  - [Examples](#examples)
    - [Output formats](#output-formats-1)
      - [Text](#text)
      - [JSON](#json)
      - [Raw](#raw)
  - [Support](#support)
  - [Contributing](#contributing)
  - [License](#license)
  - [Legal Notice](#legal-notice)

## Installation

You can install the SMS Gateway CLI in two ways:

### Option 1: Download from GitHub Releases

1. Go to the [Releases page](https://github.com/android-sms-gateway/cli/releases) of this repository.
2. Download the appropriate binary for your operating system and architecture.
3. Rename the downloaded file to `smsgate` (or `smsgate.exe` for Windows).
4. Move the binary to a directory in your system's PATH.

For example, on Linux or macOS:

```bash
mv /path/to/downloaded/binary /usr/local/bin/smsgate
chmod +x /usr/local/bin/smsgate
```

### Option 2: Install using Go

If you have Go installed on your system (version 1.23 or later), you can use the go install command:

```bash
go install github.com/android-sms-gateway/cli/cmd/smsgate@latest
```

This will download, compile, and install the latest version of the CLI tool. Make sure your Go bin directory is in your system's PATH.

After installation, you can run the CLI tool using the `smsgate` command.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Configuration

The CLI can be configured using environment variables or command-line flags. You can also use a `.env` file to set these variables.

### Environment Variables

All environment variables are prefixed with `ASG_`:

- `ASG_ENDPOINT`: The endpoint URL (default: `https://api.sms-gate.app/3rdparty/v1`)
- `ASG_USERNAME`: Your username (required)
- `ASG_PASSWORD`: Your password (required)

### Command-line Flags

- `--endpoint`, `-e`: Endpoint URL
- `--username`, `-u`: Username
- `--password`, `-p`: Password
- `--format`, `-f`: Output format (supported: text, json, raw; default: text)

#### Output Formats

The CLI supports three output formats:

1. `text`: Human-readable text output (default)
2. `json`: Pretty printed JSON-formatted output
3. `raw`: One-line JSON-formatted output

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

```bash
smsgate [global options] command [command options] [arguments...]
```

### Commands

### Send Message

Send an SMS message.

```bash
smsgate send [options] Message content
```

Options:
- `--id value`: Message ID (optional, generated if not specified)
- `--phones value, -p value, --phone value`: Phone numbers (can be specified multiple times or comma-separated)
- `--sim value`: SIM card index (1-3) (default: 0)
- `--ttl value`: Time to live as duration string like "1h30m" (default: unlimited)
- `--validUntil value`: Valid until (format: YYYY-MM-DD HH:MM:SS in local timezone)

#### Get Message Status

Retrieve the status of a sent message.

```bash
smsgate status [options] Message ID
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Examples

```bash
# Send a message
smsgate send --phone '+19162255887' 'Hello, Dr. Turk!'

# Send a message to multiple numbers
smsgate send --phone '+19162255887' --phone '+19162255888' 'Hello, Dr. Turk!'
# or
smsgate send --phones '+19162255887,+19162255888' 'Hello, Dr. Turk!'

# Get the status of a sent message
smsgate status zXDYfTmTVf3iMd16zzdBj
```

### Output formats

#### Text

```
ID: zXDYfTmTVf3iMd16zzdBj
State: Pending
IsHashed: false
IsEncrypted: false
Recipients: 
        +19162255887    Pending
        +19162255888    Pending
```

#### JSON

```json
{
  "id": "zXDYfTmTVf3iMd16zzdBj",
  "state": "Pending",
  "isHashed": false,
  "isEncrypted": false,
  "recipients": [
    {
      "phoneNumber": "+19162255887",
      "state": "Pending"
    },
    {
      "phoneNumber": "+19162255888",
      "state": "Pending"
    }
  ],
  "states": {}
}
```

#### Raw

```json
{"id":"zXDYfTmTVf3iMd16zzdBj","state":"Pending","isHashed":false,"isEncrypted":false,"recipients":[{"phoneNumber":"+19162255887","state":"Pending"},{"phoneNumber":"+19162255888","state":"Pending"}],"states":{}}
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Support

For support, please contact support@sms-gate.app

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## License

Distributed under the Apache-2.0 license. See [LICENSE](LICENSE) for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Legal Notice

Android is a trademark of Google LLC.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
