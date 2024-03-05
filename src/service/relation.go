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

// UpdateRelation update relation desc
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

	desc := ctx.PostForm("desc")

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

// DelFriendByName Delete friend by name
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

// GetGroupList return group list that user has joined in
func GetGroupList(ctx *gin.Context) {
	// try to get owner id
	ownerId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "Failed to Get OwnerId"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	// Get Group List by ownerId
	communities, err := dao.GetGroupList(uint(ownerId))
	if err != nil {
		zap.S().Info(err.Error())
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.SendNormalResp(ctx.Writer, "Successfully find group!", nil, *communities, len(*communities))
}

// CreateGroup create a group by userId
func CreateGroup(ctx *gin.Context) {
	// try to get community information
	ownerId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "Failed to Get OwnerId"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	tp, err := strconv.Atoi(ctx.PostForm("type"))
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "Failed to Get type"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	image := ctx.PostForm("image")
	desc := ctx.PostForm("desc")

	// create community record
	community := models.Community{
		Name:    "",
		OwnerId: uint(ownerId),
		Type:    tp,
		Image:   image,
		Desc:    desc,
	}
	if err = dao.CreateCommunity(community); err != nil {
		zap.S().Info(err.Error())
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.SendNormalResp(ctx.Writer, "Successfully Create group!", nil, nil, 0)
}

// SearchGroup return group list that has target name or gid
func SearchGroup(ctx *gin.Context) {
	groupName := ctx.PostForm("group_name")
	gid := ctx.PostForm("group_id")

	if gid != "" {
		community, err := dao.FindGroupByGid(gid)
		if err != nil {
			zap.S().Info(err.Error())
			common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		common.SendNormalResp(ctx.Writer, "Successfully find group!", nil, community, 1)
	} else if groupName != "" {
		communities, err := dao.FindGroupByName(groupName)
		if err != nil {
			zap.S().Info(err.Error())
			common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, err.Error(), nil)
			return
		}
		common.SendNormalResp(ctx.Writer, "Successfully find group!", nil, *communities, len(*communities))
	} else {
		zap.S().Info("Don't have necessary params")
		errMsg := "please add necessary params, such as group_name or group_id"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
	}
}

// JoinGroup Join in Group by GID
func JoinGroup(ctx *gin.Context) {
	// try to get owner id and gid
	ownerId, err := strconv.Atoi(ctx.PostForm("id"))
	if err != nil {
		zap.S().Info(err.Error())
		errMsg := "Failed to Get OwnerId"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	gid := ctx.PostForm("group_id")
	if gid == "" {
		zap.S().Info("Don't have necessary params")
		errMsg := "please add necessary params: group_id"
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, errMsg, nil)
		return
	}
	// Join in group
	err = dao.JoinInCommunityByGId(uint(ownerId), gid)
	if err != nil {
		zap.S().Info(err.Error())
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	common.SendNormalResp(ctx.Writer, "Successfully join group!", nil, nil, 0)
}
