package support

import (
	"fmt"
	"strconv"
	"strings"

	"../def"
	"../glob"
)

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
			player.OLCEdit.Sector = player.Location.Sector
			player.OLCEdit.ID = player.Location.ID
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
		OLCExits(player, input, command, cmdB, cmdl, cmdBl, argTwoThrough, argThreeThrough)
	} else if player.OLCEdit.Mode == def.OLC_ROOM {
		OLCRoom(player, input, command, cmdB, cmdl, cmdBl, argTwoThrough, argThreeThrough)
	} else if player.OLCEdit.Mode == def.OLC_OBJECT {
		OLCObject(player, input, command, cmdB, cmdl, cmdBl, argTwoThrough, argThreeThrough)
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
