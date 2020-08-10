package support

import (
	"fmt"
	"strconv"
	"strings"

	"../glob"
)

func CmdRoomList(player *glob.PlayerData, args string) {
	sec := glob.SectorsList[player.OLCEdit.Room.Sector]
	rooms := sec.Rooms

	buf := ""
	lastS := -1
	for s, r := range rooms {
		if r.Valid {
			if (s + 1) != lastS {
				buf = buf + fmt.Sprintf("%v, ", s)
			}
			lastS = s
		}
	}
	WriteToPlayer(player, buf)
}

func CmdSectorList(player *glob.PlayerData, args string) {
	for x := 1; x <= glob.SectorsListEnd; x++ {
		sec := glob.SectorsList[x]
		if sec.Valid {
			buf := fmt.Sprintf("Sector ID: %4v, Name: %-20v, Area: %-20v, Rooms: %-7v", x, sec.Name, sec.Area, len(sec.Rooms))
			WriteToPlayer(player, buf)
		}
	}
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

func CmdDig(player *glob.PlayerData, input string) {
	if player.Location.RoomLink == nil {
		WriteToPlayer(player, "You need to be in a room, to dig.")
		return
	}

	command, _ := SplitArgsTwo(input, " ")
	dirOne, dirTwo := SplitArgsTwo(input, ":")

	sector := player.Location.Sector

	if player.Location.RoomLink.Exits[strings.Title(dirOne)] != nil {
		WriteToPlayer(player, "That exit is already occupied.")
		return
	}

	rooms := glob.SectorsList[sector].Rooms

	//Find first available slot
	found := 0
	for x := 0; ; x = x + 1 {
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
