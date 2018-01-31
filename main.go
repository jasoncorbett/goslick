package main

import (
	"github.com/jasoncorbett/goslick/certs"
	"github.com/jasoncorbett/goslick/slickqa"
	"net/http"
	"strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"crypto/tls"
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"net"
	"log"
	"os"
	"google.golang.org/grpc/grpclog"
)

// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func main() {
	if len(os.Args) == 1 || (len(os.Args) > 1 && os.Args[1] == "serve") {
		opts := []grpc.ServerOption{
			grpc.Creds(credentials.NewClientTLSFromCert(certs.DemoCertPool, "localhost:8888"))}

		grpcServer := grpc.NewServer(opts...)
		slickqa.RegisterAuthServer(grpcServer, &slickqa.SlickAuthService{})
		ctx := context.Background()

		dcreds := credentials.NewTLS(&tls.Config{
			ServerName: "localhost:8888",
			RootCAs:    certs.DemoCertPool,
		})
		dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
		mux := http.NewServeMux()
		gwmux := runtime.NewServeMux()
		err := slickqa.RegisterAuthHandlerFromEndpoint(ctx, gwmux, "localhost:8888", dopts)
		if err != nil {
			fmt.Printf("serve: %v\n", err)
			return
		}

		mux.Handle("/", gwmux)

		conn, err := net.Listen("tcp", "127.0.0.1:8888")
		if err != nil {
			panic(err)
		}

		srv := &http.Server{
			Addr:    "localhost:8888",
			Handler: grpcHandlerFunc(grpcServer, mux),
			TLSConfig: &tls.Config{
				Certificates: []tls.Certificate{*certs.DemoKeyPair},
				NextProtos:   []string{"h2"},
			},
		}

		fmt.Printf("grpc on port: %d\n", 8888)
		fmt.Println("Try calling:\ncurl -k https://localhost:8888/api/v1/isAuthorized/foo")
		err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))

		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	} else if len(os.Args) > 1 && os.Args[1] == "-h" {
		fmt.Println("Usage: slickgo [command]")
		fmt.Println("Commands:")
		fmt.Println("\tserve:  Serve tells it to run the grpc server")
		fmt.Println("\tclient: Call this with a permission to test the permission against the server")
	} else if len(os.Args) > 1 && os.Args[1] == "client" {
		if len(os.Args) <= 2 {
			fmt.Println("ERROR: you must supply a permission to client")
			os.Exit(1)
		}
		permission := os.Args[2]

		var opts []grpc.DialOption
		creds := credentials.NewClientTLSFromCert(certs.DemoCertPool, "localhost:8888")
		opts = append(opts, grpc.WithTransportCredentials(creds))
		conn, err := grpc.Dial("localhost:8888", opts...)
		if err != nil {
			grpclog.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()
		client := slickqa.NewAuthClient(conn)

		msg, err := client.IsAuthorized(context.Background(), &slickqa.IsAuthorizedRequest{permission})
		fmt.Printf("Authorised for %s: %t\n", permission, msg.Allowed)
	}
}
