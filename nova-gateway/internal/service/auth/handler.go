package auth

import (
	authv1 "github.com/miiy/goc-quickstart/nova-gateway/gen/go/nova/auth/v1"
	"github.com/miiy/goc-quickstart/nova-gateway/internal/transport"

	"github.com/miiy/goc/gin"
	ginauth "github.com/miiy/goc/gin/middleware/auth"
)

func (m *Module) register(c *gin.Context) {
	var req authv1.RegisterRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.authClient.Register(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) usernameCheck(c *gin.Context) {
	var req authv1.UsernameCheckRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.authClient.UsernameCheck(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) emailCheck(c *gin.Context) {
	var req authv1.EmailCheckRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.authClient.EmailCheck(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) phoneCheck(c *gin.Context) {
	var req authv1.PhoneCheckRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.authClient.PhoneCheck(c.Request.Context(), &req)
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
	resp, err := m.authClient.Login(c.Request.Context(), &req)
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
	resp, err := m.authClient.SendSmsCode(c.Request.Context(), &req)
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
	resp, err := m.authClient.PhoneAuth(c.Request.Context(), &req)
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
	resp, err := m.authClient.MpLogin(c.Request.Context(), &req)
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
	resp, err := m.authClient.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}

func (m *Module) changePassword(c *gin.Context) {
	var req authv1.ChangePasswordRequest
	if !transport.BindProto(c, &req) {
		return
	}
	resp, err := m.authClient.ChangePassword(c.Request.Context(), &req)
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
		req.AccessToken, _ = ginauth.BearerToken(c)
	}
	resp, err := m.authClient.Logout(c.Request.Context(), &req)
	if err != nil {
		transport.WriteError(c, err)
		return
	}
	transport.WriteProto(c, resp)
}
