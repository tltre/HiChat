package service

import (
	"HiChat/common"
	"HiChat/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RedisMsg Get message from Redis
func RedisMsg(ctx *gin.Context) {
	userIdA, _ := strconv.Atoi(ctx.Query("userId"))
	userIdB, _ := strconv.Atoi(ctx.Query("targetId"))
	start, _ := strconv.Atoi(ctx.PostForm("start"))
	end, _ := strconv.Atoi(ctx.PostForm("end"))
	isRev, _ := strconv.ParseBool(ctx.PostForm("isRev"))
	res := models.GetMsgFromRedis(uint(userIdA), uint(userIdB), int64(start), int64(end), isRev)
	if res == nil {
		common.SendErrorResp(ctx.Writer, http.StatusInternalServerError, "Failed to get records", nil)
	} else {
		common.SendNormalResp(ctx.Writer, "Success to get records", nil, res, len(res))
	}
}

// SendMsg user send message to friend/group
func SendMsg(ctx *gin.Context) {
	models.Chat(ctx.Writer, ctx.Request)
}
