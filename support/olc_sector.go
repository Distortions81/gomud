package support

import (
	"fmt"
	"os"
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
		if sector == nil {
			sector = &glob.SectorsList[player.Location.Sector]
		}
		if cmdl == "done" {
			player.OLCEdit.Mode = def.OLC_NONE
			WriteToPlayer(player, "Exiting OLC.")
			player.OLCEdit.Active = false
			sector.Dirty = true
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
			if sector.Valid {
				sector.Valid = false
				sector.Dirty = true
				WriteToPlayer(player, "Sector set to invalid / inactive.")
			} else {
				sector.Valid = true
				sector.Dirty = true
				WriteToPlayer(player, "Sector set to valid / active.")

			}
		} else if strings.EqualFold(cmdl, "delete") {
			WriteToPlayer(player, "If you are COMPLETELY CERTAIN you want to PERMENATELY DLETE the ROOMS AND OBJECTS in this ENTIRE SECTOR....\r\nType: confirm-delete-sector")
			return
		} else if strings.EqualFold(cmdl, "confirm-delete-sector") {
			fileName := def.DATA_DIR + def.SECTOR_DIR + def.SECTOR_PREFIX + fmt.Sprintf("%v", sector.ID) + def.FILE_SUFFIX
			os.Remove(fileName)

			if sector.ID == glob.SectorsListEnd {
				glob.SectorsListEnd--
			}
			sector.Name = ""
			sector.Description = ""
			sector.Area = ""
			sector.Rooms = nil
			sector.Objects = nil
			sector.Valid = false
			sector.Dirty = false
			WriteToPlayer(player, "Sector is deleted.")

			return
		} else if strings.EqualFold(cmdl, "id") {
			idNum, err := strconv.Atoi(argTwoThrough)

			if err == nil && idNum > 0 {
				oldSector := player.OLCEdit.Sector
				oldSectorData := glob.SectorsList[oldSector]
				WriteToPlayer(player, "Sector moved.")
				player.OLCEdit.Sector = idNum
				if glob.SectorsListEnd < idNum {
					glob.SectorsListEnd = idNum
				}
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
			sector.Dirty = true
			if sector.Fingerprint == "" {
				sector.Fingerprint = MakeFingerprint(sector.Name)
			}
			WriteToPlayer(player, "Name set.")
		} else if strings.EqualFold(cmdl, "desc") || strings.EqualFold(cmdl, "description") {
			//Todo, editor
			sector.Description = argTwoThrough
			sector.Dirty = true
			WriteToPlayer(player, "Description set.")
		} else if strings.EqualFold(cmdl, "area") {
			sector.Area = argTwoThrough
			sector.Dirty = true
			WriteToPlayer(player, "Area set.")
		} else if strings.EqualFold(cmdl, "create") {

			newSector := CreateSector()
			newSector.Valid = true
			glob.SectorsListEnd++
			glob.SectorsList[glob.SectorsListEnd] = *newSector
			sector.ID = glob.SectorsListEnd
			player.OLCEdit.Sector = glob.SectorsListEnd
			sector.Dirty = true
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
