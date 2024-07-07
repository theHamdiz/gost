package api

import (
	"fmt"

	"github.com/theHamdiz/gost/codegen/general"
	"github.com/theHamdiz/gost/config"
	"github.com/theHamdiz/gost/installer"
	"github.com/theHamdiz/gost/runner"
)

type GenApiPlugin struct {
	Files map[string]func() string
	Data  config.ProjectData
}

func (g *GenApiPlugin) Init() error {
	g.Files = map[string]func() string{
		"app/api/http/v1/api.go": func() string {
			return `package httpAPI
		
import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.AppName}}/app/api/http/v1/helpers.go"
	{{if eq .BackendPkg "chi"}}
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	{{end}}
	{{if eq .BackendPkg "echo"}}
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	{{end}}
	{{if eq .BackendPkg "gin"}}
	"github.com/gin-gonic/gin"
	{{end}}
)

// StartHTTPServer starts the HTTP server based on the BackendPkg variable.
func StartHTTPServer() {
	{{if eq .BackendPkg "stdlib"}}
	log.Println("Starting HTTP server using stdlib")
	{{end}}
	{{if eq .BackendPkg "chi"}}
	log.Println("Starting HTTP server using chi")
	{{end}}
	{{if eq .BackendPkg "echo"}}
	log.Println("Starting HTTP server using echo")
	{{end}}
	{{if eq .BackendPkg "gin"}}
	log.Println("Starting HTTP server using gin")
	{{end}}	
	startServer()
}

// startServer starts the HTTP server based on the BackendPkg variable.
func startServer() {
	executeStartupdownHooks()

	{{if eq .BackendPkg "stdlib"}}
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	mux.HandleFunc("/", rootHandler)

	server := &http.Server{
		Addr:    ":{{.Port}}",
		Handler: mux,
	}
	{{end}}


	{{if eq .BackendPkg "chi"}}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/hello", helloHandler)
	r.Get("/", rootHandler)

	server := &http.Server{
		Addr:    ":{{.Port}}",
		Handler: r,
	}
	{{end}}


	{{if eq .BackendPkg "echo"}}
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/hello", helloHandler)
	e.GET("/", rootHandler)

	server := &http.Server{
		Addr:    ":{{.Port}}",
		Handler: e,
	}
	{{end}}

	
	{{if eq .BackendPkg "gin"}}
	r := gin.Default()
	r.GET("/hello", helloHandler)
	r.GET("/", rootHandler)

	server := &http.Server{
		Addr:    ":{{.Port}}",
		Handler: r,
	}
	{{end}}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()
	gracefulShutdown(server)
}

{{if eq .BackendPkg "stdlib"}}
// helloHandler handles the /hello endpoint for your stdlib server.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!"))
}
// rootHandler handles the / endpoint for your stdlib server.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to your gost app!"))
}
{{end}}

{{if eq .BackendPkg "chi"}}
// helloHandler handles the /hello endpoint for your chi server.
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
// rootHandler handles the / endpoint for your chi server.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to your gost app!"))
}
{{end}}


{{if eq .BackendPkg "echo"}}
// helloHandler handles the /hello endpoint for your echo server.
func helloHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
// rootHandler handles the / endpoint for your echo server.
func rootHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to your gost app!")
}
{{end}}

{{if eq .BackendPkg "gin"}}
// helloHandler handles the /hello endpoint for your gin server.
func helloHandler(c *gin.Context) {
	c.String(http.StatusOK, "Hello, World!")
}
// rootHandler handles the / endpoint for your gin server.
func rootHandler(c *gin.Context) {
	c.String(http.StatusOK, "Welcome to your gost app!")
}
{{end}}


// OnServerShutdown registers middleware functions to be called upon server shutdown.
func OnServerShutdown(middleware ...func()) {
    shutdownHooks = append(shutdownHooks, middleware...)
}

// OnServerShutdown registers middleware functions to be called upon server shutdown.
func OnBeforeServerStart(middleware ...func()) {
    startHooks = append(startHooks, middleware...)
}



		`
		},
		"app/api/http/v1/helpers.go": func() string {
			return `package httpApi

var (
	startupHooks []func()
	shutdownHooks []func()
)
// executeShutdownHooks executes all registered shutdown middleware functions.
func executeShutdownHooks() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if len(shutdownHooks) < 1{
		return
	}

	for _, hook := range shutdownHooks {
        hook()
    }
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}

// executeStartupdownHooks executes all registered startup middleware functions.
func executeStartupdownHooks() {
	if len(startupHooks) < 1{
		return
	}
	for _, hook := range startupHooks {
        hook()
    }
}
`
		},
		"app/api/grpc/v1/proto/service.proto": func() string {
			content := `syntax = "proto3";

option go_package = "app/api/grpc/v1/proto;helloworld";

// The greeting service definition.
service Greeter {
	// Sends a greeting
	rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
	string name = 1;
}

// The response message containing the greeting.
message HelloReply {
	string message = 1;
}
`
			if installer.IsCommandAvailable("protoc") {
				err := runner.RunCommandWithDir("app/api/grpc/v1/proto", "protoc", "--go_out=.", "--go-grpc_out=.", "hello.proto")
				if err != nil {
					fmt.Println("Error generating proto files:", err)
				}
			}
			return content
		},
		"app/api/grpc/v1/server/server.go": func() string {
			return `package grpcServer

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
    grpcLogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
    grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
)

// Server represents a gRPC server.
type Server struct {
	server     *grpc.Server
	listenAddr string
	enableTLS  bool
	certFile   string
	keyFile    string
}

// NewServer creates a new gRPC server.
func NewServer(listenAddr string, enableTLS bool, certFile, keyFile string) *Server {
	var opts []grpc.ServerOption

	// Add middleware
    opts = append(opts, grpc.UnaryInterceptor(
        grpcRecovery.UnaryServerInterceptor(),
        grpcLogrus.UnaryServerInterceptor(grpcLogrus.NewEntry(log.New())),
        grpcAuth.UnaryServerInterceptor(authenticate),
    ))

	// Setup TLS if enabled
    if enableTLS {
        creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
        if err != nil {
            log.Fatalf("Failed to generate credentials %v", err)
        }
        opts = append(opts, grpc.Creds(creds))
    } else {
        opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
    }

    server := grpc.NewServer(opts...)

	return &Server{
		server:     server,
		listenAddr: listenAddr,
		enableTLS:  enableTLS,
		certFile:   certFile,
		keyFile:    keyFile,
	}
}

// Start starts the gRPC server.
func (s *Server) Start() {
	lis, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	go func() {
		if err := s.server.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	log.Printf("gRPC server is running on %s", s.listenAddr)
}

// Stop stops the gRPC server gracefully.
func (s *Server) Stop() {
	log.Println("Stopping gRPC server...")
	s.server.GracefulStop()
}

// RegisterService registers a gRPC service to the server.
func (s *Server) RegisterService(registerFunc func(server *grpc.Server)) {
	registerFunc(s.server)
}

// authenticate is a sample authentication middleware.
func authenticate(ctx context.Context) (context.Context, error) {
	// Implement authentication logic here.
	return ctx, nil
}

// Run starts the server and handles graceful shutdown.
func (s *Server) Run() {
	// Start the server
	s.Start()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Attempt graceful shutdown
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		s.Stop()
		close(quit)
	}()

	<-timeoutCtx.Done()
}

`
		},
		"app/api/grpc/v1/client/client.go": func() string {
			return `package grpcClient

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/middleware"
	"google.golang.org/grpc/middleware/auth"
	"google.golang.org/grpc/middleware/logging"
	"google.golang.org/grpc/middleware/retry"
)

// Client represents a gRPC client.
type Client struct {
	conn       *grpc.ClientConn
	targetAddr string
	enableTLS  bool
	certFile   string
}

// NewClient creates a new gRPC client.
func NewClient(targetAddr string, enableTLS bool, certFile string) *Client {
	return &Client{
		targetAddr: targetAddr,
		enableTLS:  enableTLS,
		certFile:   certFile,
	}
}

// Connect establishes a connection to the gRPC server.
func (c *Client) Connect() error {
	var opts []grpc.DialOption

	// Add middleware
	opts = append(opts, grpc.WithUnaryInterceptor(
		middleware.ChainUnaryClient(
			retry.UnaryClientInterceptor(retry.WithMax(3), retry.WithPerRetryTimeout(1*time.Second)),
			logging.UnaryClientInterceptor(logging.DefaultLogger),
			auth.UnaryClientInterceptor(authenticate),
		),
	))

	// Setup TLS if enabled
	if c.enableTLS {
		creds, err := credentials.NewClientTLSFromFile(c.certFile, "")
		if err != nil {
			return err
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.Dial(c.targetAddr, opts...)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

// Close closes the connection to the gRPC server.
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// GetConn returns the gRPC client connection.
func (c *Client) GetConn() *grpc.ClientConn {
	return c.conn
}

// authenticate is a sample authentication middleware.
func authenticate(ctx context.Context) (context.Context, error) {
	// Implement authentication logic here.
	return ctx, nil
}		
`
		},
	}
	return nil
}

func (g *GenApiPlugin) Execute() error {
	return g.Generate(g.Data)
}

func (g *GenApiPlugin) Shutdown() error {
	// Any cleanup logic for the plugin
	return nil
}

func (g *GenApiPlugin) Name() string {
	return "GenApiPlugin"
}

func (g *GenApiPlugin) Version() string {
	return "1.0.0"
}

func (g *GenApiPlugin) Dependencies() []string {
	return []string{}
}

func (g *GenApiPlugin) AuthorName() string {
	return "Ahmad Hamdi"
}

func (g *GenApiPlugin) AuthorEmail() string {
	return "contact@hamdiz.me"
}

func (g *GenApiPlugin) Website() string {
	return "https://hamdiz.me"
}

func (g *GenApiPlugin) GitHub() string {
	return "https://github.com/theHamdiz/gost/gen/api"
}

func (g *GenApiPlugin) Generate(data config.ProjectData) error {
	return general.GenerateFiles(data, g.Files)
}

func NewGenApiPlugin(data config.ProjectData) *GenApiPlugin {
	return &GenApiPlugin{
		Data: data,
	}
}
