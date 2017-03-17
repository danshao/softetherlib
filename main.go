package main

import (
	"fmt"
	"gitlab.ecoworkinc.com/subspace/softetherlib/softether"
)

func main() {
	s := softether.SoftEther{IP: "52.199.244.30", Password: "ecowork"}
	serverStatus, _ := s.GetServerStatus()

	fmt.Print(serverStatus)
}
