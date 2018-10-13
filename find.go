package assets

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"time"
)

func findByUrl(at assetType, url *url.URL) (*os.File, error) {
	filename := path.Base(url.Path)
	return find(at, filename)
}

func find(at assetType, filename string) (*os.File, error) {
	var directory string
	switch at {
	case memberPictureType:
		directory = memberDir
	case titlePictureType:
		directory = titleDir
	case faviconPictureType:
		directory = assetsDir
	}
	return os.OpenFile(fmt.Sprintf("%s/%s", directory, filename), os.O_RDWR|os.O_CREATE, 0666)
}

func dateSuffix() string {
	return time.Now().Format("_2006-01-02-03-04-05")
}
