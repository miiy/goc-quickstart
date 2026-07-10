# nova-contracts

Nova API contracts.

## Layout

```text
proto/       gRPC/protobuf source files
openapi/     hand-written OpenAPI source files
gen/go/rpc/  buf/protobuf and gRPC generated output
gen/go/http/ OpenAPI generated output
```

`proto/` is only for service-to-service gRPC contracts. Do not add
`google.api.http` or gRPC OpenAPI annotations there. Gateway HTTP paths,
request/response DTOs, and public API docs are defined in `openapi/`.

OpenAPI source files are organized by API module and version. The root
`openapi/openapi.yaml` is the aggregate entrypoint used by all generators.

```text
openapi/
  openapi.yaml
  common/v1/schemas/
  auth/v1/paths/
  auth/v1/schemas/
  file/v1/paths/
  file/v1/schemas/
  post/v1/paths/
  post/v1/schemas/
  user/v1/paths/
  user/v1/schemas/
```

## Generate Protobuf

```bash
make proto
```

## OpenAPI Tooling

Install Node dependencies once:

```bash
make openapi-deps
```

Validate the aggregate OpenAPI document:

```bash
make openapi-validate
```

Generate all OpenAPI outputs:

```bash
make openapi-generate
```

Generate one target at a time:

```bash
make openapi-generate-go-gin-server
make openapi-generate-ts-client
make openapi-generate-swagger-json
```

Generated outputs:

```text
gen/go/http/go-gin-server/      Gin server contract, DTOs, and route interfaces
gen/ts/http/ts-client/          TypeScript fetch frontend client
gen/go/http/swagger-json/       bundled swagger.json
```

The OpenAPI generator version and targets are pinned in `openapitools.json`.
Use `JAVA_HOME` if Java is not already on `PATH`.
