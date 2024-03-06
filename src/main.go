package main

import (
	"HiChat/src/global"
	"HiChat/src/initialize"
	"HiChat/src/router"
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
