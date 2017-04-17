package main

import (
	"fmt"

	"gitlab.ecoworkinc.com/subspace/softetherlib/softether"
)

var printMap = func(m map[string]string) {
	for a, b := range m {
		fmt.Printf("%s: %s\n", a, b)
	}
	fmt.Println("")
}

func main() {
	s := softether.SoftEther{IP: "34.201.9.166", Password: "subspace", Hub: "subspace"}

	// Get Server Status
	serverStatus, _ := s.GetServerStatus()
	fmt.Println("Server Status")
	fmt.Println("-------------")
	printMap(serverStatus)

	// Get Session List
	sessionList, _ := s.GetSessionList()
	fmt.Println("Session List")
	fmt.Println("-------------")
	for a, session := range sessionList {
		fmt.Printf("[Session %d]\n", a)
		for c, d := range session {
			fmt.Printf("%s: %s\n", c, d)
		}
		fmt.Println("")
	}
	fmt.Println("")

	// Get Session Info
	fmt.Println("Session Info")
	fmt.Println("-------------")
	sessionInfo, _ := s.GetSessionInfo("SID-SECURENAT-1")
	printMap(sessionInfo)

	// Create User, Set Password and Get User Info
	s.CreateUser("1", "test@ecoworkinc.com", "New Account")
	s.SetUserPassword("1", "abcde")
	createdUser, _ := s.GetUserInfo("1")
	fmt.Println("Created User")
	fmt.Println("------------")
	printMap(createdUser)

	// Update User information
	s.SetUserInfo("1", "modifiedtest@ecoworkinc.com", "Modified Account Name")
	updatedUser, _ := s.GetUserInfo("1")
	fmt.Println("Updated User Alias")
	fmt.Println("------------------")
	printMap(updatedUser)

	// Revoke user
	s.ExpireUser("1")
	revokedUser, _ := s.GetUserInfo("1")
	fmt.Println("Revoked User")
	fmt.Println("------------------")
	printMap(revokedUser)

	// Delete the user
	s.DeleteUser("1")
	_, err := s.GetUserInfo("1")
	fmt.Println("User Deleted")
	fmt.Println("------------")
	fmt.Println("Return Code: ", softether.Strerror(err))

}
