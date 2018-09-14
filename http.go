package assets

import (
	"encoding/json"
	"fmt"
	"github.com/mvl-at/model"
	"io"
	"net/http"
	"os"
	"strings"
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
		writer.Header().Set("access-control-allow-origin", "*")
		writer.Header().Set("access-control-expose-headers", "access-token")
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

		if request.Method == http.MethodPost {
			jwt := request.Header.Get("access-token")

			if hasRole(fetchRoles(jwt), conf.TitleRole) {
				var directory string
				var path string

				if at == memberPictureType {
					directory = "."
					path = request.URL.Path
				}

				if at == titlePictureType {
					directory = "."
					path = "title"
				}

				file, err := os.OpenFile(fmt.Sprintf("%s/%s", directory, path), os.O_RDWR|os.O_CREATE, 0666)
				if err != nil {
					errLogger.Println(err.Error())
					return
				}
				defer file.Close()
				io.Copy(file, request.Body)
			}
		}
	}
}

func fetchRoles(jwt string) []model.Role {
	req, err := http.NewRequest(http.MethodGet, conf.RestHost+"/userinfo", nil)

	if err != nil {
		errLogger.Println(err.Error())
		return make([]model.Role, 0)
	}

	req.Header.Set("access-token", jwt)
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		errLogger.Println(err.Error())
		return make([]model.Role, 0)
	}

	if response.StatusCode != http.StatusOK {
		return make([]model.Role, 0)
	}

	userInfo := &UserInfo{}
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(userInfo)
	return userInfo.Roles
}

func hasRole(roles []model.Role, role string) bool {
	for _, v := range roles {
		if strings.ToLower(role) == strings.ToLower(v.Id) || strings.ToLower(v.Id) == "root" {
			return true
		}
	}
	return false
}
