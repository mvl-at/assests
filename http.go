package assets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mvl-at/model"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"strings"
)

const (
	memberRedirect = "/member/"
	memberPicture  = "/memberPicture/"
	titleRedirect  = "/title"
	titlePicture   = "/titlePicture/"
	defaultTitle   = "/defaultTitle"
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
	http.HandleFunc(memberRedirect, picture(memberPictureRedirectType))
	http.HandleFunc(memberPicture, picture(memberPictureType))
	http.HandleFunc(titleRedirect, picture(titlePictureRedirectType))
	http.HandleFunc(titlePicture, picture(titlePictureType))
	http.HandleFunc(defaultTitle, picture(defaultTitleType))
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
			if at == memberPictureRedirectType {

				if !assetIndex.memberHasPicture(path.Base(request.URL.Path)) {
					writer.WriteHeader(http.StatusNotFound)
					return
				}

				http.Redirect(writer, request, memberPicture+assetIndex.getMemberPictureName(path.Base(request.URL.Path)), http.StatusSeeOther)
				return
			}

			if at == titlePictureRedirectType {
				pictureName := assetIndex.getTitlePictureName()
				if assetIndex.getIsDefaultTitle() {
					pictureName = assetIndex.getDefaultTitlePictureName()
				}
				if request.URL.Query().Get("thumb") == "true" {
					pictureName += "-thumb"
				}
				http.Redirect(writer, request, titlePicture+pictureName, http.StatusSeeOther)
				return
			}

			if at == defaultTitleType {
				err := json.NewEncoder(writer).Encode(assetIndex.getIsDefaultTitle())
				if err != nil {
					errLogger.Println(err.Error())
				}
			}
			picture, err := findByUrl(at, request.URL)
			defer picture.Close()
			if err != nil {
				writer.WriteHeader(http.StatusNotFound)
				return
			}
			io.Copy(writer, picture)
		}

		if request.Method == http.MethodPost {
			jwt := request.Header.Get("access-token")

			roles := fetchRoles(jwt)
			if !hasRole(roles, at) {
				writer.WriteHeader(http.StatusForbidden)
				return
			}

			var filename string
			var persistAssetType assetType

			if at == defaultTitleType {
				var isDefault bool
				err := json.NewDecoder(request.Body).Decode(&isDefault)
				if err != nil {
					errLogger.Println(err.Error())
					return
				}
				assetIndex.setIsDefaultTitle(isDefault)
				return
			}
			if at == memberPictureRedirectType {
				id, _ := strconv.ParseInt(path.Base(request.URL.Path), 10, 64)
				filename = fetchUsername(id) + dateSuffix()
				persistAssetType = memberPictureType
				assetIndex.setMemberPictureName(path.Base(request.URL.Path), filename)
			}
			if at == titlePictureRedirectType {
				filename = "title" + dateSuffix()
				persistAssetType = titlePictureType
				if request.URL.Query().Get("default") == "true" {
					assetIndex.setDefaultTitlePictureName(filename)
				} else {
					assetIndex.setTitlePictureName(filename)
				}
				data, err := ioutil.ReadAll(request.Body)
				if err != nil {
					errLogger.Println(err.Error())
					return
				}
				img, typ, err := image.Decode(bytes.NewReader(data))
				if err != nil {
					errLogger.Printf("%s occured when trying to resize %s", err.Error(), typ)
					return
				}
				resized := resize.Resize(0, 500, img, resize.Bicubic)
				resizedFile, err := find(persistAssetType, filename+"-thumb")
				switch typ {
				case "jpeg":
					err = jpeg.Encode(resizedFile, resized, &jpeg.Options{Quality: 100})
				case "png":
					err = png.Encode(resizedFile, resized)
				default:
					errLogger.Printf("image type %s is not supported", typ)
				}
				if err != nil {
					errLogger.Println(err.Error())
				}
				request.Body = ioutil.NopCloser(bytes.NewReader(data))
			}
			file, err := find(persistAssetType, filename)
			if err != nil {
				errLogger.Println(err.Error())
				return
			}
			defer file.Close()
			io.Copy(file, request.Body)
		}
	}
}

func fetchUsername(id int64) string {
	resp, err := http.Get(conf.RestHost + "/members")
	if err != nil {
		errLogger.Println(err.Error())
		return ""
	}
	decoder := json.NewDecoder(resp.Body)
	members := make([]model.Member, 0)
	decoder.Decode(&members)
	for _, member := range members {
		if member.Id == id {
			return member.Username
		}
	}
	return ""
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

func hasRole(roles []model.Role, at assetType) bool {
	for _, v := range roles {
		if strings.ToLower(string(at)) == strings.ToLower(v.Id) || strings.ToLower(v.Id) == "root" {
			return true
		}
	}
	return false
}
