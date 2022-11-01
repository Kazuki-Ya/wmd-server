package agent_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	api "github.com/Kazuki-Ya/wmd-server/api/v1"
	"github.com/Kazuki-Ya/wmd-server/log-server/agent"
	"github.com/Kazuki-Ya/wmd-server/log-server/dynaport"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestAgent(t *testing.T) {
	var agents []*agent.Agent
	for i := 0; i < 3; i++ {
		ports := dynaport.Get(2)
		bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", ports[0])
		rpcPort := ports[1]

		dataDir, err := os.MkdirTemp("", "agent-test-log")
		require.NoError(t, err)

		var startJoinAddrs []string

		if i != 0 {
			startJoinAddrs = append(startJoinAddrs,
				agents[0].Config.BindAddr)
		}

		agent, err := agent.New(agent.Config{
			NodeName:       fmt.Sprintf("%d", i),
			StartJoinAddrs: startJoinAddrs,
			BindAddr:       bindAddr,
			RPCPort:        rpcPort,
			DataDir:        dataDir,
			Bootstrap:      i == 0,
		})
		require.NoError(t, err)

		agents = append(agents, agent)
	}

	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t,
				os.RemoveAll(agent.Config.DataDir),
			)
		}
	}()

	leaderClient := client(t, agents[0])
	produceResponse, err := leaderClient.Produce(
		context.Background(),
		&api.ProduceRequest{
			Record: &api.Record{
				Text: []byte("hello"),
			},
		},
	)
	require.NoError(t, err)

	time.Sleep(3 * time.Second)

	consumeResponse, err := leaderClient.Consume(
		context.Background(),
		&api.ConsumeRequest{
			Offset: produceResponse.Offset,
		},
	)
	require.NoError(t, err)
	require.Equal(t, consumeResponse.Record.Text, []byte("hello"))

	time.Sleep(3 * time.Second)

	followerClient := client(t, agents[1])
	consumeResponse, err = followerClient.Consume(
		context.Background(),
		&api.ConsumeRequest{
			Offset: produceResponse.Offset,
		},
	)
	require.NoError(t, err)
	require.Equal(t, consumeResponse.Record.Text, []byte("hello"))

	consumeResponse, err = leaderClient.Consume(
		context.Background(),
		&api.ConsumeRequest{
			Offset: produceResponse.Offset + 1,
		},
	)
	require.Nil(t, consumeResponse)
	require.Error(t, err)
	got := status.Code(err)
	want := status.Code(api.ErrOffsetOutOfRange{}.GRPCStatus().Err())
	require.Equal(t, got, want)
}

func client(
	t *testing.T,
	agent *agent.Agent,
) api.LogClient {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	rpcAddr, err := agent.Config.RPCAddr()
	require.NoError(t, err)

	conn, err := grpc.Dial(rpcAddr, opts...)
	require.NoError(t, err)

	client := api.NewLogClient(conn)
	return client
}
