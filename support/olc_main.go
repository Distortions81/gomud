package support

import (
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

	if strings.EqualFold(command, "olc") {
		WriteToPlayer(player, "You are in noPrefix mode, you don't need to type OLC before commands.\r\ntype settings to turn this mode off... Or type STOP to exit.")
		return
	} else if strings.EqualFold("stop", input) {
		player.OLCEdit.Active = false
		WriteToPlayer(player, "Exiting OLC!")
		return
	} else if cmdl == "settings" {
		OLCConfig(player, input, command, cmdB, cmdl, cmdBl, argTwoThrough, argThreeThrough)
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
			player.OLCEdit.Room.Sector = player.Location.Sector
			player.OLCEdit.Room.ID = player.Location.ID
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
		OLCSector(player, input, command, cmdB, cmdl, cmdBl, argTwoThrough, argThreeThrough)

	}
}
