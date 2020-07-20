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
			_, err = desc.Write([]byte("\r\nTo create a new login, type: new\r\nLogin: \r\n"))
			checkError(err, def.ERROR_NONFATAL)

			go newDescriptor(desc)
		}

	}
}

func newDescriptor(desc *net.TCPConn) {

	/*--- LOCK ---*/
	glob.ConnectionListLock.Lock()
	defer glob.ConnectionListLock.Unlock()
	/*--- LOCK ---*/

	for x := 1; x <= glob.ConnectionListMax; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con.Valid == true {
			continue
		} else {
			/*Replace*/
			con.Name = def.STRING_UNKNOWN
			con.Desc = desc
			con.Address = desc.LocalAddr().String()
			con.State = def.CON_STATE_WELCOME
			con.ConnectedTime = time.Now()
			con.IdleTime = time.Now()
			con.Id = x
			con.BytesOut = 0
			con.BytesIn = 0
			con.Player = nil
			con.Valid = true

			buf := fmt.Sprintf("Recycling connection #%d.", x)
			log.Println(buf)

			go readConnection(&glob.ConnectionList[x])
			return
		}
	}

	/*Generate new descriptor data*/
	if glob.ConnectionListMax >= def.MAX_DESCRIPTORS-1 {
		log.Println("MAX_DESCRIPTORS REACHED!")
		desc.Write([]byte("Sorry, something has gone wrong (MAX_DESCRIPTORS)!\r\nGoodbye!\r\n"))
		return
	}

	/*Create*/
	glob.ConnectionListMax++
	newConnection := glob.ConnectionData{
		Name:          def.STRING_UNKNOWN,
		Desc:          desc,
		Address:       desc.LocalAddr().String(),
		State:         def.CON_STATE_WELCOME,
		ConnectedTime: time.Now(),
		IdleTime:      time.Now(),
		Id:            glob.ConnectionListMax,
		BytesOut:      0,
		BytesIn:       0,
		Player:        nil,
		Valid:         true}
	glob.ConnectionList[glob.ConnectionListMax] = newConnection
	buf := fmt.Sprintf("Created new connection #%d.", glob.ConnectionListMax)
	log.Println(buf)

	go readConnection(&glob.ConnectionList[glob.ConnectionListMax])
	return

}

