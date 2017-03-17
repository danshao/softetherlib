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
}

// GetServerStatus executes vpncmd and gets the server status info from the SoftEther server.
// It returns a map of relevant information for the Subspace project and any error encountered.
func (s *SoftEther) GetServerStatus() (map[string]string, int) {
	cmd := exec.Command("vpncmd", "/server", s.IP, "/password:"+s.Password, "/cmd", "serverstatusget")
	statusMap := make(map[string]string)
	reFindIntegers := regexp.MustCompile("[0-9]+")
	cmdOutput := &bytes.Buffer{} // Stdout buffer

	cmd.Stdout = cmdOutput // Attach buffer to command
	// Execute command
	// printCommand(cmd)
	err := cmd.Run() // will wait for command to return
	// printError(err)
	if err != nil {
		errorCode, _ := strconv.Atoi(reFindIntegers.FindAllString(err.Error(), -1)[0])
		return make(map[string]string), errorCode
	}
	// printOutput(cmdOutput.Bytes())

	// Open iostream for parsing
	outputScanner := bufio.NewScanner(bytes.NewReader(cmdOutput.Bytes()))

	for outputScanner.Scan() {
		if strings.Contains(outputScanner.Text(), "|") {
			s := strings.Split(outputScanner.Text(), "|")
			s[0] = strings.Trim(s[0], " ")
			statusMap[s[0]] = s[1]
		}
	}

	// Calculate Outgoing and Incoming
	outgoingUnicastBytes, _ := strconv.Atoi(strings.Join(reFindIntegers.FindAllString(statusMap["Outgoing Unicast Total Size"], -1), ""))
	outgoingBroadcastBytes, _ := strconv.Atoi(strings.Join(reFindIntegers.FindAllString(statusMap["Outgoing Broadcast Total Size"], -1), ""))
	outgoingTotalBytes := strconv.Itoa(outgoingBroadcastBytes + outgoingUnicastBytes)
	incomingUnicastBytes, _ := strconv.Atoi(strings.Join(reFindIntegers.FindAllString(statusMap["Incoming Unicast Total Size"], -1), ""))
	incomingBroadcastBytes, _ := strconv.Atoi(strings.Join(reFindIntegers.FindAllString(statusMap["Incoming Broadcast Total Size"], -1), ""))
	incomingTotalBytes := strconv.Itoa(incomingBroadcastBytes + incomingUnicastBytes)

	// Create final return map
	status := map[string]string{
		"numberOfSessions":  statusMap["Number of Sessions"],
		"numberOfUsers":     statusMap["Number of Users"],
		"currentServerTime": statusMap["Current Time"][0:19],
		"serverStartTime":   statusMap["Server Started at"][0:11] + statusMap["Server Started at"][17:],
		"incomingBytes":     incomingTotalBytes,
		"outgoingBytes":     outgoingTotalBytes,
	}

	return status, 0
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
