package router

import (
	"HiChat/src/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	// initial a router
	router := gin.Default()

	// create a router group "v1"
	// the router group in this layer divided by VERSION
	// offer service for the URL likes "/v1/....."
	v1 := router.Group("v1")

	// a router group of User Module, store the User API
	// the router group in this layer divided by MODULE
	// offer service for the URL likes "/v1/user/....."
	user := v1.Group("user")
	{
		user.GET("/list", service.UserList)
		user.POST("/login", service.UserLoginByNameAndPwd)
		user.POST("/new", service.UserRegister)
		user.POST("/update", service.UpdateUserInformation)
		user.DELETE("/delete", service.DeleteUser)
	}

	return router
}
