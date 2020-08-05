package support

import (
	"fmt"
	"strconv"
	"strings"

	"../def"
	"../glob"
)

func getLocationFromString(player *glob.PlayerData, input string) (glob.LocationData, bool) {
	loc := strings.Split(input, ":")
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

		if erra == nil && errb == nil {
			room, roomFound := LocationDataFromID(sector, id)
			if roomFound {
				return room, true
			}
		}
	}
	return glob.LocationData{}, false
}

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
		player.Dirty = true
		CmdLook(player, "")
	} else {
		WriteToPlayer(player, "That location doesn't exist.")
	}
}

func CmdOLC(player *glob.PlayerData, input string) {

	//TODO, IN EDITOR HELP
	if player.OLCSettings.NoOLCPrefix == false && player.OLCSettings.NoHint == false {
		defer WriteToPlayer(player, "cmd <command> for non-olc commands, or type STOP to exit OLC.")
	}
	command, argTwoThrough := SplitArgsTwo(input, " ")
	cmdB, argThreeThrough := SplitArgsTwo(argTwoThrough, " ")
	cmdl := strings.ToLower(command)
	cmdBl := strings.ToLower(cmdB)

	player.OLCEdit.Active = true

	if player.OLCSettings.NoOLCPrefix == false && strings.EqualFold(command, "olc") {
		WriteToPlayer(player, "You are in noPrefix mode, you don't need to type OLC before commands.\r\ntype settings to turn this mode off... Or type STOP to exit.")
		return
	} else if strings.EqualFold("stop", input) {
		player.OLCEdit.Active = false
		WriteToPlayer(player, "Exiting OLC!")
		return
	} else if cmdl == "settings" {
		OLCSettings := []glob.ConfigData{
			{ID: 1, Name: "follow", Help: "If on: you are always editing the room you are standing in.",
				Ref: &player.OLCSettings.OLCRoomFollow},
			{ID: 2, Name: "showCodes", Help: "If on: Show color codes in names / descriptions / etc",
				Ref: &player.OLCSettings.OLCShowCodes},
			//{ID: 3, Name: "showAllCodes", Help: "If on: Show color codes, instead of color for the whOLC mud.",
			//Ref: player.OLCSettings.OLCShowAllCodes},
			{ID: 4, Name: "prompt", Help: "If on: Change your prompt to OLC information while in editor.",
				Ref: &player.OLCSettings.OLCPrompt},
			//{ID: 5, Name: "promptString", Help: "Customize OLC prompt.",
			//Ref: player.OLCSettings.OLCPromptString},
			{ID: 6, Name: "noOLCPrefix", Help: "If on: When in editor, all input goes to olc by default.",
				Ref: &player.OLCSettings.NoOLCPrefix},
			{ID: 7, Name: "noHint", Help: "If on: Turn off message explaining STOP/CMD for NoOLCPrefix",
				Ref: &player.OLCSettings.NoHint},
		}

		cmdNames := []string{}
		for _, c := range OLCSettings {
			cmdNames = append(cmdNames, strings.ToLower(c.Name))
		}
		match, _ := FindClosestMatch(cmdNames, argTwoThrough)
		player.Dirty = true

		if match == "follow" {
			if player.OLCSettings.OLCRoomFollow {
				player.OLCSettings.OLCRoomFollow = false
				WriteToPlayer(player, "OLC will no longer change the room you are editing when you move.")
			} else {
				player.OLCSettings.OLCRoomFollow = true
				WriteToPlayer(player, "OLC will automatically edit whatever room you move to.")
			}
		} else if match == "showcodes" {
			if player.OLCSettings.OLCShowCodes {
				player.OLCSettings.OLCShowCodes = false
				WriteToPlayer(player, "OLC will now just show normal color.")
			} else {
				player.OLCSettings.OLCShowCodes = true
				WriteToPlayer(player, "OLC will show color codes in names and descriptions.")
			}
		} else if match == "prompt" {
			if player.OLCSettings.OLCPrompt {
				player.OLCSettings.OLCPrompt = false
				WriteToPlayer(player, "Your prompt will no longer change to OLC prompt while editing.")
			} else {
				player.OLCSettings.OLCPrompt = true
				WriteToPlayer(player, "Your prompt will now be OLC information.")
			}
		} else if match == "noolcprefix" {
			if player.OLCSettings.NoOLCPrefix {
				player.OLCSettings.NoOLCPrefix = false
				WriteToPlayer(player, "Your input will NOT be sent directly to olc. Prefix all OLC commands with: olc <command>.")
			} else {
				player.OLCSettings.NoOLCPrefix = true
				WriteToPlayer(player, "All your input will be directed to OLC, until you exit it.\r\ncmd <command> will pass-through commands.")
			}
		} else if match == "nohint" {
			if player.OLCSettings.NoHint {
				player.OLCSettings.NoHint = false
				WriteToPlayer(player, "After every line, remind you of STOP/CMD commands")
			} else {
				player.OLCSettings.NoHint = true
				WriteToPlayer(player, "Do not show reminder for STOP/CMD commands")
			}
		}

		//Show settings avaialble
		for _, cmd := range OLCSettings {
			WriteToPlayer(player, fmt.Sprintf("%10v:%5v --  %v", cmd.Name, boolToOnOff(*cmd.Ref), cmd.Help))
		}
		return
	}
	if player.OLCEdit.Mode == def.OLC_NONE {
		if cmdl == "done" {
			player.OLCEdit.Mode = def.OLC_NONE
			WriteToPlayer(player, "Exiting OLC.")
			player.OLCEdit.Active = false
			return
		} else if cmdl == "" {
			WriteToPlayer(player, "Possible types:")
			WriteToPlayer(player, "Room, object, trigger, mobile, quest or sector.")
			WriteToPlayer(player, "OLC <type>, or OLC <type>")
			WriteToPlayer(player, "To exit editor: OLC done, for settings: OLC settings")
			return
		}

		if cmdl == "room" {
			player.OLCEdit.Mode = def.OLC_ROOM
			player.OLCEdit.Room = player.Location
		} else if cmdl == "object" {
			player.OLCEdit.Mode = def.OLC_OBJECT
		} else if cmdl == "trigger" {
			player.OLCEdit.Mode = def.OLC_TRIGGER
		} else if cmdl == "mobile" {
			player.OLCEdit.Mode = def.OLC_MOBILE
		} else if cmdl == "quest" {
			player.OLCEdit.Mode = def.OLC_QUEST
		} else if cmdl == "sector" {
			player.OLCEdit.Mode = def.OLC_SECTOR
		}

	}

	/* ROOM-EXITS EDITOR */
	if player.OLCEdit.Mode == def.OLC_EXITS {

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
				exitName := player.OLCEdit.ExitName
				exitDoor := player.OLCEdit.Exit.Door
				exitToRoom := player.OLCEdit.Exit.ToRoom

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

				player.OLCEdit.Exit = player.OLCEdit.Room.RoomLink.Exits[argTwoThrough]
				player.OLCEdit.ExitName = argTwoThrough
				player.OLCEdit.Mode = def.OLC_EXITS
				glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave
				CmdOLC(player, "")
				return
			} else {
				WriteToPlayer(player, "OLC exit <exit name>")
			}

		} else if cmdl == "door" {
			if player.OLCEdit.Exit.Door.Door {
				player.OLCEdit.Exit.Door.Door = false
			} else {
				player.OLCEdit.Exit.Door.Door = true
			}
		} else if cmdl == "autoopen" {
			if player.OLCEdit.Exit.Door.AutoOpen {
				player.OLCEdit.Exit.Door.AutoOpen = false
			} else {
				player.OLCEdit.Exit.Door.AutoOpen = true
			}
		} else if cmdl == "autoclose" {
			if player.OLCEdit.Exit.Door.AutoClose {
				player.OLCEdit.Exit.Door.AutoClose = false
			} else {
				player.OLCEdit.Exit.Door.AutoClose = true
			}
		} else if cmdl == "keyed" {
			if player.OLCEdit.Exit.Door.AutoClose {
				player.OLCEdit.Exit.Door.AutoClose = false
			} else {
				player.OLCEdit.Exit.Door.AutoClose = true
			}
		} else if cmdl == "delete" {
			exitName := player.OLCEdit.ExitName
			delete(player.OLCEdit.Room.RoomLink.Exits, exitName)
			player.OLCEdit.Exit = nil
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
			player.OLCEdit.Exit.ToRoom = loc
		}
		buf := fmt.Sprintf("OLC EDIT EXITS:\r\n%10v: %v\r\n\r\n%10v: %v:%v\r\n%10v: %v:%v\r\n",
			"Name",
			player.OLCEdit.ExitName,
			"FromRoom",
			player.OLCEdit.Room.Sector,
			player.OLCEdit.Room.ID,
			"ToRoom",
			player.OLCEdit.Exit.ToRoom.Sector,
			player.OLCEdit.Exit.ToRoom.ID)
		WriteToBuilder(player, buf)
		buf = fmt.Sprintf("%10v: %v\r\n%10v: %v\r%10v: %v\r\n%10v: %v\r%10v: %v",
			"Door",
			boolToYesNo(player.OLCEdit.Exit.Door.Door),
			"AutoOpen",
			boolToYesNo(player.OLCEdit.Exit.Door.AutoOpen),
			"AutoClose",
			boolToYesNo(player.OLCEdit.Exit.Door.AutoClose),
			"Hidden",
			boolToYesNo(player.OLCEdit.Exit.Door.Hidden),
			"Keyed",
			boolToYesNo(player.OLCEdit.Exit.Door.Keyed))
		WriteToPlayer(player, buf)
		WriteToPlayer(player, "Syntax for OLC exits: olc ToRoom <location>, door, autoOpen, autoClose, keyed, delete, done")

		/* ROOM EDITOR */
	} else if player.OLCEdit.Mode == def.OLC_ROOM {
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
						player.OLCEdit.Room = editRoom
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
						player.OLCEdit.Room = editRoom
						CmdOLC(player, "")
						WriteToPlayer(player, "Room already exists.")
						return
					} else {
						glob.SectorsList[sector].Rooms[id] = CreateRoom()
						editRoom, _ := LocationDataFromID(sector, id)
						player.OLCEdit.Room = editRoom
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
						player.OLCEdit.Room.RoomLink.Exits[argThreeThrough].ToRoom = player.OLCEdit.Room
						player.OLCEdit.Exit = player.OLCEdit.Room.RoomLink.Exits[argThreeThrough]
						player.OLCEdit.ExitName = argThreeThrough
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
							player.OLCEdit.Exit = exit
							player.OLCEdit.ExitName = exitName
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
					exits = exits + fmt.Sprintf("%v, ToRoom: %v:%v, Door: %v, AutoOpen: %v, AutoClose: %v, Hidden: %v, Keyed: %v\r\n",
						name, exit.ToRoom.Sector, exit.ToRoom.ID, exit.Door.Door, exit.Door.AutoOpen, exit.Door.AutoClose,
						exit.Door.Hidden, exit.Door.Keyed)
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
				WriteToPlayer(player, "No room selected in editor")
			}

		}

	} else if player.OLCEdit.Mode == def.OLC_OBJECT {
		if cmdl == "done" {
			player.OLCEdit.Mode = def.OLC_NONE
			WriteToPlayer(player, "Exiting OLC.")
			player.OLCEdit.Active = false
			return
		}
		WriteToPlayer(buf)
	} else if player.OLCEdit.Mode == def.OLC_TRIGGER {
		if cmdl == "done" {
			player.OLCEdit.Mode = def.OLC_NONE
			WriteToPlayer(player, "Exiting OLC.")
			player.OLCEdit.Active = false
			return
		}
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLCEdit.Mode == def.OLC_MOBILE {
		if cmdl == "done" {
			player.OLCEdit.Mode = def.OLC_NONE
			WriteToPlayer(player, "Exiting OLC.")
			player.OLCEdit.Active = false
			return
		}
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLCEdit.Mode == def.OLC_QUEST {
		if cmdl == "done" {
			player.OLCEdit.Mode = def.OLC_NONE
			WriteToPlayer(player, "Exiting OLC.")
			player.OLCEdit.Active = false
			return
		}
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLCEdit.Mode == def.OLC_SECTOR {

		if player.OLCEdit.Sector == 0 {
			sid := player.Location.Sector
			player.OLCEdit.Sector = sid
		}
		sector := &glob.SectorsList[player.OLCEdit.Sector]

		if cmdl == "" {
			buf := fmt.Sprintf("Name: %v\r\nID %v\r\nFingerprint: %v\r\nDescription: %v\r\nArea: %v\r\nRoom count: %v",
				sector.Name, sector.ID, sector.Fingerprint, sector.Description, sector.Area, sector.NumRooms)
			WriteToBuilder(player, buf)
		} else {
			/* If sector specified, use it, otherwise use player location */

			sector = &glob.SectorsList[player.OLCEdit.Sector]

			if cmdl == "done" {
				player.OLCEdit.Mode = def.OLC_NONE
				WriteToPlayer(player, "Exiting OLC.")
				player.OLCEdit.Active = false
				return
			} else if strings.EqualFold(cmdl, "sector") {
				psid, err := strconv.Atoi(cmdBl)

				if err == nil {
					if glob.SectorsList[psid].Valid {
						player.OLCEdit.Sector = psid
						WriteToPlayer(player, "Sector "+cmdBl+" selected")
					} else {
						WriteToPlayer(player, "Invalid sector, use sector create.")
						return
					}
				}
				CmdOLC(player, "")
				return
			} else if strings.EqualFold(cmdl, "name") {
				sector.Name = argTwoThrough
				sector.Valid = true
				if sector.Fingerprint == "" {
					sector.Fingerprint = MakeFingerprint(sector.Name)
				}
				WriteToPlayer(player, "Name set.")
			} else if strings.EqualFold(cmdl, "desc") || strings.EqualFold(cmdl, "description") {
				//Todo, editor
				sector.Description = argTwoThrough
				WriteToPlayer(player, "Description set.")
			} else if strings.EqualFold(cmdl, "area") {
				sector.Area = argTwoThrough
				WriteToPlayer(player, "Area set.")
			} else if strings.EqualFold(cmdl, "create") {
				glob.SectorsListEnd++

				newSector := CreateSector()
				newSector.Valid = true
				glob.SectorsList[glob.SectorsListEnd] = *newSector

				player.OLCEdit.Sector = glob.SectorsListEnd
				WriteToPlayer(player, "Sector created.")
				CmdOLC(player, "")
				return
			} else {
				WriteToPlayer(player, "That isn't a valid command.\r\nCommands: <name/description/area> <text>")
				return
			}
			CmdOLC(player, "")
			return
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
			glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave

		} else if dirTwo != "" {
			doDigCustom(player, rooms, found, dirOne, dirTwo, sector)
			WriteToPlayer(player, fmt.Sprintf("Digging %v:%v", dirOne, dirTwo))
			glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave
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

	if player.OLCSettings.OLCShowCodes {
		bytes, err = player.Connection.Desc.Write([]byte(text + "\r\n"))
	} else {
		bytes, err = player.Connection.Desc.Write([]byte(ANSIColor(text) + "\r\n"))
	}
	player.Connection.BytesOut += bytes
	trackBytesOut(player.Connection)

	DescWriteError(player.Connection, err)
}
