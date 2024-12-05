package client

import (
	"context"
	"time"

	"github.com/Nexell-AI-Network/nexelliad/v2/cmd/nexelliawallet/daemon/server"

	"github.com/pkg/errors"

	"github.com/Nexell-AI-Network/nexelliad/v2/cmd/nexelliawallet/daemon/pb"
	"google.golang.org/grpc"
)

// Connect connects to the nexelliawalletd server, and returns the client instance
func Connect(address string) (pb.KaspawalletdClient, func(), error) {
	// Connection is local, so 1 second timeout is sufficient
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(server.MaxDaemonSendMsgSize)))
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, nil, errors.New("nexelliawallet daemon is not running, start it with `nexelliawallet start-daemon`")
		}
		return nil, nil, err
	}

	return pb.NewKaspawalletdClient(conn), func() {
		conn.Close()
	}, nil
}
