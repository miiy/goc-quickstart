package auth

import (
	authv1 "github.com/miiy/goc-quickstart/api-gateway/gen/go/blog/auth/v1"
	"github.com/miiy/goc-quickstart/api-gateway/internal/transport"

	"github.com/golang-jwt/jwt/v5/request"
	"github.com/miiy/goc/gin"
)

func (m *Module) register(c *gin.Context) {
	var req authv1.RegisterRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.client.Register(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) usernameCheck(c *gin.Context) {
	var req authv1.FieldCheckRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.client.UsernameCheck(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) emailCheck(c *gin.Context) {
	var req authv1.FieldCheckRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.client.EmailCheck(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) phoneCheck(c *gin.Context) {
	var req authv1.FieldCheckRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.client.PhoneCheck(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) login(c *gin.Context) {
	var req authv1.LoginRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.client.Login(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) sendSMSCode(c *gin.Context) {
	var req authv1.SendSmsCodeRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.client.SendSmsCode(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) phoneAuth(c *gin.Context) {
	var req authv1.PhoneAuthRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.client.PhoneAuth(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) mpLogin(c *gin.Context) {
	var req authv1.MpLoginRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.client.MpLogin(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) refreshToken(c *gin.Context) {
	var req authv1.RefreshTokenRequest
	if !transport.BindProto(c, &req) {
		return
	}
	if req.AccessToken == "" {
		req.AccessToken = bearerAccessToken(c)
	}
	resp, err := m.client.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) logout(c *gin.Context) {
	var req authv1.LogoutRequest
	if !transport.BindProto(c, &req) {
		return
	}
	if req.AccessToken == "" {
		req.AccessToken = bearerAccessToken(c)
	}
	resp, err := m.client.Logout(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func bearerAccessToken(c *gin.Context) string {
	token, err := request.BearerExtractor{}.ExtractToken(c.Request)
	if err != nil {
		return ""
	}
	return token
}
