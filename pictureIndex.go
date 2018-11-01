package assets

import (
	"encoding/json"
	"os"
)

const (
	memberDir = "members"
	titleDir  = "title"
	index     = "index.json"
)

var assetIndex = loadAssesIndex()

type AssetIndex struct {
	Members        map[string]string `json:"members"`
	Title          string            `json:"title"`
	DefaultTitle   string            `json:"defaultTitle"`
	IsDefaultTitle bool              `json:"isDefaultTitle"`
	lastUpdate     int64             `json:"-"`
	fileName       string            `json:"-"`
}

//loads filenames from disk if necessary
func (a *AssetIndex) load() {
	file, err := a.openIndex()
	if err != nil {
		errLogger.Println(err.Error())
		return
	}
	defer file.Close()
	stat, _ := file.Stat()
	if stat.ModTime().Unix() > a.lastUpdate {
		decoder := json.NewDecoder(file)
		decoder.Decode(a)
		a.lastUpdate = stat.ModTime().Unix()
	}
}

//saves filenames to disk
func (a *AssetIndex) save() {
	file, err := a.openIndex()
	if err != nil {
		errLogger.Println(err.Error())
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(a)
}

//returns the filename if a members
func (a *AssetIndex) getMemberPictureName(id string) string {
	a.load()
	return a.Members[id]
}

//checks if a member has a picture
func (a *AssetIndex) memberHasPicture(id string) bool {
	a.load()
	return a.Members[id] != ""
}

//saves the filename of a members
func (a *AssetIndex) setMemberPictureName(id string, file string) {
	a.Members[id] = file
	go a.save()
}

//returns the filename of the title picture
func (a *AssetIndex) getTitlePictureName() string {
	a.load()
	return assetIndex.Title
}

//saves the filename of the title picture
func (a *AssetIndex) setTitlePictureName(file string) {
	assetIndex.Title = file
	go a.save()
}

//returns the filename of the default title picture
func (a *AssetIndex) getDefaultTitlePictureName() string {
	a.load()
	return assetIndex.DefaultTitle
}

//saves the filename of the default title picture
func (a *AssetIndex) setDefaultTitlePictureName(file string) {
	assetIndex.DefaultTitle = file
	go a.save()
}

//returns if the title picture is currently the default one
func (a *AssetIndex) getIsDefaultTitle() bool {
	a.load()
	return a.IsDefaultTitle
}

//sets if the title picture is currently the default one
func (a *AssetIndex) setIsDefaultTitle(isDefault bool) {
	a.IsDefaultTitle = isDefault
	go a.save()
}

//returns the members index
func (a *AssetIndex) openIndex() (*os.File, error) {
	return os.OpenFile(a.fileName, os.O_CREATE|os.O_RDWR, 0666)
}

//load the default assets file
func loadAssesIndex() *AssetIndex {
	index := &AssetIndex{fileName: index, lastUpdate: 0, Members: make(map[string]string)}
	index.load()
	return index
}
