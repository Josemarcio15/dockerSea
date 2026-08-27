# DockSea

DockSea is a modern, agentless Docker manager for local environments and multiple VPSs. It provides a single interface for monitoring and operating Docker hosts without installing a management agent or exposing the Docker API to the public internet.

> **Status:** Alpha — currently tested only as a native desktop application.

[![Wails 3](https://img.shields.io/badge/Built%20with-Wails%20v3-007acc.svg)](https://v3.wails.io/)
[![Go](https://img.shields.io/badge/Backend-Go%201.25-00ADD8.svg)](https://go.dev/)
[![Svelte](https://img.shields.io/badge/Frontend-Svelte%205-FF3E00.svg)](https://svelte.dev/)
[![Tailwind CSS](https://img.shields.io/badge/Styling-Tailwind%20v4-38B2AC.svg)](https://tailwindcss.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Overview

DockSea connects to Docker hosts through the local Unix socket or an SSH tunnel to the remote host's Docker socket. Remote servers do not require DockSea, a Docker agent, or an additional listening port. This makes it suitable for privately managed VPS infrastructure running on providers such as AWS, Hetzner, DigitalOcean, Contabo, Linode, or Oracle Cloud.

The application is designed to run as a native desktop application.

## Features

- Container management: list, start, stop, restart, remove, inspect, and stream logs from containers.
- Real-time Docker events: react to daemon events without periodic polling.
- Dashboard and metrics: view CPU, memory, swap, disk, and uptime information for the active host.
- Nginx tools: edit "sites-available" and "sites-enabled" configurations, validate them with "nginx -t", reload Nginx, and inspect logs.
- Port scanner: inspect listening TCP and UDP ports, including associated processes and PIDs.
- Local and remote connections: manage Docker Desktop or local Unix sockets alongside SSH-based VPS connections.
- SSH authentication: support for private keys, passphrases, passwords, custom SSH ports, and optional sudo access.
- Internationalization: English and Brazilian Portuguese translations.
- Light and dark themes: a responsive interface for desktop usage.

## How remote connections work

1. DockSea opens an SSH connection to the configured VPS.
2. Commands and Docker API requests are sent through the encrypted SSH channel.
3. The remote Docker Unix socket is used directly; no public Docker API port is required.

For production use, create a dedicated Linux user with the minimum permissions required to manage Docker. Protect SSH with key-based authentication, disable password authentication where possible, and restrict access with your firewall or VPN.

## Requirements

- Go 1.25 or newer
- Node.js 20 or newer and npm
- Wails 3 CLI
- [Task](https://taskfile.dev/), recommended for repository commands

Install the Wails CLI and Task if they are not already available:

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
go install github.com/go-task/task/v3/cmd/task@latest
```

## Installation and development

Clone the repository and install frontend dependencies:

```bash
git clone git@github.com:Josemarcio15/dockSea.git
cd dockSea
cd frontend
npm install
cd ..
```

Start the Wails development environment:

```bash
task dev
```

Run frontend checks independently:

```bash
npm run check --prefix frontend
```

Run Go unit tests:

```bash
go test ./...
```

Build the native desktop application:

```bash
task build
```

The generated application is written to "bin/" for the current platform.

## Configuration and data

DockSea uses SQLite for local application data, including saved connection configuration. The database is created by the application at runtime.

Do not commit private keys, passwords, ".env" files, database files, or production configuration to the repository. Use SSH keys, a secrets manager, or your deployment platform's secret mechanism.

## Security checklist

- Use SSH key authentication and a dedicated management user.
- Restrict inbound traffic with a firewall, VPN, or private network.
- Rotate credentials and remove unused VPS connections regularly.

## Useful commands

```bash
task dev # Start Wails development mode
task build # Build the native desktop application
task package # Package the native application for the current platform
```

## License

DockSea is distributed under the MIT License. See "LICENSE" (LICENSE) for the full license text.
