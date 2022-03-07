# FattarielloDB
a cluster db for your fattariellis

## build
```
go build -a -installsuffix cgo -ldflags '-w -extldflags "-static"' -o ./bin/client_run ./client/
go build -a -installsuffix cgo -ldflags '-w -extldflags "-static"' -o ./bin/server_run ./server/
```
### protobuf + gRPC

```
protoc \
    --go_out=plugins=grpc:. \
    --go_opt=paths=source_relative \
    ./proto/node.proto
```
