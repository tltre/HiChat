package router

import (
	"HiChat/src/middleware"
	"HiChat/src/service"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	// initial a router
	// the router default use Logger and Recovery MiddleWare
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
		user.GET("/list", middleware.Authentication(), service.UserList)
		user.POST("/login", service.UserLoginByNameAndPwd)
		user.POST("/new", service.UserRegister)
		user.POST("/update", middleware.Authentication(), service.UpdateUserInformation)
		user.DELETE("/delete", middleware.Authentication(), service.DeleteUser)
	}

	// Relation Module
	relation := v1.Group("relation").Use(middleware.Authentication())
	{
		// Friends API
		relation.GET("/list", service.FriendList)
		relation.POST("/add", service.AddFriendByName)
		relation.POST("/update", service.UpdateRelation)
		relation.DELETE("/delete", service.DelFriendByName)

		// Group API
		relation.GET("/group-list", service.GetGroupList)
		relation.POST("/new", service.CreateGroup)
		relation.GET("/search", service.SearchGroup)
		relation.POST("/join", service.JoinGroup)
		relation.POST("/update-group", service.UpdateGroup)
		relation.DELETE("/delete-group", service.DelGroup)
	}

	// Message Module
	message := v1.Group("message").Use(middleware.Authentication())
	{
		message.GET("/get-records", service.RedisMsg)
		message.POST("/send", service.SendMsg)
	}

	// File Upload Module
	v1.POST("/upload", service.UploadFile)

	return router
}
