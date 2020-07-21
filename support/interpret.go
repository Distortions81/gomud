package support

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"../def"
	"../glob"
)

func interpretInput(con *glob.ConnectionData, input string) {

	/*Clean up user input*/
	alphaChar := AlphaCharOnly(input)
	alphaCharLen := len(alphaChar)
	inputc := strings.ReplaceAll(input, "\n", "")
	inputc = strings.ReplaceAll(inputc, "\r", "")
	inputc = strings.ReplaceAll(inputc, "\t", "")
	inputc = strings.TrimSpace(inputc)
	msg := StripCtlAndExtFromBytes(inputc)

	msgLen := len(msg)
	command := ""
	aargs := ""
	arglen := -1

	args := strings.Split(msg, " ")

	//If we have arguments
	if msgLen > 0 {

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
			_, found := ReadPlayer(alphaChar, false)
			if found == false {
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
		player, _ := ReadPlayer(con.Name, true)
		con.Player = player

		err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(inputc))

		if err == nil {
			con.Player.Connection = con
			WriteToDesc(con, "Welcome back, "+player.Name+"!")
			WriteToAll(player.Name + " has joined.")
			con.State = def.CON_STATE_PLAYING
		} else {
			log.Println("Invalid password attempt: " + player.Name + " ip: " + con.Address)
			time.Sleep(5 * time.Second)
			WriteToDesc(con, "Invalid password.")
			time.Sleep(5 * time.Second)
			WriteToDesc(con, "Reconnect, and try again...")
			con.State = def.CON_STATE_DISCONNECTING
		}
	} else if con.State == def.CON_STATE_NEW_LOGIN {
		if alphaCharLen > def.MIN_PLAYER_NAME_LENGTH &&
			alphaCharLen < def.MAX_PLAYER_NAME_LENGTH &&
			alphaChar != def.STRING_UNKNOWN {
			con.Name = alphaChar
			_, found := ReadPlayer(alphaChar, false)
			if found {
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
			con.Player = CreatePlayerFromDesc(con)
			con.Player.Connection = con
			WriteToDesc(con, "You shall be called "+con.Name+", then...")
			WriteToDesc(con, "Passwords must be between 9 and 72 characters long, and contain at least 2 numbers/symbols.")
			WriteToDesc(con, "Password:")
			con.State = def.CON_STATE_NEW_PASSWORD
		} else {
			con.State = def.CON_STATE_NEW_LOGIN
			WriteToDesc(con, "What name would you like to go by then?")
		}

	} else if con.State == def.CON_STATE_NEW_PASSWORD {
		symbolCount := len(NonAlphaCharOnly(inputc))
		inputcLen := len(inputc)
		if inputcLen >= 8 && inputcLen <= 72 && symbolCount >= 2 {
			con.Temp = inputc
			WriteToDesc(con, "Type again to confirm:")
			con.State = def.CON_STATE_NEW_PASSWORD_CONFIRM
		} else {
			WriteToDesc(con, "That isn't an acceptable password!")
			debug := fmt.Sprintf("len: %d, sym: %d", inputcLen, symbolCount)
			log.Println(debug)
			WriteToDesc(con, "Password:")
		}
	} else if con.State == def.CON_STATE_NEW_PASSWORD_CONFIRM {

		/*Hash password*/
		if inputc == con.Temp {
			WriteToDesc(con, "Encrypting password... One second please!")
			hash, err := bcrypt.GenerateFromPassword([]byte(msg), def.PASSWORD_HASH_COST)
			if err != nil {
				CheckError("interp: password hash", err, def.ERROR_NONFATAL)
				WriteToDesc(con, "Password encryption failed, sorry something is wrong.")

				//TODO disconnect/invalidate and report
				return
			}
			con.Temp = ""
			con.Player.Password = string(hash)
			WriteToDesc(con, "Done, logging in!")
			con.State = def.CON_STATE_PLAYING
			WritePlayer(con.Player, true)
		} else {
			con.Temp = ""
			WriteToDesc(con, "Passwords didn't match, try again.")
			WriteToDesc(con, "Password:")
		}

	} else if con.State == def.CON_STATE_PLAYING {
		/***************/
		/*Commands area*/
		/***************/
		if command == "quit" {
			WritePlayer(con.Player, true)
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
		} else if command == "save" {
			okay := WritePlayer(con.Player, true)
			if okay == false {
				WriteToPlayer(con.Player, "Saving character failed!!!")
			}

		} else {
			WriteToDesc(con, "That isn't a valid command.")
		}

	}
	if con.State == def.CON_STATE_DISCONNECTING {
		con.Desc.Close()
		con.Valid = false
		if con.Player != nil && con.Player.Valid {
			con.Player.Valid = false
			con.Player.Connection = nil
			con.Player = nil
		}
		con = nil
		return
	}
}
