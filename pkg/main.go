package pkg

import (
	"database/sql"
	"gorm.io/gorm"
)

type RoleAndPermission interface {
	GetResourceName() string
	GetResourceId() string
	GetTeamId() sql.NullString
	GetTeamIdValue() string
}

type GrpcAcl struct {
	DB *gorm.DB
}

func (acl *GrpcAcl) NewPermission(title, name, description string) *Permission {
	permission := &Permission{
		Title: title,
		Name:  name,
		Description: sql.NullString{
			String: description,
		},
	}

	acl.DB.Create(&permission)
	return permission
}

func (acl *GrpcAcl) GetAllPermissions() ([]Permission, error) {
	var permissions []Permission
	if err := acl.DB.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

func (acl *GrpcAcl) GetPermissionByName(name string) *Permission {
	var permission Permission
	if err := acl.DB.Find(&permission).
		Where("name", name).Error; err != nil {
			return nil
	}

	return &permission
}

func (acl *GrpcAcl) NewRole(title, name, description string) *Role {
	role := &Role{
		Title: title,
		Name:  name,
		Description: sql.NullString{
			String: description,
		},
	}

	acl.DB.Create(&role)
	return role
}

func (acl *GrpcAcl) GetAllRoles() ([]Role, error) {
	var roles []Role
	if err := acl.DB.Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (acl *GrpcAcl) GetRoleByName(name string) *Role{
	var role Role
	if err := acl.DB.Find(&role).
		Where("name", name).Error; err != nil {
		return nil
	}

	return &role
}

func (acl *GrpcAcl) AssignPermission(resource RoleAndPermission, permission *Permission, actions ...string) error {
	return acl.AssignPermissionById(resource, permission.ID, actions...)
}

func (acl *GrpcAcl) AssignPermissionById(resource RoleAndPermission, permissionId int64, actions ...string) error {
	var assignedPermissions []AssignedPermission

	for _, action := range actions {
		assignedPermissions = append(assignedPermissions, AssignedPermission{
			PermissionId: sql.NullInt64{
				Valid: true,
				Int64: permissionId,
			},
			TeamId: resource.GetTeamId(),
			ResourceId:   resource.GetResourceId(),
			ResourceName: resource.GetResourceName(),
			Action:       action,
		})
	}

	if err := acl.DB.Create(&assignedPermissions).Error; err != nil {
		return err
	}

	return nil
}

func (acl *GrpcAcl) RetractPermission(resource RoleAndPermission, permissions ...Permission) error {
	var permissionIds []int64
	for _, permission := range permissions {
		permissionIds = append(permissionIds, permission.ID)
	}

	return acl.RetractPermissionById(resource, permissionIds...)
}

func (acl *GrpcAcl) RetractPermissionById(resource RoleAndPermission, permissionIds ...int64) error {
	if err := acl.DB.Delete(&AssignedPermission{}).
		Where("resource_name", resource.GetResourceName()).
		Where("resource_id", resource.GetResourceId()).
		Where("permission_id", &permissionIds).Error; err != nil {
		return err
	}

	return nil
}

func (acl *GrpcAcl) GetModelRoles(model RoleAndPermission) []*Role {
	var assignedPermissions []AssignedPermission
	var roles []*Role
	if err := acl.DB.Find(&assignedPermissions).
		Where("resource_name", model.GetResourceName()).
		Where("resource_id", model.GetResourceId()).
		Where("role_id", "<> NULL").
		Select("role_id").Error; err != nil {
		return roles
	}

	var roleIds []int64
	for _, assignedPermission := range assignedPermissions {
		roleIds = append(roleIds, assignedPermission.ID)
	}

	if err := acl.DB.Find(&roles).
		Where("id", roleIds).Error; err != nil {
		return roles
	}

	return roles
}

func (acl *GrpcAcl) GetModelPermissions(model RoleAndPermission) ([]Permission, error) {
	var assignedPermissions []AssignedPermission
	var permissionIds []int64
	var permissions []Permission

	if err := acl.DB.Find(&assignedPermissions).
		Where("resource_name", model.GetResourceName()).
		Where("resource_id", model.GetResourceId()).
		Where("team_id", model.GetTeamId()).Error; err != nil {
		return nil, err
	}

	for _, assignedPermission := range assignedPermissions {
		permissionIds = append(permissionIds, assignedPermission.PermissionId.Int64)
	}

	if err := acl.DB.Find(&permissions).
		Where("id", permissionIds).Error; err != nil {
			return nil, err
	}

	return permissions, nil
}

func (acl *GrpcAcl) CheckPermissionInRole(permission *Permission, role *Role, actions ...string) (bool) {
	return acl.CheckPermissionInModel(permission, role, actions...)
}

func (acl *GrpcAcl) CheckPermissionInModel(permission *Permission, model RoleAndPermission, actions ...string) bool {
	var assignedPermissions []AssignedPermission
	acl.DB.Where("")
	if err := acl.DB.Find(&assignedPermissions).
		Where("permission_id", permission.ID).
		Where("resource_name", model.GetResourceName()).
		Where("resource_id", model.GetResourceId()).
		Where("action", actions).Error; err != nil {
		return false
	}

	return true
}

func New(db *gorm.DB) *GrpcAcl {
	return &GrpcAcl{
		DB: db,
	}
}
