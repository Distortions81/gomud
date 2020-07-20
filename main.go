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
	ServerListener.Close()
}

func mainLoop() {

	/*Seperate thread, wait for new connections*/
	for glob.ServerState == def.SERVER_RUNNING {
		/*Check for new connections*/
		desc, err := glob.ServerListener.AcceptTCP()

		if err == nil {
			buf := fmt.Sprintf("New connection from %s.", desc.LocalAddr().String())
			log.Println(buf)

			/*Change connections settings*/
			err := desc.SetLinger(-1)
			checkError(err, def.ERROR_NONFATAL)
			err = desc.SetNoDelay(true)
			checkError(err, def.ERROR_NONFATAL)
			err = desc.SetReadBuffer(10000) //10k, 10 seconds of insanely-fast typing
			checkError(err, def.ERROR_NONFATAL)
			err = desc.SetWriteBuffer(12500000) //12.5MB, 10 second buffer at 10mbit
			checkError(err, def.ERROR_NONFATAL)

			//TODO Add full greeting/info
			/*Respond here, so we don't have to wait for lock*/
			_, err = desc.Write([]byte("To create a new login, type: new\r\nLogin:"))
			checkError(err, def.ERROR_NONFATAL)

			go newDescriptor(desc)
		}

	}
}

func newDescriptor(desc *net.TCPConn) {

	/*--- LOCK ---*/
	glob.ConnectionListLock.Lock()
	/*--- LOCK ---*/

	/*Add to descriptor list*/
	//TODO re-use old ConnectionData on reconnect, re-atach to old player

	/*Generate new descriptor data*/
	if glob.ConnectionListMax >= def.MAX_DESCRIPTORS {
		log.Println("MAX_DESCRIPTORS REACHED!")
		desc.Write([]byte("Sorry, something has gone wrong (MAX_DESCRIPTORS)!\r\nPlease report this error!\r\n"))
		return
	}
	glob.ConnectionListMax++
	id := glob.ConnectionListMax
	addr := desc.LocalAddr()

	/*Append*/
	newConnection := glob.ConnectionData{
		Name:          def.STRING_UNKNOWN,
		Desc:          desc,
		Address:       addr.String(),
		State:         def.CON_STATE_WELCOME,
		ConnectedTime: time.Now(),
		IdleTime:      time.Now(),
		Id:            id,
		BytesOut:      0,
		BytesIn:       0,
		Player:        nil,
		Valid:         true}

	glob.ConnectionList[id] = newConnection

	go readConnection(&glob.ConnectionList[id])

	/*--- UNLOCK ---*/
	glob.ConnectionListLock.Unlock()
	/*--- UNLOCK ---*/
}

func readConnection(con *glob.ConnectionData) {

	/*--- LOCK ---*/
	glob.ConnectionListLock.Lock()
	reader := bufio.NewReader(con.Desc)
	glob.ConnectionListLock.Unlock()
	/*--- UNLOCK ---*/

	for {

		umes, err := reader.ReadString('\n')

		//TODO max line length and max lines/sec
		if err == nil && umes != "" {

			/*Clean up user input*/
			//TODO, strip non-printable, space and telnet but not unicode.
			message := support.StripCtlAndExtFromBytes(umes)
			msg := strings.ReplaceAll(message, "\n", "")
			msg = strings.ReplaceAll(msg, "\r", "")
			msg = strings.ReplaceAll(msg, "\t", "")
			msg = strings.TrimSpace(msg)

			//TODO, split into normal commands and fight/round commands, if we do rounds.
			if msg != "" {
				/*--- LOCK ---*/
				glob.ConnectionListLock.Lock()
				/*--- LOCK ---*/

				slen := len(msg)
				command := ""
				aargs := ""
				arglen := -1

				args := strings.Split(msg, " ")

				//If we have arguments
				if slen > 1 {

					arglen = len(args)
					if arglen > 0 {
						//Command name, tolower
						command = strings.ToLower(args[0])
						//arguments
						if arglen > 1 {
							aargs = strings.Join(args[1:arglen], " ")
						}
					}
				}

				//Move all this to handlers, to get rid of if/elseif mess,
				//and to enable autocomplete and shortcuts/aliases.

				if con.State == def.CON_STATE_WELCOME {
					if slen > 3 && slen < 128 {

						con.Name = fmt.Sprintf("%s", msg)
						con.State = def.CON_STATE_PLAYING //needs locks

						WriteToDesc(con, "Okay, you will be called "+msg)
						showCommands(con)
						buf := fmt.Sprintf("%s has joined!", msg)
						WriteToAll(buf)
					} else {
						WriteToDesc(con, "That isn't a valid name.")
					}
				} else if con.State == def.CON_STATE_PLAYING {
					if command == "quit" {
						WriteToDesc(con, "Goodbye")
						buf := fmt.Sprintf("%s has quit.", con.Name)
						WriteToAll(buf)

						con.State = def.CON_STATE_DISCONNECTING
						return
					} else if command == "who" {
						output := "Players online:\n"

						max := len(glob.ConnectionList)
						for x, p := range glob.ConnectionList {
							if x >= def.MAX_DESCRIPTORS {
								break
							}
							if p.Valid == false {
								continue
							}
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
						WriteToDesc(con, output)
					} else if command == "say" {
						if arglen > 0 {
							out := fmt.Sprintf("%s says: %s", con.Name, aargs)
							us := fmt.Sprintf("You say: %s", aargs)

							WriteToOthers(con, out)
							WriteToDesc(con, us)
						} else {
							WriteToDesc(con, "But, what do you want to say?")
						}

					} else {
						WriteToDesc(con, "That isn't a valid command.")
						showCommands(con)
					}
				}
				/*--- UNLOCK ---*/
				glob.ConnectionListLock.Unlock()
				/*--- UNLOCK ---*/
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

	for x, con := range glob.ConnectionList {
		if x >= def.MAX_DESCRIPTORS {
			break
		}
		if con.Valid == false {
			continue
		}
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

	for x, con := range glob.ConnectionList {
		if x >= def.MAX_DESCRIPTORS {
			break
		}
		if con.Valid == false {
			continue
		}
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
