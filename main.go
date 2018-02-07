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
	"os"
	"github.com/jasoncorbett/goslick/jwtauth"
	"github.com/serussell/logxi/v1"
	"mime"
	"io/ioutil"
	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/google"
	"golang.org/x/oauth2"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	googleOAuth2 "golang.org/x/oauth2/google"
	"net/http/httputil"
	"net/url"
	"github.com/koding/websocketproxy"
)


// grpcHandlerFunc returns an http.Handler that delegates to grpcServer on incoming gRPC
// connections or otherHandler otherwise. Copied from cockroachdb.
func grpcHandlerFunc(grpcServer *grpcweb.WrappedGrpcServer, otherHandler http.Handler) http.Handler {
	logger := log.New("http")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Header.Get("Authorization"), "Bearer") {
			token := r.Header.Get("Authorization")[7:]
			logger.Info("Auth token: %s", token)
		}
		// TODO(tamird): point to merged gRPC code rather than a PR.
		// This is a partial recreation of gRPC's internal checks https://github.com/grpc/grpc-go/pull/514/files#diff-95e9a25b738459a2d3030e1e6fa2a718R61
		if strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func serveSwagger(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")

	// Expose files in third_party/swagger-ui/ on <host>/swagger-ui
	fileServer := http.FileServer(http.Dir("third_party/swagger-ui"))
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}

func serveStaticAssets(mux *http.ServeMux) {
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
}


func loginHandlers(mux *http.ServeMux) {
	oauth2Config := &oauth2.Config{
		ClientID:     "",
		ClientSecret: "",
		RedirectURL:  "https://localhost:8888/login/google-callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	mux.Handle("/login/google", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil)))
	mux.Handle("/login/google-callback", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, issueSession(), nil)))
}


// issueSession issues a cookie session after successful Google login
func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		token, err := jwtauth.CreateJWT()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf(`
<html>
	<head>
		<title>Slick Authentication Callback</title>
		<script lang="javascript">
			localStorage.token="%s";
			localStorage.userName="%s";
			localStorage.userFirstName="%s";
			localStorage.userFamilyName="%s";
			localStorage.userGender="%s";
			localStorage.userPicture="%s";
			window.location.replace("https://localhost:8888/index.html");
		</script>
	</head>
	<body>
	</body>
</html>
`, token, googleUser.Name, googleUser.GivenName, googleUser.FamilyName, googleUser.Gender, googleUser.Picture)))
	}
	return http.HandlerFunc(fn)
}


func rootPage(mux *http.ServeMux) {
	mux.HandleFunc("/index.html", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<html>
	<head>
		<title>Jason's Awesome Authentication Page</title>
		<script src="/static/jquery-3.3.1.min.js"></script>
		<script lang="javascript">
			function loggedOut() {
				$('#auth-message').text("You are NOT authenticated.");
				$('#logout').hide();
			}

			$(function() {
				if(localStorage.token) {
					$('#auth-message').text("You are authenticated with token: " + localStorage.token);
					$('#login-button').hide()
					$('#logout').click(function() {
						delete localStorage.token;
						$('#login-button').show();
						loggedOut();
					});
				} else {
					loggedOut()
				}
			});
		</script>
	</head>
	<body>
		<h1 id="auth-message">You are NOT Authenticated</h1>
		<a id="login-button" href="/login/google"><img src="static/google.png"></img></a>
		<a id="logout">Logout</a>
	</body>
</html>
`))
	})
}

func reverseProxy(mux *http.ServeMux) {
	proxyTarget, _ := url.Parse("http://localhost:3000/")
	wsProxyTarget, _ := url.Parse("ws://localhost:3000/")
	proxy := httputil.NewSingleHostReverseProxy(proxyTarget)
	socketproxy := websocketproxy.NewProxy(wsProxyTarget)
	mux.Handle("/index.html", proxy)
	mux.Handle("/index.css", proxy)
	mux.Handle("/index.js", proxy)
	mux.Handle("/index.js.map", proxy)
	 mux.HandleFunc("/sockjs-node/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.Proto, "ws") || strings.Contains(r.Header.Get("Upgrade"), "websocket") {
			socketproxy.ServeHTTP(w, r)
		} else {
			proxy.ServeHTTP(w, r)
		}
	})
	mux.Handle("/__webpack_dev_server__/", proxy)
	mux.Handle("/img/", proxy)
}

func main() {
	if len(os.Args) == 1 || (len(os.Args) > 1 && os.Args[1] == "serve") {
		opts := []grpc.ServerOption{
			grpc.Creds(credentials.NewClientTLSFromCert(certs.DemoCertPool, "localhost:8888"))}

		grpcServer := grpc.NewServer(opts...)
		slickqa.RegisterAuthServer(grpcServer, &slickqa.SlickAuthService{})
		ctx := context.Background()

		wrappedServer := grpcweb.WrapServer(grpcServer)
		dcreds := credentials.NewTLS(&tls.Config{
			ServerName: "localhost:8888",
			RootCAs:    certs.DemoCertPool,
		})
		dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}
		mux := http.NewServeMux()
		gwmux := runtime.NewServeMux()

		swaggerJsonContents, err := ioutil.ReadFile("slickqa/slick.swagger.json")
		mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Mimetype", "application/json")
			w.Write(swaggerJsonContents)
		})

		err = slickqa.RegisterAuthHandlerFromEndpoint(ctx, gwmux, "localhost:8888", dopts)
		if err != nil {
			fmt.Printf("serve: %v\n", err)
			return
		}

		mux.Handle("/", gwmux)
		serveSwagger(mux)
		serveStaticAssets(mux)
		//rootPage(mux)
		reverseProxy(mux)
		loginHandlers(mux)

		conn, err := net.Listen("tcp", "127.0.0.1:8888")
		if err != nil {
			panic(err)
		}

		srv := &http.Server{
			Addr:    "localhost:8888",
			Handler: grpcHandlerFunc(wrappedServer, mux),
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
		fmt.Println("\tauth:   Create an authentication token for use with curl")
	} else if len(os.Args) > 1 && os.Args[1] == "client" {
		if len(os.Args) <= 2 {
			fmt.Println("ERROR: you must supply a permission to client")
			os.Exit(1)
		}
		permission := os.Args[2]

		var opts []grpc.DialOption
		creds := credentials.NewClientTLSFromCert(certs.DemoCertPool, "localhost:8888")
		jwtCreds, _ := jwtauth.NewCredential()
		opts = append(opts, grpc.WithTransportCredentials(creds))
		opts = append(opts, grpc.WithPerRPCCredentials(jwtCreds))
		conn, err := grpc.Dial("localhost:8888", opts...)
		if err != nil {
			fmt.Printf("fail to dial: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close()
		client := slickqa.NewAuthClient(conn)

		msg, err := client.IsAuthorized(context.Background(), &slickqa.IsAuthorizedRequest{permission})
		if err == nil {
			fmt.Printf("Authorised for %s: %t\n", permission, msg.Allowed)
		} else {
			fmt.Println("ERROR: ", err)
			os.Exit(1)
		}
	} else if len(os.Args) > 1 && os.Args[1] == "auth" {
		token, err :=jwtauth.CreateJWT()
		if err != nil {
			fmt.Printf("Error occured: %#v", err)
			os.Exit(1)
		}
		fmt.Printf("Authorization: Bearer %s", token)
	}
}
