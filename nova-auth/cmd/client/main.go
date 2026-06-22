package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/miiy/goc-quickstart/nova-auth/gen/go/nova/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var addr = flag.String("addr", "localhost:50051", "the address to connect to")
	flag.Parse()

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("dit not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewAuthServiceClient(conn)
	//register(c)
	//login(c)
	mpLogin(c)
}

func register(c pb.AuthServiceClient) {
	regReq := pb.RegisterRequest{
		Email:                "test@test.com",
		Username:             "username",
		Password:             "123456",
		PasswordConfirmation: "123456",
	}
	rResp, err := callRegister(c, &regReq)
	if err != nil {
		log.Fatalf("client.callRegister(_) = _, %v", err)
	}
	log.Println("SignUp:", rResp)
}

func login(c pb.AuthServiceClient) {
	regReq := pb.LoginRequest{
		Username: "username",
		Password: "123456",
	}
	rResp, err := callLogin(c, &regReq)
	if err != nil {
		log.Fatalf("client.callRegister(_) = _, %v", err)
	}
	log.Println("Login:", rResp)
}

func mpLogin(c pb.AuthServiceClient) {
	code := ""
	req := pb.MpLoginRequest{
		Code: code,
	}
	resp, err := callMpLogin(c, &req)
	if err != nil {
		log.Fatalf("client.callMpLogin(_) = _, %v", err)
	}
	log.Println("MpLogin:", resp)
}

func callMpLogin(client pb.AuthServiceClient, req *pb.MpLoginRequest) (*pb.MpLoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	return client.MpLogin(ctx, req)
}

func callRegister(client pb.AuthServiceClient, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return client.Register(ctx, req)
}

func callLogin(client pb.AuthServiceClient, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return client.Login(ctx, req)
}
