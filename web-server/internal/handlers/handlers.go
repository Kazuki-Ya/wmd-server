package handlers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"

	api "github.com/Kazuki-Ya/wmd-server/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Index(writer http.ResponseWriter, request *http.Request) {

	t, err := template.ParseFiles("./web-server/template/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(writer, nil)
}

func Inference(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	input_sentence := request.FormValue("input_sentence")
	inferenceClient, err := InferenceClient()
	if err != nil {
		log.Fatal("failed to create client")
	}

	output, err := inferenceClient.InferenceCall(ctx, &api.InputDataForInference{
		Text: input_sentence,
	})
	if err != nil {
		log.Fatal("failed to recieve output")
	}

	t, err := template.ParseFiles("./web-server/template/inference.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(writer, output.Label)
}

func Store(writer http.ResponseWriter, request *http.Request) {
	ctx := context.Background()
	sent_sentence := request.FormValue("failed_sentence")
	logClient, err := LogClient()
	if err != nil {
		log.Fatal(err)
	}

	response, err := logClient.Produce(ctx, &api.ProduceRequest{
		Record: &api.Record{
			Text: []byte(sent_sentence),
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintln(writer, response.Offset)
}

func LogClient() (api.LogClient, error) {
	addr := "127.0.0.1:8400"

	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cc, err := grpc.Dial(addr, clientOptions...)
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewLogClient(cc)

	return client, nil
}

func InferenceClient() (api.InferenceClient, error) {
	addr := "127.0.0.1:8403"

	clientOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	cc, err := grpc.Dial(addr, clientOptions...)
	if err != nil {
		log.Fatal(err)
	}

	client := api.NewInferenceClient(cc)

	return client, nil
}

func generateHTML(writer http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("template/%s.html", file))
	}

	templates := template.Must(template.ParseFiles(files...))
	templates.ExecuteTemplate(writer, "layout", data)
}
