package support

import (
	"fmt"
	"io/ioutil"
	"strings"

	"../def"
	"../glob"
)

func interpretInput(con *glob.ConnectionData, input string) {

	/*Clean up user input*/
	alphaChar := AlphaCharOnly(input)
	alphaCharLen := len(alphaChar)
	message := StripCtlAndExtFromBytes(input)
	msg := strings.ReplaceAll(message, "\n", "")
	msg = strings.ReplaceAll(msg, "\r", "")
	msg = strings.ReplaceAll(msg, "\t", "")
	msg = strings.TrimSpace(msg)

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
	/*NEW/Login/Password area*/
	if con.State == def.CON_STATE_WELCOME {
		if command == "new" {
			buf := fmt.Sprintf("Names must be between %d and %d letters long, A-z only.", def.MIN_PLAYER_NAME_LENGTH, def.MAX_PLAYER_NAME_LENGTH)
			WriteToDesc(con, buf)
			WriteToDesc(con, "What name would you like to go by?")
			con.State = def.CON_STATE_NEW_LOGIN
		} else {
			_, err := ioutil.ReadFile(def.DATA_DIR + def.PLAYER_DIR + alphaChar)
			if err != nil {
				WriteToDesc(con, "Couldn't find a player by that name.")
				WriteToDesc(con, "Try again, or type 'NEW' to create a new character.")
				WriteToDesc(con, "Name:")
			} else {
				/* Login check goes here alphaChar*/
				con.State = def.CON_STATE_PASSWORD
				con.Name = alphaChar
				WriteToDesc(con, "Password:")
			}
		}
	} else if con.State == def.CON_STATE_PASSWORD {
		WriteToDesc(con, "Welcome back, "+con.Name+"!")
		con.State = def.CON_STATE_PLAYING
	} else if con.State == def.CON_STATE_NEW_LOGIN {
		if alphaCharLen > def.MIN_PLAYER_NAME_LENGTH && alphaCharLen < def.MAX_PLAYER_NAME_LENGTH {
			con.Name = alphaChar
			_, err := ioutil.ReadFile(def.DATA_DIR + def.PLAYER_DIR + alphaChar)
			if err != nil {
				WriteToDesc(con, "Player name is already taken! Try again.")
				WriteToDesc(con, "Name:")
			} else {
				WriteToDesc(con, "Are you sure you want your name to be known as '"+alphaChar+"'? (y/n)")
				con.State = def.CON_STATE_NEW_LOGIN_CONFIRM
			}
		} else {
			WriteToDesc(con, "That isn't a acceptable name... Try again:")
		}

	} else if con.State == def.CON_STATE_NEW_LOGIN_CONFIRM {
		if command == "y" || command == "yes" {
			con.Player = CreatePlayer(con)
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
			con.State = def.CON_STATE_PLAYING
			//support.WritePlayer()
		} else {
			WriteToDesc(con, "Passwords didn't match, try again.")
			WriteToDesc(con, "Password:")
		}

	} else if con.State == def.CON_STATE_PLAYING {
		/***************/
		/*Commands area*/
		/***************/
		if command == "quit" {
			WriteToDesc(con, "Goodbye!")
			buf := fmt.Sprintf("%s has quit.", con.Name)
			WriteToAll(buf)

			con.State = def.CON_STATE_DISCONNECTING
		} else if command == "who" {
			output := "Players online:\n"

			for x := 0; x <= glob.ConnectionListMax; x++ {
				var p *glob.ConnectionData = &glob.ConnectionList[x]
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
		} else if command == "writetest" {
			if arglen > 0 {
				WritePlayer(con.Player)
				WriteToDesc(con, "Wrote test.")
			} else {
				WriteToDesc(con, "But, what do you want to say?")
			}

		} else {
			WriteToDesc(con, "That isn't a valid command.")
		}

	} else if con.State == def.CON_STATE_DISCONNECTING {
		con.Valid = false
		con.Desc.Close()
		return
	}
}
