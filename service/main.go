package main

import (
	model2 "shorturl/modle"
	"shorturl/routers"
)

func main() {

	model2.InitRedis()
	model2.InitDb()
	routers.InitRouter()
}
