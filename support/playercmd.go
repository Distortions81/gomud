package support

import (
	"fmt"
	"strings"
	"time"

	"../def"
	"../glob"
)

func CmdNews(player *glob.PlayerData, input string) {
	WriteToPlayer(player, glob.News)
}

//Hard coded aliases
func CmdNorth(player *glob.PlayerData, input string) {
	CmdGo(player, "north")
}
func CmdSouth(player *glob.PlayerData, input string) {
	CmdGo(player, "south")
}
func CmdEast(player *glob.PlayerData, input string) {
	CmdGo(player, "east")
}
func CmdWest(player *glob.PlayerData, input string) {
	CmdGo(player, "west")
}
func CmdUp(player *glob.PlayerData, input string) {
	CmdGo(player, "up")
}
func CmdDown(player *glob.PlayerData, input string) {
	CmdGo(player, "down")
}

func CmdRelog(player *glob.PlayerData, input string) {
	c := player.Connection
	CmdSave(player, "")
	RemovePlayer(player)
	c.Input.BufferInPos = 0
	c.Input.BufferInCount = 0
	c.Input.BufferOutPos = 0
	c.Input.BufferOutCount = 0

	for p := 0; p < def.MAX_INPUT_LINES; p++ {
		c.Input.InputBuffer[p] = ""
	}

	c.State = def.CON_STATE_RELOG
	c.Valid = true

	WriteToDesc(c, "")
	WriteToDesc(c, "")
	WriteToDesc(c, "")
	WriteToDesc(c, glob.Greeting)

	WriteToDesc(c, "(Type NEW to create character) Name:")
	return
}

func CmdRecall(player *glob.PlayerData, input string) {

	if input == "set" {
		//I love how elegant this is
		player.Recall = player.Location
		WriteToPlayer(player, "Recall set!")
		return
	}

	if player.Location.ID == player.Recall.ID && player.Location.Sector == player.Recall.Sector {
		WriteToPlayer(player, "You try to recall, but strain to remember... wait, does this place look familiar?")
		return
	} else {
		WriteToPlayer(player, "You recall, and are suddenly transported, in a bright blue {Cflash{x of {Ylight.")
	}

	WriteToRoom(player, fmt.Sprintf("%v {Kvanishes{x with a bright blue {Cflash{x of {Ylight.", player.Name))
	if player.Recall.Sector != 0 || player.Recall.ID != 0 {
		PlayerToRoom(player, player.Recall.Sector, player.Recall.ID)
	} else {
		PlayerToRoom(player, def.PLAYER_START_SECTOR, def.PLAYER_START_ROOM)
	}
	WriteToRoom(player, fmt.Sprintf("%v suddenly {mappears{x, with a bright blue {Cflash{x of {Ylight.", player.Name))
	player.Dirty = true
	CmdLook(player, "")
}

func CmdAlias(player *glob.PlayerData, input string) {
	command, longArgs := SplitArgsTwo(input, " ")
	firstArg, lastArgs := SplitArgsTwo(longArgs, " ")

	if command == "list" {
		aliases := ""
		for key, value := range player.Aliases {
			aliases = aliases + fmt.Sprintf("%s: %s\r\n", key, value)
		}
		if aliases == "" {
			aliases = "None"
		}
		WriteToPlayer(player, "Aliases:\r\n"+aliases)
	} else if command == "add" {

		/*Prevent problems*/
		if firstArg == "" {
			WriteToPlayer(player, "The alias needs a name.")
			return
		} else if firstArg == "alias" {
			WriteToPlayer(player, "That would be a bad idea.")
			return
		}
		if len(lastArgs) > 80 {
			WriteToPlayer(player, "That output is too long.")
			return
		}

		/*Write data to player*/
		player.Aliases[firstArg] = lastArgs
		WriteToPlayer(player, "Alias added.")
		player.Dirty = true
	} else if command == "delete" {
		found := false
		for key, _ := range player.Aliases {
			if strings.EqualFold(key, firstArg) {
				found = true
				break
			}
		}
		if found {
			delete(player.Aliases, firstArg)
			WriteToPlayer(player, "Alias deleted")
			player.Dirty = true
		}
	} else {
		WriteToPlayer(player, "Aliases can be the same name as commands,")
		WriteToPlayer(player, "but you can still alias the orignal command to something else.")
		WriteToPlayer(player, "Aliases can not call other aliases, or the alias command,")
		WriteToPlayer(player, "and must use full-length command names (no shorthand).")
		WriteToPlayer(player, "")
		WriteToPlayer(player, "alias add <shortcut> <output> (max 80 characters)")
		WriteToPlayer(player, "alias delete <shortcut>")
		WriteToPlayer(player, "alias list")
	}
}

func CmdCommands(player *glob.PlayerData, args string) {
	WriteToPlayer(player, glob.QuickHelp)
}

func movePlayerExit(player *glob.PlayerData, arg string, exit *glob.ExitData) {
	WriteToPlayer(player, "You go "+arg+".")
	WriteToRoom(player, player.Name+" goes "+arg+".")
	PlayerToRoom(player, exit.ToRoom.Sector, exit.ToRoom.ID)

	WriteToRoom(player, player.Name+" arrives.")
	player.Dirty = true

	CmdLook(player, "")
}

