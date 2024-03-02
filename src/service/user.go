package service

import (
	"HiChat/src/common"
	"HiChat/src/dao"
	"HiChat/src/middleware"
	"HiChat/src/models"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// UserList Get Method, Provided for Admin
func UserList(ctx *gin.Context) {
	list, err := dao.GetUserList()
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1, // -1 means error
			"msg":  "GetUserList Error",
		})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

// UserLoginByNameAndPwd Post method
func UserLoginByNameAndPwd(ctx *gin.Context) {
	username := ctx.PostForm("name")
	pwd := ctx.PostForm("password")
	data, err := dao.GetUserByNameForLoginIn(username)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1, // -1 means error
			"msg":  "Failed to Login in",
		})
		return
	}

	if data.Name == "" {
		ctx.JSON(200, gin.H{
			"code": -1, // -1 means error
			"msg":  "Username didn't exist",
		})
		return
	}

	// password in database is encrypted, thus we should encrypt again to compare
	ok := common.CheckPassword(pwd, data.Salt, data.PassWord)
	if !ok {
		ctx.JSON(200, gin.H{
			"code": -1, // -1 means error
			"msg":  "Password incorrect",
		})
		return
	}

	Rsp, err := dao.GetUserByNameAndPwd(username, data.PassWord)
	if err != nil {
		zap.S().Info("Failed to Login in")
		return
	}

	// JWT identify
	token, err := middleware.GenerateToken(Rsp.ID, "xy")
	if err != nil {
		zap.S().Info("Failed to generate token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":   0,
		"msg":    "Success to Login in",
		"token":  token,
		"userId": Rsp.ID,
	})
}

// UserRegister Post Method
func UserRegister(ctx *gin.Context) {
	user := models.UserBasic{}
	username := ctx.PostForm("name")
	pwd := ctx.PostForm("password")
	rePwd := ctx.PostForm("identity")

	// Username and Password can not be empty
	if username == "" || pwd == "" || rePwd == "" {
		ctx.JSON(200, gin.H{
			"code": -1, // -1 means error
			"msg":  "Username and Password can not be empty",
			"data": username,
		})
		return
	}

	//Check if the Identity matches the password
	if pwd != rePwd {
		ctx.JSON(200, gin.H{
			"code": -1, // -1 means error
			"msg":  "Password is not equal to Identity",
			"data": username,
		})
		return
	}

	// Check if the username exists
	_, err := dao.GetUserByNameForRegister(username)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": -1, // -1 means error
			"msg":  "Username is exist",
			"data": username,
		})
		return
	}

	// Generate Salt value
	salt := fmt.Sprintf("%d", rand.Int31())
	user.Salt = salt

	// encrypted password
	user.PassWord = common.SaltPassword(pwd, salt)

	// Fill in User information
	user.Name = username

	t := time.Now()
	user.LoginTime = &t
	user.LoginOutTime = &t
	user.HeartBeatTime = &t

	// create user in DB
	err = dao.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, // -1 means error
			"msg":  "Create User in DB error",
			"data": username,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "Success to Register",
		"data": user,
	})
}

// UpdateUserInformation 更新用户信息
func UpdateUserInformation(ctx *gin.Context) {
	user := models.UserBasic{}

	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info("Failed to Exchange type")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, // -1 means error
			"msg":  "Failed to Update User Information",
		})
		return
	}

	user.ID = uint(id)

	Name := ctx.Request.FormValue("name")
	PassWord := ctx.Request.FormValue("password")
	Email := ctx.Request.FormValue("email")
	Phone := ctx.Request.FormValue("phone")
	avatar := ctx.Request.FormValue("icon")
	gender := ctx.Request.FormValue("gender")

	if Name != "" {
		user.Name = Name
	}
	if PassWord != "" {
		salt := fmt.Sprintf("%d", rand.Int31())
		user.Salt = salt
		user.PassWord = common.SaltPassword(PassWord, salt)
	}
	if Email != "" {
		user.Email = Email
	}
	if Phone != "" {
		user.Phone = Phone
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if gender != "" {
		user.Gender = gender
	}

	// ValidateStruct, which will check if the Params in Struct with tag "valid" is validated
	_, err = govalidator.ValidateStruct(user)
	if err != nil {
		zap.S().Info("Params didn't match")
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": -1, // -1 means error
			"msg":  "Params didn't validated",
		})
		return
	}

	// Update User Information in DB
	err = dao.UpdateUser(user)
	if err != nil {
		zap.S().Info("Update User in DB error")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, // -1 means error
			"msg":  "Failed to update User",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0, // -1 means error
		"msg":  "Success to Update User",
		"data": user,
	})
}

// DeleteUser Del Method
func DeleteUser(ctx *gin.Context) {
	user := models.UserBasic{}

	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info("Failed to Exchange type")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, // -1 means error
			"msg":  "Failed to Delete User",
		})
		return
	}

	user.ID = uint(id)
	err = dao.DeleteUser(user)
	if err != nil {
		zap.S().Info("Delete User in DB error")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": -1, // -1 means error
			"msg":  "Failed to delete User",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code": 0, // -1 means error
		"msg":  "Success to Delete User",
	})
}
