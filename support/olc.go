package support

import (
	"fmt"
	"strconv"
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

	if command == "done" {
		player.OLCEdit.Mode = def.OLC_NONE
		WriteToPlayer(player, "Exiting OLC.")
		player.OLCEdit.Active = false
		return
	}
	player.OLCEdit.Active = true

	if command == "" && player.OLCEdit.Mode == def.OLC_NONE {
		WriteToPlayer(player, "Possible types:")
		WriteToPlayer(player, "Room, object, trigger, mobile, quest or sector.")
		WriteToPlayer(player, "OLC <type>, or OLC <type> <sector:id> (for a specfic item)")
		WriteToPlayer(player, "Other commands: DONE (to exit olc), and Settings.")
		WriteToPlayer(player, "Typing the command OLC (by itself) will show the editor again.")
		return
	} else {

		if command == "settings" {
			olcSettings := []glob.CommandArgData{
				{ID: 1, Name: "follow", Help: "If on: you are always editing the room you are standing in.",
					Ref: player.OLCSettings.OlcRoomFollow},
				{ID: 2, Name: "showCodes", Help: "If on: Show color codes in names / descriptions / etc",
					Ref: player.OLCSettings.OlcShowCodes},
				//{ID: 3, Name: "showAllCodes", Help: "If on: Show color codes, instead of color for the whole mud.",
				//Ref: player.OLCSettings.OlcShowAllCodes},
				{ID: 4, Name: "prompt", Help: "If on: Change your prompt to OLC information while in editor.",
					Ref: player.OLCSettings.OlcPrompt},
				//{ID: 5, Name: "promptString", Help: "Customize OLC prompt.",
				//Ref: player.OLCSettings.OlcPromptString},
			}

			cmdNames := []string{}
			for _, c := range olcSettings {
				cmdNames = append(cmdNames, strings.ToLower(c.Name))
			}
			match, _ := FindClosestMatch(cmdNames, longArg)

			if match == "follow" {
				if player.OLCSettings.OlcRoomFollow {
					player.OLCSettings.OlcRoomFollow = false
					WriteToPlayer(player, "OLC will no longer change the room you are editing when you move.")
					return
				} else {
					player.OLCSettings.OlcRoomFollow = true
					WriteToPlayer(player, "OLC will automatically edit whatever room you move to.")
					return
				}
			} else if match == "showcodes" {
				if player.OLCSettings.OlcShowCodes {
					player.OLCSettings.OlcShowCodes = false
					WriteToPlayer(player, "OLC will now just show normal color.")
					return
				} else {
					player.OLCSettings.OlcShowCodes = true
					WriteToPlayer(player, "OLC will show color codes in names and descriptions.")
					return
				}
			} else if match == "prompt" {
				if player.OLCSettings.OlcPrompt {
					player.OLCSettings.OlcPrompt = false
					WriteToPlayer(player, "Your prompt will no longer change to OLC prompt while editing.")
					return
				} else {
					player.OLCSettings.OlcPrompt = true
					WriteToPlayer(player, "Your prompt will now be OLC information.")
					return
				}
			}

			//Show settings avaialble
			for _, cmd := range olcSettings {
				WriteToPlayer(player, fmt.Sprintf("%10v:%5v --  %v", cmd.Name, boolToOnOff(cmd.Ref), cmd.Help))
			}
			return
		}
		if command == "room" {
			WriteToPlayer(player, "OLC EDIT: ROOM")
			player.OLCEdit.Mode = def.OLC_ROOM
			if longArg == "" {
				player.OLCEdit.Room = player.Location
			} else {
				loc := strings.Split(longArg, ":")
				locLen := len(loc)

				if locLen == 2 {
					sector, errb := strconv.Atoi(loc[0])
					id, erra := strconv.Atoi((loc[1]))

					if erra != nil || errb != nil {
						WriteToPlayer(player, "Syntax: olc room sector:id\r\nExample: olc room 1:1")
						return
					}

					editRoom, roomFound := LocationDataFromID(sector, id)
					if roomFound {
						player.OLCEdit.Room = editRoom
						CmdOLC(player, "")
					} else {
						WriteToPlayer(player, "I couldn't find that room.")
					}

				} else {
					WriteToPlayer(player, "Syntax: olc room sector:id\r\nExample: olc room 1:1")
					return
				}

			}
		} else if command == "object" {
			player.OLCEdit.Mode = def.OLC_OBJECT
		} else if command == "trigger" {
			player.OLCEdit.Mode = def.OLC_TRIGGER
		} else if command == "mobile" {
			player.OLCEdit.Mode = def.OLC_MOBILE
		} else if command == "quest" {
			player.OLCEdit.Mode = def.OLC_QUEST
		} else if command == "sector" {
			player.OLCEdit.Mode = def.OLC_SECTOR
		}

	}

	/* ROOM-EXITS EDITOR */
	if player.OLCEdit.Mode == def.OLC_EXITS {
		buf := fmt.Sprintf("OLC EDIT EXITS:\r\nName: %v\r\n\r\nFromRoom: %v:%v\r\nToRoom: %v:%v\r\n",
			player.OLCEdit.ExitName,
			player.OLCEdit.Room.Sector,
			player.OLCEdit.Room.ID,
			player.OLCEdit.Exit.ToRoom.Sector,
			player.OLCEdit.Exit.ToRoom.ID)
		WriteToBuilder(player, buf)
		buf = fmt.Sprintf("Door: %v\r\nAutoOpen: %v, AutoClose %v\r\nHidden: %v, Keyed %v.",
			boolToYesNo(player.OLCEdit.Exit.Door.Door),
			boolToYesNo(player.OLCEdit.Exit.Door.AutoOpen),
			boolToYesNo(player.OLCEdit.Exit.Door.AutoClose),
			boolToYesNo(player.OLCEdit.Exit.Door.Hidden),
			boolToYesNo(player.OLCEdit.Exit.Door.Keyed))
		WriteToPlayer(player, buf)

		/* ROOM EDITOR */
	} else if player.OLCEdit.Mode == def.OLC_ROOM {
		if longArg != "" {
			if command == "name" {
				player.OLCEdit.Room.RoomLink.Name = longArg
				WriteToPlayer(player, "Name set")
				glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave

			} else if command == "description" || command == "desc" {
				player.OLCEdit.Room.RoomLink.Description = longArg
				WriteToPlayer(player, "Description set")
				glob.SectorsList[player.OLCEdit.Room.Sector].Dirty = true //Autosave
			} else if command == "exit" || command == "exits" {
				for exitName, exit := range player.OLCEdit.Room.RoomLink.Exits {
					if strings.EqualFold(exitName, longArg) {
						player.OLCEdit.Exit = exit
						player.OLCEdit.ExitName = exitName
						player.OLCEdit.Mode = def.OLC_EXITS
						WriteToPlayer(player, "Exit found, switching to exit editor.")
						CmdOLC(player, "")
						return
					}
				}
				WriteToPlayer(player, "I didn't find an exit by that name.")

			}
		} else {

			buf := ""
			exits := ""
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
			WriteToPlayer(player, "Syntax for OLC room: olc <name, description, exit> <item>")

		}

	} else if player.OLCEdit.Mode == def.OLC_OBJECT {
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLCEdit.Mode == def.OLC_TRIGGER {
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLCEdit.Mode == def.OLC_MOBILE {
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLCEdit.Mode == def.OLC_QUEST {
		WriteToPlayer(player, "Not available yet (WIP).")
	} else if player.OLCEdit.Mode == def.OLC_SECTOR {
		WriteToPlayer(player, "Not available yet (WIP).")
	}

}

func CmdDig(player *glob.PlayerData, input string) {
	if player.Location.RoomLink == nil {
		WriteToPlayer(player, "You need to be in a room, to dig.")
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

	if player.OLCSettings.OlcShowCodes {
		bytes, err = player.Connection.Desc.Write([]byte(text + "\r\n"))
	} else {
		bytes, err = player.Connection.Desc.Write([]byte(ANSIColor(text) + "\r\n"))
	}
	player.Connection.BytesOut += bytes
	trackBytesOut(player.Connection)

	DescWriteError(player.Connection, err)
}
