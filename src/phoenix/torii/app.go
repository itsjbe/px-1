package torii

import (
	pb "../message"
	"database/sql"
	"fmt"
	"github.com/braintree/manners"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"jbe/iris/core/errors"
	"km/phoenix/src/phoenix/torii/handlers"
	"km/phoenix/src/phoenix/torii/health"
	"km/phoenix/src/phoenix/torii/user"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const version = "1.0.0"

func Run() error {

	log.Println("starting torii")

	vaultToken := os.Getenv("VAULT_TOKEN")
	if vaultToken == "" {
		log.Fatal("VAULT_TOKEN must be set and non-empty")
	}

	vaultAddr := os.Getenv("VAULT_ADDR")
	if vaultAddr == "" {
		log.Fatal("VAULT_ADDR must be set and non-empty")
	}

	vc, err := newVaultClient(vaultAddr, vaultToken)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Getting JWT shared secret...")
	secret, err := vc.getJWTSecret("secret/torii")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Getting database credentials...")
	username, password, err := vc.getDatabaseCredentials("mysql/creds/torii")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initializing database connection pool...")
	dbAddr := os.Getenv("TORIIAPP_DB_HOST")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/torii", username, password, dbAddr)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	httpAddr := os.Getenv("NOMAD_ADDR_http")
	if httpAddr == "" {
		log.Fatal("NOMAD_ADDR_http must be set and non-empty")
	}
	log.Printf("HTTP service listening on %s", httpAddr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.Handle("/login", handlers.LoginHandler(secret, user.DB))
	mux.Handle("/secure", handlers.JWTAuthHandler(handlers.HelloHandler))
	mux.Handle("/version", handlers.VersionHandler(version))
	mux.HandleFunc("/healthz", health.HealthzHandler)
	mux.HandleFunc("/healthz/status", health.HealthzStatusHandler)

	httpServer := manners.NewServer()
	httpServer.Addr = httpAddr
	httpServer.Handler = handlers.LoggingHandler(mux)

	errChan := make(chan error, 10)

	go func() {
		errChan <- http.ListenAndServe("", mux)
	}()

	go func() {
		errChan <- vc.renewDatabaseCredentials()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("captured %v. exiting...", s))

			httpServer.BlockingClose()
			os.Exit(0)
		}
	}

	log.Println("server starting on port 50050...")

	var opts = []grpc.ServerOption{
		grpc.StreamInterceptor(streamInterceptor),
		grpc.UnaryInterceptor(unaryInterceptor),
	}

	s := grpc.NewServer(opts...)

	pb.RegisterHelloServiceServer(s, &RegisterWphServer{})
	pb.RegisterPositionServiceServer(s, &PositionServer{})

	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}

	return nil
}

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := authorize(stream.Context()); err != nil {
		return err
	}

	return handler(srv, stream)
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := authorize(ctx); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

var (
	EmptyMetadataErr = errors.New("empty metadata")
	AccessDeniedErr  = errors.New("access denied")
)

func authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return EmptyMetadataErr
	}

	if len(md["token"]) > 0 && md["token"][0] == "123" {
		return nil
	}

	return AccessDeniedErr
}

type RegisterWphServer struct {
}

func (*RegisterWphServer) Greet(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		WphID: req.WphID,
	}, nil
}

type PositionServer struct {
}

func (PositionServer) Position(coord *pb.Coordinate, srv pb.PositionService_PositionServer) error {

	log.Println("handling position request from from req")
	ctx := srv.Context()
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < 10; i++ {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		for i := 0; i < 10; i++ {
			res := pb.CoordinateResponse{
				Lat: coord.GetLat() + r.Float64(),
				Lon: coord.GetLon() + r.Float64(),
			}
			err := srv.Send(&res)

			if err != nil {
				log.Printf("send error %v", err)
				return err
			}

			log.Printf("send new coords %v", res)
		}
	}
	return nil
}

func HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func HomeHandler(w http.ResponseWriter, req *http.Request) {

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Hello from home to %s.", req.RemoteAddr)))
}
