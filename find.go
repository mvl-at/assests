package assets

import (
	"fmt"
	"net/url"
	"os"
	"path"
)

func find(at assetType, url *url.URL) (*os.File, error) {
	var directory string
	switch at {
	case memberPictureType:
		directory = "member"
	case titlePictureType:
		directory = "."
		url, _ = url.Parse("title")
	case faviconPictureType:
		directory = "."
		url, _ = url.Parse("icon")
	}
	return os.Open(fmt.Sprintf("%s/%s", directory, path.Base(url.Path)))
}
