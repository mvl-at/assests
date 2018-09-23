package assets

import (
	"fmt"
	"net/url"
	"os"
	"path"
)

func find(at assetType, url *url.URL) (*os.File, error) {
	var directory string
	filename := path.Base(url.Path)
	switch at {
	case memberPictureType:
		directory = "member"
	case titlePictureType:
		directory = "."
	case faviconPictureType:
		directory = "."
	}
	return os.OpenFile(fmt.Sprintf("%s/%s", directory, filename), os.O_RDWR|os.O_CREATE, 0666)
}
