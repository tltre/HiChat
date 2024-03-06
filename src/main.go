package main

import (
	"HiChat/src/initialize"
	"HiChat/src/router"
)

func main() {
	// initialize work
	initialize.InitLogger()
	initialize.InitDB()
	initialize.InitRedis()
	println("successfully initialize!")

	// start the router (gin service)
	r := router.Router()
	r.Run(":8000")
}
