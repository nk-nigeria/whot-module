# Whot Game Module

A card game module for Nakama server using the generic bot management system.

## Features

- **Whot Card Game** - Traditional Nigerian card game
- **Bot Integration** - Uses cgp-common bot system for AI players
- **Multiplayer Support** - Up to 4 players per match
- **Configurable Rules** - Bot behavior configurable via database

## Bot Integration

This module uses the generic bot management system from `cgp-common/bot`:

- **WhotBotIntegration** - Implements BotIntegration interface for Whot game
- **Database Configuration** - Bot rules stored in database tables
- **Delayed Bot Joining** - Bots join with random delays for natural behavior

## Usage

The bot system is integrated through `WhotBotIntegration` in `usecase/service/whot_bot_integration.go`.

For detailed bot system documentation, see `cgp-common/bot/README.md`.