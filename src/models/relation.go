package models

import "gorm.io/gorm"

// Relation describes the relations between users
/*
the params are:
	* OwnerId is the user id of the relationship owner
	* TargetId is the user id of the target user
	* Type = 1 means Friends relationship; while Type = 2 means Group relationship
	* Desc store the description message
*/
type Relation struct {
	gorm.Model
	OwnerId  uint
	TargetId uint
	Type     int
	Desc     string
}

func (r *Relation) RelTableName() string {
	return "relation"
}
