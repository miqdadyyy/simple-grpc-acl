package test

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"simple-grpc-acl/pkg"
	"testing"
)

/*
- Create Role
	- Admin
	- Member
- Create Permission
	- CRUD Blog
	- CRUD User
- Assign Permission To Role
	- Admin
		- CRUD Blog
		- CRUD User
	- Member
		- R Blog
		- R User
- Assign Permission To User
	- Member
		- CUD Blog
- Check Permission
	- Admin can create blog
	- Admin can delete user
	- Member can read blog
	- Member can't create blog
*/

func getAclInstance() *pkg.GrpcAcl {
	db, _ := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	return pkg.New(db)
}

func TestClearDatabase(t *testing.T) {
	_ = os.Remove("test.sqlite")
	t.Log("Database Cleared")
}

func TestMigrate(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("test.sqlite"), &gorm.Config{})
	err := db.AutoMigrate(&pkg.Role{}, &pkg.Permission{}, &pkg.AssignedPermission{})
	if err != nil {
		t.Error("Migration Failed")
		t.Failed()
	}
}

func TestCreateAdminRole(t *testing.T) {
	acl := getAclInstance()
	role := acl.NewRole("Admin", "admin", "")
	if role == nil {
		t.Error("Create Admin Failed")
		t.Failed()
	}
}

func TestCreateMemberRole(t *testing.T) {
	acl := getAclInstance()
	role := acl.NewRole("Member", "member", "")
	if role == nil {
		t.Error("Create Member Failed")
		t.Failed()
	}
}

func TestCreateBlogPermission(t *testing.T) {
	acl := getAclInstance()
	permission := acl.NewPermission("Blog", "blog", "")
	if permission == nil {
		t.Error("Failed to create blog permission")
		t.Failed()
	}
}

func TestCreteUserPermission(t *testing.T) {
	acl := getAclInstance()
	permission := acl.NewPermission("User", "user", "")
	if permission == nil {
		t.Error("Failed to create user permission")
		t.Failed()
	}
}

func TestAssignBlogPermissionToAdmin(t *testing.T) {
	acl := getAclInstance()
	// Get admin role
	role := acl.GetRoleByName("admin")
	permission := acl.GetPermissionByName("blog")

	_, err := role.AssignPermission(acl, permission, "create", "read", "update", "delete")
	if err != nil {
		t.Error("Failed to assign blog permission to admin")
		t.Failed()
	}
}

func TestAssignUserPermissionToAdmin(t *testing.T) {
	acl := getAclInstance()
	// Get admin role
	role := acl.GetRoleByName("admin")
	permission := acl.GetPermissionByName("user")

	_, err := role.AssignPermission(acl, permission, "create", "read", "update", "delete")
	if err != nil {
		t.Error("Failed to assign blog permission to admin")
		t.Failed()
	}
}

func TestAssignBlogPermissionToUser(t *testing.T) {
	acl := getAclInstance()
	// Get admin role
	role := acl.GetRoleByName("member")
	permission := acl.GetPermissionByName("blog")

	_, err := role.AssignPermission(acl, permission, "read")
	if err != nil {
		t.Error("Failed to assign blog permission to admin")
		t.Failed()
	}
}

func TestAssignUserPermissionToUser(t *testing.T) {
	acl := getAclInstance()
	// Get admin role
	role := acl.GetRoleByName("member")
	permission := acl.GetPermissionByName("user")

	_, err := role.AssignPermission(acl, permission, "read")
	if err != nil {
		t.Error("Failed to assign blog permission to admin")
		t.Failed()
	}
}

func TestAdminCanCreateBlog(t *testing.T){
	acl := getAclInstance()
	role := acl.GetRoleByName("admin")
	permission := acl.GetPermissionByName("blog")
	res := acl.CheckPermissionInModel(permission, role, "create")
	if res != true {
		t.Failed()
	}
}

func TestAdminCanDeleteUser(t *testing.T) {
	acl := getAclInstance()
	role := acl.GetRoleByName("admin")
	permission := acl.GetPermissionByName("user")
	res := acl.CheckPermissionInModel(permission, role, "delete")
	if res != true {
		t.Failed()
	}
}

func TestMemberCanReadBlog(t *testing.T) {
	acl := getAclInstance()
	role := acl.GetRoleByName("member")
	permission := acl.GetPermissionByName("blog")
	res := acl.CheckPermissionInModel(permission, role, "read")
	if res != true {
		t.Failed()
	}
}

func TestMemberCantCreateBlog(t *testing.T) {
	acl := getAclInstance()
	role := acl.GetRoleByName("member")
	permission := acl.GetPermissionByName("blog")
	res := acl.CheckPermissionInModel(permission, role, "create")
	if res == true {
		t.Failed()
	}
}