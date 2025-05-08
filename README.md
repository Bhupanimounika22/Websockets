# Go WebSocket Server

A simple WebSocket server implementation in Go using the Gorilla WebSocket package.

## Prerequisites

- Go 1.18 or later

## Installation

1. Clone this repository
2. Install dependencies:

```bash
go mod download
```

## Running the Server

Start the WebSocket server:

```bash
go run main.go
```

The server will start on port 8080. You can access the chat interface by opening http://localhost:8080 in your browser.

## Features

- Real-time chat functionality
- Multiple client support
- Automatic reconnection on disconnection

## How It Works

The server uses the following components:

- Gorilla WebSocket package for WebSocket implementation
- Client struct to represent each connected WebSocket client
- ClientManager to handle multiple client connections
- Goroutines for concurrently handling multiple clients 