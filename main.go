package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"./def"
	"./support"
)

func checkError(err error, fatal bool) {
	if err != nil {
		buf := fmt.Sprintf("Fatal error: %s", err.Error())
		log.Println(buf)
		if fatal {
			os.Exit(1)
		}
	}
}

func main() {

	/*Find Network*/
	addr, err := net.ResolveTCPAddr("tcp", def.DEFAULT_PORT)
	checkError(err, def.ERROR_FATAL)

	/*Open Listener*/
	ServerListener, err := net.ListenTCP("tcp", addr)
	checkError(err, def.ERROR_FATAL)

	/*Print Connection*/
	buf := fmt.Sprint("Server online at: %s", addr.String())
	log.Println(buf)

	/*Listen for connections*/
	mainLoop()
	ServerListener.Close()
}

func mainLoop() {

	/*Check for new connections*/
	conn, err := ServerListener.Accept()
	if err == nil {
		log.Println("[INFO] New connection.")
		conn.newDescriptor(conn)

		/*Change connections settings*/
		conn.SetLinger(-1)
		conn.SetNoDelay(true)
		conn.SetReadBuffer(10000)     //10k, 10 seconds of insanely-fast typing
		conn.SetWriteBuffer(12500000) //12.5MB, 10 second buffer at 10mbit
	} else {
		log.Println("[WARN] Attempted to accept invalid connection.")
		return
	}

	for p := range ConnectionList {
		readConnection(p)
	}
}

func newDescriptor(desc net.Conn) {
	WriteToDesc(desc, "Connected!")
	newConnection := ConnectionData{connection: desc, address: def.STRING_UNKNOWN, state: def.CON_STATE_WELCOME, connectedTime: time.Now(), idleTime: time.Now(), player: nil, valid: true}
	ConnectionList = append(ConnectionList, newConnection)
	WriteToDesc(desc, "To create a new login, type: new\nLogin:")
}

func readConnection(player Connection) {
	reader := bufio.NewReader(player.desc)

	umes, err := reader.ReadString('\n')

	if err != nil {
		lostConnection(desc)
		desc.Close()
		return
	}
	message := support.StripCtlAndExtFromBytes(umes)

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

		if player.state == def.CON_STATE_WELCOME {
			if slen > 3 && slen < 128 {
				player.name = fmt.Sprintf("%s", msg)
				player.state = CON_STATE_PLAYING

				syncPlayer(player)

				WriteToDesc(desc, "Okay, you will be called "+msg)
				showCommands(desc)
				buf := fmt.Sprintf("%s has joined!", msg)
				WriteToAll(buf)
			} else {
				WriteToDesc(desc, "That isn't a valid name.")
			}
		} else if player.state == def.CON_STATE_PLAYING {
			if command == "quit" {
				WriteToDesc(desc, "Goodbye")
				buf := fmt.Sprintf("%s has quit.", player.name)
				desc.Close()

				pnum := findPlayer(desc)
				removePlayer(pnum)

				WriteToAll(buf)
				return
			} else if command == "who" {
				output := "Players online:\n"

				max := len(ConnectionList)
				for x, p := range ConnectionList {
					buf := ""

					if p.state == def.CON_STATE_PLAYING {
						buf = fmt.Sprintf("%d: %s", x+1, p.name)
					} else {
						buf = fmt.Sprintf("%d: %s", x+1, "(Connecting)")
					}
					output = output + buf
					if x <= max {
						output = output + "\r\n"
					}
				}
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

	for _, con := range ConnectionList {
		if con.state == CON_STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			con.desc.Write([]byte(message))
		}
	}
	log.Println("[ALL] " + text)
}

func WriteToOthers(desc net.Conn, text string) {

	for _, con := range ConnectionList {
		if con.desc != desc && con.state == def.CON_STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			con.desc.Write([]byte(message))
		}
	}
	fmt.Println(text)
}

func lostConnection(desc net.Conn) {

	pnum := findPlayer(desc)
	if pnum >= 0 {
		if ConnectionList[pnum].state == def.CON_STATE_PLAYING {
			msg := fmt.Sprintf("%s disconnected.", ConnectionList[pnum].name)
			go WriteToAll(msg)
			removePlayer(pnum)
			return
		}
	}

	removePlayer(pnum)
	fmt.Println("Connection dropped at login.")
}

func removePlayer(i int) {
	if i >= 0 {
		ConnectionList = append(ConnectionList[:i], ConnectionList[i+1:]...)
		fmt.Println("Player removed.")
	}
}
