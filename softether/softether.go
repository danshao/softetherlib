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
	"time"
)

// SoftEther is a struct which holds the IP, Password, and Hub of the SoftEther server.
type SoftEther struct {
	IP       string
	Password string
	Hub      string
}

const SOFT_ETHER_TABLE_HEADER_KEY = "Item"

var reFindIntegers = regexp.MustCompile("[0-9]+")
var cleanBytesOutput = func(aString string) (int, error) {
	return strconv.Atoi(strings.Join(reFindIntegers.FindAllString(aString, -1), ""))
}

// GetServerStatus executes vpncmd and gets the server status info from the SoftEther server.
func (s SoftEther) GetServerStatus() (status map[string]string, returnCode int) {

	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /cmd ServerStatusGet
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/cmd",
		"ServerStatusGet",
	)

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

// GetSessionList executes vpncmd and gets the session list from the SoftEther server for a specific Hub.
func (s SoftEther) GetSessionList() (sessionListMap map[int]map[string]string, returnCode int) {

	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd SessionList
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"SessionList",
	)

	// Local variables
	sessionListMap = make(map[int]map[string]string)
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
	pos := 0
	for outputScanner.Scan() {
		if strings.Contains(outputScanner.Text(), "|") {
			s := strings.Split(outputScanner.Text(), "|")
			s[0] = strings.Trim(s[0], " ")

			if SOFT_ETHER_TABLE_HEADER_KEY == s[0] {
				continue // Skip table header
			}

			if _, ok := sessionListMap[pos]; !ok {
				sessionListMap[pos] = make(map[string]string)
			}

			if s[0] == "Transfer Bytes" || s[0] == "Transfer Packets" {
				aBytes, _ := cleanBytesOutput(s[1])
				s[1] = strconv.Itoa(aBytes)
			}

			sessionListMap[pos][s[0]] = s[1]

			if s[0] == "Transfer Packets" {
				pos++
			}
		}
	}

	return
}

// GetSessionInfo executes vpncmd and gets the session information for a specific Session Name
func (s SoftEther) GetSessionInfo(sessionName string) (sessionInfo map[string]string, returnCode int) {
	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd SessionGet [SESSION_NAME]
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"SessionGet", sessionName,
	)

	// Local variables
	sessionInfo = make(map[string]string)
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

			if SOFT_ETHER_TABLE_HEADER_KEY == s[0] {
				continue // Skip table header
			}

			sessionInfo[s[0]] = s[1]
		}
	}

	for key, value := range sessionInfo {
		switch key {
		case
			"Outgoing Data Size",
			"Incoming Data Size",
			"Outgoing Unicast Packets",
			"Outgoing Unicast Total Size",
			"Outgoing Broadcast Packets",
			"Outgoing Broadcast Total Size",
			"Incoming Unicast Packets",
			"Incoming Unicast Total Size",
			"Incoming Broadcast Packets",
			"Incoming Broadcast Total Size":
			// Convert "4,734,874 bytes" to "4734874"
			formattedValue, _ := cleanBytesOutput(value)
			sessionInfo[key] = strconv.Itoa(formattedValue)
		case
			"First Session has been Established since",
			"Current Session has been Established since":
			// Convert "2017-04-19 (Wed) 02:05:16" to "2017-04-19 02:05:16"
			sessionInfo[key] = value[0:11] + value[17:]
		}
	}

	return
}

