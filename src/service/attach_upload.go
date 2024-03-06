package service

import (
	"HiChat/src/common"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

func UploadFile(ctx *gin.Context) {
	req := ctx.Request

	// Get File
	srcFile, head, err := req.FormFile("file")
	if err != nil {
		zap.S().Info("Failed to Get File")
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, "failed to get file", nil)
		return
	}

	// Get File Suffix
	var suffix string
	filename := head.Filename
	terms := strings.Split(filename, ".")
	if len(terms) > 1 {
		suffix = "." + terms[len(terms)-1]
	} else {
		zap.S().Info("Failed to get suffix")
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, "Cannot get suffix", nil)
		return
	}

	// Store the File in Project Server
	newFileName := fmt.Sprintf("%d_%04d%s", time.Now().Unix(), rand.Int31(), suffix)
	dstFile, err := os.Create("./src/asset/upload" + newFileName)
	if err != nil {
		zap.S().Info("Failed to Create new File")
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, "failed to store file", nil)
		return
	}
	if _, err = io.Copy(dstFile, srcFile); err != nil {
		zap.S().Info("Failed to Create new File")
		common.SendErrorResp(ctx.Writer, http.StatusBadRequest, "failed to store file", nil)
		return
	}

	data := make(map[string]string, 0)
	data["url"] = "./src/asset/upload" + newFileName
	common.SendNormalResp(ctx.Writer, "Success to Upload File", data, nil, 0)
}
