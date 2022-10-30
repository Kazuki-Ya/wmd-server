package handlers

import (
	"context"
	"os"
	"testing"
	"time"

	api "github.com/Kazuki-Ya/wmd-server/api/v1"
	"github.com/Kazuki-Ya/wmd-server/log-server/agent"

	"github.com/stretchr/testify/require"
)

// before run this code, please ensure that inference-server is listen at port 8403
func TestInference(t *testing.T) {
	inferenceClient, err := inferenceClient()
	require.NoError(t, err)

	ctx := context.Background()
	input_sentence := "Today is a good day"
	output, err := inferenceClient.InferenceCall(ctx, &api.InputDataForInference{
		Text: input_sentence,
	})
	require.NoError(t, err)

	label := output.Label
	require.Equal(t, uint32(1), label)
}

func TestLogstore(t *testing.T) {
	// boot agent
	bindAddr := "127.0.0.1:8401" // for serf port
	dataDir, err := os.MkdirTemp("", "web-log-communication-test")
	require.NoError(t, err)
	var startJoinAddrs []string

	agent, err := agent.New(agent.Config{
		NodeName:       "0",
		StartJoinAddrs: startJoinAddrs,
		BindAddr:       bindAddr,
		RPCPort:        8400,
		DataDir:        dataDir,
		Bootstrap:      true,
	})
	require.NoError(t, err)

	defer func() {
		err := agent.Shutdown()
		require.NoError(t, err)
		require.NoError(t, os.RemoveAll(agent.Config.DataDir))
	}()

	logClient, err := logClient()
	require.NoError(t, err)

	produceResponse, err := logClient.Produce(
		context.Background(),
		&api.ProduceRequest{
			Record: &api.Record{
				Text: []byte("hello"),
			},
		},
	)
	require.NoError(t, err)

	time.Sleep(1 * time.Second)

	consumeResponse, err := logClient.Consume(
		context.Background(),
		&api.ConsumeRequest{
			Offset: produceResponse.Offset,
		},
	)
	require.NoError(t, err)
	require.Equal(t, consumeResponse.Record.Text, []byte("hello"))

}
