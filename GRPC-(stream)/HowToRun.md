# How To Run — gRPC Streaming Server + Client

This guide explains how to run the **server** and **client** folders of this
project. The code is written in the **Go** language (the `go.mod`-based Go
toolchain) and uses **gRPC streaming** (no TLS — plaintext on localhost).

```
GRPC-(stream)/
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
    └── proto/
        ├── main.proto
        └── gen/          # generated code (protoc output)
```

---

## 1. Prerequisites

Make sure the following tools are installed and available on your `PATH`:

| Tool            | Purpose                                    | Check with         |
|-----------------|--------------------------------------------|--------------------|
| **Go**          | The language toolchain (`go` command)      | `go version`       |
| **protoc**      | Protocol Buffers compiler                  | `protoc --version` |

You also need the two gRPC code‑generation plugins for Go installed
(the project uses these to generate the `proto/gen/*.pb.go` files):

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

> Make sure your `GOPATH/bin` is added to `PATH` so `protoc` can find the
> `protoc-gen-go` and `protoc-gen-go-grpc` binaries. If `protoc` can’t find
> them, export the path first:
> `export PATH="$PATH:$(go env GOPATH)/bin"`

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
If they are missing or you changed the `.proto` file, regenerate it.

Run this **inside both `server` and `client`** folders (each project has its
own `proto/main.proto`):

```bash
cd server
protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto
```

```bash
cd client
protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto
```

> The stream server uses **plaintext** gRPC (`grpc.NewServer()` with no TLS
> credentials), so there is no need to generate certificates for this
> project.

---

## 4. Run the Server

From the `server` folder:

```bash
cd server
go run .
```

The server starts a TCP listener on port **50051** (no TLS) and registers the
`Calculator` service, which implements **all three streaming RPCs** from
`main.proto`:

- **GenerateFibonacci** (server‑side streaming) — given `N`, streams `N`
  Fibonacci numbers back to the client, one per second.
- **SendNumbers** (client‑side streaming) — receives numbers from the client,
  accumulates them, and once the client closes the stream returns the total
  sum.
- **Chat** (bidirectional streaming) — echoes each client message back,
  prefixed with `"Server replied: "`, while printing what the client said.

> Leave the server running in its own terminal window.

---

## 5. Run the Client

In a **separate terminal**, from the `client` folder:

```bash
cd client
go run .
```

The client connects to `localhost:50051` **using insecure/plaintext
credentials** (no certificate validation, matching the plaintext server) and
exercises all three streaming RPCs from the `Calculator` service:

1. **GenerateFibonacci** — sends `N = 10` and prints each Fibonacci number as
   the server streams it back (`0, 1, 1, 2, 3, 5, …`).
2. **SendNumbers** — sends numbers `0..8` one per second, then closes the
   stream and prints the server‑computed sum.
3. **Chat** — concurrently sends `"Hi"`, `"How are you?"`, `"Bye"` while
   printing whatever the server echoes back.

Example client output:

```
Fibonacci number: 0
Fibonacci number: 1
Fibonacci number: 1
...
End of Stream
Server resp after stream: 36
Server:  Server replied: Hi
Server:  Server replied: How are you?
Server:  Server replied: Bye
```

Example server output (in the server terminal):

```
Client said: Hi
Client said: How are you?
Client said: Bye
```

> Notes:
> - The sum is `0+1+2+…+8 = 36`.
> - The `Chat` stream may take a few seconds to print all replies because
>   each message is sent ~1 second apart.

---

## 6. Order of Operations (Quick Summary)

```bash
# 1. Install gRPC proto plugins (once)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# 2. Install gRPC dependency (server + client)
cd server && go get google.golang.org/grpc
cd ../client && go get google.golang.org/grpc

# 3. Regenerate code (only if proto/gen is missing / proto changed)
cd ../server && protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto
cd ../client && protoc -I=proto --go_out=. --go-grpc_out=. proto/main.proto

# 4. Run server, then client (two terminals)
cd ../server && go run .
cd ../client && go run .
```

---

## Troubleshooting

- **`protoc` can’t find `protoc-gen-go`** → add your `GOPATH/bin` to `PATH`
  (`export PATH="$PATH:$(go env GOPATH)/bin"`).
- **Client can’t connect** → make sure the server terminal is still running
  and that the client host/port (`localhost:50051`) matches the server.
- **`go get` fails** → run it from the folder that contains the `go.mod`
  for the module you are updating.
- **Missing generated files** → re-run the `protoc` commands in step 3.
- **No output on the client** → the streaming calls print with delays
  (1 second per Fibonacci number / chat message); give it a few seconds.