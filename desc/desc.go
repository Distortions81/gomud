package desc

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
	"time"

	"../def"
	"../glob"
	"../support"
)

func NewDescriptor(desc *net.TCPConn) {

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
			con.ConnectedFor = time.Now()
			con.IdleTime = time.Now()
			con.Id = x
			con.BytesOut = 0
			con.BytesIn = 0
			con.Player = nil
			con.Valid = true

			buf := fmt.Sprintf("Recycling connection #%d.", x)
			log.Println(buf)

			go ReadConnection(&glob.ConnectionList[x])
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
		Name:         def.STRING_UNKNOWN,
		Desc:         desc,
		Address:      desc.LocalAddr().String(),
		State:        def.CON_STATE_WELCOME,
		ConnectedFor: time.Now(),
		IdleTime:     time.Now(),
		Id:           glob.ConnectionListMax,
		BytesOut:     0,
		BytesIn:      0,
		Player:       nil,
		Valid:        true}
	glob.ConnectionList[glob.ConnectionListMax] = newConnection
	buf := fmt.Sprintf("Created new connection #%d.", glob.ConnectionListMax)
	log.Println(buf)

	go desc.ReadConnection(&glob.ConnectionList[glob.ConnectionListMax])
	return

}

func ReadConnection(con *glob.ConnectionData) {

	/*--- LOCK ---*/
	glob.ConnectionListLock.Lock()
	/*Create reader*/
	reader := bufio.NewReader(con.Desc)
	glob.ConnectionListLock.Unlock()
	/*--- UNLOCK ---*/

	for {

		umes, err := reader.ReadString('\n')

		if err != nil {
			glob.ConnectionListLock.Lock()
			DescWriteError(con, err)
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
						desc.WriteToDesc(con, buf)
						desc.WriteToDesc(con, "What name would you like to go by?")
						con.State = def.CON_STATE_NEW_LOGIN
					} else {
						filedata, err := ioutil.ReadFile(def.PLAYER_DIR + alphaChar)
						if err != nil {
							desc.WriteToDesc(con, "Couldn't find a player by that name.")
							desc.WriteToDesc(con, "Try again, or type 'NEW' to create a new character.")
							desc.WriteToDesc(con, "Name:")
						} else {
							/* Login check goes here alphaChar*/
							con.State = def.CON_STATE_PASSWORD
							con.Name = alphaChar
							desc.WriteToDesc(con, "Password:")
						}
					}
				} else if con.State == def.CON_STATE_PASSWORD {
					desc.WriteToDesc(con, "Welcome back, "+con.Name+"!")
					con.State = def.CON_STATE_PLAYING
				} else if con.State == def.CON_STATE_NEW_LOGIN {
					if alphaCharLen > def.MIN_PLAYER_NAME_LENGTH && alphaCharLen < def.MAX_PLAYER_NAME_LENGTH {
						con.Name = alphaChar
						desc.WriteToDesc(con, "Are you sure you want your name to be known as '"+alphaChar+"'? (y/n)")
						con.State = def.CON_STATE_NEW_LOGIN_CONFIRM
					} else {
						desc.WriteToDesc(con, "That isn't a acceptable name... Try again:")
					}

				} else if con.State == def.CON_STATE_NEW_LOGIN_CONFIRM {
					if command == "y" || command == "yes" {
						desc.WriteToDesc(con, "You shall be called "+alphaChar+", then...")
						desc.WriteToDesc(con, "Password:")
						con.State = def.CON_STATE_NEW_PASSWORD
					} else {
						con.State = def.CON_STATE_NEW_LOGIN
						desc.WriteToDesc(con, "What name would you like to go by then?")
					}

				} else if con.State == def.CON_STATE_NEW_PASSWORD {
					desc.WriteToDesc(con, "Type again to confirm:")
					con.State = def.CON_STATE_NEW_PASSWORD_CONFIRM
				} else if con.State == def.CON_STATE_NEW_PASSWORD_CONFIRM {

					/*Check password*/
					if 1 == 1 {
						desc.WriteToDesc(con, "Password confirmed, logging in!")
						showCommands(con)
						con.State = def.CON_STATE_PLAYING
					} else {
						desc.WriteToDesc(con, "Passwords didn't match, try again.")
						desc.WriteToDesc(con, "Password:")
					}

					/*Commands area*/
				} else if con.State == def.CON_STATE_PLAYING {
					if command == "quit" {
						desc.WriteToDesc(con, "Goodbye!")
						buf := fmt.Sprintf("%s has quit.", con.Name)
						desc.WriteToAll(buf)

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
						desc.WriteToDesc(con, output)
					} else if command == "say" {
						if arglen > 0 {
							out := fmt.Sprintf("%s says: %s", con.Name, aargs)
							us := fmt.Sprintf("You say: %s", aargs)

							desc.WriteToOthers(con, out)
							desc.WriteToDesc(con, us)
						} else {
							desc.WriteToDesc(con, "But, what do you want to say?")
						}

					} else {
						desc.WriteToDesc(con, "That isn't a valid command.")
						showCommands(con)
					}
				}
				if con.State == def.CON_STATE_DISCONNECTING {
					con.Valid = false
					con.Desc.Close()
					/*--- UNLOCK ---*/
					glob.ConnectionListLock.Unlock()
					/*--- UNLOCK ---*/
					return /*EXIT*/
				}
				/*--- UNLOCK ---*/
				glob.ConnectionListLock.Unlock()
				/*--- UNLOCK ---*/
			}
		}
	}
}
func DescWriteError(c *glob.ConnectionData, err error) {
	if err != nil {
		support.CheckError(err, def.ERROR_NONFATAL)

		if c.Name != def.STRING_UNKNOWN && c.State == def.CON_STATE_PLAYING {
			buf := fmt.Sprintf("%s lost their connection.", c.Name)
			desc.WriteToOthers(c, buf)
		} else {
			buf := fmt.Sprintf("%s disconnected.", c.Address)
			log.Println(buf)
		}

		c.State = def.CON_STATE_DISCONNECTED
		c.Valid = false
		c.Desc.Close()
	}
}

func WriteToDesc(c *glob.ConnectionData, text string) {
	message := fmt.Sprintf("%s\r\n", text)
	bytes, err := c.Desc.Write([]byte(message))
	c.BytesOut += bytes

	desc.DescWriteError(c, err)
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

			desc.DescWriteError(con, err)
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

			desc.DescWriteError(c, err)
		}
	}
	log.Println(text)
}
