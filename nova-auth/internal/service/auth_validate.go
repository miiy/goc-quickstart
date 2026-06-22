package service

import (
	"strings"

	"buf.build/go/protovalidate"
	pb "github.com/miiy/goc-quickstart/nova-auth/gen/go/nova/auth/v1"
	"github.com/miiy/goc/utils/password"
)

func registerValidate(req *pb.RegisterRequest) error {
	if err := protovalidate.Validate(req); err != nil {
		return err
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	req.PasswordConfirmation = strings.TrimSpace(req.PasswordConfirmation)

	if req.Username == "" || req.Email == "" || req.Password == "" || req.PasswordConfirmation == "" {
		return ErrInvalidArgument
	}
	if req.Password != req.PasswordConfirmation {
		return ErrPasswordsDiffer
	}
	if err := password.Validate(req.Password); err != nil {
		return err
	}
	return nil
}

func loginValidate(req *pb.LoginRequest) error {
	return protovalidate.Validate(req)
}

func mpLoginValidate(req *pb.MpLoginRequest) error {
	return protovalidate.Validate(req)
}

func changePasswordValidate(req *pb.ChangePasswordRequest) error {
	if err := protovalidate.Validate(req); err != nil {
		return err
	}

	req.OldPassword = strings.TrimSpace(req.OldPassword)
	req.NewPassword = strings.TrimSpace(req.NewPassword)
	req.NewPasswordConfirmation = strings.TrimSpace(req.NewPasswordConfirmation)

	if req.OldPassword == "" || req.NewPassword == "" || req.NewPasswordConfirmation == "" {
		return ErrInvalidArgument
	}
	if req.NewPassword != req.NewPasswordConfirmation {
		return ErrPasswordsDiffer
	}
	if err := password.Validate(req.NewPassword); err != nil {
		return err
	}
	return nil
}
