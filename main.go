package main

import (
	"HiChat/global"
	"HiChat/initialize"
	"HiChat/router"
	"fmt"
)

func main() {
	// initialize work
	initialize.InitLogger()
	initialize.InitConfig("debug")
	initialize.InitDB()
	initialize.InitRedis()
	println("successfully initialize!")

	// start the router (gin service)
	r := router.Router()
	r.Run(fmt.Sprintf(":%d", global.ServiceConfig.Port))
}
