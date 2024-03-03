package dao

import (
	"HiChat/src/common"
	"HiChat/src/global"
	"HiChat/src/models"
	"errors"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// GetUserList list All Users
func GetUserList() ([]*models.UserBasic, error) {
	var users []*models.UserBasic
	if tx := global.DB.Find(&users); tx.RowsAffected == 0 {
		return nil, errors.New("get User List Failed")
	}
	return users, nil
}

// GetUserByNameAndPwd Query User by name and pwd, Used in verify when Login in
func GetUserByNameAndPwd(name string, pwd string) (*models.UserBasic, error) {
	var user models.UserBasic
	if tx := global.DB.Where("name = ? and pass_word = ?", name, pwd).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("query User by name and pwd Failed")
	}

	// Login in Identify
	curTime := strconv.Itoa(int(time.Now().UnixNano()))
	md5time := common.Md5Encoder(curTime)
	if tx := global.DB.Model(&user).Where("id = ?", user.ID).Update("identity", md5time); tx.RowsAffected == 0 {
		return nil, errors.New("update Identity Failed")
	}

	return &user, nil
}

// GetUserByNameForLoginIn Query User by Name, Used in Login in
func GetUserByNameForLoginIn(name string) (*models.UserBasic, error) {
	var user models.UserBasic
	if tx := global.DB.Where("name = ?", name).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("didn't find the user")
	}
	return &user, nil
}

// GetUserByNameForRegister Query User by Name, Used in Register
func GetUserByNameForRegister(name string) (*models.UserBasic, error) {
	var user models.UserBasic
	if tx := global.DB.Where("name = ?", name).First(&user); tx.RowsAffected != 0 {
		return nil, errors.New("username has existed")
	}
	return &user, nil
}

// GetUserById Query User by id, Used in Login in
func GetUserById(id uint) (*models.UserBasic, error) {
	var user models.UserBasic
	if tx := global.DB.Where("id = ?", id).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("didn't find the user")
	}
	return &user, nil
}

// GetUserByPhone Query User by Phone, Used in Login in
func GetUserByPhone(phone string) (*models.UserBasic, error) {
	var user models.UserBasic
	if tx := global.DB.Where("phone = ?", phone).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("didn't find the user")
	}
	return &user, nil
}

// GetUserByEmail Query User by Email, Used in Login in
func GetUserByEmail(email string) (*models.UserBasic, error) {
	var user models.UserBasic
	if tx := global.DB.Where("email = ?", email).First(&user); tx.RowsAffected == 0 {
		return nil, errors.New("didn't find the user")
	}
	return &user, nil
}

// CreateUser create User
func CreateUser(user models.UserBasic) error {
	tx := global.DB.Create(&user)
	if tx.RowsAffected == 0 {
		// Log the Error
		zap.S().Info("Create User Failed")
		return errors.New("create User Failed")
	}
	return nil
}

// UpdateUser modifier the User Information
func UpdateUser(user models.UserBasic) error {
	tx := global.DB.Model(&user).Updates(models.UserBasic{
		Name:     user.Name,
		PassWord: user.PassWord,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Phone:    user.Phone,
		Email:    user.Email,
		Salt:     user.Salt,
	})
	if tx.RowsAffected == 0 {
		// Log the Error
		zap.S().Info("Update User Failed")
		return errors.New("update User Failed")
	}
	return nil
}

// DeleteUser Delete User
func DeleteUser(user models.UserBasic) error {
	tx := global.DB.Delete(&user)
	if tx.RowsAffected == 0 {
		// Log the Error
		zap.S().Info("Delete User Failed")
		return errors.New("delete User Failed")
	}
	return nil
}
