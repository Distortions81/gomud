package support

import (
	"fmt"
	"strconv"
	"strings"

	"../def"
	"../glob"
)

func OLCSector(player *glob.PlayerData,
	input string, command string, cmdB string, cmdl string, cmdBl string,
	argTwoThrough string, argThreeThrough string) {

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
