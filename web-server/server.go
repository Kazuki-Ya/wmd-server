package web

import (
	"net/http"

	"github.com/Kazuki-Ya/wmd-server/web-server/internal/handlers"
)

func WebInit() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./public"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", handlers.Index)
	mux.HandleFunc("/inference", handlers.Inference)
	mux.HandleFunc("/store", handlers.Store)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
