package assets

import (
	"encoding/json"
	"os"
)

const (
	assets  = "assets/"
	members = assets + "members/"
	index   = "index.json"
)

type AssetIndex struct {
	Members map[string]string `json:"members"`
	Title   string            `json:"title"`
}

var pictureIndex = make(map[string]string)

const pictureIndexName = "member-index.json"

var lastUpdate int64 = 0

//loads filenames from disk if necessary
func loadIndex() {
	file, err := openIndex()
	if err != nil {
		errLogger.Println(err.Error())
		return
	}
	defer file.Close()
	stat, _ := file.Stat()
	if stat.ModTime().Unix() > lastUpdate {
		decoder := json.NewDecoder(file)
		decoder.Decode(&pictureIndex)
		lastUpdate = stat.ModTime().Unix()
	}
}

//saves filenames to disk
func saveIndex() {
	file, err := openIndex()
	if err != nil {
		errLogger.Println(err.Error())
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(&pictureIndex)
}

//returns the filename if a members
func getMemberPictureName(id string) string {
	loadIndex()
	return pictureIndex[id]
}

//saves the filename of a members
func setMemberPictureName(id string, file string) {
	pictureIndex[id] = file
	go saveIndex()
}

//returns the members index
func openIndex() (*os.File, error) {
	return os.OpenFile(pictureIndexName, os.O_CREATE|os.O_RDWR, 0666)
}