func readConnection(con *glob.ConnectionData) {

	/*--- LOCK ---*/
	glob.ConnectionListLock.Lock()
	reader := bufio.NewReader(con.Desc)
	glob.ConnectionListLock.Unlock()
	/*--- UNLOCK ---*/

	for {

		umes, err := reader.ReadString('\n')

		if err != nil {
			glob.ConnectionListLock.Lock()
			descWriteError(con, err)
			glob.ConnectionListLock.Unlock()
			return
		}

		//TODO max line length and max lines/sec
		if err == nil && umes != "" {

			/*Clean up user input*/
			alphaChar := support.AlphaCharOnly(umes)
			alphaCharLen := len(alphaChar)
			message := support.StripCtlAndExtFromBytes(umes)
			msg := strings.ReplaceAll(message, "\n", "")
			msg = strings.ReplaceAll(msg, "\r", "")
			msg = strings.ReplaceAll(msg, "\t", "")
			msg = strings.TrimSpace(msg)

			//TODO, split into normal commands and fight/round commands, if we do rounds.
			if msg != "" {

				slen := len(msg)
				command := ""
				aargs := ""
				arglen := -1

				args := strings.Split(msg, " ")

				//If we have arguments
				if slen > 0 {

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
				/*--- LOCK ---*/
				glob.ConnectionListLock.Lock()
				/*--- LOCK ---*/

				/*NEW/Login/Password area*/
				if con.State == def.CON_STATE_WELCOME {
					if command == "new" {
						buf := fmt.Sprintf("Names must be between %d and %d letters long, A-z only.", def.MIN_PLAYER_NAME_LENGTH, def.MAX_PLAYER_NAME_LENGTH)
						WriteToDesc(con, buf)
						WriteToDesc(con, "What name would you like to go by?")
						con.State = def.CON_STATE_NEW_LOGIN
					} else {
						/* Login check goes here alphaChar*/
						con.State = def.CON_STATE_PASSWORD
						con.Name = alphaChar
						WriteToDesc(con, "Password:")
					}
				} else if con.State == def.CON_STATE_PASSWORD {
					WriteToDesc(con, "Welcome back, "+con.Name+"!")
					con.State = def.CON_STATE_PLAYING
				} else if con.State == def.CON_STATE_NEW_LOGIN {
					if alphaCharLen > def.MIN_PLAYER_NAME_LENGTH && alphaCharLen < def.MAX_PLAYER_NAME_LENGTH {
						con.Name = alphaChar
						WriteToDesc(con, "Are you sure you want your name to be known as '"+alphaChar+"'? (y/n)")
						con.State = def.CON_STATE_NEW_LOGIN_CONFIRM
					} else {
						WriteToDesc(con, "That isn't a acceptable name... Try again:")
					}

				} else if con.State == def.CON_STATE_NEW_LOGIN_CONFIRM {
					if command == "y" || command == "yes" {
						WriteToDesc(con, "You shall be called "+alphaChar+", then...")
						WriteToDesc(con, "Password:")
						con.State = def.CON_STATE_NEW_PASSWORD
					} else {
						con.State = def.CON_STATE_NEW_LOGIN
						WriteToDesc(con, "What name would you like to go by then?")
					}

				} else if con.State == def.CON_STATE_NEW_PASSWORD {
					WriteToDesc(con, "Type again to confirm:")
					con.State = def.CON_STATE_NEW_PASSWORD_CONFIRM
				} else if con.State == def.CON_STATE_NEW_PASSWORD_CONFIRM {

					/*Check password*/
					if 1 == 1 {
						WriteToDesc(con, "Password confirmed, logging in!")
						showCommands(con)
						con.State = def.CON_STATE_PLAYING
					} else {
						WriteToDesc(con, "Passwords didn't match, try again.")
						WriteToDesc(con, "Password:")
					}

					/*Commands area*/
				} else if con.State == def.CON_STATE_PLAYING {
					if command == "quit" {
						WriteToDesc(con, "Goodbye!")
						buf := fmt.Sprintf("%s has quit.", con.Name)
						WriteToAll(buf)

						con.State = def.CON_STATE_DISCONNECTING
					} else if command == "who" {
						output := "Players online:\n"

						for x := 0; x <= glob.ConnectionListMax; x++ {
							var p *glob.ConnectionData
							p = &glob.ConnectionList[x]
							if p.Valid == false {
								continue
							}
							buf := ""

							if p.State == def.CON_STATE_PLAYING {
								buf = fmt.Sprintf("%d: %s", x, p.Name)
							} else {
								buf = fmt.Sprintf("%d: %s", x, "(Connecting)")
							}
							output = output + buf
							if x <= glob.ConnectionListMax {
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
				if con.State == def.CON_STATE_DISCONNECTING {
					con.Valid = false
					con.Desc.Close()
					/*--- UNLOCK ---*/
					glob.ConnectionListLock.Unlock()
					/*--- UNLOCK ---*/
					return
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

func descWriteError(c *glob.ConnectionData, err error) {
	if err != nil {
		checkError(err, def.ERROR_NONFATAL)

		buf := fmt.Sprintf("%s lost their connection.", c.Name)
		WriteToOthers(c, buf)

		c.State = def.CON_STATE_DISCONNECTED
		c.Valid = false
		c.Desc.Close()
	}
}

func WriteToDesc(c *glob.ConnectionData, text string) {
	message := fmt.Sprintf("%s\r\n", text)
	bytes, err := c.Desc.Write([]byte(message))
	c.BytesOut += bytes

	descWriteError(c, err)
}

func WriteToAll(text string) {

	for x := 0; x <= glob.ConnectionListMax; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con.Valid == false {
			continue
		}
		if con.State == def.CON_STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			bytes, err := con.Desc.Write([]byte(message))
			con.BytesOut += bytes

			descWriteError(con, err)
		}
	}
	log.Println(text)
}

func WriteToOthers(c *glob.ConnectionData, text string) {

	for x := 0; x <= glob.ConnectionListMax; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con.Valid == false {
			continue
		}
		if con.Desc != c.Desc && con.State == def.CON_STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			bytes, err := con.Desc.Write([]byte(message))
			con.BytesOut += bytes

			descWriteError(c, err)
		}
	}
	log.Println(text)
}
