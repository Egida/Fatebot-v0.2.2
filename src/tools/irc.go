package tools

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

/*
You can config bot profile in "IRC_Login" function.
*/

type logFormation struct {
	user, nick string
}

func IRC_Conn(set_serv string) net.Conn {
	conn, err := net.Dial("tcp", set_serv)
	for err != nil {
		continue
	}
	return conn
}

func IRC_Find(read, msg string) bool {
	return strings.Contains(read, msg)
}

func IRC_Recv(cmd string, arg int) string {
	return strings.Split(cmd, " ")[arg]
}

func IRC_Send(sendIRC net.Conn, data string) {
	fmt.Fprintf(sendIRC, "%s\r\n", data)
}

func IRC_Login(log_serv net.Conn, set_chan, set_chan_pass string) {
	rand.Seed(time.Now().UnixNano())
	alphabet := 'A' + rune(rand.Intn(26))
	sAlphabet := string(alphabet)
	botID := rand.Intn(1000000)

	formation := logFormation{
		user: fmt.Sprint("USER [FATE][", sAlphabet, "][", botID, "]", " 8 * :bot"),
		nick: fmt.Sprint("NICK [FATE][", sAlphabet, "][", botID, "]"),
	}

	//login to server
	IRC_Send(log_serv, formation.user)
	IRC_Send(log_serv, formation.nick)
}
