package softether

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// SoftEther is a struct which holds the IP and Password of the SoftEther server
type SoftEther struct {
	IP       string
	Password string
	Hub      string
}

var reFindIntegers *regexp.Regexp = regexp.MustCompile("[0-9]+")
var cleanBytesOutput = func(aString string) (int, error) {
	return strconv.Atoi(strings.Join(reFindIntegers.FindAllString(aString, -1), ""))
}

// GetServerStatus executes vpncmd and gets the server status info from the SoftEther server.
// It returns a map of relevant information for the Subspace project and any error encountered.
// @returns status map[string]string
// @returns returnCode int
func (s *SoftEther) GetServerStatus() (status map[string]string, returnCode int) {

	// Command to execute
	// vpncmd /server [IP] /password:[PASSWORD] /cmd ServerStatusGet
	cmd := exec.Command("vpncmd", "/server", s.IP, "/password:"+s.Password, "/cmd", "ServerStatusGet")

	// Local variables
	statusMap := make(map[string]string)
	cmdOutput := &bytes.Buffer{} // Stdout buffer

	// Attach buffer to command output and execute
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return
	}

	// Prepare iostream and extract data
	outputScanner := bufio.NewScanner(bytes.NewReader(cmdOutput.Bytes()))
	for outputScanner.Scan() {
		if strings.Contains(outputScanner.Text(), "|") {
			s := strings.Split(outputScanner.Text(), "|")
			s[0] = strings.Trim(s[0], " ")
			statusMap[s[0]] = s[1]
		}
	}

	// Perform calculations for outgoing and incoming traffic
	outgoingUnicastBytes, _ := cleanBytesOutput(statusMap["Outgoing Unicast Total Size"])
	outgoingBroadcastBytes, _ := cleanBytesOutput(statusMap["Outgoing Broadcast Total Size"])
	incomingUnicastBytes, _ := cleanBytesOutput(statusMap["Incoming Unicast Total Size"])
	incomingBroadcastBytes, _ := cleanBytesOutput(statusMap["Incoming Broadcast Total Size"])
	outgoingTotalBytes := strconv.Itoa(outgoingBroadcastBytes + outgoingUnicastBytes)
	incomingTotalBytes := strconv.Itoa(incomingBroadcastBytes + incomingUnicastBytes)

	// Put it all together
	status = map[string]string{
		"numberOfSessions":  statusMap["Number of Sessions"],
		"numberOfUsers":     statusMap["Number of Users"],
		"currentServerTime": statusMap["Current Time"][0:19],
		"serverStartTime":   statusMap["Server Started at"][0:11] + statusMap["Server Started at"][17:],
		"incomingBytes":     incomingTotalBytes,
		"outgoingBytes":     outgoingTotalBytes,
	}

	return
}

// GetUserInfo executes vpncmd and gets the details of a specific user
// It returns a map of relevant information for the Subspace project and any error encountered.
// @param id string
// @returns userInfo map[string]string
// @returns returnCode int
func (s *SoftEther) GetUserInfo(id string) (userInfo map[string]string, returnCode int) {

	// Command to execute
	// vpncmd /server [IP] /password:[PASSWORD] /hub:[HUB] /cmd UserGet [NAME]
	cmd := exec.Command("vpncmd", "/server", s.IP, "/password:"+s.Password, "/hub:"+s.Hub, "/cmd", "UserGet", id)

	// Local variables
	userInfoMap := make(map[string]string)
	cmdOutput := &bytes.Buffer{} // Stdout buffer

	// Attach buffer to command output and execute
	cmd.Stdout = cmdOutput
	err := cmd.Run()
	if err != nil {
		returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return
	}

	// Prepare iostream and extract data
	outputScanner := bufio.NewScanner(bytes.NewReader(cmdOutput.Bytes()))
	for outputScanner.Scan() {
		if strings.Contains(outputScanner.Text(), "|") {
			s := strings.Split(outputScanner.Text(), "|")
			s[0] = strings.Trim(s[0], " ")
			userInfoMap[s[0]] = s[1]
		}
	}

	// Perform calculations for outgoing and incoming traffic
	outgoingUnicastBytes, _ := cleanBytesOutput(userInfoMap["Outgoing Unicast Total Size"])
	outgoingBroadcastBytes, _ := cleanBytesOutput(userInfoMap["Outgoing Broadcast Total Size"])
	incomingUnicastBytes, _ := cleanBytesOutput(userInfoMap["Incoming Unicast Total Size"])
	incomingBroadcastBytes, _ := cleanBytesOutput(userInfoMap["Incoming Broadcast Total Size"])
	outgoingTotalBytes := strconv.Itoa(outgoingBroadcastBytes + outgoingUnicastBytes)
	incomingTotalBytes := strconv.Itoa(incomingBroadcastBytes + incomingUnicastBytes)

	// Put it all together
	userInfo = map[string]string{
		"email":          userInfoMap["Description"],
		"alias":          userInfoMap["Full Name"],
		"numberOfLogins": userInfoMap["Number of Logins"],
		"expirationDate": userInfoMap["Expiration Date"],
		"creationDate":   userInfoMap["Created on"][0:11] + userInfoMap["Created on"][17:],
		"updatedDate":    userInfoMap["Updated on"][0:11] + userInfoMap["Updated on"][17:],
		"incomingBytes":  incomingTotalBytes,
		"outgoingBytes":  outgoingTotalBytes,
	}

	return
}

// CreateUser executes vpncmd and creates a user
// @param id string
// @param email string
// @param alias string ""
func (s *SoftEther) CreateUser(args ...interface{}) {

	// // Mandatory parameters
	// var id string
	// var email string

	// // Optional parameters
	// var alias string = ""

	// // Ensure that we have at least 2 parameters
	// if 2 > len(args) {
	// 	panic("Not enough parameters.")
	// }

	// // Attach buffer to command output and execute
	// cmd.Stdout = cmdOutput
	// printCommand(cmd)
	// err := cmd.Run() // will wait for command to return
	// printError(err)
	// if err != nil {
	// 	returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
	// 	return
	// }
	// printOutput(cmdOutput.Bytes())

}

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, " "))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}
