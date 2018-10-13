package assets

import "os"

func Setup() {
	tree()
	routes()
	run()
}

func tree() {
	os.MkdirAll("assets/members", os.ModeDir)
	os.MkdirAll("assets/title", os.ModeDir)
}
