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
	http.HandleFunc("/member/", member)
}

func member(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		picture, err := find(memberPictureType, r.URL)
		defer picture.Close()
		if err != nil {
			rw.WriteHeader(http.StatusNotFound)
			errLogger.Println(err.Error())
			return
		}
		io.Copy(rw, picture)
	}
}
