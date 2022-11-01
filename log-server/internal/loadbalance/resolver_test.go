package loadbalance_test

import (
	"net"
	"testing"

	api "github.com/Kazuki-Ya/wmd-server/api/v1"
	"github.com/Kazuki-Ya/wmd-server/log-server/internal/loadbalance"
	"github.com/Kazuki-Ya/wmd-server/log-server/internal/server"
	"github.com/stretchr/testify/require"

	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestResolver(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	srv, err := server.NewGRPCServer(&server.Config{
		GetServerer: &getServers{},
	})
	require.NoError(t, err)

	go srv.Serve(l)

	conn := &clientConn{}
	require.NoError(t, err)

	clientCreds := insecure.NewCredentials()

	opts := resolver.BuildOptions{
		DialCreds: clientCreds,
	}
	r := &loadbalance.Resolver{}
	_, err = r.Build(
		resolver.Target{
			Endpoint: l.Addr().String(),
		},
		conn,
		opts,
	)
	require.NoError(t, err)

	wantState := resolver.State{
		Addresses: []resolver.Address{{
			Addr:       "localhost:9001",
			Attributes: attributes.New("is_leader", true),
		}, {
			Addr:       "localhost:9002",
			Attributes: attributes.New("is_leader", false),
		}},
	}
	require.Equal(t, wantState, conn.state)

	conn.state.Addresses = nil
	r.ResolveNow(resolver.ResolveNowOptions{})
	require.Equal(t, wantState, conn.state)
}

type getServers struct{}

func (s *getServers) GetServers() ([]*api.Server, error) {
	return []*api.Server{{
		Id:       "leader",
		RpcAddr:  "localhost:9001",
		IsLeader: true,
	}, {
		Id:      "follower",
		RpcAddr: "localhost:9002",
	}}, nil
}

type clientConn struct {
	resolver.ClientConn
	state resolver.State
}

func (c *clientConn) UpdateState(state resolver.State) error {
	c.state = state
	return nil
}

func (c *clientConn) ReportError(err error) {}

func (c *clientConn) NewAddress(addrs []resolver.Address) {}

func (c *clientConn) NewServiceConfig(config string) {}

func (c *clientConn) ParseServiceConfig(
	config string,
) *serviceconfig.ParseResult {
	return nil
}
