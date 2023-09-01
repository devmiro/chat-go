# Chat Application in Go

This is a simple chat application implemented in Go.

## Setup

- Clone the repository.
- Install the required dependencies listed in `requirements.txt`.
- Run the application using `go run cmd/main.go`.

## Features

- Allows registered users to log in and chat with others.
- Supports posting stock commands in the chatroom.
- Includes a decoupled bot to fetch stock quotes from an API.
- Orders chat messages by timestamp and shows only the last 50 messages.
- Provides unit tests for selected functionality.

## Bonus Features

- Supports multiple chatrooms.
- Handles messages not understood by the bot.

## Considerations

- This project is backend-focused; keep the frontend simple.
- Ensure the security of confidential information.
- Monitor resource consumption.
- Use Git for version control.
- Feel free to use small helper libraries.
