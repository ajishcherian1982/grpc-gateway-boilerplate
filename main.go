package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"

	"github.com/ajishcherian1982/grpc-gateway-boilerplate/db"
	"github.com/ajishcherian1982/grpc-gateway-boilerplate/gateway"
	"github.com/ajishcherian1982/grpc-gateway-boilerplate/handler"
	"github.com/ajishcherian1982/grpc-gateway-boilerplate/insecure"
	pbExample "github.com/ajishcherian1982/grpc-gateway-boilerplate/proto"
	"github.com/ajishcherian1982/grpc-gateway-boilerplate/server"
	"github.com/ajishcherian1982/grpc-gateway-boilerplate/store"
)

func main() {
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	addr := "0.0.0.0:10000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}

	s := grpc.NewServer(
		// TODO: Replace with your own certificate!
		grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
	)
	pbExample.RegisterUserServiceServer(s, server.New())

	d, err := db.New()
	if err != nil {
		err = fmt.Errorf("failed to connect database: %w", err)
		log.Fatalln("Failed to connect to database:", err)
	}
	us := store.NewUserStore(d)
	as := store.NewArticleStore(d)
	h := handler.New(&l, us, as)
	pbExample.RegisterUsersServer(s)

	// Serve gRPC Server
	log.Info("Serving gRPC on https://", addr)
	go func() {
		log.Fatal(s.Serve(lis))
	}()

	err = gateway.Run("dns:///" + addr)
	log.Fatalln(err)
}
