package grpc

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps a gRPC client connection
type Client struct {
	conn   *grpc.ClientConn
	logger *zap.Logger
}

// NewClient creates a new gRPC client
func NewClient(address string, logger *zap.Logger, opts ...grpc.DialOption) (*Client, error) {
	// Add default options
	defaultOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(10 * time.Second),
		grpc.WithChainUnaryInterceptor(
			ClientLoggingInterceptor(logger),
		),
	}

	allOpts := append(defaultOpts, opts...)

	conn, err := grpc.Dial(address, allOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
	}

	logger.Info("gRPC client connected", zap.String("address", address))

	return &Client{
		conn:   conn,
		logger: logger,
	}, nil
}

// GetConn returns the underlying gRPC connection
func (c *Client) GetConn() *grpc.ClientConn {
	return c.conn
}

// Close closes the gRPC client connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ClientLoggingInterceptor logs gRPC client requests
func ClientLoggingInterceptor(logger *zap.Logger) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()

		err := invoker(ctx, method, req, reply, cc, opts...)

		duration := time.Since(start)

		if err != nil {
			logger.Error("gRPC client request failed",
				zap.String("method", method),
				zap.Duration("duration", duration),
				zap.Error(err),
			)
		} else {
			logger.Debug("gRPC client request",
				zap.String("method", method),
				zap.Duration("duration", duration),
			)
		}

		return err
	}
}
