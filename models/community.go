package models

import (
	"HiChat/common"
	"HiChat/global"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Community describes the information of group
/*
the params are:
	* Name: the name of group
	* GroupId: the user-visible group id
	* OwnerId: the group owner's userId
	* Type: type of group
	* Image: icon of group
	* Desc: describe of group
*/
type Community struct {
	gorm.Model
	Name    string
	GroupId string
	OwnerId uint
	Type    int
	Image   string
	Desc    string
}

// AfterCreate Hook function, generate group id by ID
func (c *Community) AfterCreate(tx *gorm.DB) error {
	if t := tx.Model(c).Update("group_id", common.GenerateId(c.ID)); t.RowsAffected == 0 {
		zap.S().Info("failed to add group id")
		return errors.New("failed to add group id")
	}
	return nil
}

// FindMembersId find member id by community id
func FindMembersId(id uint) (*[]uint, error) {
	relation := make([]Relation, 0)
	if tx := global.DB.Where("target_id = ? and type = 2", id).Find(&relation); tx.RowsAffected == 0 {
		zap.S().Info("Didn't Find any Member")
		return nil, errors.New("didn't Find any Member")
	}
	membersId := make([]uint, 0)
	for _, r := range relation {
		membersId = append(membersId, r.OwnerId)
	}
	return &membersId, nil
}
