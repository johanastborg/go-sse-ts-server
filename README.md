# Go SSE Time Series Server

A high-performance Go server for streaming simulated time-series data over HTTP/2 using Server-Sent Events (SSE) with a fan-out pattern.

## Features

- **50Hz Streaming**: Generates a simulated sine wave with added noise at 50Hz (20ms intervals).
- **Fan-out Architecture**: A central hub manages multiple concurrent client connections, broadcasting data from a single producer.
- **SSE/HTTP2**: Uses Server-Sent Events for efficient real-time data delivery over HTTP/2.
- **Bazel & Gazelle**: Fully managed build system with automatic `BUILD` file generation for Go packages.

## Prerequisites

- [Bazel](https://bazel.build/install) (version 7.x or later recommended)
- [Go](https://go.dev/doc/install) (version 1.23.0+)

## Getting Started

### 1. Clone the repository
```bash
git clone github.com/johanastborg/go-sse-ts-server
cd go-sse-ts-server
```

### 2. Update BUILD files (Gazelle)
If you add new packages or change dependencies, run:
```bash
bazel run //:gazelle
```

### 3. Build the Project
```bash
bazel build //...
```

### 4. Run the Server
```bash
bazel run //cmd/server:server
```
The server will start on `http://localhost:8080`.

## Testing the Stream

Open a new terminal and use `curl` to watch the real-time data feed:

```bash
curl -N -H "Accept: text/event-stream" http://localhost:8080/stream
```

You should see an output like this:
```text
data: {"timestamp":1772957403107,"value":5.3295302516544005}
data: {"timestamp":1772957403127,"value":5.686293395940926}
data: {"timestamp":1772957403147,"value":7.21783367995546}
...
```

## Project Structure

- `cmd/server/`: Main application entry point.
- `internal/feed/`: Sine wave data producer (50Hz).
- `internal/hub/`: Generic fan-out hub for broadcasting data.
- `internal/sse/`: HTTP handler for SSE protocol.
- `MODULE.bazel`: Bazel module definitions and dependencies.
- `BUILD.bazel`: Root build file with Gazelle configuration.
