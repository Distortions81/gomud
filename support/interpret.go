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

type Commands struct {
	CommandList []Command
}

type Command struct {
	name    string
	Command func()
	Level   int
}

var CL Commands

func AddCommands() {

}

func interpretInput(con *glob.ConnectionData, input string) {

	if con == nil && !con.Valid {
		return
	}

	/*********************/
	/*Clean up user input*/
	/*********************/
	overflow := false

	input = strings.TrimSpace(input)
	input = StripControl(input)
	input, overflow = TruncateString(input, def.MAX_INPUT_LENGTH)
	inputLen := len(input)

	if overflow {
		WriteToDesc(con, "That line was too long, truncating...")
	}

	command := ""
	longArg := ""
	argNum := 0
	//If we have arguments
	if inputLen > 0 {
		args := strings.Split(input, " ")
		argNum = len(args)

		if argNum > 0 {
			//Command name, tolower
			command = strings.ToLower(args[0])

			//all arguments after command
			if argNum > 1 {
				longArg = strings.Join(args[1:argNum], " ")
			}
		}
	}

	//Set player as no longer idle
	con.IdleTime = time.Now()

	/*Skip out if we are playing the game*/
	if con.State == def.CON_STATE_PLAYING {
		if con.Player != nil && con.Player.Valid &&
			con.Player.Connection != nil &&
			con.Player.Connection.Valid {

			PlayerCommand(con.Player, command, longArg)
			return
		}
	}

	//For names, only letters allowed
	alphaChar := AlphaOnly(input)
	alphaCharLen := len(alphaChar)

	if con.State == def.CON_STATE_WELCOME {
		if command == "new" {
			buf := fmt.Sprintf("Names must be between %d and %d letters long, A-z only.", def.MIN_PLAYER_NAME_LENGTH, def.MAX_PLAYER_NAME_LENGTH)
			WriteToDesc(con, buf)
			WriteToDesc(con, "What name would you like to go by?")
			con.State = def.CON_STATE_NEW_LOGIN
		} else {
			if alphaCharLen > def.MIN_PLAYER_NAME_LENGTH &&
				alphaCharLen < def.MAX_PLAYER_NAME_LENGTH &&
				alphaChar != def.STRING_UNKNOWN {

				_, found := ReadPlayer(alphaChar, false)

				if found == false {
					WriteToDesc(con, "Couldn't find a player by that name.")
					WriteToDesc(con, "Try again, or type 'NEW' to create a new character.")
					WriteToDesc(con, "Name:")
				} else {
					con.State = def.CON_STATE_PASSWORD
					con.Name = alphaChar
					WriteToDesc(con, "Password:")
				}
			} else {
				if con != nil && con.Valid {
					buf := fmt.Sprintf("Illegal characters or length: login attempt by %v", con.Address)
					log.Println(buf)
					con.State = def.CON_STATE_DISCONNECTING
				} else {
					buf := fmt.Sprintf("Illegal characters or length: login attempt by unknown.")
					log.Println(buf)
				}
			}
		}
	} else if con.State == def.CON_STATE_PASSWORD {
		player, _ := ReadPlayer(con.Name, true)
		con.Player = player

		err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(input))

		if err == nil {
			con.State = def.CON_STATE_PLAYING

			LinkPlayerConnection(player, con)
			WriteToDesc(con, "Welcome back, "+player.Name+"!")
		} else {
			log.Println("Invalid password attempt: " + player.Name + " ip: " + con.Address)
			WriteToDesc(con, "Invalid password.")
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

			WriteToDesc(con, "You shall be called "+con.Name+", then...")
			WriteToDesc(con, "Passwords must be between 9 and 72 characters long, and contain at least 2 numbers/symbols.")
			WriteToDesc(con, "Password:")
			con.State = def.CON_STATE_NEW_PASSWORD
		} else {
			con.State = def.CON_STATE_NEW_LOGIN
			WriteToDesc(con, "What name would you like to go by then?")
		}

	} else if con.State == def.CON_STATE_NEW_PASSWORD {
		symbolCount := len(NonAlpha(input))
		inputLen := len(input)
		if inputLen >= 8 && inputLen <= 72 && symbolCount >= 2 {
			con.TempPass = input
			WriteToDesc(con, "Type again to confirm:")
			con.State = def.CON_STATE_NEW_PASSWORD_CONFIRM
		} else {
			WriteToDesc(con, "That isn't an acceptable password!")
			WriteToDesc(con, "Password:")
		}
	} else if con.State == def.CON_STATE_NEW_PASSWORD_CONFIRM {

		/*Hash password*/
		if input == con.TempPass {
			WriteToDesc(con, "Hashing password... One second please!")
			hash, err := bcrypt.GenerateFromPassword([]byte(input), def.PASSWORD_HASH_COST)
			if err != nil {
				CheckError("interp: password hash", err, def.ERROR_NONFATAL)
				WriteToDesc(con, "Password encryption failed, sorry something is wrong.")

				con.State = def.CON_STATE_DISCONNECTING
				return
			}
			con.TempPass = ""
			con.Player.Password = string(hash)

			SetupNewCharacter(con.Player)
			con.State = def.CON_STATE_PLAYING
			LinkPlayerConnection(con.Player, con)

			buf := fmt.Sprintf("%s has just arrived to the world.", con.Player.Name)
			WriteToRoom(con.Player, buf)

			okay := WritePlayer(con.Player)
			if okay == false {
				WriteToPlayer(con.Player, "Saving character failed!!!")
			} else {
				WriteToPlayer(con.Player, "Character saved.")
			}

			showHelp(con.Player)

		} else {
			con.TempPass = ""
			WriteToDesc(con, "Passwords didn't match, try again.")
			WriteToDesc(con, "Password:")
			con.State = def.CON_STATE_NEW_PASSWORD
		}

	}
}