func CmdGo(player *glob.PlayerData, args string) {
	found := false

	if args == "" {
		WriteToPlayer(player, "Go where?")
		return
	}
	for exitName, exit := range player.Location.RoomLink.Exits {
		if strings.HasPrefix(strings.ToLower(exitName), strings.ToLower(args)) {
			found = true

			if IsStandardDirection(exitName) {
				WriteToPlayer(player, "You go "+exitName+".")
				WriteToRoom(player, player.Name+" goes "+exitName+".")
			} else {
				WriteToPlayer(player, "You '"+exitName+"'.")
				WriteToRoom(player, player.Name+" went '"+exitName+"'.")
			}

			PlayerToRoom(player, exit.ToRoom.Sector, exit.ToRoom.ID)

			if IsStandardDirection(exitName) {
				WriteToRoom(player, player.Name+" arrives from '"+GetStandardDirectionMirror(exitName)+"'.")
			} else {
				WriteToRoom(player, player.Name+" went '"+exitName+"'.")
			}

			CmdLook(player, "")
			player.Dirty = true
			return
		}
	}
	if !found {

		var exitsList []string
		for exitName, _ := range player.Location.RoomLink.Exits {
			exitsList = append(exitsList, exitName)
		}
		result, _ := FindClosestMatch(exitsList, args)
		if result != "" {

		} else {
			WriteToPlayer(player, "Go where?")
		}
	}
}

func CmdQuit(player *glob.PlayerData, args string) {
	okay := WritePlayer(player)
	if okay == false {
		WriteToPlayer(player, "Saving character failed!!!")
		return //Don't quit if we couldn't save
	} else {
		WriteToPlayer(player, "Character saved.")
	}
	buf := fmt.Sprintf("%s has quit.", player.Name)
	WriteToAll(buf)
	player.Connection.State = def.CON_STATE_DISCONNECTING
}

func CmdWho(player *glob.PlayerData, args string) {
	buf := "Players online:"
	WriteToPlayer(player, buf)

	pos := 0
	for x := 1; x <= glob.ConnectionListEnd; x++ {
		var p *glob.ConnectionData = &glob.ConnectionList[x]
		if p.Valid == false {
			continue
		}

		if p.State == def.CON_STATE_PLAYING {
			idleString := ""
			connectedString := ""

			if time.Since(p.IdleTime) > time.Minute {
				idleString = " (idle " + RoundSinceTime("m", p.IdleTime) + ")"
			}
			if time.Since(p.ConnectedFor) > time.Minute {
				connectedString = " (on " + RoundSinceTime("m", p.ConnectedFor) + ")"
			}
			pos++
			buf = fmt.Sprintf("%d: %s%s%s", pos, p.Name, connectedString, idleString)
			WriteToPlayer(player, buf)
		} else {
			pos++
			buf = fmt.Sprintf("%d: %s", pos, "(Connecting)")
			WriteToPlayer(player, buf)
		}
	}
}

func CmdSay(player *glob.PlayerData, args string) {
	if len(args) > 0 {
		out := fmt.Sprintf("%s says: %s", player.Name, args)
		us := fmt.Sprintf("You say: %s", args)

		WriteToRoom(player, out)
		WriteToPlayer(player, us)
	} else {
		WriteToPlayer(player, "But, what do you want to say?")
	}
}

func CmdEmote(player *glob.PlayerData, args string) {
	if len(args) > 0 {
		out := fmt.Sprintf("%s %s", player.Name, args)

		WriteToRoom(player, out)
		WriteToPlayer(player, out)
	} else {
		WriteToPlayer(player, "But, what do you want to say?")
	}
}

func CmdChat(player *glob.PlayerData, args string) {
	if len(args) > 0 {
		out := fmt.Sprintf("%s chats: %s", player.Name, args)
		us := fmt.Sprintf("You chat: %s", args)

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

func CmdLook(player *glob.PlayerData, args string) {

	err := true
	sector := glob.SectorsList[player.Location.Sector]
	if sector.Valid {
		if sector.Rooms[player.Location.ID] != nil && sector.Rooms[player.Location.ID].Valid {
			room := sector.Rooms[player.Location.ID]
			roomName := room.Name
			roomDesc := room.Description
			buf := fmt.Sprintf("%s:\r\n%s", roomName, roomDesc)
			WriteToPlayer(player, buf)
			err = false
		}

		if player.Location.RoomLink != nil {
			exits := "["
			l := len(player.Location.RoomLink.Exits)
			x := 0
			for name, _ := range player.Location.RoomLink.Exits {
				x++
				exits = exits + name
				if x < l {
					exits = exits + ", "
				}
			}
			if exits == "[" {
				exits = exits + " None... "
			}
			exits = exits + "]"

			WriteToPlayer(player, "\r\nExits: "+exits)

			names := ""
			unlinked := ""
			for _, target := range player.Location.RoomLink.Players {
				if target != nil && target != player {
					if target.Connection != nil && target.Connection.Valid == false {
						unlinked = " (lost connection)"
					}
					names = names + fmt.Sprintf("%s is here.%s", target.Name, unlinked)
				}
			}
			//Newline if there are players here.
			if names != "" {
				WriteToPlayer(player, "\r\n"+names)
			}
		} else {
			err = true
		}
	}
	if err {
		WriteToPlayer(player, "You are floating in the VOID...")
		PlayerToRoom(player, def.PLAYER_START_SECTOR, def.PLAYER_START_ROOM)
	}

}
