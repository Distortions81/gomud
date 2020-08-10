package support

import (
	"fmt"
	"strconv"
	"strings"

	"../def"
	"../glob"
)

func OLCRoom(player *glob.PlayerData,
	input string, command string, cmdB string, cmdl string, cmdBl string,
	argTwoThrough string, argThreeThrough string) {

	if cmdl == "room" {
		loc := strings.Split(argTwoThrough, ":")
		locLen := len(loc)

		//TODO replace with existing function
		if locLen >= 1 {
			sector := 0
			id := 0
			var erra error
			var errb error
			if locLen == 1 {
				sector = player.Location.Sector
				id, erra = strconv.Atoi((loc[0]))
			} else if locLen == 2 {
				sector, errb = strconv.Atoi(loc[0])
				id, erra = strconv.Atoi((loc[1]))
			}

			if erra == nil && errb == nil {
				editRoom, roomFound := LocationDataFromID(sector, id)
				if roomFound {
					player.OLCEdit.Room.ID = editRoom.ID
					player.OLCEdit.Room.Sector = editRoom.Sector
					player.OLCEdit.Room.RoomLink = editRoom.RoomLink
					CmdOLC(player, "")
				}
			}
		} else if cmdl == "create" {
			loc := strings.Split(argTwoThrough, ":")
			locLen := len(loc)

			if locLen >= 1 {
				sector := 0
				id := 0
				var erra error
				var errb error
				if locLen == 1 {
					sector = player.Location.Sector
					id, erra = strconv.Atoi((loc[0]))
				} else if locLen == 2 {
					sector, errb = strconv.Atoi(loc[0])
					id, erra = strconv.Atoi((loc[1]))
				}

				if erra != nil || errb != nil {
					WriteToPlayer(player, "Syntax: OLC room create <location>\r\nExample: OLC room create 1:1")
					return
				}

				editRoom, roomFound := LocationDataFromID(sector, id)
				if roomFound {
					player.OLCEdit.Room.ID = editRoom.ID
					player.OLCEdit.Room.Sector = editRoom.Sector
					player.OLCEdit.Room.RoomLink = editRoom.RoomLink
					CmdOLC(player, "")
					WriteToPlayer(player, "Room already exists.")
					return
				} else {
					glob.SectorsList[sector].Rooms[id] = CreateRoom()
					editRoom, _ := LocationDataFromID(sector, id)
					player.OLCEdit.Room.ID = editRoom.ID
					player.OLCEdit.Room.Sector = editRoom.Sector
					player.OLCEdit.Room.RoomLink = editRoom.RoomLink
					CmdOLC(player, "")
					WriteToPlayer(player, fmt.Sprintf("Room %v:%v created!", sector, id))
					glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave
					return
				}
			} else {
				WriteToPlayer(player, "Syntax: OLC room create <location>\r\nExample: OLC room create 1:1")
				return
			}
		}
	} else if cmdl == "done" {
		player.OLCEdit.Mode = def.OLC_NONE
		WriteToPlayer(player, "Exiting OLC.")
		player.OLCEdit.Active = false
		return
	}
	if argTwoThrough != "" {
		if cmdl == "name" {
			player.OLCEdit.Room.RoomLink.Name = argTwoThrough
			WriteToPlayer(player, "Name set")
			glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave

		} else if cmdl == "description" || cmdl == "desc" {
			if cmdBl == "editor" {
				player.CurEdit.Active = true
				player.CurEdit.CallBack = "olc room"
				player.CurEdit.CallBackP = &player.OLCEdit.Room.RoomLink.Description

				dLines := strings.Split(player.OLCEdit.Room.RoomLink.Description, "\r\n")
				dLen := len(dLines)
				player.CurEdit.NumLines = 0
				player.CurEdit.CurLine = 0
				if player.CurEdit.Lines == nil {
					player.CurEdit.Lines = make(map[int]string)
				}
				for x := 0; x < dLen; x++ {
					player.CurEdit.Lines[x] = dLines[x]
					player.CurEdit.NumLines++
					player.CurEdit.CurLine++
				}
				player.CurEdit.NumLines--
				player.CurEdit.CurLine--
				MleEditor(player, argThreeThrough)
				WriteToPlayer(player, "Description sent to editor.")
				return
			} else if cmdBl == "paste" {
				newDesc := ""
				for x := 0; x <= player.CurEdit.NumLines; x++ {
					newDesc = newDesc + player.CurEdit.Lines[x] + "\r\n"
				}
				player.OLCEdit.Room.RoomLink.Description = newDesc
				WriteToPlayer(player, "Text transfered from editor.")
				CmdOLC(player, "")
				return
			}
			player.OLCEdit.Room.RoomLink.Description = argTwoThrough
			WriteToPlayer(player, "Description set")
			glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave
		} else if cmdl == "exit" || cmdl == "exits" {
			if cmdB == "" {
				WriteToPlayer(player, "OLC exit <exit name>")
			} else if cmdBl == "create" {
				if argThreeThrough != "" {
					for exitName, _ := range player.OLCEdit.Room.RoomLink.Exits {
						if strings.EqualFold(exitName, argThreeThrough) {
							WriteToPlayer(player, "That exit already exists.")
							return
						}
					}
					if player.OLCEdit.Room.RoomLink.Exits == nil {
						player.OLCEdit.Room.RoomLink.Exits = make(map[string]*glob.ExitData)
					}
					player.OLCEdit.Room.RoomLink.Exits[argThreeThrough] = CreateExit()
					player.OLCEdit.Room.RoomLink.Exits[argThreeThrough].ToRoom.RoomLink = player.OLCEdit.Room.RoomLink
					player.OLCEdit.Room.RoomLink.Exits[argThreeThrough].ToRoom.ID = player.OLCEdit.Room.ID
					player.OLCEdit.Room.RoomLink.Exits[argThreeThrough].ToRoom.Sector = player.OLCEdit.Room.Sector
					player.OLCEdit.Exit.ExitLink = player.OLCEdit.Room.RoomLink.Exits[argThreeThrough]
					player.OLCEdit.Exit.Name = argThreeThrough
					player.OLCEdit.Mode = def.OLC_EXITS
					glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave
					CmdOLC(player, "")
					return
				} else {
					WriteToPlayer(player, "OLC exit <exit name>")
				}
			} else {

				for exitName, exit := range player.OLCEdit.Room.RoomLink.Exits {
					if strings.EqualFold(exitName, argTwoThrough) {
						player.OLCEdit.Exit.ExitLink = exit
						player.OLCEdit.Exit.Name = exitName
						player.OLCEdit.Mode = def.OLC_EXITS
						WriteToPlayer(player, "Exit found, switching to exit editor.")
						CmdOLC(player, "")
						return
					}
				}
				WriteToPlayer(player, "I didn't find an exit by that name, to create: OLC exit create <Exit Name>")
			}

		}
	} else {

		buf := ""
		exits := ""
		if player.OLCEdit.Room.RoomLink != nil {
			for name, exit := range player.OLCEdit.Room.RoomLink.Exits {
				if name != "" {
					exits = exits + fmt.Sprintf("%v, ToRoom: %v:%v, Door: %v, AutoOpen: %v, AutoClose: %v, Hidden: %v, Keyed: %v\r\n",
						name, exit.ToRoom.Sector, exit.ToRoom.ID, exit.Door.Door, exit.Door.AutoOpen, exit.Door.AutoClose,
						exit.Door.Hidden, exit.Door.Keyed)
				}
			}
			if exits == "" {
				exits = "None"
			}
			buf = buf + fmt.Sprintf("Room: %v:%v (sector/id)\r\nName: %v\r\nDescription: \r\n\r\n%v\r\n\r\nExits:\r\n%v",
				player.OLCEdit.Room.Sector, player.OLCEdit.Room.ID,
				player.OLCEdit.Room.RoomLink.Name, player.OLCEdit.Room.RoomLink.Description, exits)
			WriteToBuilder(player, buf)
			WriteToPlayer(player, "Syntax for OLC room:\r\rolc <name/description> <text>\r\nolc exit <exit name>\r\nolc exit create <direction name>\r\nolc room <location>")
		} else {
			WriteToPlayer(player, "No room selected in editor, selecting current room.")
			player.OLCEdit.Room.ID = player.Location.ID
			player.OLCEdit.Room.Sector = player.Location.Sector
			player.OLCEdit.Room.RoomLink = player.Location.RoomLink
			CmdOLC(player, "")
			return
		}

	}
}
