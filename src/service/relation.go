package service

import (
	"HiChat/src/common"
	"HiChat/src/dao"
	"HiChat/src/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

// a data model that define the User information return to User
type user struct {
	Name   string
	Avatar string
	Gender string
	Phone  string
	Email  string
}

// FriendList Get one's friend list by his userID
func FriendList(ctx *gin.Context) {
	// try to get user id
	userId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "Failed to Get UserId"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	// search friend list by id in DAO
	friendList, err := dao.GetFriendList(uint(userId))
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "friends list is empty"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}

	friends := make([]user, 0)
	for _, f := range *friendList {
		friends = append(friends, user{
			Name:   f.Name,
			Avatar: f.Avatar,
			Gender: f.Gender,
			Phone:  f.Phone,
			Email:  f.Email,
		})
	}

	common.SendNormalResp(ctx.Writer, "Success to Get Friend List", nil, friends, len(friends))
}

// AddFriendByName call DAO to create a relationship between currentUser and targetUser
func AddFriendByName(ctx *gin.Context) {
	// try to get ownerId and targetName
	ownerId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "Failed to Get OwnerId"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	targetName := ctx.PostForm("target_name")

	// Add Friend in DAO
	err = dao.AddFriendByName(uint(ownerId), targetName)
	if err != nil {
		zap.S().Info(err.Error())
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.SendNormalResp(ctx.Writer, "Success to Add Friend", nil, nil, 0)
}

func UpdateRelation(ctx *gin.Context) {
	// try to get ownerId and targetName
	ownerId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "Failed to Get OwnerId"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	targetName := ctx.PostForm("target_name")

	// Get update Information
	r := models.Relation{}

	typeStr := ctx.PostForm("type")
	desc := ctx.PostForm("desc")

	if typeStr != "" {
		typeInt, err := strconv.Atoi(typeStr)
		if err != nil {
			zap.S().Info(err.Error())
			errMsg := "Failed to Get Type"
			common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
			return
		}
		r.Type = typeInt
	}

	if desc != "" {
		r.Desc = desc
	}

	if err = dao.UpdateRelation(uint(ownerId), targetName, r); err != nil {
		zap.S().Info(err.Error())
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.SendNormalResp(ctx.Writer, "Successfully Update", nil, nil, 0)
}

func DelFriendByName(ctx *gin.Context) {
	// try to get ownerId and targetName
	ownerId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "Failed to Get OwnerId"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	targetName := ctx.PostForm("target_name")

	if err := dao.DeleteRelation(uint(ownerId), targetName); err != nil {
		zap.S().Info(err.Error())
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.SendNormalResp(ctx.Writer, "Successfully delete", nil, nil, 0)
}
