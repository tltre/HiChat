package dao

import (
	"HiChat/global"
	"HiChat/models"
	"errors"
	"go.uber.org/zap"
)

// GetGroupList return a group list which user has joined in
func GetGroupList(userId uint) (*[]models.Community, error) {
	// Find all group ID
	relations := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id = ? and type = 2", userId).Find(&relations); tx.RowsAffected == 0 {
		zap.S().Info("User didn't join in any community")
		return nil, nil
	}
	communitiesId := make([]uint, 0)
	for _, r := range relations {
		communitiesId = append(communitiesId, r.TargetId)
	}
	// Get Communities Record
	communities := make([]models.Community, 0)
	if tx := global.DB.Where("id in ?", communitiesId).Find(&communities); tx.RowsAffected == 0 {
		zap.S().Info("Cannot Find Communities Record")
		return nil, errors.New("cannot Find Communities Record")
	}

	return &communities, nil
}

// CreateCommunity Create a community if not exist
func CreateCommunity(community models.Community) error {
	// check if Community has existed
	if tx := global.DB.Where("name = ? and owner_id = ?", community.Name, community.OwnerId).First(&community); tx.RowsAffected != 0 {
		zap.S().Info("Cannot create another group which has same Name and Owner")
		return errors.New("cannot create another group which has same Name and Owner")
	}
	// create new community record in Table Community and Relation
	tx := global.DB.Begin()
	if t := tx.Create(&community); t.RowsAffected == 0 {
		tx.Rollback()
		zap.S().Info("Failed to create community")
		return errors.New("failed to create community")
	}

	relation := models.Relation{}
	relation.OwnerId = community.OwnerId
	relation.TargetId = community.ID
	relation.Type = 2
	relation.Desc = community.Desc

	if t := tx.Create(&relation); t.RowsAffected == 0 {
		tx.Rollback()
		zap.S().Info("Failed to create community")
		return errors.New("failed to create community")
	}
	tx.Commit()
	return nil
}

// FindGroupByName Find group list by group name
func FindGroupByName(name string) (*[]models.Community, error) {
	// Get Communities Record
	communities := make([]models.Community, 0)
	if tx := global.DB.Where("name = ?", name).Find(&communities); tx.RowsAffected == 0 {
		zap.S().Info("Cannot Find Communities Record")
		return nil, errors.New("cannot Find Communities Record")
	}

	return &communities, nil
}

// FindGroupByGid Find group list by group Gid
func FindGroupByGid(groupId string) (*models.Community, error) {
	// Get Community Record
	community := models.Community{}
	if tx := global.DB.Where("group_id = ?", groupId).Find(&community); tx.RowsAffected == 0 {
		zap.S().Info("Cannot Find Community Record")
		return nil, errors.New("cannot Find Community Record")
	}

	return &community, nil
}

// JoinInCommunityByGId join in group by gid
func JoinInCommunityByGId(userId uint, groupId string) error {
	community := models.Community{}
	relation := models.Relation{}

	// check if group exist
	if tx := global.DB.Where("group_id = ?", groupId).First(&community); tx.RowsAffected == 0 {
		zap.S().Info("Target Community did not exist")
		return errors.New("target Community did not exist")
	}

	// check if user has join in before
	cid := community.ID
	if tx := global.DB.Where("owner_id = ? and target_id = ? and type = 2", userId, cid).First(&relation); tx.RowsAffected != 0 {
		zap.S().Info("User had join in before")
		return errors.New("user had join in before")
	}

	// add record in relation table
	relation = models.Relation{}
	relation.OwnerId = userId
	relation.TargetId = cid
	relation.Type = 2

	if tx := global.DB.Create(&relation); tx.RowsAffected == 0 {
		zap.S().Info("failed to join in group")
		return errors.New("failed to join in group")
	}

	return nil
}

// UpdateCommunityInformation update ownerId, name, type, image and desc information by gid
func UpdateCommunityInformation(userId uint, community models.Community) (*models.Community, error) {
	// check if community exist
	group, err := FindGroupByGid(community.GroupId)
	if group.ID == 0 || err != nil {
		zap.S().Info("Group is not exist")
		return nil, errors.New("group is not exist")
	}
	// update
	newCommunity := models.Community{
		Name:  community.Name,
		Type:  community.Type,
		Image: community.Image,
		Desc:  community.Desc,
	}
	// Only group owner can modify OwnerId
	if group.OwnerId == userId && community.OwnerId != 0 {
		newCommunity.OwnerId = community.OwnerId
	}
	tx := global.DB.Model(&community).Where("group_id = ?", community.GroupId).Updates(&newCommunity)
	if tx.RowsAffected == 0 {
		zap.S().Info("Failed to update")
		return nil, errors.New("failed to update")
	}
	return FindGroupByGid(community.GroupId)
}

// DelGroup Delete the group record if user is group owner, otherwise Quit the group
func DelGroup(userId uint, gid string) (string, error) {
	// check if community exist
	group, err := FindGroupByGid(gid)
	if group.ID == 0 || err != nil {
		zap.S().Info("Group is not exist")
		return "", errors.New("group is not exist")
	}
	// check if the user is owner
	if group.OwnerId == userId {
		// delete record in Community & Relation
		tx := global.DB.Begin()
		if tx := global.DB.Where("id = ? ", group.ID).Delete(&models.Community{}); tx.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Info("Failed to delete in Table Community")
			return "", errors.New("failed to delete")
		}

		if tx := global.DB.Where("target_id = ? and type = 2", group.ID).Delete(&models.Relation{}); tx.RowsAffected == 0 {
			tx.Rollback()
			zap.S().Info("Failed to delete in Table Relation")
			return "", errors.New("failed to delete")
		}
		tx.Commit()
		return "Successfully Delete the group", nil
	} else {
		// delete record in Relation
		if tx := global.DB.Where("owner_id = ? and target_id = ? and type = 2", userId, group.ID).Delete(&models.Relation{}); tx.RowsAffected == 0 {
			zap.S().Info("Failed to delete")
			return "", errors.New("failed to delete")
		}
		return "Successfully quit the group", nil
	}
}
