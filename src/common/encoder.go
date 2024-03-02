package common

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
)

// Md5encoder Encode the PassWord by Md5
func Md5encoder(code string) string {
	m := md5.New()
	io.WriteString(m, code)
	return hex.EncodeToString(m.Sum(nil))
}

// CheckPassword Encode the input password and compare the password in DB
func CheckPassword(curPwd string, salt string, dbPwd string) bool {
	pwd := SaltPassword(curPwd, salt)
	return pwd == dbPwd
}

// SaltPassword return the password with salt
func SaltPassword(pwd string, salt string) string {
	saltPwd := fmt.Sprintf("%s%s", Md5encoder(pwd), salt)
	return saltPwd
}
