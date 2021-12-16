// Package servo contains a gRPC bases servo client
package servo

import (
	"context"

	rpcclient "go.viam.com/utils/rpc/client"

	"github.com/edaniels/golog"
	"go.viam.com/utils/rpc/dialer"

	"go.viam.com/core/grpc"
	pb "go.viam.com/core/proto/api/component/v1"
)

// serviceClient is a client that satisfies the servo.proto contract
type serviceClient struct {
	conn   dialer.ClientConn
	client pb.ServoServiceClient
	logger golog.Logger
}

// newServiceClient returns a new serviceClient served at the given address
func newServiceClient(ctx context.Context, address string, opts rpcclient.DialOptions, logger golog.Logger) (*serviceClient, error) {
	conn, err := grpc.Dial(ctx, address, opts, logger)
	if err != nil {
		return nil, err
	}
	sc := newSvcClientFromConn(conn, logger)
	return sc, nil
}

// newSvcClientFromConn constructs a new serviceClient using the passed in connection.
func newSvcClientFromConn(conn dialer.ClientConn, logger golog.Logger) *serviceClient {
	client := pb.NewServoServiceClient(conn)
	sc := &serviceClient{
		conn:   conn,
		client: client,
		logger: logger,
	}
	return sc
}

// Close cleanly closes the underlying connections
func (sc *serviceClient) Close() error {
	return sc.conn.Close()
}

// client is a servo client
type client struct {
	*serviceClient
	name string
}

// NewClient constructs a new client that is served at the given address.
func NewClient(ctx context.Context, name string, address string, opts rpcclient.DialOptions, logger golog.Logger) (Servo, error) {
	sc, err := newServiceClient(ctx, address, opts, logger)
	if err != nil {
		return nil, err
	}
	return clientFromSvcClient(sc, name), nil
}

// NewClientFromConn constructs a new Client from connection passed in.
func NewClientFromConn(conn dialer.ClientConn, name string, logger golog.Logger) Servo {
	sc := newSvcClientFromConn(conn, logger)
	return clientFromSvcClient(sc, name)
}

func clientFromSvcClient(sc *serviceClient, name string) Servo {
	return &client{sc, name}
}

func (c *client) Move(ctx context.Context, angle uint8) error {
	req := &pb.ServoServiceMoveRequest{AngleDeg: uint32(angle), Name: c.name}
	_, err := c.client.Move(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (c *client) AngularOffset(ctx context.Context) (uint8, error) {
	req := &pb.ServoServiceAngularOffsetRequest{Name: c.name}
	resp, err := c.client.AngularOffset(ctx, req)
	if err != nil {
		return 0, err
	}
	return uint8(resp.GetAngleDeg()), nil
}