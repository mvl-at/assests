package assets

import "os"

func Setup() {
	tree()
	routes()
	run()
}

func tree() {
	os.MkdirAll("members", os.ModePerm)
	os.MkdirAll("title", os.ModePerm)
}
