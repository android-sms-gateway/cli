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
    - [Available Options](#available-options)
    - [Output Formats](#output-formats)
  - [Usage](#usage)
    - [Commands](#commands)
      - [Send a message](#send-a-message)
      - [Get the status of a sent message](#get-the-status-of-a-sent-message)
  - [Usage Examples](#usage-examples)
    - [Output formats](#output-formats-1)
      - [Text](#text)
      - [JSON](#json)
      - [Raw](#raw)
  - [Exit codes](#exit-codes)
  - [Support](#support)
  - [Contributing](#contributing)
  - [License](#license)
  - [Legal Notice](#legal-notice)

## Installation

You can install the SMS Gateway CLI in two ways:

### Option 1: Download from GitHub Releases

1. Go to the [Releases page](https://github.com/android-sms-gateway/cli/releases/latest) of this repository.
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

The CLI can be configured using environment variables or command-line flags. You can also use a `.env` file in the working directory to set these variables.

### Available Options

| Option             | Env Var        | Description      | Default value                          |
| ------------------ | -------------- | ---------------- | -------------------------------------- |
| `--endpoint`, `-e` | `ASG_ENDPOINT` | The endpoint URL | `https://api.sms-gate.app/3rdparty/v1` |
| `--username`, `-u` | `ASG_USERNAME` | Your username    | **required**                           |
| `--password`, `-p` | `ASG_PASSWORD` | Your password    | **required**                           |
| `--format`, `-f`   | n/a            | Output format    | `text`                                 |

### Output Formats

The CLI supports three output formats:

1. `text`: Human-readable text output (default)
2. `json`: Pretty printed JSON-formatted output
3. `raw`: One-line JSON-formatted output

Please note that when the exit code is not `0`, the error description is printed to stderr without any formatting.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage

```bash
smsgate [global options] command [command options] [arguments...]
```

### Commands

The CLI supports the following commands:

- `send` - send a message with single or multiple recipients
- `status` - get the status of a sent message by message ID

#### Send a message

Syntax:
```bash
smsgate send [options] 'Message content'
```

| Option                      | Description                                                                                | Default value | Example                 |
| --------------------------- | ------------------------------------------------------------------------------------------ | ------------- | ----------------------- |
| `--id`                      | Message ID, will be generated if not provided                                              | empty         | `zXDYfTmTVf3iMd16zzdBj` |
| `--phone`, `--phones`, `-p` | Phone number, can be used multiple times or with comma-separated values                    | **required**  | `+19162255887`          |
| `--sim`                     | SIM card slot number, if empty, the default SIM card will be used                          | empty         | `2`                     |
| `--ttl`                     | Time-to-live (TTL), if empty, the message will not expire<br>Conflicts with `--validUntil` | empty         | `1h30m`                 |
| `--validUntil`              | Valid until, if empty, the message will not expire<br>Conflicts with `--ttl`               | empty         | `2024-12-31 23:59:59`   |

#### Get the status of a sent message

Syntax:
```bash
smsgate status 'Message ID'
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Usage Examples

For security reasons, it is recommended to pass credentials using environment variables or a `.env` file.

```bash
# Send a message
smsgate send --phone '+19162255887' 'Hello, Dr. Turk!'

# Send a message to multiple numbers
smsgate send --phone '+19162255887' --phone '+19162255888' 'Hello, doctors!'
# or
smsgate send --phones '+19162255887,+19162255888' 'Hello, doctors!'

# Get the status of a sent message
smsgate status zXDYfTmTVf3iMd16zzdBj
```

Credentials can also be passed via CLI options:

```bash
# Pass credentials by options
smsgate send -u <username> -p <password> --phone '+19162255887' 'Hello, Dr. Turk!'
```

If you prefer not to install the CLI tool locally, you can use Docker to run it:

```bash
docker run -it --rm --env-file .env ghcr.io/android-sms-gateway/cli send --phone '+19162255887' 'Hello, Dr. Turk!'
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

## Exit codes

The CLI uses exit codes to indicate the outcome of operations:

- `0`: success
- `1`: invalid options or arguments
- `2`: server request error
- `3`: output formatting error

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
