package support

import (
	"fmt"

	"../def"
	"../glob"
)

func OLCObject(player *glob.PlayerData,
	input string, command string, cmdB string, cmdl string, cmdBl string,
	argTwoThrough string, argThreeThrough string) {
	if cmdl == "done" {
		player.OLCEdit.Mode = def.OLC_NONE
		WriteToPlayer(player, "Exiting OLC.")
		player.OLCEdit.Active = false
		return
	} else if cmdl == "create" {
		sector, id, err := ParseVnum(player, argThreeThrough)
		if err == false {
			glob.SectorsList[sector].Objects[id] = CreateObject()
		} else {
			objs := glob.SectorsList[sector].Objects

			found := 0
			for x := 0; ; x++ {
				if objs[x] != nil && objs[x].Valid == false {
					found = x
					break
				}
				if objs[x] == nil {
					found = x
					break
				}
			}
			objs[found] = CreateObject()
			objs[found].ID = found
		}
	} else if cmdl == "" {
		if player.OLCEdit.Object.ID != 0 {
			objId := player.OLCEdit.Object.ID
			objSec := player.OLCEdit.Sector
			obj := glob.SectorsList[objSec].Objects[objId]
			buf := fmt.Sprintf("Name: %v\r\nID: %v\r\nDesc: %v\r\n", obj.Name, obj.ID, obj.Description)
			WriteToPlayer(player, buf)
		} else {
			WriteToPlayer(player, "No object selected.")
		}
	} else {
		sector, id, err := ParseVnum(player, argTwoThrough)
		obj, found := GetObjectFromID(sector, id)
		if err == false && found == true {
			obj = player.OLCEdit.Object.ObjectLink
			buf := fmt.Sprintf("Name: %v, ID, %v", obj.Name, obj.ID)
			WriteToPlayer(player, buf)
		} else {
			player.OLCEdit.Object.ObjectLink, found = GetObjectFromID(sector, id)
			player.OLCEdit.Object.Sector = sector
			player.OLCEdit.Object.ID = id
			WriteToPlayer(player, "Object selected.")
		}
	}
}
