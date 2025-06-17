# Reminder

Reminder is a simple web-based message board and reminder system written in Go. It serves a list of messages over HTTPS and provides an admin page to update them.

## Features

- Displays date and user-defined messages with optional colours
- Admin interface for editing messages using basic authentication
- Stores messages in `messages.json` and plays an alarm when they are saved
- Supports configurable TLS certificate and key paths

## Requirements

- Go 1.21 or later
- A TLS certificate (`cert.pem`) and key (`key.pem`)
- Optional: an MP3 file (`Alarm05.mp3`) to play when saving messages

## Running

```bash
go run main.go -cert path/to/cert.pem -key path/to/key.pem
```

You can also set `CERT_PATH` and `KEY_PATH` environment variables instead of using command-line flags.

The admin page is available at `/admin` with the default credentials `hasting` / `holidays`.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
