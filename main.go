package main

import (
	"fmt"
	"gitlab.ecoworkinc.com/subspace/softetherlib/softether"
)

func main() {
	s := softether.SoftEther{IP: "52.199.244.30", Password: "ecowork", Hub: "ecowork-aws"}

	// Get Server Status
	s.GetServerStatus()

	// Create User and Get User Info
	s.CreateUser("1", "test@ecoworkinc.com", "Test Account")
	createdUser, _ := s.GetUserInfo("1")
	fmt.Println(createdUser)

	// Update the user alias
	s.UpdateUserAlias("1", "Dan")
	updatedUser, _ := s.GetUserInfo("1")
	fmt.Println(updatedUser)

	// Delete the user
	s.DeleteUser("1")

}
