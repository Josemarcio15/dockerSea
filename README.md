# DockSea

DockSea is a modern, agentless Docker manager for local environments and multiple VPSs. It provides a single interface for monitoring and operating Docker hosts without installing a management agent or exposing the Docker API to the public internet.

[![Wails 3](https://img.shields.io/badge/Built%20with-Wails%20v3-007acc.svg)](https://v3.wails.io/)
[![Go](https://img.shields.io/badge/Backend-Go%201.25-00ADD8.svg)](https://go.dev/)
[![Svelte](https://img.shields.io/badge/Frontend-Svelte%205-FF3E00.svg)](https://svelte.dev/)
[![Tailwind CSS](https://img.shields.io/badge/Styling-Tailwind%20v4-38B2AC.svg)](https://tailwindcss.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Overview

DockSea connects to Docker hosts through the local Unix socket or an SSH tunnel to the remote host's Docker socket. Remote servers do not require DockSea, a Docker agent, or an additional listening port. This makes it suitable for privately managed VPS infrastructure running on providers such as AWS, Hetzner, DigitalOcean, Contabo, Linode, or Oracle Cloud.

The application can be used as a native desktop application or built in server mode and deployed as an HTTP service with Docker.

## Features

- **Container management:** list, start, stop, restart, remove, inspect, and stream logs from containers.
- **Real-time Docker events:** react to daemon events without periodic polling.
- **Dashboard and metrics:** view CPU, memory, swap, disk, and uptime information for the active host.
- **Nginx tools:** edit `sites-available` and `sites-enabled` configurations, validate them with `nginx -t`, reload Nginx, and inspect logs.
- **Port scanner:** inspect listening TCP and UDP ports, including associated processes and PIDs.
- **Local and remote connections:** manage Docker Desktop or local Unix sockets alongside SSH-based VPS connections.
- **SSH authentication:** support for private keys, passphrases, passwords, custom SSH ports, and optional sudo access.
- **Internationalization:** English and Brazilian Portuguese translations.
- **Light and dark themes:** a responsive interface for desktop and browser-based usage.

## How Remote Connections Work

1. DockSea opens an SSH connection to the configured VPS.
2. Commands and Docker API requests are sent through the encrypted SSH channel.
3. The remote Docker Unix socket is used directly; no public Docker API port is required.

For production use, create a dedicated Linux user with the minimum permissions required to manage Docker. Protect SSH with key-based authentication, disable password authentication where possible, and restrict access with your firewall or VPN.


## Requirements

- Go 1.25 or newer
- Node.js 20 or newer and npm
- Wails 3 CLI
- Task, recommended for the repository build commands
- Docker, when building the server image or cross-compiling

Install the Wails CLI and Task if they are not already available:

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
go install github.com/go-task/task/v3/cmd/task@latest
```

## Local Development

Clone the repository and install frontend dependencies:

```bash
git clone git@github.com:Josemarcio15/dockerSea.git
cd dockerSea
cd frontend && npm install && cd ..
```

Start the Wails development environment:

```bash
task dev
```

Run frontend checks independently:

```bash
cd frontend
npm run check
```

Build the native desktop application:

```bash
task build
```

The generated application is written to `bin/` for the current platform.

## VPS Deployment with Docker

Server mode exposes DockSea as an HTTP service on container port `8080`. The included Dockerfile builds the frontend and Go backend in a multi-stage build and uses a minimal distroless runtime image by default.

### Build the image

Run this from the repository root on a machine with Docker:

```bash
task build:docker TAG=docksea:latest
```

### Run the service

```bash
docker run -d \
  --name docksea \
  --restart unless-stopped \
  -p 8080:8080 \
  docksea:latest
```

To use another host port, change only the host-side mapping:

```bash
docker run -d --name docksea --restart unless-stopped -p 3000:8080 docksea:latest
```

The server binds to `0.0.0.0` inside the container. This can be overridden with `WAILS_SERVER_HOST` when required by the deployment environment.

### Run locally with Task

The following command builds the image and starts it on port `8080`:

```bash
task run:docker
```

To use a different host port:

```bash
task run:docker PORT=3000
```

### Reverse proxy and TLS

For a public deployment, place DockSea behind Nginx, Caddy, or another TLS-terminating reverse proxy. Expose only ports `80` and `443` publicly, keep port `8080` bound to localhost or an internal Docker network, and forward requests to `http://127.0.0.1:8080`.

Example Nginx location:

```nginx
server {
    listen 443 ssl http2;
    server_name docksea.example.com;

    # Configure ssl_certificate and ssl_certificate_key here.

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

## Configuration and Data

DockSea uses SQLite for local application data, including saved connection configuration. The database is created by the application at runtime. Before deploying to production, identify the database location for your target runtime and mount a persistent volume if connection data must survive container replacement.

Do not commit private keys, passwords, `.env` files, database files, or production configuration to the repository. Use SSH keys, a secrets manager, or your deployment platform's secret mechanism.

## Security Checklist

- Use SSH key authentication and a dedicated management user.
- Avoid exposing the Docker daemon TCP socket to the internet.
- Put the HTTP service behind HTTPS and an access-control layer when publicly reachable.
- Restrict inbound traffic with a firewall, VPN, or private network.
- Use a persistent, protected volume for application data when required.
- Rotate credentials and remove unused VPS connections regularly.
- Review container and host logs after deployment.

## Useful Commands

```bash
task dev             # Start Wails development mode
task build           # Build the native desktop application
task build:server    # Build a server-mode binary
task build:docker    # Build the production Docker image
task run:docker      # Build and run the Docker image
task package         # Package the native application for the current platform
```

## License

DockSea is distributed under the MIT License. See [LICENSE](LICENSE) for the full license text.

