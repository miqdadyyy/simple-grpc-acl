package pkg

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"time"
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
	if err := acl.DB.Where("name", name).
		Find(&permission).Error; err != nil {
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

func (acl *GrpcAcl) GetRoleByName(name string) *Role {
	var role Role
	if err := acl.DB.
		Where("name", name).
		Find(&role).Error; err != nil {
		return nil
	}

	return &role
}

func (acl *GrpcAcl) AssignPermissionByName(resource RoleAndPermission, permissionName string, actions ...string) error {
	permission := acl.GetPermissionByName(permissionName)
	if permission == nil {
		return errors.New("Permission not found")
	}

	return acl.AssignPermission(resource, permission, actions...)
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
			TeamId:       resource.GetTeamId(),
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

func (acl *GrpcAcl) AssignRoleByName(resource RoleAndPermission, roleName string, teamId string) error {
	role := acl.GetRoleByName(roleName)
	if role == nil {
		return errors.New("Role not found")
	}

	return acl.AssignRole(resource, role, teamId, "")
}

func (acl *GrpcAcl) AssignRoleByNameWithStatus(resource RoleAndPermission, roleName string, teamId string, status string) error {
	role := acl.GetRoleByName(roleName)
	if role == nil {
		return errors.New("Role not found")
	}

	return acl.AssignRole(resource, role, teamId, status)
}

func (acl *GrpcAcl) AssignRole(resource RoleAndPermission, role *Role, teamId string, status string) error {
	return acl.AssignRoleById(resource, role.ID, teamId, status)
}

func (acl *GrpcAcl) AssignRoleWithRemarks(resource RoleAndPermission, role *Role, teamId, status, remarks string) error {
	var team sql.NullString
	if teamId != "" {
		team = sql.NullString{
			Valid:  true,
			String: teamId,
		}
	} else {
		team = sql.NullString{
			Valid: false,
		}
	}

	assignedPermission := AssignedPermission{
		RoleId: sql.NullInt64{
			Valid: true,
			Int64: role.ID,
		},
		TeamId:       team,
		Status:       status,
		ResourceName: resource.GetResourceName(),
		ResourceId:   resource.GetResourceId(),
		Remarks:      remarks,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := acl.DB.Create(&assignedPermission).Error; err != nil {
		return err
	}

	return nil
}

func (acl *GrpcAcl) AssignRoleById(resource RoleAndPermission, roleId int64, teamId string, status string) error {
	var team sql.NullString
	if teamId != "" {
		team = sql.NullString{
			Valid:  true,
			String: teamId,
		}
	} else {
		team = sql.NullString{
			Valid: false,
		}
	}

	assignedPermission := AssignedPermission{
		RoleId: sql.NullInt64{
			Valid: true,
			Int64: roleId,
		},
		TeamId:       team,
		Status:       status,
		ResourceName: resource.GetResourceName(),
		ResourceId:   resource.GetResourceId(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := acl.DB.Create(&assignedPermission).Error; err != nil {
		return err
	}

	return nil
}

func (acl *GrpcAcl) GetModelPermissions(model RoleAndPermission) ([]Permission, error) {
	var assignedPermissions []AssignedPermission
	var permissionIds []int64
	var permissions []Permission

	if err := acl.DB.Find(&assignedPermissions).
		Joins("LEFT JOIN assigned_permissions as AP on assigned_permissions.resource_id = AP.role_id AND assigned_permissions.resource_name = 'role'").
		Where("AP.resource_name = ?", model.GetResourceName()).
		Where("AP.resource_id = ?", model.GetResourceId()).
		Or("assigned_permissions.resource_name = ?", model.GetResourceName()).
		Where("assigned_permissions.resource_id = ?", model.GetResourceId()).
		Where("assigned_permissions.role_id IS NULL").Error; err != nil {
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

func (acl *GrpcAcl) CheckPermissionInRole(permission *Permission, role *Role, actions ...string) bool {
	return acl.CheckPermissionInModel(permission, role, actions...)
}

func (acl *GrpcAcl) CheckPermissionInModel(permission *Permission, model RoleAndPermission, actions ...string) bool {
	var assignedPermission AssignedPermission
	if err := acl.DB.Model(&assignedPermission).
		Joins("LEFT JOIN assigned_permissions as AP on assigned_permissions.resource_id = AP.role_id AND assigned_permissions.resource_name = 'role'").
		Where("AP.resource_name = ?", model.GetResourceName()).
		Where("AP.resource_id = ?", model.GetResourceId()).
		Where("assigned_permissions.permission_id = ?", permission.ID).
		Where("assigned_permissions.action", actions).
		Or("assigned_permissions.resource_name = ?", model.GetResourceName()).
		Where("assigned_permissions.resource_id = ?", model.GetResourceId()).
		Where("assigned_permissions.role_id IS NULL").
		Where("assigned_permissions.permission_id = ?", permission.ID).
		Where("assigned_permissions.action", actions).
		First(&assignedPermission).
		Error; err != nil {
		return false
	}

	return true
}

func (acl *GrpcAcl) CheckPermissionInModelWithTeam(permission *Permission, model RoleAndPermission, teamId string, actions ...string) bool {
	var assignedPermission AssignedPermission
	if err := acl.DB.Model(&assignedPermission).
		Joins("LEFT JOIN assigned_permissions AS AP "+
			"ON assigned_permissions.resource_id = AP.role_id "+
			"AND assigned_permissions.resource_name = 'role' "+
			"AND AP.team_id = ?", teamId).
		Where("AP.resource_name = ?", model.GetResourceName()).
		Where("AP.resource_id = ?", model.GetResourceId()).
		Where("AP.team_id = ?", teamId).
		Where("assigned_permissions.permission_id = ?", permission.ID).
		Where("assigned_permissions.action", actions).
		Or("assigned_permissions.resource_name = ?", model.GetResourceName()).
		Where("assigned_permissions.resource_id = ?", model.GetResourceId()).
		Where("assigned_permissions.role_id IS NULL").
		Where("assigned_permissions.permission_id = ?", permission.ID).
		Where("assigned_permissions.action", actions).
		Where("assigned_permissions.team_id = ?", teamId).
		First(&assignedPermission).
		Error; err != nil {
		return false
	}

	return true
}

func (acl *GrpcAcl) CheckModelRoleWithTeam(model RoleAndPermission, role *Role, teamId string) bool {
	var permission AssignedPermission
	if err := acl.DB.Where("role_id = ?", role.ID).
		Where("resource_name = ?", model.GetResourceName()).
		Where("resource_id = ?", model.GetResourceId()).
		Where("team_id = ?", teamId).First(&permission).Error; err != nil {
		return false
	}

	return true
}

func (acl *GrpcAcl) UpdatePermissionStatus(model RoleAndPermission, teamId string, status string) error {
	if err := acl.DB.Model(&AssignedPermission{}).
		Where("resource_name = ?", model.GetResourceName()).
		Where("resource_id = ?", model.GetResourceId()).
		Where("team_id = ?", teamId).
		Update("status", status).Error; err != nil {
		return err
	}

	return nil
}

func (acl *GrpcAcl) RetractModelFromTeam(model RoleAndPermission, teamId string) error {
	if err := acl.DB.
		Where("resource_name = ?", model.GetResourceName()).
		Where("resource_id = ?", model.GetResourceId()).
		Where("team_id = ?", teamId).
		Delete(&AssignedPermission{}).Error; err != nil {
		return err
	}

	return nil
}

func (acl *GrpcAcl) RetractAllTeamMember(teamId string) error {
	if err := acl.DB.Model(&AssignedPermission{}).
		Where("team_id = ?", teamId).
		Update("status", "inactive").Error; err != nil {
		return err
	}

	return nil
}

func New(db *gorm.DB) *GrpcAcl {
	return &GrpcAcl{
		DB: db,
	}
}