// GetUserList executes vpncmd and gets the user list from the SoftEther server for a specific Hub.
func (s SoftEther) GetUserList() (userListMap map[int]map[string]string, returnCode int) {

	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd SessionList
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"UserList",
	)

	// Local variables
	userListMap = make(map[int]map[string]string)
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
	pos := 0
	for outputScanner.Scan() {
		if strings.Contains(outputScanner.Text(), "|") {
			s := strings.Split(outputScanner.Text(), "|")
			s[0] = strings.Trim(s[0], " ")

			if SOFT_ETHER_TABLE_HEADER_KEY == s[0] {
				continue // Skip table header
			}

			if _, ok := userListMap[pos]; !ok {
				userListMap[pos] = make(map[string]string)
			}

			userListMap[pos][s[0]] = s[1]

			if s[0] == "Transfer Packets" {
				pos++
			}
		}
	}
	for _, user := range userListMap {
		fmt.Println(user)
		for key, value := range user {
			switch key {
			case
				"Transfer Bytes",
				"Transfer Packets":
				// Convert "4,734,874 bytes" to "4734874"
				formattedValue, _ := cleanBytesOutput(value)
				user[key] = strconv.Itoa(formattedValue)
			case
				"Last Login":
				// Convert "2017-04-19 (Wed) 02:05:16" to "2017-04-19 02:05:16"
				if "(None)" != value {
					user[key] = value[0:11] + value[17:]
				}
			}
		}
	}
	return
}

// GetUserInfo executes vpncmd and gets the details of a specific User for a specific Hub.
func (s SoftEther) GetUserInfo(id string) (userInfo map[string]string, returnCode int) {

	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd UserGet [NAME]
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"UserGet", id,
	)

	// Local variables
	userInfo = make(map[string]string)
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

			if SOFT_ETHER_TABLE_HEADER_KEY == s[0] {
				continue // Skip table header
			}

			userInfo[s[0]] = s[1]
		}
	}

	for key, value := range userInfo {
		switch key {
		case
			"Outgoing Unicast Packets",
			"Outgoing Unicast Total Size",
			"Outgoing Broadcast Packets",
			"Outgoing Broadcast Total Size",
			"Incoming Unicast Packets",
			"Incoming Unicast Total Size",
			"Incoming Broadcast Packets",
			"Incoming Broadcast Total Size":
			// Convert "4,734,874 bytes" to "4734874"
			formattedValue, _ := cleanBytesOutput(value)
			userInfo[key] = strconv.Itoa(formattedValue)
		case
			"Created on",
			"Updated on":
			// Convert "2017-04-19 (Wed) 02:05:16" to "2017-04-19 02:05:16"
			userInfo[key] = value[0:11] + value[17:]
		}
	}

	return
}

// CreateUser executes vpncmd and creates a User for a specific Hub.
func (s SoftEther) CreateUser(args ...interface{}) (returnCode int) {

	// Mandatory parameters
	var id string
	var email string

	// Optional parameters
	var description string

	// Ensure that we have at least 2 parameters
	if 2 > len(args) {
		panic("Not enough parameters.")
	}

	// Get any parameters passed out of args
	for i, p := range args {
		switch i {
		case 0: // id
			param, ok := p.(string)
			if !ok {
				panic("First paramter (id) not type string.")
			}
			id = param

		case 1: // email
			param, ok := p.(string)
			if !ok {
				panic("Second parameter (email) not type string.")
			}
			email = param

		case 2: // description
			param, ok := p.(string)
			if !ok {
				panic("Third parameter (description) not type string.")
			}
			description = param

		default:
			panic("Wrong parameter count.")
		}
	}

	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd UserCreate [NAME] /GROUP:[GROUP] /REALNAME:[EMAIL] /NOTE:[DESCRIPTION]
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"UserCreate", id,
		"/REALNAME:"+email,
		"/NOTE:"+description,
		"/GROUP:",
	)
	cmdOutput := &bytes.Buffer{} // Stdout buffer

	// Attach buffer to command output and execute
	cmd.Stdout = cmdOutput
	err := cmd.Run() // will wait for command to return
	if err != nil {
		returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return
	}

	return
}

// SetUserPassword executes vpncmd and updates a specific User's password in a specific Hub.
func (s SoftEther) SetUserPassword(id string, password string) (returnCode int) {
	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd UserPasswordSet [NAME] /GROUP:[GROUP] /REALNAME:[ALIAS] /NOTE:[EMAIL]
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"UserPasswordSet", id,
		"/PASSWORD:"+password,
	)
	cmdOutput := &bytes.Buffer{} // Stdout buffer

	// Attach buffer to command output and execute
	cmd.Stdout = cmdOutput
	err := cmd.Run() // will wait for command to return
	if err != nil {
		returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return
	}

	return
}

