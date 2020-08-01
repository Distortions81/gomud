package support

import (
	"fmt"
	"strconv"
	"strings"

	"../def"
	"../glob"
)

func CmdAsave(player *glob.PlayerData, args string) {
	WriteSectorList()
	WriteToPlayer(player, "All sectors saving.")
}

func CmdGoto(player *glob.PlayerData, input string) {
	a, b := SplitArgsTwo(input, ":")
	sector := 0
	id := 0

	if b == "" {
		sector = player.Location.Sector
		id, _ = strconv.Atoi(a)
	} else {
		sector, _ = strconv.Atoi(a)
		id, _ = strconv.Atoi(b)
	}

	if glob.SectorsList[sector].Valid &&
		glob.SectorsList[sector].Rooms != nil &&
		glob.SectorsList[sector].Rooms[id] != nil {
		WriteToPlayer(player, fmt.Sprintf("Going to %v:%v...", sector, id))
		WriteToRoom(player, fmt.Sprintf("%v vanishes in a puff of smoke.", player.Name))
		PlayerToRoom(player, sector, id)
		WriteToRoom(player, fmt.Sprintf("A puff of smoke appears, and %v emerges from it.", player.Name))
		CmdLook(player, "")
	} else {
		WriteToPlayer(player, "That location doesn't exist.")
	}
}

