package router

import (
	"regexp"
	"strings"
	"testing"

	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	postv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/post/v1"
	userv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/user/v1"
	authsvc "github.com/miiy/goc-quickstart/nova-gateway/internal/service/auth"
	filesvc "github.com/miiy/goc-quickstart/nova-gateway/internal/service/file"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/service/health"
	postsvc "github.com/miiy/goc-quickstart/nova-gateway/internal/service/post"
	usersvc "github.com/miiy/goc-quickstart/nova-gateway/internal/service/user"
	"github.com/miiy/goc/gin"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type route struct {
	method string
	path   string
}

func TestGinRoutesMatchProtoHTTPAnnotations(t *testing.T) {
	got := ginRoutesForTest()
	want := protoRoutesForTest(t)

	for r := range want {
		if _, ok := got[r]; !ok {
			t.Fatalf("missing gin route for proto annotation: %s %s", r.method, r.path)
		}
	}

	for r := range got {
		// File uploads are multipart endpoints assembled by the gateway, not proto JSON routes.
		if r.path == "/healthz" || r.path == "/api/v1/files/upload" || r.path == "/api/v1/files/upload/avatar" {
			continue
		}
		if _, ok := want[r]; !ok {
			t.Fatalf("gin route has no proto annotation: %s %s", r.method, r.path)
		}
	}
}

func ginRoutesForTest() map[route]struct{} {
	r := gin.New()

	health.NewModule(r).RegisterRouter()

	authModule := authsvc.NewModule(nil)
	authModule.RegisterPublicRouter(r)

	public := r.Group("/api/v1")
	postModule := postsvc.NewModule(nil, nil)
	postModule.RegisterPublicRouter(public)

	protected := r.Group("/api/v1")
	protected.Use(func(c *gin.Context) {
		c.Next()
	})

	authModule.RegisterProtectedRouter(protected)
	postModule.RegisterProtectedRouter(protected)
	filesvc.NewModule(nil, nil).RegisterRouter(protected)
	usersvc.NewModule(nil).RegisterRouter(protected)

	routes := make(map[route]struct{})
	for _, item := range r.Routes() {
		routes[route{method: item.Method, path: item.Path}] = struct{}{}
	}
	return routes
}

func protoRoutesForTest(t *testing.T) map[route]struct{} {
	t.Helper()

	files := []protoreflect.FileDescriptor{
		authv1.File_nova_auth_v1_auth_proto,
		postv1.File_nova_post_v1_post_proto,
		userv1.File_nova_user_v1_user_proto,
	}

	routes := make(map[route]struct{})
	for _, file := range files {
		services := file.Services()
		for i := 0; i < services.Len(); i++ {
			methods := services.Get(i).Methods()
			for j := 0; j < methods.Len(); j++ {
				method := methods.Get(j)
				opts, ok := method.Options().(*descriptorpb.MethodOptions)
				if !ok || !proto.HasExtension(opts, annotations.E_Http) {
					continue
				}

				ext := proto.GetExtension(opts, annotations.E_Http)
				rule, ok := ext.(*annotations.HttpRule)
				if !ok {
					t.Fatalf("unexpected http option type for %s: %T", method.FullName(), ext)
				}
				addHTTPRule(routes, rule)
			}
		}
	}
	return routes
}

func addHTTPRule(routes map[route]struct{}, rule *annotations.HttpRule) {
	if rule == nil {
		return
	}

	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		routes[route{method: "GET", path: protoPathToGin(pattern.Get)}] = struct{}{}
	case *annotations.HttpRule_Post:
		routes[route{method: "POST", path: protoPathToGin(pattern.Post)}] = struct{}{}
	case *annotations.HttpRule_Put:
		routes[route{method: "PUT", path: protoPathToGin(pattern.Put)}] = struct{}{}
	case *annotations.HttpRule_Delete:
		routes[route{method: "DELETE", path: protoPathToGin(pattern.Delete)}] = struct{}{}
	case *annotations.HttpRule_Patch:
		routes[route{method: "PATCH", path: protoPathToGin(pattern.Patch)}] = struct{}{}
	case *annotations.HttpRule_Custom:
		routes[route{method: strings.ToUpper(pattern.Custom.Kind), path: protoPathToGin(pattern.Custom.Path)}] = struct{}{}
	}

	for _, binding := range rule.AdditionalBindings {
		addHTTPRule(routes, binding)
	}
}

var protoPathParam = regexp.MustCompile(`\{([^}=]+)(=[^}]*)?\}`)

func protoPathToGin(path string) string {
	return protoPathParam.ReplaceAllString(path, `:$1`)
}
