package pkg

import (
	"database/sql"
	"time"
)

type AssignedPermission struct {
	ID           int64          `gorm:"column:id;primarykey" json:"id"`
	PermissionId sql.NullInt64  `gorm:"column:permission_id" json:"permission_id"`
	TeamId       sql.NullString `gorm:"column:team_id" json:"team_id"`
	RoleId       sql.NullInt64  `gorm:"column:role_id" json:"role_id"`
	ResourceId   string         `gorm:"column:resource_id;not null;" json:"resource_id"`
	ResourceName string         `gorm:"column:resource_name;not null;" json:"resource_name"`
	Action       string         `gorm:"column:action;not null;" json:"action"`
	Status       string         `gorm:"column:status" json:"status"`
	Remarks      string         `gorm:"column:remarks" json:"remarks"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    sql.NullTime   `gorm:"column:deleted_at" json:"deleted_at"`
}
