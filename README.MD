# Gopher Status 🔍

Lightweight website monitoring service with gRPC API and React dashboard.

## Features

- 🚀 **gRPC API** with HTTP/JSON gateway
- 🔄 **Concurrent monitoring** using worker pool pattern
- 📊 **Real-time dashboard** (React + Vite)
- 🐳 **Single Docker container** for both backend and frontend
- ⚡ **In-memory storage** for fast operations

## Tech Stack

**Backend:**
- Go 1.25+
- gRPC + gRPC-Gateway
- Protocol Buffers

**Frontend:**
- React
- Vite

**Infrastructure:**
- Docker (multi-stage build)

## Quick Start

### Using Docker

```bash
# Build the image
docker build -t gopher-status .

# Run the container
docker run -p 8080:8080 -p 50051:50051 gopher-status
