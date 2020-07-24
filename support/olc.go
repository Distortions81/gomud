package support

import (
	"fmt"
	"strings"

	"../def"
	"../glob"
)

func CmdOLC(player *glob.PlayerData, input string) {

	inputLen := len(input)
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

	if command == "exit" {
		player.OLCEdit.Mode = def.OLC_NONE
		WriteToPlayer(player, "Exiting OLC.")
		player.OLCEdit.Active = false
		return
	}
	player.OLCEdit.Active = true

	if player.OLCEdit.Mode == def.OLC_NONE {
		if command == "" {
			WriteToPlayer(player, "Possible types:")
			WriteToPlayer(player, "Exit, room, object, trigger, mobile or quest.")
			WriteToPlayer(player, "OLC <type>, or OLC <type> <sector:id> (for a specfic item)")
			WriteToPlayer(player, "sector can be omitted if the item is in the same sector as your character")
		} else {
			if command == "room" {
				player.OLCEdit.Mode = def.OLC_ROOM
				if longArg == "" {
					player.OLCEdit.Room = player.Location
					player.OLCEdit.Room.Valid = true
				}
				WriteToPlayer(player, "Room edit mode.\r\n")
			} else if command == "object" {
				WriteToPlayer(player, "Object edit mode.\r\n")
				player.OLCEdit.Mode = def.OLC_OBJECT
			} else if command == "trigger" {
				WriteToPlayer(player, "Trigger edit mode.\r\n")
				player.OLCEdit.Mode = def.OLC_TRIGGER
			} else if command == "mobile" {
				WriteToPlayer(player, "Mobile edit mode.\r\n")
				player.OLCEdit.Mode = def.OLC_MOBILE
			} else if command == "quest" {
				WriteToPlayer(player, "Quest edit mode.\r\n")
				player.OLCEdit.Mode = def.OLC_QUEST
			} else {
				WriteToPlayer(player, "Edit "+command+"? What is that?")
			}
			CmdOLC(player, "")

		}

	} else if player.OLCEdit.Mode == def.OLC_ROOM {
		if player.OLCEdit.Room.Valid {
			buf := ""
			exits := ""
			for name, exit := range player.OLCEdit.Room.RoomLink.Exits {
				exits = exits + fmt.Sprintf("Name: %v, ToRoom: %v:%v\r\nDoor: %v, AutoOpen: %v, AutoClose: %v\r\nHidden: %v, Keyed: %v",
					name, exit.ToRoom.Sector, exit.ToRoom.ID, exit.Door.Door, exit.Door.AutoOpen, exit.Door.AutoClose,
					exit.Door.Hidden, exit.Door.Keyed)
			}
			if exits == "" {
				exits = "None"
			}
			buf = buf + fmt.Sprintf("Sector: %v, ID %v\n\rName: %v\n\rDescription: \n\r\r\n%v\r\n\r\nExits: %v",
				player.OLCEdit.Room.Sector, player.OLCEdit.Room.ID,
				player.OLCEdit.Room.RoomLink.Name, player.OLCEdit.Room.RoomLink.Description, exits)
			WriteToPlayer(player, buf)
		} else {
			WriteToPlayer(player, "No room selected, olc <sector:id> (sector optional)")
		}

	} else if player.OLCEdit.Mode == def.OLC_OBJECT {

	} else if player.OLCEdit.Mode == def.OLC_TRIGGER {

	} else if player.OLCEdit.Mode == def.OLC_MOBILE {

	} else if player.OLCEdit.Mode == def.OLC_QUEST {

	}

}
