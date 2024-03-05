package common

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

// GenerateId Generate user-visible ID from raw ID
/* Just use 5 chars as user-visible ID, because of little amount of groups and users */
func GenerateId(id uint) string {
	hash := sha256.Sum256([]byte(strconv.Itoa(int(id))))
	return hex.EncodeToString(hash[:])[:5]
}
