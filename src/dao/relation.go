package dao

import (
	"HiChat/src/global"
	"HiChat/src/models"
	"errors"
	"go.uber.org/zap"
)

// GetFriendList return Friend List by user ID
func GetFriendList(userId uint) (*[]models.UserBasic, error) {
	relations := make([]models.Relation, 0)
	if tx := global.DB.Where("OwnerId = ? and type = 1", userId).Find(&relations); tx.RowsAffected == 0 {
		zap.S().Info("Didn't Find relation data")
		return nil, errors.New("failed to get friend list")
	}

	// Get Friends ID List
	friendsId := make([]uint, 0)
	for _, r := range relations {
		friendsId = append(friendsId, r.TargetId)
	}

	var friends *[]models.UserBasic
	if tx := global.DB.Where("id in ?", friendsId).Find(&friends); tx.RowsAffected == 0 {
		zap.S().Info("Didn't Find friends data")
		return nil, errors.New("failed to get friend list")
	}

	return friends, nil
}

// AddFriendById should open transaction to guarantee consistent
func AddFriendById(userId uint, targetId uint) error {
	if userId == targetId {
		zap.S().Info("userId cannot equal to targetId")
		return errors.New("cannot add yourself as friends")
	}

	// check if targetId exist
	targetUser, err := GetUserById(targetId)
	if err != nil || targetUser.ID == 0 {
		zap.S().Info("Target User is not exist")
		return errors.New("target User is not exist")
	}

	// check if relation had existed
	relation := models.Relation{}
	if tx := global.DB.Where("OwnerId = ? and TargetId = ? and type = 1", userId, targetId).First(&relation); tx.RowsAffected != 0 {
		zap.S().Info("Relation has already existed")
		return errors.New("relation has already existed")
	}

	if tx := global.DB.Where("OwnerId = ? and TargetId = ? and type = 1", targetId, userId).First(&relation); tx.RowsAffected != 0 {
		zap.S().Info("Relation has already existed")
		return errors.New("relation has already existed")
	}

	// open transaction
	tx := global.DB.Begin()

	relation.OwnerId = userId
	relation.TargetId = targetId
	relation.Type = 1
	if t := tx.Create(&relation); t.RowsAffected == 0 {
		tx.Rollback()
		zap.S().Info("Failed to Create Relation")
		return errors.New("failed to Create Relation")
	}

	relation = models.Relation{}
	relation.OwnerId = targetId
	relation.TargetId = userId
	relation.Type = 1
	if t := tx.Create(&relation); t.RowsAffected == 0 {
		tx.Rollback()
		zap.S().Info("Failed to Create Relation")
		return errors.New("failed to Create Relation")
	}

	tx.Commit()
	return nil
}

// AddFriendByName Find Target User Id and call AddFriendById
func AddFriendByName(userId uint, targetName string) error {
	// check if targetUser exist
	targetUser, err := GetUserByNameForLoginIn(targetName)
	if err != nil || targetUser.ID == 0 {
		zap.S().Info("Target User is not exist")
		return errors.New("target User is not exist")
	}

	return AddFriendById(userId, targetUser.ID)
}
