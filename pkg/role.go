package pkg

import (
	"database/sql"
	"strconv"
)

type Role struct {
	ID                int64          `gorm:"column:id;primarykey" json:"id"`
	Title             string         `gorm:"column:title;not null;unique;type:varchar;size:255" json:"title"`
	Name              string         `gorm:"column:name;not null;unique;type:varchar;size:255" json:"name"`
	Description       sql.NullString `gorm:"column:description;type:varchar;size:255" json:"description"`
	RoleAndPermission `gorm:"-"`
}

func (role *Role) GetResourceName() string {
	return "role"
}

func (role *Role) GetResourceId() string {
	return strconv.FormatInt(role.ID, 10)
}

func (role *Role) GetTeamId() sql.NullString {
	return sql.NullString{Valid: false}
}

func (role *Role) GetTeamIdValue() string {
	return ""
}

func (role *Role) Update(acl *GrpcAcl, title, name, description string) (*Role, error) {
	role.Title = title
	role.Name = name
	role.Description = sql.NullString{
		String: description,
	}

	if err := acl.DB.Save(&role).Error; err != nil {
		return nil, err
	}

	return role, nil
}

func (role *Role) Delete(acl *GrpcAcl) error {
	if err := acl.DB.Delete(&role).Error; err != nil {
		return err
	}
	return nil
}

func (role *Role) GetPermissions(acl *GrpcAcl) ([]Permission, error) {
	return acl.GetModelPermissions(role)
}

func (role *Role) AssignPermission(acl *GrpcAcl, permission *Permission, actions ...string) (*Role, error) {
	if err := acl.AssignPermission(role, permission, actions...); err != nil {
		return nil, err
	}
	return role, nil
}

func (role *Role) RetractPermissions(acl *GrpcAcl, permissions ...Permission) (*Role, error) {
	if err := acl.RetractPermission(role, permissions...); err != nil {
		return nil, err
	}

	return role, nil
}
