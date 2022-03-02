# FattarielloDB
a cluster db for your fattariellis

### protobuf + gRPC

```
protoc --go_out=./cluster --go_opt=paths=source_relative \
    --go-grpc_out=./cluster --go-grpc_opt=paths=source_relative \
    node.proto
```
