package main

import (
	"gitlab.ecoworkinc.com/subspace/softetherlib/softether"
)

func main() {
	s := softether.SoftEther{IP: "52.199.244.30", Password: "ecowork", Hub: "ecowork-aws"}
	// serverStatus, err := s.GetServerStatus()

	s.GetUserInfo("dan")
}
