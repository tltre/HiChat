package dao

import (
	"HiChat/global"
	"HiChat/models"
	"errors"
	"go.uber.org/zap"
)

// GetFriendList return Friend List by user ID
func GetFriendList(userId uint) (*[]models.UserBasic, error) {
	relations := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_Id = ? and type = 1", userId).Find(&relations); tx.RowsAffected == 0 {
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

// GetRelationId will return the id of the relation (0 if relation didn't exist) by userId and targetId
func GetRelationId(userId uint, targetId uint) uint {
	relation := models.Relation{}
	tx := global.DB.Where("owner_id = ? and target_id = ? and type = 1", userId, targetId).First(&relation)
	if tx.RowsAffected == 0 {
		zap.S().Info("Relation did not exist")
		return 0
	}
	return relation.ID
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
	if GetRelationId(userId, targetId) != 0 {
		zap.S().Info("Relation has already existed")
		return errors.New("relation has already existed")
	}

	if GetRelationId(targetId, userId) != 0 {
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

// UpdateFriendRelation Update type and desc of target friend relation record
func UpdateFriendRelation(userId uint, targetName string, r models.Relation) error {
	// check if targetUser exist
	targetUser, err := GetUserByNameForLoginIn(targetName)
	if err != nil || targetUser.ID == 0 {
		zap.S().Info("Target User is not exist")
		return errors.New("target User is not exist")
	}

	targetId := targetUser.ID
	if userId == targetId {
		zap.S().Info("userId cannot equal to targetId")
		return errors.New("userId cannot equal to targetId")
	}

	// check if relation had existed
	r1 := GetRelationId(userId, targetId)
	r2 := GetRelationId(targetId, userId)
	if r1 == 0 || r2 == 0 {
		zap.S().Info("Relation didn't exist")
		return errors.New("relation didn't exist")
	}

	// open transaction
	tx := global.DB.Begin()

	t := tx.Model(&r).Where("id = ?", r1).Updates(&r)
	if t.RowsAffected == 0 {
		tx.Rollback()
		zap.S().Info("Failed to Update Relation")
		return errors.New("failed to Update Relation")
	}

	t = tx.Model(&r).Where("id = ?", r2).Updates(&r)
	if t.RowsAffected == 0 {
		tx.Rollback()
		zap.S().Info("Failed to Update Relation")
		return errors.New("failed to Update Relation")
	}

	tx.Commit()
	return nil
}

// UpdateGroupRelation update relation by gid
func UpdateGroupRelation(gid string, r models.Relation) error {
	group, err := FindGroupByGid(gid)
	if err != nil {
		zap.S().Info("Target Group is not exist")
		return errors.New("target group is not exist")
	}
	if tx := global.DB.Model(&r).Where("id = ?", group.ID).Updates(&r); tx.RowsAffected == 0 {
		zap.S().Info("Failed to Update Relation")
		return errors.New("failed to Update Relation")
	}
	return nil
}

// DeleteFriendRelation delete friend relation
func DeleteFriendRelation(userId uint, targetName string) error {
	// check if targetUser exist
	targetUser, err := GetUserByNameForLoginIn(targetName)
	if err != nil || targetUser.ID == 0 {
		zap.S().Info("Target User is not exist")
		return errors.New("target User is not exist")
	}

	targetId := targetUser.ID
	if userId == targetId {
		zap.S().Info("userId cannot equal to targetId")
		return errors.New("userId cannot equal to targetId")
	}

	// check if relation had existed
	r1 := GetRelationId(userId, targetId)
	r2 := GetRelationId(targetId, userId)
	if r1 == 0 && r2 == 0 {
		zap.S().Info("Relations didn't exist")
		return errors.New("relations didn't exist")
	}

	var delIds = []uint{r1, r2}
	if tx := global.DB.Delete(&models.Relation{}, delIds); tx.RowsAffected == 0 {
		zap.S().Info("Failed to Delete")
		return errors.New("failed to Delete")
	}

	return nil
}
