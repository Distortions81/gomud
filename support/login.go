package support

import (
	"fmt"
	"log"
	"strings"
	"time"

	"../def"
	"../glob"
	"../mlog"
	"golang.org/x/crypto/bcrypt"
)

func interpretInput(con *glob.ConnectionData, input string, isAlias bool) {

	if con == nil {
		log.Println("interpretInput: nil connection")
		return
	}
	//May have to bypass this for "force" command
	if con.Valid == false {
		return
	}

	/*********************/
	/*Clean up user input*/
	/*********************/
	overflow := false

	input = strings.TrimSpace(input)
	input, overflow = TruncateString(input, def.MAX_INPUT_LENGTH)

	if overflow {
		WriteToDesc(con, "Your input was too long, rejecting.")
		con.State = def.CON_STATE_DISCONNECTING
		return
	}

	/*Split up arguments*/
	command, argTwoThrough := SplitArgsTwo(input, " ")
	cmdB, argThreeThrough := SplitArgsTwo(argTwoThrough, " ")
	cmdl := strings.ToLower(command)
	//cmdBl := strings.ToLower(cmdB)

	//Set player as no longer idle
	con.IdleTime = time.Now()

	if con.State == def.CON_STATE_DISCONNECTED || con.State == def.CON_STATE_DISCONNECTING {
		/*Connections marked disconnect can't input*/

		CloseConnection(con)
		return
	} else if con.State == def.CON_STATE_RELOG {
		/*Player relog*/

		con.Player = nil
		con.TempPlayer = nil
		con.TempPlayer = nil
		con.State = def.CON_STATE_WELCOME
	} else if con.State == def.CON_STATE_NEWS {
		con.State = def.CON_STATE_PLAYING
		LinkPlayerConnection(con.Player, con)
		if con.Player.Level == 0 {
			WriteToPlayer(con.Player, "NEW PLAYER HELP:")
			CmdHelp(con.Player, "")
			con.Player.Level = 1
		}
		return
	} else if con.State == def.CON_STATE_PLAYING && con.Player != nil && con.Player.Valid {
		/*If we are playing the game, this is a command */

		if con.Player.OLCSettings.NoOLCPrefix == true && con.Player.OLCEdit.Active {
			if cmdl == "cmd" {
				if strings.EqualFold(cmdB, "olc") == false {
					PlayerCommand(con.Player, cmdB, argThreeThrough, false)
					return
				} else {
					WriteToPlayer(con.Player, "DOES NOT COMPUTE, DOES NOT COMPUTE! <EXPLODES>...\r\nFor real though, you can't do that.")
					return
				}
			} else {
				CmdOLC(con.Player, input)
				return
			}
		} else if command == "" && con.Player.OLCEdit.Active {
			CmdOLC(con.Player, "")
			return
		}

		PlayerCommand(con.Player, command, argTwoThrough, isAlias)
		if con.Player.OLCEdit.Active && con.Player.OLCSettings.OLCPrompt {
			OLCPrompt := ""
			if con.Player.OLCEdit.Mode == def.OLC_NONE {
				OLCPrompt = "<OLC: none>:"
			} else if con.Player.OLCEdit.Mode == def.OLC_ROOM {
				OLCPrompt = fmt.Sprintf("<OLC ROOM: %v:%v>:", con.Player.OLCEdit.Room.Sector, con.Player.OLCEdit.Room.ID)
			} else if con.Player.OLCEdit.Mode == def.OLC_OBJECT {
				OLCPrompt = fmt.Sprintf("<OLC OBJECT: WIP>:")
			} else if con.Player.OLCEdit.Mode == def.OLC_TRIGGER {
				OLCPrompt = fmt.Sprintf("<OLC TRIGGER: WIP>:")
			} else if con.Player.OLCEdit.Mode == def.OLC_MOBILE {
				OLCPrompt = fmt.Sprintf("<OLC MOBILE: WIP>: ")
			} else if con.Player.OLCEdit.Mode == def.OLC_QUEST {
				OLCPrompt = fmt.Sprintf("<OLC QUEST: WIP>: ")
			} else if con.Player.OLCEdit.Mode == def.OLC_SECTOR {
				OLCPrompt = fmt.Sprintf("<OLC SECTOR: WIP>: ")
			} else if con.Player.OLCEdit.Mode == def.OLC_EXITS {
				OLCPrompt = fmt.Sprintf("<OLC EXITS: [%v] Room: %v:%v>",
					con.Player.OLCEdit.ExitName,
					con.Player.OLCEdit.Room.Sector,
					con.Player.OLCEdit.Room.ID)
			}
			defer WriteToDesc(con, OLCPrompt)
			return
		} else {
			defer WriteToDesc(con, ">")
			return
		}
	}

	//Aliases can only be used for commands
	if isAlias {
		return
	}

	//For names, only letters allowed
	alphaChar := AlphaOnly(input)
	alphaCharLen := len(alphaChar)

	/*Inital connection*/
	if con.State == def.CON_STATE_PLAYING {
		/*Players shouldn't be here*/
		return
	} else if con.State == def.CON_STATE_WELCOME {
		if command == "new" {
			buf := fmt.Sprintf("Names must be between %d and %d letters long, A-z only.", def.MIN_PLAYER_NAME_LENGTH, def.MAX_PLAYER_NAME_LENGTH)
			WriteToDesc(con, buf)
			WriteToDesc(con, "What name would you like to go by?")
			con.State = def.CON_STATE_NEW_LOGIN
		} else {
			/*Login Name */
			if alphaCharLen > def.MIN_PLAYER_NAME_LENGTH &&
				alphaCharLen < def.MAX_PLAYER_NAME_LENGTH &&
				alphaChar != def.STRING_UNKNOWN {

				for x := 1; x <= glob.PlayerListEnd; x++ {
					target := glob.PlayerList[x]

					if target != nil && target.Valid &&
						target.Connection != nil &&
						target.Connection.Valid &&
						target.Connection.Desc != nil &&
						target.Connection.State == def.CON_STATE_PLAYING {

						if strings.EqualFold(target.Name, alphaChar) {
							WriteToDesc(con, "That character is already online!")
							WriteToDesc(con, "Login anyway? (y/n)")
							if target != nil {
								WriteToPlayer(target, "Someone is attempting to login to this character.")
							}
							con.TempPlayer = target
							con.Name = alphaChar
							con.State = def.CON_STATE_RECONNECT_CONFIRM
							return
						}
					}
				}
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
				WriteToDesc(con, "...That isn't a valid name, try again.")
				WriteToDesc(con, "Name:")
			}
		}
	} else if con.State == def.CON_STATE_RECONNECT_CONFIRM {
		if command == "y" || command == "yes" {
			WriteToDesc(con, "Password:")
			con.Name = con.TempPlayer.Name
			con.State = def.CON_STATE_PASSWORD
		} else {
			con.State = def.CON_STATE_DISCONNECTING
		}

		/* Player's password */
	} else if con.State == def.CON_STATE_PASSWORD {
		player, _ := ReadPlayer(con.Name, true)
		con.Player = player

		err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(input))

		if err == nil {
			con.State = def.CON_STATE_NEWS

			if con.TempPlayer != nil && con.TempPlayer.Connection.Valid {

				WriteToDesc(con.TempPlayer.Connection, "You logged in from another connection!")
				CloseConnection(con.TempPlayer.Connection)
				WriteToDesc(con, "Closing other connection to character...")
				WriteToDesc(con, "")
			}
			WriteToDesc(con, glob.News)
			WriteToDesc(con, "{K[Press enter or return to continue]")
		} else {

			mlog.Write("Invalid password attempt: " + player.Name + " ip: " + con.Address)
			WriteToDesc(con, "Invalid password.")
			con.State = def.CON_STATE_DISCONNECTING
		}
		/*New player*/
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
		/*Confirm new player*/
	} else if con.State == def.CON_STATE_NEW_LOGIN_CONFIRM {
		if command == "y" || command == "yes" {
			con.Player = CreatePlayerFromDesc(con)

			WriteToDesc(con, "You shall be called "+con.Name+", then...")
			WriteToDesc(con, "Passwords must be between 9 and 72 characters long,\r\nand contain at least 2 numbers/symbols.")
			WriteToDesc(con, "Password:")
			con.State = def.CON_STATE_NEW_PASSWORD
		} else {
			con.State = def.CON_STATE_NEW_LOGIN
			WriteToDesc(con, "What name would you like to go by then?")
		}

		/*New player password*/
	} else if con.State == def.CON_STATE_NEW_PASSWORD {
		symbOLCount := len(NonAlpha(input))
		inputLen := len(input)
		if inputLen >= 8 && inputLen <= 72 && symbOLCount >= 2 {
			con.TempPass = input
			WriteToDesc(con, "Type again to confirm:")
			con.State = def.CON_STATE_NEW_PASSWORD_CONFIRM
		} else {
			WriteToDesc(con, "That isn't an acceptable password!")
			WriteToDesc(con, "Password:")
		}
		/*Confirm new player password */
	} else if con.State == def.CON_STATE_NEW_PASSWORD_CONFIRM {

		/*Hash password*/
		if input == con.TempPass {
			WriteToDesc(con, "Hashing password...")
			hash, err := bcrypt.GenerateFromPassword([]byte(input), def.PASSWORD_HASH_COST)
			if err != nil {
				CheckError("interp: password hash", err, def.ERROR_NONFATAL)
				WriteToDesc(con, "Password encryption failed, sorry something is wrong.")

				con.State = def.CON_STATE_DISCONNECTING
				return
			}
			WriteToDesc(con, "...done! Welcome to GoMud!")
			WriteToDesc(con, "")

			con.TempPass = ""
			con.Player.Password = string(hash)

			SetupNewCharacter(con.Player)
			con.State = def.CON_STATE_NEWS

			okay := WritePlayer(con.Player)
			if okay == false {
				WriteToPlayer(con.Player, "Saving character failed!!!")
			} else {
				WriteToPlayer(con.Player, "Character saved.")
			}

			WriteToDesc(con, glob.News)
			WriteToDesc(con, "{K[Press enter or return to continue]")
		} else {
			con.TempPass = ""
			WriteToDesc(con, "Passwords didn't match, try again.")
			WriteToDesc(con, "Password:")
			con.State = def.CON_STATE_NEW_PASSWORD
		}

	} else {
		/* Player in weird mode? */
		WriteToDesc(con, "Your connection is in an unknown mode... Please reconnect.")
		CloseConnection(con)
		return
	}
}