func CmdOLE(player *glob.PlayerData, input string) {

	command, argTwoThrough := SplitArgsTwo(input, " ")
	cmdB, argThreeThrough := SplitArgsTwo(argTwoThrough, " ")
	command = strings.ToLower(command)
	cmdB = strings.ToLower(cmdB)

	player.OLEEdit.Active = true

	if player.OLEEdit.Mode == def.OLE_NONE {
		if command == "done" {
			player.OLEEdit.Mode = def.OLE_NONE
			WriteToPlayer(player, "Exiting OLE.")
			player.OLEEdit.Active = false
			return
		} else if command == "" {
			WriteToPlayer(player, "Possible types:")
			WriteToPlayer(player, "Room, object, trigger, mobile, quest or sector.")
			WriteToPlayer(player, "OLE <type>, or OLE <type> <sector:id> (for a specfic item), or just <id> (for sector you are standing in)")
			WriteToPlayer(player, "Other commands: DONE (to exit OLE), and Settings.")
			WriteToPlayer(player, "Typing the command OLE (by itself) will show the editor again, so will enter/return on a blank line.")
			return
		}

		if command == "settings" {
			OLESettings := []glob.ConfigData{
				{ID: 1, Name: "follow", Help: "If on: you are always editing the room you are standing in.",
					Ref: &player.OLESettings.OLERoomFollow},
				{ID: 2, Name: "showCodes", Help: "If on: Show color codes in names / descriptions / etc",
					Ref: &player.OLESettings.OLEShowCodes},
				//{ID: 3, Name: "showAllCodes", Help: "If on: Show color codes, instead of color for the whole mud.",
				//Ref: player.OLESettings.OLEShowAllCodes},
				{ID: 4, Name: "prompt", Help: "If on: Change your prompt to OLE information while in editor.",
					Ref: &player.OLESettings.OLEPrompt},
				//{ID: 5, Name: "promptString", Help: "Customize OLE prompt.",
				//Ref: player.OLESettings.OLEPromptString},
			}

			cmdNames := []string{}
			for _, c := range OLESettings {
				cmdNames = append(cmdNames, strings.ToLower(c.Name))
			}
			match, _ := FindClosestMatch(cmdNames, argTwoThrough)

			if match == "follow" {
				if player.OLESettings.OLERoomFollow {
					player.OLESettings.OLERoomFollow = false
					WriteToPlayer(player, "OLE will no longer change the room you are editing when you move.")
					return
				} else {
					player.OLESettings.OLERoomFollow = true
					WriteToPlayer(player, "OLE will automatically edit whatever room you move to.")
					return
				}
			} else if match == "showcodes" {
				if player.OLESettings.OLEShowCodes {
					player.OLESettings.OLEShowCodes = false
					WriteToPlayer(player, "OLE will now just show normal color.")
					return
				} else {
					player.OLESettings.OLEShowCodes = true
					WriteToPlayer(player, "OLE will show color codes in names and descriptions.")
					return
				}
			} else if match == "prompt" {
				if player.OLESettings.OLEPrompt {
					player.OLESettings.OLEPrompt = false
					WriteToPlayer(player, "Your prompt will no longer change to OLE prompt while editing.")
					return
				} else {
					player.OLESettings.OLEPrompt = true
					WriteToPlayer(player, "Your prompt will now be OLE information.")
					return
				}
			}

			//Show settings avaialble
			for _, cmd := range OLESettings {
				WriteToPlayer(player, fmt.Sprintf("%10v:%5v --  %v", cmd.Name, boolToOnOff(*cmd.Ref), cmd.Help))
			}
			return
		}
		if command == "room" {
			WriteToPlayer(player, "OLE EDIT: ROOM")
			player.OLEEdit.Mode = def.OLE_ROOM
			if argTwoThrough == "" {
				player.OLEEdit.Room = player.Location
				CmdOLE(player, "")
			} else if cmdB == "create" {
				loc := strings.Split(argThreeThrough, ":")
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
						WriteToPlayer(player, "Syntax: OLE room create <sector:id> or just <id> (for sector you are standing in)\r\nExample: OLE room create 1:1")
						return
					}

					editRoom, roomFound := LocationDataFromID(sector, id)
					if roomFound {
						player.OLEEdit.Room = editRoom
						CmdOLE(player, "")
						WriteToPlayer(player, "Room already exists.")
					} else {
						glob.SectorsList[sector].Rooms[id] = CreateRoom()
						editRoom, _ := LocationDataFromID(sector, id)
						player.OLEEdit.Room = editRoom
						CmdOLE(player, "")
						WriteToPlayer(player, fmt.Sprintf("Room %v:%v created!", sector, id))
						glob.SectorsList[player.OLEEdit.Room.Sector].Dirty = true //Autosave
					}
				} else {
					WriteToPlayer(player, "Syntax: OLE room create <sector:id> or just <id> (for sector you are standing in)\r\nExample: OLE room create 1:1")
				}
			} else {
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
						WriteToPlayer(player, "Syntax: OLE room <sector:id> or just <id> (for sector you are standing in\r\nExample: OLE room 1:1")
						return
					}

					editRoom, roomFound := LocationDataFromID(sector, id)
					if roomFound {
						player.OLEEdit.Room = editRoom
						CmdOLE(player, "")
					} else {
						WriteToPlayer(player, "I couldn't find that room, to create: OLE room create <sector:id> or just <id> (for sector you are standing in)")
					}
				} else {
					WriteToPlayer(player, "Syntax: OLE room <sector:id> or just <id> (for sector you are standing in)\r\nExample: OLE room 1:1")
					return
				}

			}
		} else if command == "object" {
			player.OLEEdit.Mode = def.OLE_OBJECT
		} else if command == "trigger" {
			player.OLEEdit.Mode = def.OLE_TRIGGER
		} else if command == "mobile" {
			player.OLEEdit.Mode = def.OLE_MOBILE
		} else if command == "quest" {
			player.OLEEdit.Mode = def.OLE_QUEST
		} else if command == "sector" {
			player.OLEEdit.Mode = def.OLE_SECTOR
		}

	}

	/* ROOM-EXITS EDITOR */
	if player.OLEEdit.Mode == def.OLE_EXITS {

		if command == "done" {
			player.OLEEdit.Mode = def.OLE_ROOM
			WriteToPlayer(player, "Going back to room editor..")
			CmdOLE(player, "")
			return
		}
		buf := fmt.Sprintf("OLE EDIT EXITS:\r\nName: %v\r\n\r\nFromRoom: %v:%v\r\nToRoom: %v:%v\r\n",
			player.OLEEdit.ExitName,
			player.OLEEdit.Room.Sector,
			player.OLEEdit.Room.ID,
			player.OLEEdit.Exit.ToRoom.Sector,
			player.OLEEdit.Exit.ToRoom.ID)
		WriteToBuilder(player, buf)
		buf = fmt.Sprintf("Door: %v\r\nAutoOpen: %v, AutoClose %v\r\nHidden: %v, Keyed %v.",
			boolToYesNo(player.OLEEdit.Exit.Door.Door),
			boolToYesNo(player.OLEEdit.Exit.Door.AutoOpen),
			boolToYesNo(player.OLEEdit.Exit.Door.AutoClose),
			boolToYesNo(player.OLEEdit.Exit.Door.Hidden),
			boolToYesNo(player.OLEEdit.Exit.Door.Keyed))
		WriteToPlayer(player, buf)

		/* ROOM EDITOR */
	} else if player.OLEEdit.Mode == def.OLE_ROOM {
		if command == "room" {
			WriteToPlayer(player, "Already in room editor.")
		} else if command == "done" {
			player.OLEEdit.Mode = def.OLE_NONE
			WriteToPlayer(player, "Exiting OLE.")
			player.OLEEdit.Active = false
			return
		}
		if argTwoThrough != "" {
			if command == "name" {
				player.OLEEdit.Room.RoomLink.Name = argTwoThrough
				WriteToPlayer(player, "Name set")
				glob.SectorsList[player.OLEEdit.Room.Sector].Dirty = true //Autosave

			} else if command == "description" || command == "desc" {
				player.OLEEdit.Room.RoomLink.Description = argTwoThrough
				WriteToPlayer(player, "Description set")
				glob.SectorsList[player.OLEEdit.Room.Sector].Dirty = true //Autosave
			} else if command == "exit" || command == "exits" {
				if cmdB == "" {
					WriteToPlayer(player, "OLE exit <exit name>")
				} else if cmdB == "create" {
					if argThreeThrough != "" {
						for exitName, _ := range player.OLEEdit.Room.RoomLink.Exits {
							if strings.EqualFold(exitName, argThreeThrough) {
								WriteToPlayer(player, "That exit already exists.")
								return
							}
						}
						player.OLEEdit.Room.RoomLink.Exits[argThreeThrough] = CreateExit()
						player.OLEEdit.Room.RoomLink.Exits[argThreeThrough].ToRoom = player.OLEEdit.Room
						player.OLEEdit.Exit = player.OLEEdit.Room.RoomLink.Exits[argThreeThrough]
						player.OLEEdit.ExitName = argThreeThrough
						player.OLEEdit.Mode = def.OLE_EXITS
						CmdOLE(player, "")
						return
					} else {
						WriteToPlayer(player, "OLE exit <exit name>")
					}
				} else {

					for exitName, exit := range player.OLEEdit.Room.RoomLink.Exits {
						if strings.EqualFold(exitName, argTwoThrough) {
							player.OLEEdit.Exit = exit
							player.OLEEdit.ExitName = exitName
							player.OLEEdit.Mode = def.OLE_EXITS
							WriteToPlayer(player, "Exit found, switching to exit editor.")
							CmdOLE(player, "")
							return
						}
					}
					WriteToPlayer(player, "I didn't find an exit by that name, to create: OLE exit create <Exit Name>")
				}

			}
		} else {

			buf := ""
			exits := ""
			if player.OLEEdit.Room.RoomLink != nil {
				for name, exit := range player.OLEEdit.Room.RoomLink.Exits {
					exits = exits + fmt.Sprintf("%v, ToRoom: %v:%v, Door: %v, AutoOpen: %v, AutoClose: %v, Hidden: %v, Keyed: %v\r\n",
						name, exit.ToRoom.Sector, exit.ToRoom.ID, exit.Door.Door, exit.Door.AutoOpen, exit.Door.AutoClose,
						exit.Door.Hidden, exit.Door.Keyed)
				}
				if exits == "" {
					exits = "None"
				}
				buf = buf + fmt.Sprintf("Room: %v:%v (sector/id)\r\nName: %v\r\nDescription: \r\n\r\n%v\r\n\r\nExits:\r\n%v",
					player.OLEEdit.Room.Sector, player.OLEEdit.Room.ID,
					player.OLEEdit.Room.RoomLink.Name, player.OLEEdit.Room.RoomLink.Description, exits)
				WriteToBuilder(player, buf)
				WriteToPlayer(player, "Syntax for OLE room: OLE <name, description, exit> <item>")
			} else {
				WriteToPlayer(player, "No room selected in editor")
			}

		}

	} else if player.OLEEdit.Mode == def.OLE_OBJECT {
		if command == "done" {
			player.OLEEdit.Mode = def.OLE_NONE
			WriteToPlayer(player, "Exiting OLE.")
			player.OLEEdit.Active = false
			return
		}
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLEEdit.Mode == def.OLE_TRIGGER {
		if command == "done" {
			player.OLEEdit.Mode = def.OLE_NONE
			WriteToPlayer(player, "Exiting OLE.")
			player.OLEEdit.Active = false
			return
		}
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLEEdit.Mode == def.OLE_MOBILE {
		if command == "done" {
			player.OLEEdit.Mode = def.OLE_NONE
			WriteToPlayer(player, "Exiting OLE.")
			player.OLEEdit.Active = false
			return
		}
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLEEdit.Mode == def.OLE_QUEST {
		if command == "done" {
			player.OLEEdit.Mode = def.OLE_NONE
			WriteToPlayer(player, "Exiting OLE.")
			player.OLEEdit.Active = false
			return
		}
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLEEdit.Mode == def.OLE_SECTOR {
		if command == "done" {
			player.OLEEdit.Mode = def.OLE_NONE
			WriteToPlayer(player, "Exiting OLE.")
			player.OLEEdit.Active = false
			return
		}
		if argTwoThrough == "" {
			sid := player.Location.Sector
			sector := glob.SectorsList[sid]
			player.OLEEdit.Sector = sid

			buf := fmt.Sprintf("Name: %v\r\nID %v\r\nFingerprint: %v\r\nDescription: %v\r\nArea: %v\r\nRoom count: %v\r\nValid: %v",
				sector.Name, sector.ID, sector.Fingerprint, sector.Description, sector.Area, sector.NumRooms, sector.Valid)
			WriteToBuilder(player, buf)
		}
	}

}

func CmdDig(player *glob.PlayerData, input string) {
	if player.Location.RoomLink == nil {
		WriteToPlayer(player, "You need to be in a room, to dig.")
		return
	}

	command, _ := SplitArgsTwo(input, " ")
	dirOne, dirTwo := SplitArgsTwo(input, ":")

	curID := player.Location.ID
	sector := player.Location.Sector

	if player.Location.RoomLink.Exits[strings.Title(dirOne)] != nil {
		WriteToPlayer(player, "That exit is already occupied.")
		return
	}

	rooms := glob.SectorsList[sector].Rooms

	//Find first available slot
	found := 0
	for x := curID; ; x++ {
		/* Re-use old room */
		if rooms[x] != nil && rooms[x].Valid == false {
			found = x
			break
		}
		/* New room */
		if rooms[x] == nil {
			found = x
			break
		}
	}

	if rooms[found] == nil {
		rooms[found] = CreateRoom()
	}

	if command != "" {
		if IsStandardDirection(command) {
			doDig(player, rooms, found, command, sector)

		} else if dirTwo != "" {
			doDigCustom(player, rooms, found, dirOne, dirTwo, sector)
			WriteToPlayer(player, fmt.Sprintf("Digging %v:%v", dirOne, dirTwo))
		} else {
			WriteToPlayer(player, "Custom directions require names for both sides of the direction. dig climb up:slide down")
		}
	} else {
		WriteToPlayer(player, "dig <direction> (north,south,east,west), or dig <enter:exit>.\r\nExample: dig climb up:slide down.")
	}

}

func doDig(player *glob.PlayerData, rooms map[int]*glob.RoomData, found int, dir string, sector int) {
	room := rooms[found]

	room.Valid = true
	room.Name = "new room"

	room.Exits[GetStandardDirectionMirror(dir)] = CreateExit()
	room.Exits[GetStandardDirectionMirror(dir)].ToRoom = player.Location

	player.Location.RoomLink.Exits[strings.Title(dir)] = CreateExit()
	player.Location.RoomLink.Exits[strings.Title(dir)].ToRoom.ID = found
	player.Location.RoomLink.Exits[strings.Title(dir)].ToRoom.Sector = sector

	CmdGo(player, dir)
	glob.SectorsList[sector].Dirty = true //Autosave
}

func doDigCustom(player *glob.PlayerData, rooms map[int]*glob.RoomData, found int, dirOne string, dirTwo string, sector int) {
	room := rooms[found]

	room.Valid = true
	room.Name = "new room"

	room.Exits[dirTwo] = CreateExit()
	room.Exits[dirTwo].ToRoom = player.Location

	player.Location.RoomLink.Exits[dirOne] = CreateExit()
	player.Location.RoomLink.Exits[dirOne].ToRoom.ID = found
	player.Location.RoomLink.Exits[dirOne].ToRoom.Sector = sector

	CmdGo(player, dirOne)
	glob.SectorsList[sector].Dirty = true //Autosave
}

func WriteToBuilder(player *glob.PlayerData, text string) {

	var bytes int
	var err error

	if player == nil || !player.Valid || player.Connection == nil || !player.Connection.Valid {
		return
	}

	if player.OLESettings.OLEShowCodes {
		bytes, err = player.Connection.Desc.Write([]byte(text + "\r\n"))
	} else {
		bytes, err = player.Connection.Desc.Write([]byte(ANSIColor(text) + "\r\n"))
	}
	player.Connection.BytesOut += bytes
	trackBytesOut(player.Connection)

	DescWriteError(player.Connection, err)
}
