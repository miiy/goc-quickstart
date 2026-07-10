module github.com/miiy/goc-quickstart/nova-gateway

go 1.26.3

require (
	buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go v1.36.11-20260415201107-50325440f8f2.1
	github.com/google/wire v0.7.0
	github.com/miiy/goc v0.1.1
	github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server v0.0.0
	google.golang.org/genproto/googleapis/api v0.0.0-20260523011958-0a33c5d7ca68
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)

replace github.com/miiy/goc-quickstart/nova-contracts/gen/go/http/go-gin-server => ../nova-contracts/gen/go/http/go-gin-server

replace github.com/miiy/goc => ../../goc

require (
	github.com/boj/redistore v1.4.1 // indirect
	github.com/bytedance/gopkg v0.1.4 // indirect
	github.com/bytedance/sonic v1.15.1 // indirect
	github.com/bytedance/sonic/loader v0.5.1 // indirect
	github.com/cloudwego/base64x v0.1.7 // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/gin-contrib/cors v1.7.7 // indirect
	github.com/gin-contrib/sessions v1.1.0 // indirect
	github.com/gin-contrib/sse v1.1.1 // indirect
	github.com/gin-contrib/zap v1.1.7 // indirect
	github.com/gin-gonic/gin v1.12.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.2 // indirect
	github.com/goccy/go-json v0.10.6 // indirect
	github.com/goccy/go-yaml v1.19.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-isatty v0.0.22 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.3.1 // indirect
	github.com/quic-go/qpack v0.6.0 // indirect
	github.com/quic-go/quic-go v0.59.1 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/sqids/sqids-go v0.4.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.1 // indirect
	go.mongodb.org/mongo-driver/v2 v2.6.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.28.0 // indirect
	golang.org/x/arch v0.27.0 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/net v0.55.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
	golang.org/x/text v0.37.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260523011958-0a33c5d7ca68 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
