package handlers_test

import (
	"context"
	"testing"

	api "github.com/Kazuki-Ya/wmd-server/api/v1"
	"github.com/Kazuki-Ya/wmd-server/web-server/internal/handlers"
	"github.com/stretchr/testify/require"
)

// before test this code, please ensure that inference-server is listen at port 8403
func TestInference(t *testing.T) {
	inferenceClient, err := handlers.InferenceClient()
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
