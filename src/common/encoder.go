package common

func Md5encoder(str string) string {
	return ""
}

func CheckPassword(curPwd string, salt string, dbPwd string) bool {
	return true
}

func SaltPassword(pwd string, salt string) string {
	return pwd
}
