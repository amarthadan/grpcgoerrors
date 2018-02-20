# gRPCGoErrors

## How to run
Install client and server
```bash
go get github.com/misenko/grpcgoerrors/ggeclient
go get github.com/misenko/grpcgoerrors/ggeserver
```

Run server and client (expects `$GOPATH/bin` in your `$PATH`)
```bash
ggeserver &
ggeclient
```
