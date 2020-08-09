package support

import (
	"fmt"
	"strings"

	"../def"
	"../glob"
)

func OLCExits(player *glob.PlayerData,
	input string, command string, cmdB string, cmdl string, cmdBl string,
	argTwoThrough string, argThreeThrough string) {

	if cmdl == "done" {
		player.OLCEdit.Mode = def.OLC_ROOM
		WriteToPlayer(player, "Going back to room editor...")
		CmdOLC(player, "")
		return
	} else if cmdl == "name" {
		if cmdB != "" {
			for exitName, _ := range player.OLCEdit.Room.RoomLink.Exits {
				if strings.EqualFold(exitName, cmdB) {
					WriteToPlayer(player, "That exit already exists.")
					return
				}
			}

			//Delete old
			exitName := player.OLCEdit.Exit.Name
			exitDoor := player.OLCEdit.Exit.ExitLink.Door
			exitToRoom := player.OLCEdit.Exit.ExitLink.ToRoom

			delete(player.OLCEdit.Room.RoomLink.Exits, exitName)
			if player.OLCEdit.Room.RoomLink.Exits == nil {
				player.OLCEdit.Room.RoomLink.Exits = make(map[string]*glob.ExitData)
			}

			//Make new
			if player.OLCEdit.Room.RoomLink.Exits == nil {
				player.OLCEdit.Room.RoomLink.Exits = make(map[string]*glob.ExitData)
			}
			player.OLCEdit.Room.RoomLink.Exits[argTwoThrough] = CreateExit()

			//Copy data over
			player.OLCEdit.Room.RoomLink.Exits[argTwoThrough].Door = exitDoor
			player.OLCEdit.Room.RoomLink.Exits[argTwoThrough].ToRoom = exitToRoom

			player.OLCEdit.Exit.ExitLink = player.OLCEdit.Room.RoomLink.Exits[argTwoThrough]
			player.OLCEdit.Exit.Name = argTwoThrough
			player.OLCEdit.Mode = def.OLC_EXITS
			glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave
			CmdOLC(player, "")
			return
		} else {
			WriteToPlayer(player, "OLC exit <exit name>")
		}

	} else if cmdl == "door" {
		if player.OLCEdit.Exit.ExitLink.Door.Door {
			player.OLCEdit.Exit.ExitLink.Door.Door = false
		} else {
			player.OLCEdit.Exit.ExitLink.Door.Door = true
		}
	} else if cmdl == "autoopen" {
		if player.OLCEdit.Exit.ExitLink.Door.AutoOpen {
			player.OLCEdit.Exit.ExitLink.Door.AutoOpen = false
		} else {
			player.OLCEdit.Exit.ExitLink.Door.AutoOpen = true
		}
	} else if cmdl == "autoclose" {
		if player.OLCEdit.Exit.ExitLink.Door.AutoClose {
			player.OLCEdit.Exit.ExitLink.Door.AutoClose = false
		} else {
			player.OLCEdit.Exit.ExitLink.Door.AutoClose = true
		}
	} else if cmdl == "keyed" {
		if player.OLCEdit.Exit.ExitLink.Door.AutoClose {
			player.OLCEdit.Exit.ExitLink.Door.AutoClose = false
		} else {
			player.OLCEdit.Exit.ExitLink.Door.AutoClose = true
		}
	} else if cmdl == "delete" {
		exitName := player.OLCEdit.Exit.Name
		delete(player.OLCEdit.Room.RoomLink.Exits, exitName)
		player.OLCEdit.Mode = def.OLC_ROOM
		CmdOLC(player, "")
		WriteToPlayer(player, "Exit deleted, returning to room editor.")
		return
	} else if cmdl == "toroom" {
		loc, found := getLocationFromString(player, cmdB)
		if found == false {
			WriteToPlayer(player, "Invalid location. <sector:id>, or <id> for sector you are standing in.")
			return
		}
		player.OLCEdit.Exit.ExitLink.ToRoom = loc
	}
	buf := fmt.Sprintf("OLC EDIT EXITS:\r\n%10v: %v\r\n\r\n%10v: %v:%v\r\n%10v: %v:%v\r\n",
		"Name",
		player.OLCEdit.Exit.Name,
		"FromRoom",
		player.OLCEdit.Room.Sector,
		player.OLCEdit.Room.ID,
		"ToRoom",
		player.OLCEdit.Exit.ExitLink.ToRoom.Sector,
		player.OLCEdit.Exit.ExitLink.ToRoom.ID)
	WriteToBuilder(player, buf)
	buf = fmt.Sprintf("%10v: %v\r\n%10v: %v\r%10v: %v\r\n%10v: %v\r%10v: %v",
		"Door",
		boolToYesNo(player.OLCEdit.Exit.ExitLink.Door.Door),
		"AutoOpen",
		boolToYesNo(player.OLCEdit.Exit.ExitLink.Door.AutoOpen),
		"AutoClose",
		boolToYesNo(player.OLCEdit.Exit.ExitLink.Door.AutoClose),
		"Hidden",
		boolToYesNo(player.OLCEdit.Exit.ExitLink.Door.Hidden),
		"Keyed",
		boolToYesNo(player.OLCEdit.Exit.ExitLink.Door.Keyed))
	WriteToPlayer(player, buf)
	WriteToPlayer(player, "Syntax for OLC exits: olc ToRoom <location>, door, autoOpen, autoClose, keyed, delete, done")

}
