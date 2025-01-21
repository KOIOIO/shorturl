package main

import (
	model "shorturl/modle"
	"shorturl/routers"
)

func main() {

	model.InitRedis()
	model.InitDb()
	routers.InitRouter()
}