// SetUserAlias executes vpncmd and updates a specific User's information in a specific Hub.
func (s SoftEther) SetUserInfo(id, email, description string) (returnCode int) {

	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd UserSet [NAME] /GROUP:[GROUP] /REALNAME:[EMAIL] /NOTE:[DESCRIPTION]
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"UserSet", id,
		"/REALNAME:"+email,
		"/NOTE:"+description,
		"/GROUP:",
	)
	cmdOutput := &bytes.Buffer{} // Stdout buffer

	// Attach buffer to command output and execute
	cmd.Stdout = cmdOutput
	err := cmd.Run() // will wait for command to return
	if err != nil {
		returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return
	}

	return
}

// DeleteUser executes vpncmd and deletes a specific User in a specific Hub.
func (s SoftEther) DeleteUser(id string) (returnCode int) {
	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd UserDelete [NAME]
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"UserDelete", id,
	)
	cmdOutput := &bytes.Buffer{} // Stdout buffer

	// Attach buffer to command output and execute
	cmd.Stdout = cmdOutput
	err := cmd.Run() // will wait for command to return
	if err != nil {
		returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return
	}

	return
}

// DisconnectSession executes vpncmd and disconnects a specific session
func (s SoftEther) DisconnectSession(sessionName string) (returnCode int) {
	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd SessionDisconnect [SESSION_NAME]
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"SessionDisconnect", sessionName,
	)
	cmdOutput := &bytes.Buffer{} // Stdout buffer

	// Attach buffer to command output and execute
	cmd.Stdout = cmdOutput
	err := cmd.Run() // will wait for command to return
	if err != nil {
		returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return
	}

	return
}

// SetUserEnabled executes vpncmd to enable/disable a specified Username
func (s SoftEther) SetUserEnabled(username string, enabled bool) (returnCode int) {
	var (
		expirationDate string
	)

	if enabled == false {
		// Get current time, set to one day before
		t := time.Now()
		expirationDate = t.AddDate(0, 0, -1).Format("2006/01/02 15:04:05")
	} else {
		expirationDate = "none"
	}

	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /hub:[HUB] /cmd UserExpiresSet [SESSION_NAME] /EXPIRES:[EXPIRATION_DATE}]
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/hub:"+s.Hub,
		"/cmd",
		"UserExpiresSet", username,
		"/expires:"+expirationDate,
	)
	cmdOutput := &bytes.Buffer{} // Stdout buffer

	// Attach buffer to command output and execute
	cmd.Stdout = cmdOutput
	err := cmd.Run() // will wait for command to return
	if err != nil {
		returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return
	}

	return
}

// SetPreSharedKey executes vpncmd to modify the preshared key
func (s SoftEther) SetPreSharedKey(preSharedKey string) (returnCode int) {
	// Command to execute
	// vpncmd /server [IP]:992 /password:[PASSWORD] /cmd IPsecEnable [/L2TP:yes|no] [/L2TPRAW:yes|no] [/ETHERIP:yes|no] [/PSK:pre-shared-key] [/DEFAULTHUB:default_hub]
	cmd := exec.Command(
		"vpncmd",
		"/server",
		s.IP+":992",
		"/password:"+s.Password,
		"/cmd",
		"IPsecEnable",
		"/L2TP:yes",
		"/L2TPRAW:no",
		"/ETHERIP:no",
		"/DEFAULTHUB:"+s.Hub,
		"/PSK:"+preSharedKey,
	)
	cmdOutput := &bytes.Buffer{} // Stdout buffer
	printCommand(cmd)

	// Attach buffer to command output and execute
	cmd.Stdout = cmdOutput
	err := cmd.Run() // will wait for command to return
	if err != nil {
		returnCode, _ = strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return
	}

	return
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
