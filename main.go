package main

import (
	"net/http"
	"netdisk-go/handler"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))

	mux.HandleFunc("/file/upload", handler.UploadHandler)
	mux.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	mux.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	mux.HandleFunc("/file/query", handler.FileQueryHandler)
	mux.HandleFunc("/file/download", handler.DownloadHandler)
	mux.HandleFunc("/file/update", handler.FileMetaUpdateHandler)

	server := http.Server{
		Addr:    ":8888",
		Handler: mux,
	}
	server.ListenAndServe()
}
