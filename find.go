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
	}
	return os.Open(fmt.Sprintf("%s/%s", directory, path.Base(url.Path)))
}
