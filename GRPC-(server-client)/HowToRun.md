# How To Run — gRPC Server + Client

This guide explains how to run the **server** and **client** folders of this
project. The code is written in the **Go** language (the `go.mod`-based Go
toolchain) and uses **gRPC** over **TLS**.

```
GRPC-(server-client)/
├── client/
│   ├── go.mod / go.sum
│   ├── main.go
│   ├── command.txt
│   └── proto/
│       ├── main.proto
│       └── gen/          # generated code (protoc output)
└── server/
    ├── go.mod / go.sum
    ├── main.go
    ├── command.txt
    ├── cert.conf / cert.pem / key.pem
    └── proto/
        ├── main.proto
        ├── greeter.proto
        └── gen/          # generated code (protoc output)
```

---

## 1. Prerequisites

Make sure the following tools are installed and available on your `PATH`:

| Tool            | Purpose                                    | Check with          |
|-----------------|--------------------------------------------|---------------------|
| **Go**          | The language toolchain (`go` command)      | `go version`        |
| **protoc**      | Protocol Buffers compiler                  | `protoc --version`  |
| **OpenSSL**     | Generate local TLS certificates            | `openssl version`   |

You also need the two gRPC code‑generation plugins for Go installed
(the project uses these to generate the `proto/gen/*.pb.go` files):

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

> Make sure your `GOPATH/bin` is added to `PATH` so `protoc` can find the
> `protoc-gen-go` and `protoc-gen-go-grpc` binaries.

---

## 2. Fetch Dependencies

Dependencies are already declared in `go.mod` (and pinned in `go.sum`).
To make sure the gRPC dependency is installed, run this **inside each folder**
(`server` and `client`):

```bash
go get google.golang.org/grpc
```

This must be run from the folder that contains the `go.mod` you want to
update (i.e. once in `./server` and once in `./client`).

---

## 3. One‑Time Code Generation

The generated files under `proto/gen/` come from the `.proto` definitions.
If they are missing or you changed a `.proto` file, regenerate them.

### Server (generates both `main.proto` and `greeter.proto`)

```bash
cd server
protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto proto/greeter.proto
```

### Client

```bash
cd client
protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto
```

---

## 4. Server — Generate TLS Certificates (one‑time)

The server serves gRPC over TLS and loads `cert.pem` and `key.pem` at startup.
They are already present, but if you need to regenerate them, run this inside
the `server` folder:

```bash
cd server
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config cert.conf
```

> `cert.conf` is already configured for local development (CN = `localhost`,
> SAN for `localhost` / `127.0.0.1`), so you do not need to edit it.

---

## 5. Run the Server

From the `server` folder:

```bash
cd server
go run .
```

You should see:

```
Server is running on port :50051
```

The server starts a TCP listener on port **50051**, loads the TLS credentials,
and registers two services:

- **Calculator** — `Add(a, b)` → sum (also demonstrates gRPC metadata /
  headers / trailers handling)
- **Greeter** — `Greet(name)` → `"Hello <name>, nice to meet you,"`

> Leave the server running in its own terminal window.

---

## 6. Run the Client

In a **separate terminal**, from the `client` folder:

```bash
cd client
go run .
```

The client connects to `localhost:50051` **using insecure/plaintext
credentials** (no certificate validation, since the server uses a
self‑signed certificate) and exercises all three streaming RPCs from the
`Calculator` service:

1. **GenerateFibonacci** (server‑side streaming) — sends `N = 10` and prints
   each Fibonacci number as it is received.
2. **SendNumbers** (client‑side streaming) — sends numbers `0..8` one per
   second, then closes the stream and prints the server‑computed sum.
3. **Chat** (bidirectional streaming) — concurrently sends `"Hi"`,
   `"How are you?"`, `"Bye"` while printing whatever the server sends back.

Example client output:

```
Fibonacci number: 0
Fibonacci number: 1
...
End of Stream
Server resp after stream: <sum>
Server:  <message>
```

---

## 7. Order of Operations (Quick Summary)

```bash
# 1. Install gRPC proto plugins (once)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 2. Install gRPC dependency (server + client)
cd server && go get google.golang.org/grpc
cd ../client && go get google.golang.org/grpc

# 3. Regenerate code (only if proto/gen is missing / proto changed)
cd ../server && protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto proto/greeter.proto
cd ../client && protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto

# 4. (Re)generate TLS certs if needed (server only)
cd ../server && openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout key.pem -out cert.pem -config cert.conf

# 5. Run server, then client (two terminals)
cd server && go run .
cd client && go run .
```

---

## Troubleshooting

- **`protoc` can’t find `protoc-gen-go`** → add your `GOPATH/bin` to `PATH`.
- **Client can’t connect** → make sure the server terminal still says
  `Server is running on port :50051` and that the client host/port
  (`localhost:50051`) matches the server.
- **TLS errors** → ensure `cert.pem` / `key.pem` exist in `server/` (see
  step 4).
- **`go get` fails** → run it from the folder that contains the `go.mod`
  for the module you are updating.
- **Missing generated files** → re-run the `protoc` commands in step 3.