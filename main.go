package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	grpcAcl "grpc-acl/pkg"
)

func main() {
	dsn := "miqdad:anone@tcp(127.0.0.1:3306)/analytics?charset=utf8mb4&parseTime=True&loc=Local"
	dbClient, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	acl := grpcAcl.New(dbClient)
	_ = acl.NewPermission("Blog", "hello", "")
	//x := acl.NewRole("Admin", "adm", "")

}