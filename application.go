package assets

import "os"

func Setup() {
	tree()
	routes()
	run()
}

func tree() {
	os.MkdirAll("members", os.ModeDir)
	os.MkdirAll("title", os.ModeDir)
}
