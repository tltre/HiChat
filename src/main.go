package main

import "HiChat/src/initialize"

func main() {
	initialize.InitLogger()
	initialize.InitDB()
	println("successfully initialize!")
}
