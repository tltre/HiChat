package service

import (
	"HiChat/common"
	"HiChat/dao"
	"HiChat/middleware"
	"HiChat/models"
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
		zap.S().Info("DB GetUserList Failed")
		errMsg := "GetUserList Error"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	common.SendNormalResp(ctx.Writer, "Success to Get User List", nil, list, len(list))
}

// UserLoginByNameAndPwd Post method
func UserLoginByNameAndPwd(ctx *gin.Context) {
	username := ctx.PostForm("name")
	pwd := ctx.PostForm("password")
	user, err := dao.GetUserByNameForLoginIn(username)
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "Failed to Login in"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}

	if user.Name == "" {
		errMsg := "Username didn't exist"
		zap.S().Info(errMsg)
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, errMsg, nil)
		return
	}

	// password in database is encrypted, thus we should encrypt again to compare
	ok := common.CheckPassword(pwd, user.Salt, user.PassWord)
	if !ok {
		errMsg := "Password incorrect"
		zap.S().Info(errMsg)
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, errMsg, nil)
		return
	}

	Rsp, err := dao.GetUserByNameAndPwd(username, user.PassWord)
	if err != nil {
		zap.S().Info(err.Error())
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, "Failed to Login in", nil)
		return
	}

	// JWT identify
	token, err := middleware.GenerateToken(Rsp.ID, "xy")
	if err != nil {
		zap.S().Info("Failed to generate token")
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, "Failed to generate token", nil)
		return
	}

	data := make(map[string]string)
	data["token"] = token
	data["userId"] = strconv.Itoa(int(Rsp.ID))
	common.SendNormalResp(ctx.Writer, "Success to Login in", data, user, 1)
}

// UserRegister Post Method
func UserRegister(ctx *gin.Context) {
	user := models.UserBasic{}
	username := ctx.PostForm("name")
	pwd := ctx.PostForm("password")
	rePwd := ctx.PostForm("Identity")

	data := make(map[string]string)

	// Username and Password can not be empty
	if username == "" || pwd == "" || rePwd == "" {
		errMsg := "Username and Password can not be empty"
		data["username"] = username
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, errMsg, data)
		return
	}

	//Check if the Identity matches the password
	if pwd != rePwd {
		errMsg := "Password is not equal to Identity"
		zap.S().Info(errMsg)
		data["username"] = username
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, errMsg, nil)
		return
	}

	// Check if the username exists
	_, err := dao.GetUserByNameForRegister(username)
	if err != nil {
		errMsg := "Username is exist"
		zap.S().Info(errMsg)
		data["username"] = username
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, errMsg, nil)
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
		errMsg := "Create User in DB error"
		zap.S().Info(errMsg)
		data["username"] = username
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}

	common.SendNormalResp(ctx.Writer, "Success to Register", nil, user, 1)
}

// UpdateUserInformation 更新用户信息
func UpdateUserInformation(ctx *gin.Context) {
	user := models.UserBasic{}

	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info("Failed to Exchange type")
		errMsg := "Failed to Update User Information"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
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
		errMsg := "Params didn't validated"
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, errMsg, nil)
		return
	}

	// Update User Information in DB
	err = dao.UpdateUser(user)
	if err != nil {
		zap.S().Info("Update User in DB error")
		errMsg := "Failed to update User"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}

	common.SendNormalResp(ctx.Writer, "Success to Update User", nil, user, 1)
}

// DeleteUser Del Method
func DeleteUser(ctx *gin.Context) {
	user := models.UserBasic{}

	id, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info("Failed to Exchange type")
		errMsg := "Failed to delete User"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}

	user.ID = uint(id)
	err = dao.DeleteUser(user)
	if err != nil {
		zap.S().Info("Delete User in DB error")
		errMsg := "Failed to delete User"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}

	common.SendNormalResp(ctx.Writer, "Success to Delete User", nil, nil, 0)
}
