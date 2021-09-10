package pkg

type ServerAuthorization struct {
	Acl *GrpcAcl
	PermissionName string
}

func NewServerAuthorization(acl *GrpcAcl, permissionName string) *ServerAuthorization{
	return &ServerAuthorization{
		Acl: acl,
		PermissionName: permissionName,
	}
}

func (s *ServerAuthorization) CanRead(resource RoleAndPermission) bool {
	// Check Permission on Role
	return false
}
