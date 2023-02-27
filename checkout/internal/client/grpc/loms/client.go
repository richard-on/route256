package loms

import (
	"gitlab.ozon.dev/rragusskiy/homework-1/loms/pkg/loms"
	"google.golang.org/grpc"
)

// Client is a wrapper for LOMS gRPC client.
type Client struct {
	lomsClient loms.LOMSClient
}

// NewClient creates new LOMS gRPC client.
func NewClient(cc *grpc.ClientConn) *Client {
	return &Client{
		lomsClient: loms.NewLOMSClient(cc),
	}
}
