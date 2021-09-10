package pkg

import (
	"database/sql"
)

type AssignedPermission struct {
	ID           int64          `gorm:"column:id;primarykey" json:"id"`
	PermissionId sql.NullInt64  `gorm:"column:permission_id" json:"permission_id"`
	TeamId       sql.NullString `gorm:"column:team_id" json:"team_id"`
	RoleId       sql.NullInt64  `gorm:"column:role_id" json:"role_id"`
	ResourceId   string         `gorm:"column:resource_id;not null;" json:"resource_id"`
	ResourceName string         `gorm:"column:resource_name;not null;" json:"resource_name"`
	Action       string         `gorm:"column:action;not null;" json:"action"`
}
