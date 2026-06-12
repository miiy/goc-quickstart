# nova-proto

nova api protocol

## generate protobuf

```bash
make proto
```

```bash
sed -i '' 's/github.com\/miiy\/goc-quickstart\/nova-proto\/gen\/go\/protoc-gen-openapiv2/github.com\/grpc-ecosystem\/grpc-gateway\/v2\/protoc-gen-openapi\/v2/g' gen/go/nova-proto/post/v1/post.pb.go
```
