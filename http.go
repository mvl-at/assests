package assets

import (
	"fmt"
	"io"
	"net/http"
)

//Runs the http Server.
func run() {
	host := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	logger.Println("Listen on " + host)
	err := http.ListenAndServe(host, nil)

	if err != nil {
		errLogger.Fatalln(err.Error())
	}
}

//Registers all http routes.
func routes() {
	http.HandleFunc("/member/", picture(memberPictureType))
	http.HandleFunc("/title", picture(titlePictureType))
}

//Modifies the http header for use with REST.
func picture(at assetType) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == http.MethodOptions {
			writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "access-token,content-type")
			return
		}

		if request.Method == http.MethodGet {
			picture, err := find(at, request.URL)
			defer picture.Close()
			if err != nil {
				writer.WriteHeader(http.StatusNotFound)
				return
			}
			io.Copy(writer, picture)
		}
	}
}
