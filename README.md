# Hotel-IO Backend

## Overview

This repository contains the backend implementation for Hotel-IO, an action-packed multiplayer.

## Features

- Real-time multiplayer combat system
- Player authentication and profiles
- Character selection and customization
- Matchmaking system
- Leaderboards and ranking system
- Battle statistics tracking
- Real-time game state synchronization
- Tournament management
- Player progression system
- Combat logging and replay system

## Tech Stack

- Go (1.21+)
- Gin (Web Framework)
- GORM (ORM)
- MongoDB (for player data and statistics)
- Gorilla WebSocket (for real-time multiplayer)
- JWT-Go for authentication
- RESTful API architecture
- Redis (for game state caching)

## Prerequisites

- Go 1.21 or higher
- Redis 6 or higher
- Make (optional, for using Makefile commands)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/yourusername/hotel-io-backend.git
cd hotel-io-backend
```

2. Install dependencies:

```bash
go mod download
```

3. Start the development server:

```bash
go run cmd/main.go
```

## Project Structure

```
├── cmd/
│   └── main.go          # Application entry point
├── internal/
│   ├── database/        # Database configurations and connections
│   ├── repositories/    # Data access layer
│   ├── models/          # Data models and entities
│   ├── v1/              # API version 1 handlers and routes
│   └── game/            # Game logic and mechanics
├── go.mod               # Go modules file
├── go.sum               # Go modules checksums
```
