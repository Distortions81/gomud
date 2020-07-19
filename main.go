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
	"./glob"
	"./support"
)

func checkError(err error, fatal bool) {
	if err != nil {
		buf := fmt.Sprintf("error: %s", err.Error())
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
	glob.ServerListener = ServerListener
	checkError(err, def.ERROR_FATAL)

	/*Print Connection*/
	buf := fmt.Sprintf("Server online at: %s", addr.String())
	log.Println(buf)

	/*Listen for connections*/
	mainLoop()
	//ServerListener.Close()
}

func mainLoop() {

	/*Seperate thread, wait for new connections*/
	go func() {
		for glob.ServerState == def.SERVER_RUNNING {
			/*Check for new connections*/
			conn, err := glob.ServerListener.AcceptTCP()

			if err == nil {
				buf := fmt.Sprintf("New connection from %s.", conn.LocalAddr().String())
				log.Println(buf)

				/*Change connections settings*/
				conn.SetLinger(-1)
				conn.SetNoDelay(true)
				conn.SetReadBuffer(10000)     //10k, 10 seconds of insanely-fast typing
				conn.SetWriteBuffer(12500000) //12.5MB, 10 second buffer at 10mbit

				newDescriptor(conn)
			}

		}
	}()

	/* Main game loop */
	for glob.ServerState == def.SERVER_RUNNING {

		/*--- LOCK ---*/
		glob.ConnectionListLock.Lock()
		/*--- LOCK ---*/

		/*Read all connections*/
		for i, p := range glob.ConnectionList {
			if p.Valid {
				readConnection(&glob.ConnectionList[i])
			}
		}

		/*Memory cleanup area, do this at the *very end* only*/
		//
		/*Connections*/
		for i, p := range glob.ConnectionList {
			/*Remove no longer valid connections*/
			if p.Valid == false {
				glob.ConnectionList = append(glob.ConnectionList[:i], glob.ConnectionList[i+1:]...)
			}
			/*Disconnect people here be to be cleaner*/
			if p.State == def.CON_STATE_DISCONNECTING {
				p.State = def.CON_STATE_DISCONNECTED
				lostConnection(&p)
				p.Desc.Close()
				p.Valid = false
			}
		}
		/*--- UNLOCK ---*/
		glob.ConnectionListLock.Unlock()
		/*--- UNLOCK ---*/

	}
	log.Println("Loop exited.")
}

func newDescriptor(desc *net.TCPConn) {

	/*Respond here, so we don't have to wait for lock*/
	_, err := desc.Write([]byte("To create a new login, type: new\nLogin:"))
	if err != nil {
		checkError(err, def.ERROR_NONFATAL)
	}

	/*--- LOCK ---*/
	glob.ConnectionListLock.Lock()
	/*--- LOCK ---*/

	/* WE NEED TO TAKE OVER, TO ACCEPT NEW INCOMING CONNECTION */

	/*Generate new descriptor data*/
	glob.LastConnectionID++
	id := glob.LastConnectionID
	reader := bufio.NewReader(desc)
	addr := desc.LocalAddr()

	/*Append*/
	newConnection := glob.ConnectionData{
		Name:          def.STRING_UNKNOWN,
		Desc:          desc,
		Address:       addr.String(),
		State:         def.CON_STATE_WELCOME,
		ConnectedTime: time.Now(),
		IdleTime:      time.Now(),
		Reader:        reader,
		Id:            id,
		BytesOut:      0,
		BytesIn:       0,
		Player:        nil,
		Valid:         true}

	glob.ConnectionList = append(glob.ConnectionList, newConnection)

	/*--- UNLOCK ---*/
	glob.ConnectionListLock.Unlock()
	/*--- UNLOCK ---*/
}

func readConnection(c *glob.ConnectionData) {
	umes, err := c.Reader.ReadString('\n')

	if err == nil && umes != "" {
		/*Clean up user input*/
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

			if c.State == def.CON_STATE_WELCOME {
				if slen > 3 && slen < 128 {
					c.Name = fmt.Sprintf("%s", msg)
					c.State = def.CON_STATE_PLAYING

					WriteToDesc(c, "Okay, you will be called "+msg)
					showCommands(c)
					buf := fmt.Sprintf("%s has joined!", msg)
					WriteToAll(buf)
				} else {
					WriteToDesc(c, "That isn't a valid name.")
				}
			} else if c.State == def.CON_STATE_PLAYING {
				if command == "quit" {
					WriteToDesc(c, "Goodbye")
					buf := fmt.Sprintf("%s has quit.", c.Name)
					WriteToAll(buf)

					c.State = def.CON_STATE_DISCONNECTING
					return
				} else if command == "who" {
					output := "Players online:\n"

					max := len(glob.ConnectionList)
					for x, p := range glob.ConnectionList {
						buf := ""

						if p.State == def.CON_STATE_PLAYING {
							buf = fmt.Sprintf("%d: %s", x+1, p.Name)
						} else {
							buf = fmt.Sprintf("%d: %s", x+1, "(Connecting)")
						}
						output = output + buf
						if x <= max {
							output = output + "\r\n"
						}
					}
					WriteToDesc(c, output)
				} else if command == "say" {
					if arglen > 0 {
						out := fmt.Sprintf("%s says: %s", c.Name, aargs)
						us := fmt.Sprintf("You say: %s", aargs)

						WriteToOthers(c, out)
						WriteToDesc(c, us)
					} else {
						WriteToDesc(c, "But, what do you want to say?")
					}

				} else {
					WriteToDesc(c, "That isn't a valid command.")
					showCommands(c)
				}
			}
		}
	}
}

func showCommands(c *glob.ConnectionData) {
	us := fmt.Sprintf("commands: say, who, quit")
	WriteToDesc(c, us)
}

func WriteToDesc(c *glob.ConnectionData, text string) {
	message := fmt.Sprintf("%s\r\n", text)
	bytes, err := c.Desc.Write([]byte(message))
	c.BytesOut += bytes

	if err != nil {
		c.State = def.CON_STATE_DISCONNECTING
	}
	checkError(err, def.ERROR_NONFATAL)
}

func WriteToAll(text string) {

	for _, con := range glob.ConnectionList {
		if con.State == def.CON_STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			bytes, err := con.Desc.Write([]byte(message))
			con.BytesOut += bytes
			if err != nil {
				con.State = def.CON_STATE_DISCONNECTING
			}
			checkError(err, def.ERROR_NONFATAL)
		}
	}
	log.Println(text)
}

func WriteToOthers(c *glob.ConnectionData, text string) {

	for _, con := range glob.ConnectionList {
		if con.Desc != c.Desc && con.State == def.CON_STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			bytes, err := con.Desc.Write([]byte(message))
			con.BytesOut += bytes
			if err != nil {
				con.State = def.CON_STATE_DISCONNECTING
			}
			checkError(err, def.ERROR_NONFATAL)
		}
	}
	log.Println(text)
}

func lostConnection(c *glob.ConnectionData) {

	if c.Name != def.STRING_UNKNOWN {
		msg := fmt.Sprintf("%s disconnected.", c.Name)
		go WriteToAll(msg)
		return
	}
	buf := fmt.Sprintf("%s disconnected.", c.Address)
	log.Println(buf)
}
