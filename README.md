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
  <h3 align="center">SMS Gateway for Android™ CLI</h3>

  <p align="center">
    A command-line interface for interacting with the SMS Gateway for Android API
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
    - [Output formats](#output-formats-1)
- [👥 Contributing](#-contributing)
- [©️ License](#️-license)
- [⚠️ Legal Notice](#️-legal-notice)


<!-- ABOUT THE PROJECT -->
## 📱 About The Project

There are two CLI tools in this repository: `smsgate` and `smsgate-ca`. The first one is for SMS Gateway for Android itself, and the second one is for the Certificate Authority.

This CLI provides a robust interface for:
- Sending and managing SMS messages
- Configuring webhook integrations
- Issuing certificates for private deployments

### ⚙️ Built With

- [![Go][Go-shield]][Go-url]
- [![Goreleaser][Goreleaser-shield]][Goreleaser-url]

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## 💻 Getting Started

### Prerequisites

- Go 1.23+ (for building from source)
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

The CLI supports three output formats:

1. `text`: Human-readable text output (default)
2. `json`: Pretty printed JSON-formatted output
3. `raw`: One-line JSON-formatted output

Please note that when the exit code is not `0`, the error description is printed to stderr without any formatting.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## 💻 Usage

```bash
smsgate [global options] command [command options] [arguments...]
```

### Commands

The CLI offers two main groups of commands:

- **Messages**: Commands for sending messages and checking their status.
- **Webhooks**: Commands for managing webhooks, including creating, updating, and deleting them.

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

```bash
# Send a message
smsgate send --phones '+12025550123' 'Hello, Dr. Turk!'

# Send a message to multiple numbers
smsgate send --phones '+12025550123' --phones '+12025550124' 'Hello, doctors!'
# or
smsgate send --phones '+12025550123,+12025550124' 'Hello, doctors!'

# Get the status of a sent message
smsgate status zXDYfTmTVf3iMd16zzdBj
```

Credentials can also be passed via CLI options:

```bash
smsgate -u <username> -p <password> send --phones '+12025550123' 'Hello, Dr. Turk!'
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
