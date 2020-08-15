package support

import (
	"fmt"

	"../def"
	"../glob"
)

func OLCReset(player *glob.PlayerData,
	input string, command string, cmdB string, cmdl string, cmdBl string,
	argTwoThrough string, argThreeThrough string) {

	found := 0
	isFound := false
	sector := player.OLCEdit.Reset.Sector
	id := player.OLCEdit.Reset.ID
	wasErr := false

	if cmdl == "done" {
		/* Exit editor */
		player.OLCEdit.Mode = def.OLC_NONE
		WriteToPlayer(player, "Exiting OLC.")
		player.OLCEdit.Active = false
		return

		/* Create new reset */
	} else if cmdl == "create" {

		/* Check if sector/id is specified */
		sector, id, wasErr = ParseVnum(player, argThreeThrough)
		if wasErr == false && sector > 0 && id > 0 {
			if glob.SectorsList[sector].Rooms[id] == nil {
				WriteToPlayer(player, "That room does not exist!")
				return
			}
		} else {
			if player.OLCEdit.Reset.ID > 0 && player.OLCEdit.Reset.Sector > 0 {
				sector = player.OLCEdit.Reset.Sector
				id = player.OLCEdit.Reset.ID
			} else if player.OLCEdit.Room.ID > 0 && player.OLCEdit.Room.Sector > 0 {
				sector = player.OLCEdit.Room.Sector
				id = player.OLCEdit.Room.ID
			} else {
				sector = player.Location.Sector
				id = player.Location.ID
				WriteToPlayer(player, "No location specified, defaulting to current room.")
			}
		}

		/* Make resets map, if it doesn't exist yet */
		if glob.SectorsList[sector].Rooms[id].Resets == nil {
			glob.SectorsList[sector].Rooms[id].Resets = make(map[int]*glob.ResetsData)
		}

		bufe := fmt.Sprintf("sector: %v, id: %v found: %v, err: %v, input: %v",
			sector, id, isFound, wasErr, input)
		WriteToPlayer(player, bufe)

		for x := 1; ; x++ {
			r := glob.SectorsList[sector].Rooms[id].Resets[x]
			if r == nil || r.Valid == false {
				found = x
				break
			}
		}

		glob.SectorsList[sector].Rooms[id].Resets[found] = CreateReset()
		glob.SectorsList[sector].Rooms[id].Resets[found].Number = found

		roomLoc, rFound := LocationDataFromID(sector, id)
		if rFound {
			player.OLCEdit.Reset.RoomLink = roomLoc.RoomLink
			glob.SectorsList[sector].Rooms[id].Resets[found].Location = roomLoc
		} else {
			WriteToPlayer(player, "We can't add resets to a room that does not exist!")
			return
		}

		buf := fmt.Sprintf("Reset #%v, for room %v:%v created.", found, sector, id)
		WriteToPlayer(player, buf)
		glob.SectorsList[sector].Dirty = true

	} else if cmdl == "name" {
		if player.OLCEdit.Reset.ResetLink == nil {
			WriteToPlayer(player, "No selected reset")
		} else {
			player.OLCEdit.Reset.ResetLink.Name = argTwoThrough
			WriteToPlayer(player, "Name set.")
			glob.SectorsList[player.OLCEdit.Reset.Sector].Dirty = true
		}
	} else if cmdl == "" {
		sector := player.OLCEdit.Reset.Sector
		id := player.OLCEdit.Reset.ID
		buf := ""
		if sector != 0 && id != 0 && glob.SectorsList[sector].Rooms[id] != nil {
			for _, res := range glob.SectorsList[sector].Rooms[id].Resets {
				buf = fmt.Sprintf("Name: %v\r\nID: %v\r\nType: %v\r\n",
					res.Name, res.Number, res.Type)
				WriteToPlayer(player, buf)
			}
		} else {
			WriteToPlayer(player, "No room selected!")
			return
		}
		if buf == "" {
			WriteToPlayer(player, "There are no resets in this room.")
			return
		}
	} else {
		sector, id, wasErr = ParseVnum(player, input)
		_, isFound = GetRoomFromID(sector, id)
		if wasErr == false && isFound == true {
			var rLoc glob.LocationData
			rLoc, isFound = LocationDataFromID(sector, id)
			if isFound {
				player.OLCEdit.Reset.Sector = rLoc.Sector
				player.OLCEdit.Reset.ID = rLoc.ID
				player.OLCEdit.Reset.RoomLink = rLoc.RoomLink
				WriteToPlayer(player, "Resets selected.")
				CmdOLC(player, "")
				return
			}
		} else {
			WriteToPlayer(player, "Didn't find a valid room.")
			buf := fmt.Sprintf("sector: %v, id: %v found: %v, err: %v, input: %v",
				sector, id, isFound, wasErr, input)
			WriteToPlayer(player, buf)
		}
	}
}
