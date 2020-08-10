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

	sec, err := strconv.Atoi(cmdl)

	if err == nil && sec > 0 {
		secData := &glob.SectorsList[sec]
		if secData != nil {
			WriteToPlayer(player, fmt.Sprintf("Sector %v selected.", sec))
			player.OLCEdit.Sector = sec
			sector = secData
			CmdOLC(player, "")
			return
		}
	} else if cmdl == "" {
		buf := fmt.Sprintf("Name: %v\r\nID %v\r\nFingerprint: %v\r\nDescription: %v\r\nArea: %v\r\nRoom count: %v\r\nVALID: %v",
			sector.Name, sector.ID, sector.Fingerprint, sector.Description, sector.Area, sector.NumRooms, sector.Valid)
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
		} else if strings.EqualFold(cmdl, "valid") {
			if glob.SectorsList[player.OLCEdit.Sector].Valid {
				glob.SectorsList[player.OLCEdit.Sector].Valid = false
				WriteToPlayer(player, "Sector set to invalid / inactive.")
			} else {
				glob.SectorsList[player.OLCEdit.Sector].Valid = true
				WriteToPlayer(player, "Sector set to valid / active.")

			}
		} else if strings.EqualFold(cmdl, "delete") {
			WriteToPlayer(player, "If you are COMPLETELY CERTAIN you want to PERMENATELY DLETE the ROOMS AND OBJECTS in this ENTIRE SECTOR....\r\nType: confirm-delete-sector")
			return
		} else if strings.EqualFold(cmdl, "confirm-delete-sector") {
			glob.SectorsList[player.OLCEdit.Sector].Rooms = nil
			glob.SectorsList[player.OLCEdit.Sector].Objects = nil
			glob.SectorsList[player.OLCEdit.Sector].Valid = false
			glob.SectorsList[player.OLCEdit.Sector].Dirty = true
			WriteToPlayer(player, "Rooms and objects in sector cleard, and sector set to invalid / inactive.")
			return
		} else if strings.EqualFold(cmdl, "id") {
			idNum, err := strconv.Atoi(argTwoThrough)

			if err == nil && idNum > 0 {
				oldSector := player.OLCEdit.Sector
				oldSectorData := glob.SectorsList[oldSector]
				WriteToPlayer(player, "Sector moved.")
				player.OLCEdit.Sector = idNum
				glob.SectorsList[idNum] = oldSectorData
				glob.SectorsList[oldSector].Valid = false
				glob.SectorsList[idNum].Dirty = true
				glob.SectorsList[oldSector].Dirty = true
			} else {
				WriteToPlayer(player, "sector <sector number>")
			}
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
