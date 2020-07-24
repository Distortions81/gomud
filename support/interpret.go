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

type Command struct {
	Name  string
	Cmd   func(player *glob.PlayerData, args string)
	Admin bool
	Help  string
}

const FOR_ADMIN = true
const FOR_USER = false

var CommandList = []Command{
	{Name: "bytes", Cmd: CmdBytes, Admin: FOR_ADMIN,
		Help: "See bandwidth usage"},
	{Name: "asave", Cmd: CmdAsave, Admin: FOR_ADMIN,
		Help: "Save game areas"},

	{Name: "help", Cmd: CmdHelp, Admin: FOR_USER,
		Help: "You are here"},
	{Name: "who", Cmd: CmdWho, Admin: FOR_USER,
		Help: "See who is online"},
	{Name: "look", Cmd: CmdLook, Admin: FOR_USER,
		Help: "Look around the room"},
	{Name: "say", Cmd: CmdSay, Admin: FOR_USER,
		Help: "Talk to other people in the room"},
	{Name: "save", Cmd: CmdSave, Admin: FOR_USER,
		Help: "Talk to other people in the room"},
	{Name: "quit", Cmd: CmdQuit, Admin: FOR_USER,
		Help: "Quit the game"},
}

func MakeQuickHelp() {
	buf := "Commands:\r\n"

	for _, cmd := range CommandList {
		admin := ""
		if cmd.Admin {
			admin = " (admin)"
		}
		buf = buf + fmt.Sprintf("%12v : %-56v%8v\r\n", cmd.Name, cmd.Help, admin)
	}
	glob.QuickHelp = buf
}

func PlayerCommand(player *glob.PlayerData, command string, args string) {

	if player != nil && player.Valid {
		for _, cmd := range CommandList {
			if strings.ToLower(cmd.Name) == strings.ToLower(command) {
				cmd.Cmd(player, args)
				return
			}
		}
		WriteToPlayer(player, "Invalid command.")
		CmdHelp(player, "")
	}
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

	/*Inital connection*/
	if con.State == def.CON_STATE_WELCOME {
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
		/* Player's password */
	} else if con.State == def.CON_STATE_PASSWORD {
		player, _ := ReadPlayer(con.Name, true)
		con.Player = player

		err := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(input))

		if err == nil {
			con.State = def.CON_STATE_PLAYING

			WriteToDesc(con, "Welcome back, "+player.Name+"!")
			LinkPlayerConnection(player, con)
		} else {
			log.Println("Invalid password attempt: " + player.Name + " ip: " + con.Address)
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

			con.TempPass = ""
			con.Player.Password = string(hash)

			SetupNewCharacter(con.Player)
			con.State = def.CON_STATE_PLAYING
			LinkPlayerConnection(con.Player, con)

			okay := WritePlayer(con.Player)
			if okay == false {
				WriteToPlayer(con.Player, "Saving character failed!!!")
			} else {
				WriteToPlayer(con.Player, "Character saved.")
			}

			CmdHelp(con.Player, "")

		} else {
			con.TempPass = ""
			WriteToDesc(con, "Passwords didn't match, try again.")
			WriteToDesc(con, "Password:")
			con.State = def.CON_STATE_NEW_PASSWORD
		}

	}
}

func CmdHelp(player *glob.PlayerData, args string) {
	WriteToPlayer(player, glob.QuickHelp)
}

func CmdQuit(player *glob.PlayerData, args string) {
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
}

func CmdWho(player *glob.PlayerData, args string) {
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
}

func CmdBytes(player *glob.PlayerData, args string) {
	output := ""

	for x := 1; x <= glob.ConnectionListEnd; x++ {
		con := &glob.ConnectionList[x]
		target := con.Player
		buf := ""
		ssl := ""

		if con.SSL {
			ssl = "(SSL)"
		}

		if target != nil {
			for key, value := range target.Connections {
				buf = buf + fmt.Sprintf("%-5s%32v: %16v(%4v) %v/%v\r\n", ssl, target.Name, key, value, ScaleBytes(target.BytesIn[key]), ScaleBytes(target.BytesOut[key]))
			}
		} else if con != nil {
			buf = buf + fmt.Sprintf("%-5s%32v: %16v(%4v) %v/%v\r\n", ssl, con.Name, "", "", ScaleBytes(con.BytesIn), ScaleBytes(con.BytesOut))
		}

		output = output + buf
	}
	WriteToPlayer(player, output)
}

func CmdSay(player *glob.PlayerData, args string) {
	if len(args) > 0 {
		out := fmt.Sprintf("%s says: %s", player.Name, args)
		us := fmt.Sprintf("You say: %s", args)

		WriteToOthers(player, out)
		WriteToPlayer(player, us)
	} else {
		WriteToPlayer(player, "But, what do you want to say?")
	}
}

func CmdSave(player *glob.PlayerData, args string) {
	okay := WritePlayer(player)
	if okay == false {
		WriteToPlayer(player, "Saving character failed!!!")
	} else {
		WriteToPlayer(player, "Character saved.")
	}
}

func CmdAsave(player *glob.PlayerData, args string) {
	okay := WriteSector(&glob.SectorsList[0])
	if okay == false {
		WriteToPlayer(player, "Saving sector failed!!!")
	} else {
		WriteToPlayer(player, "Sector saved.")
	}
}

func CmdLook(player *glob.PlayerData, args string) {

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

}
