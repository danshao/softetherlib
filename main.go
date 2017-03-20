package main

import (
	"fmt"
	"gitlab.ecoworkinc.com/subspace/softetherlib/softether"
)

func main() {
	s := softether.SoftEther{IP: "52.199.244.30", Password: "ecowork", Hub: "ecowork-aws"}

	// Get Server Status
	serverStatus, _ := s.GetServerStatus()
	fmt.Println("Server Status")
	fmt.Println("-------------")
	for a, b := range serverStatus {
		fmt.Printf("%s: %s\n", a, b)
	}
	fmt.Println("")

	// Create User and Get User Info
	s.CreateUser("1", "test@ecoworkinc.com", "New Account")
	createdUser, _ := s.GetUserInfo("1")
	fmt.Println("Created User")
	fmt.Println("------------")
	for a, b := range createdUser {
		fmt.Printf("%s: %s\n", a, b)
	}
	fmt.Println("")

	// Update the user alias
	s.UpdateUserAlias("1", "Modified Account Name")
	updatedUser, _ := s.GetUserInfo("1")
	fmt.Println("Updated User Alias")
	fmt.Println("------------------")
	for a, b := range updatedUser {
		fmt.Printf("%s: %s\n", a, b)
	}
	fmt.Println("")

	// Delete the user
	s.DeleteUser("1")
	_, err := s.GetUserInfo("1")
	fmt.Println("User Deleted")
	fmt.Println("------------")
	fmt.Println("Return Code: ", softether.Strerror(err))

}
