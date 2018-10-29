package wph

import (
	pb "../message"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

const (
	address = "localhost:50050"
)

func Run() error {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&tokenCred{
			token: "123",
		}))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return err
	}
	defer conn.Close()

	c := pb.NewHelloServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	r, err := c.Greet(ctx, &pb.HelloRequest{
		WphID: "some wph id",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
		return err
	}

	log.Printf("Greeting: %s", r.WphID)

	StreamTest(conn)

	return nil
}

type tokenCred struct {
	token string
}

func (c *tokenCred) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"token": c.token,
	}, nil
}

func (c *tokenCred) RequireTransportSecurity() bool {
	return false
}

func StreamTest(conn *grpc.ClientConn) error {

	client := pb.NewPositionServiceClient(conn)

	outCoord := pb.Coordinate{}

	stream, err := client.Position(context.Background(), &outCoord)

	if err != nil {
		return err
	}

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// read done.
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("Got message at point(%v, %v)", in.GetLat(), in.GetLon())
	}
}

type HelloHandler struct {
}

func (*HelloHandler) Greet(context.Context, *pb.HelloRequest) (*pb.HelloResponse, error) {
	panic("implement me")
}
