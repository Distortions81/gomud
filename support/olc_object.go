package support

import (
	"fmt"

	"../def"
	"../glob"
)

func OLCObject(player *glob.PlayerData,
	input string, command string, cmdB string, cmdl string, cmdBl string,
	argTwoThrough string, argThreeThrough string) {

	found := 0
	isFound := false
	sector := 0
	id := 0
	wasErr := false
	var obj *glob.ObjectData

	if cmdl == "done" {
		player.OLCEdit.Mode = def.OLC_NONE
		WriteToPlayer(player, "Exiting OLC.")
		player.OLCEdit.Active = false
		return
	} else if cmdl == "create" {
		sector, id, wasErr = ParseVnum(player, argThreeThrough)
		if wasErr == false && sector > 0 && id > 0 {
			glob.SectorsList[sector].Objects[id] = CreateObject()
		} else {
			sector = player.Location.Sector
			objs := glob.SectorsList[sector].Objects

			for x := 1; ; x++ {
				if objs[x] != nil && objs[x].Valid == false {
					found = x
					break
				}
				if objs[x] == nil {
					found = x
					break
				}
			}
		}

		/* Make object map, if it doesn't exist yet */
		if glob.SectorsList[sector].Objects == nil {
			glob.SectorsList[sector].Objects = make(map[int]*glob.ObjectData)
		}
		glob.SectorsList[sector].Objects[found] = CreateObject()
		glob.SectorsList[sector].Objects[found].ID = found

		buf := fmt.Sprintf("Object %v:%v created.", sector, found)
		WriteToPlayer(player, buf)
		glob.SectorsList[player.OLCEdit.Object.Sector].Dirty = true

	} else if cmdl == "name" {
		if player.OLCEdit.Object.ObjectLink == nil {
			WriteToPlayer(player, "No selected object")
		} else {
			player.OLCEdit.Object.ObjectLink.Name = argTwoThrough
			WriteToPlayer(player, "Name set.")
			glob.SectorsList[player.OLCEdit.Object.Sector].Dirty = true
		}

	} else if cmdl == "" {
		if player.OLCEdit.Object.ID != 0 {
			objId := player.OLCEdit.Object.ID
			objSec := player.OLCEdit.Object.Sector
			obj := glob.SectorsList[objSec].Objects[objId]
			buf := fmt.Sprintf("Name: %v\r\nID: %v\r\nDesc: %v\r\n", obj.Name, obj.ID, obj.Description)
			WriteToPlayer(player, buf)
		} else {
			WriteToPlayer(player, "No object selected.")
		}
	} else {
		sector, id, wasErr = ParseVnum(player, input)
		obj, isFound = GetObjectFromID(sector, id)
		if wasErr == false && isFound == true && obj != nil {
			player.OLCEdit.Object.ObjectLink, isFound = GetObjectFromID(sector, id)
			player.OLCEdit.Object.Sector = sector
			player.OLCEdit.Object.ID = id
			WriteToPlayer(player, "Object selected.")
		} else {
			WriteToPlayer(player, "Didn't find a valid object.")
		}
	}
}
