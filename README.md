# SysTrace Agent
A simple Go agent for collecting system data and communicating with a master server over WebSocket.
## What the Agent Does
- Collects hardware and system information (for example CPU, RAM, hostname, OS)
- Retrieves location data (GPS)
- Connects to the master server
- Sends update messages regularly
- Receives server messages (`test`, `command`, `config`)
## Requirements
- Windows 10/11 (x64)
- Go 1.26 (according to `go.mod`)
- Docker Desktop installed
- Docker Engine running (`docker info` must succeed)
- A `.env` file in the project directory
## Installer Prerequisites
The installer expects Docker to be available before setup starts:
- `docker.exe` must be available in `PATH`
- Docker Desktop must be started
- Docker Engine must be reachable (`docker info` exits with code 0)
If this is not fulfilled, the installer aborts early with a message.
## Configuration
Create a `.env` file in the project root:
```env
MASTER_SERVER_URL=http://localhost:8080
GEOLOCATION_API_KEY=your_api_key
```
Notes:
- `MASTER_SERVER_URL` is used for the WebSocket connection.
- `GEOLOCATION_API_KEY` is required for geolocation.
## Run
Install dependencies and start the agent:
```bash
go mod tidy
go run .
```
## Project Structure (Short)
- `main.go`: Application entry point
- `internal/agent`: Agent logic (collect, send, receive)
- `internal/collector`: Data collectors (CPU, RAM, hardware, GPS)
- `internal/handler`: Handles incoming server messages
- `internal/transport`: Connection and `.env` handling
- `internal/data`: Data models for transfer and storage
## Build
```bash
go build -o SysTrace_Agent.exe .
```
