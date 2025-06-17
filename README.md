# Reminder

Reminder is a simple web-based message board and reminder system written in Go. It serves a list of messages over HTTPS and provides an admin page to update them.

## Features

- Displays date and user-defined messages with optional colours
- Admin interface for editing messages using basic authentication
- Stores messages in `messages.json` and plays an alarm when they are saved
- Supports configurable TLS certificate and key paths
- Admin credentials configurable via `config.json`

## Requirements

- Go 1.21 or later
- A TLS certificate (`cert.pem`) and key (`key.pem`)
- Optional: an MP3 file (`Alarm05.mp3`) to play when saving messages
- A `config.json` file in the root directory for admin credentials.

## Configuration

Admin credentials are now managed via a `config.json` file located in the root of the repository. Create this file with the following structure:

```json
{
  "adminUsername": "YOUR_USERNAME",
  "adminPassword": "YOUR_PASSWORD"
}
```

Replace `YOUR_USERNAME` and `YOUR_PASSWORD` with your desired administrator credentials.

**Important:** If you are using real credentials, ensure that `config.json` is not committed to your version control system (e.g., add it to your `.gitignore` file).

## Running

```bash
go run main.go
```

The server will look for `cert.pem` and `key.pem` in the current directory by default. You can also specify paths using command-line flags (`-cert` and `-key`) or environment variables (`CERT_PATH` and `KEY_PATH`). If these files are not found or are invalid, the server will fail to start. For local development, you can generate self-signed certificates.

The application now runs on port `8443` by default. Access it at `https://localhost:8443`.

The admin page is available at `/admin`. Credentials for accessing the admin page are configured in `config.json` (see the Configuration section above).

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
