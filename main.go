package main

import (
	"gitlab.ecoworkinc.com/subspace/softetherlib/softether"
)

func main() {
	s := softether.SoftEther{IP: "52.199.244.30", Password: "ecowork", Hub: "ecowork-aws"}

	s.GetServerStatus()
	s.CreateUser("1", "test@ecoworkinc.com", "Test Account")
	s.GetUserInfo("1")
	s.UpdateUserAlias("1", "Dan")
}