func showHelp(player *glob.PlayerData) {
	WriteToPlayer(player, "Commands:")
	WriteToPlayer(player, "help       (you are here)")
	WriteToPlayer(player, "bytes      (bandwidth use)")
	WriteToPlayer(player, "look       (shows room)")
	WriteToPlayer(player, "who        (shows all online)")
	WriteToPlayer(player, "say <text> (Speaks)")
	WriteToPlayer(player, "save       (saves character)")
	WriteToPlayer(player, "quit       (leave game)")
	return
}

func PlayerCommand(player *glob.PlayerData, command string, args string) {
	/***************/
	/*Commands area*/ //TODO make into nice list with separate functions
	/***************/
	if player != nil && player.Valid {

		if command == "help" {
			showHelp(player)
			return
		} else if command == "quit" {
			okay := WritePlayer(player)
			if okay == false {
				WriteToPlayer(player, "Saving character failed!!!")
				return
			} else {
				WriteToPlayer(player, "Character saved.")
			}
			buf := fmt.Sprintf("%s has quit.", player.Name)
			WriteToAll(buf)
			player.Connection.State = def.CON_STATE_DISCONNECTING
			RemovePlayerWorld(player)
			return
		} else if command == "who" {
			output := "Players online:\n"

			for x := 0; x <= glob.ConnectionListEnd; x++ {
				var p *glob.ConnectionData = &glob.ConnectionList[x]
				if p.Valid == false {
					continue
				}
				buf := ""

				if p.State == def.CON_STATE_PLAYING {
					idleString := ""
					connectedString := ""

					if time.Since(p.IdleTime) > time.Minute {
						idleString = " (idle " + ToHourMinute(time.Since(p.IdleTime)) + ")"
					}
					if time.Since(p.ConnectedFor) > time.Minute {
						connectedString = " (on " + ToHourMinute(time.Since(p.ConnectedFor)) + ")"
					}

					buf = fmt.Sprintf("%d: %s%s%s", x, p.Name, connectedString, idleString)
				} else {
					buf = fmt.Sprintf("%d: %s", x, "(Connecting)")
				}
				output = output + buf
				if x <= glob.ConnectionListEnd {
					output = output + "\r\n"
				}
			}
			WriteToPlayer(player, output)
			return
		} else if command == "bytes" {
			output := ""

			for x := 0; x <= glob.ConnectionListEnd; x++ {
				var p *glob.ConnectionData = &glob.ConnectionList[x]
				if p.Valid == false {
					continue
				}
				buf := ""

				if p.Player != nil {
					output = "Connections:\r\nname: ip(count), in/out kb\r\n"
					for key, value := range player.Connections {
						buf = buf + fmt.Sprintf("%32v: %16v(%4v) %v/%v\r\n", player.Name, key, value, player.BytesIn[key]/1024, player.BytesOut[key]/1024)
					}
				} else {
					output = "Connections:\r\nname: in/out bytes\r\n"
					buf = buf + fmt.Sprintf("%32v %v/%v\r\n", p.Name, p.BytesIn, p.BytesOut)
				}

				output = output + buf
				if x <= glob.ConnectionListEnd {
					output = output + "\r\n"
				}
			}
			WriteToPlayer(player, output)
			return
		} else if command == "say" {
			if len(args) > 0 {
				out := fmt.Sprintf("%s says: %s", player.Name, args)
				us := fmt.Sprintf("You say: %s", args)

				WriteToOthers(player, out)
				WriteToPlayer(player, us)
			} else {
				WriteToPlayer(player, "But, what do you want to say?")
			}
			return
		} else if command == "save" {
			okay := WritePlayer(player)
			if okay == false {
				WriteToPlayer(player, "Saving character failed!!!")
			} else {
				WriteToPlayer(player, "Character saved.")
			}
			return
		} else if command == "asave" {
			okay := WriteSector(&glob.SectorsList[0])
			if okay == false {
				WriteToPlayer(player, "Saving sector failed!!!")
			} else {
				WriteToPlayer(player, "Sector saved.")
			}
			return

		} else if command == "look" {
			err := true
			if glob.SectorsList[player.Sector].Valid {
				sector := glob.SectorsList[player.Sector]
				if sector.Rooms[player.Room].Valid {
					room := sector.Rooms[player.Room]
					roomName := room.Name
					roomDesc := room.Description
					buf := fmt.Sprintf("%s:\r\n%s", roomName, roomDesc)
					WriteToPlayer(player, buf)
					err = false
				}

				if player.RoomLink != nil {
					names := ""
					unlinked := ""
					for _, target := range player.RoomLink.Players {
						if target != nil && target != player {
							if target.Connection != nil && target.Connection.Valid == false {
								unlinked = " (lost connection)"
							}
							names = names + fmt.Sprintf("%s is here.%s\r\n", target.Name, unlinked)
						}
					}
					//Newline if there are players here.
					if names != "" {
						names = "\r\n" + names
					}
					WriteToPlayer(player, names)
				}
			}
			if err {
				WriteToPlayer(player, "You are floating in the VOID...")
			}

			return
		} else {

			WriteToPlayer(player, "That isn't a valid command.")
			showHelp(player)
			return
		}

	}
}
