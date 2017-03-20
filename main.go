package main

import (
	"gitlab.ecoworkinc.com/subspace/softetherlib/softether"
)

func main() {
	s := softether.SoftEther{IP: "52.199.244.30", Password: "ecowork", Hub: "ecowork-aws"}

	s.GetServerStatus()
	s.GetUserInfo("dan")
	s.CreateUser("1", "test@ecoworkinc.com", "Test Account")
}
