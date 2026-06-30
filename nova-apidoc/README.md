# nova-apidoc

Serves the generated Nova OpenAPI document with Swagger UI.

The embedded document is copied from:

```text
../nova-contracts/gen/go/http/swagger-json/swagger.json
```

to:

```text
gen/openapi/swagger.json
```

Run the server:

```bash
make run
```

Swagger UI is available at `/`, and the raw OpenAPI JSON is served under
`/openapi/swagger.json`.
