# Building

## Build All Binaries

```bash
make build
```

Outputs to `.bin/`:

| Binary | Command |
|---|---|
| `cmn-server` | `./cmd/server/main.go` |
| `cmn-client-admin` | `./cmd/client/admin/main.go` |
| `cmn-client-app` | `./cmd/client/app/main.go` |
| `cmn-client-anonymous` | `./cmd/client/anonymous/main.go` |

## Individual Builds

```bash
go build -o .bin/cmn-server ./cmd/server/main.go
go build -o .bin/cmn-client-admin ./cmd/client/admin/main.go
```

## Clean

```bash
make clean
```
