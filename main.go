package main

import (
	model "shorturl/modle"
	"shorturl/routers"
)

func main() {
	// TODO: Implement me!
	model.InitRedis()
	model.InitDb()
	routers.InitRouter()
}
