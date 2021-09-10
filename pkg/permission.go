package pkg

import (
	"database/sql"
	"fmt"
)

type Permission struct {
	ID          int64          `gorm:"column:id;primarykey" json:"id"`
	Title       string         `gorm:"column:title;not null;unique;type:varchar;size:255" json:"title"`
	Name        string         `gorm:"column:name;not null;unique;type:varchar;size:255" json:"name"`
	Description sql.NullString `gorm:"column:description;type:varchar;size:255" json:"description"`
}

func (permission *Permission) Delete(acl *GrpcAcl) error {
	// Delete Assigned Permissions First
	if err := acl.DB.Delete(&AssignedPermission{}, fmt.Sprintf("permission_id = %v", permission.ID)).Error; err != nil {
		return err
	}

	// Delete permission
	if err := acl.DB.Delete(&permission).Error; err != nil {
		return err
	}

	return nil
}

func (permission *Permission) Update(acl *GrpcAcl, title, name, description string) (*Permission, error) {
	permission.Title = title
	permission.Name = name
	permission.Description = sql.NullString{String: description}

	if err := acl.DB.Save(&permission).Error; err != nil {
		return nil, err
	}

	return permission, nil
}
