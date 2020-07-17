package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"./support"
	"github.com/sasha-s/go-deadlock"
)

const STATE_WELCOME = 0
const STATE_PLAYING = 10

var ConnectionListLock deadlock.RWMutex
var ConnectionList []Connection

type Connection struct {
	desc  net.Conn
	life  time.Time
	state int
	name  string
}

func main() {
	service := ":7777"

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer listener.Close()

	fmt.Println("Online.")
	listenForConnections(listener)
}

func listenForConnections(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		//Log new desc
		fmt.Println("new descriptor.")
		newDescriptor(conn)
		time.Sleep(time.Millisecond * 100)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func newDescriptor(desc net.Conn) {
	WriteToDesc(desc, "Connected!")
	newConnection := Connection{desc: desc, life: time.Now(), state: STATE_WELCOME, name: "Unknown"}
	ConnectionListLock.Lock()
	ConnectionList = append(ConnectionList, newConnection)
	ConnectionListLock.Unlock()

	time.Sleep(time.Millisecond * 100)
	WriteToDesc(desc, "What would you like to be called?")

	go readConnection(desc) //new thread!
}

func readConnection(desc net.Conn) {
	reader := bufio.NewReader(desc)

	for {

		ConnectionListLock.Lock()
		pnum := findPlayer(desc)
		if pnum < 0 {
			ConnectionListLock.Unlock()
			return //Player is dead
		}
		player := ConnectionList[pnum]
		ConnectionListLock.Unlock()

		umes, err := reader.ReadString('\n')
		message := support.StripCtlAndExtFromBytes(umes)

		if err != nil {
			lostConnection(desc)
			desc.Close()
			return
		}

		msg := strings.ReplaceAll(message, "\n", "")
		msg = strings.ReplaceAll(msg, "\r", "")
		msg = strings.ReplaceAll(msg, "\t", "")
		msg = strings.TrimSpace(msg)

		if msg != "" {

			slen := len(msg)
			command := ""
			aargs := ""
			arglen := -1

			args := strings.Split(msg, " ")

			if slen > 1 {

				arglen = len(args)

				if arglen > 0 {
					command = strings.ToLower(args[0])
					if arglen > 1 {
						aargs = strings.Join(args[1:arglen], " ")
					}
				}
			}

			if player.state == STATE_WELCOME {
				if slen > 3 && slen < 128 {
					player.name = fmt.Sprintf("%s", msg)
					player.state = STATE_PLAYING

					syncPlayer(player)

					WriteToDesc(desc, "Okay, you will be called "+msg)
					showCommands(desc)
					buf := fmt.Sprintf("%s has joined!", msg)
					WriteToAll(buf)
				} else {
					WriteToDesc(desc, "That isn't a valid name.")
				}
			} else if player.state == STATE_PLAYING {
				if command == "quit" {
					WriteToDesc(desc, "Goodbye")
					buf := fmt.Sprintf("%s has quit.", player.name)
					desc.Close()

					ConnectionListLock.Lock()
					pnum := findPlayer(desc)
					removePlayer(pnum)
					ConnectionListLock.Unlock()

					WriteToAll(buf)
					return
				} else if command == "who" {
					output := "Players online:\n"

					max := len(ConnectionList)
					ConnectionListLock.RLock()
					for x, p := range ConnectionList {
						buf := ""

						if p.state == STATE_PLAYING {
							buf = fmt.Sprintf("%d: %s", x+1, p.name)
						} else {
							buf = fmt.Sprintf("%d: %s", x+1, "(Connecting)")
						}
						output = output + buf
						if x <= max {
							output = output + "\r\n"
						}
					}
					ConnectionListLock.RUnlock()
					WriteToDesc(desc, output)
				} else if command == "say" {
					if arglen > 0 {
						out := fmt.Sprintf("%s says: %s", player.name, aargs)
						us := fmt.Sprintf("You say: %s", aargs)

						WriteToOthers(desc, out)
						WriteToDesc(desc, us)
					} else {
						WriteToDesc(desc, "But, what do you want to say?")
					}

				} else {
					WriteToDesc(desc, "That isn't a valid command.")
					showCommands(desc)
				}
			}
		}

		time.Sleep(time.Millisecond * 10)
	}
}

func showCommands(desc net.Conn) {
	us := fmt.Sprintf("commands: say, who, quit")
	WriteToDesc(desc, us)
}

func WriteToDesc(desc net.Conn, text string) {
	message := fmt.Sprintf("%s\r\n", text)
	desc.Write([]byte(message))
}

func WriteToAll(text string) {
	ConnectionListLock.RLock()
	defer ConnectionListLock.RUnlock()

	for _, con := range ConnectionList {
		if con.state == STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			con.desc.Write([]byte(message))
		}
	}
	fmt.Println(text)
}

func WriteToOthers(desc net.Conn, text string) {
	ConnectionListLock.RLock()
	defer ConnectionListLock.RUnlock()

	for _, con := range ConnectionList {
		if con.desc != desc && con.state == STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			con.desc.Write([]byte(message))
		}
	}
	fmt.Println(text)
}

func lostConnection(desc net.Conn) {

	ConnectionListLock.Lock()
	defer ConnectionListLock.Unlock()

	pnum := findPlayer(desc)
	if pnum >= 0 {
		if ConnectionList[pnum].state == STATE_PLAYING {
			msg := fmt.Sprintf("%s disconnected.", ConnectionList[pnum].name)
			go WriteToAll(msg)
			removePlayer(pnum)
			return
		}
	}

	removePlayer(pnum)
	fmt.Println("A connection was lost.")
}

func findPlayer(desc net.Conn) int {

	for x, con := range ConnectionList {
		if con.desc == desc {
			return x
		}
	}

	return -1
}

func removePlayer(i int) {
	if i >= 0 {
		ConnectionList = append(ConnectionList[:i], ConnectionList[i+1:]...)
		fmt.Println("player removed.")
	}
}

func syncPlayer(player Connection) {
	//Sync everything back
	ConnectionListLock.Lock()
	pnum := findPlayer(player.desc)
	if pnum >= 0 {
		ConnectionList[pnum] = player
	}
	ConnectionListLock.Unlock()
}
