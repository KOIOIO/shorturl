package main

import (
	model "shorturl/modle"
	"shorturl/routers"
	"shorturl/server"
)

func main() {
	server.DeleteWithTime()
	model.InitRedis()
	model.InitDb()
	routers.InitRouter()
}
