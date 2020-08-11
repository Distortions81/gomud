package support

import (
	"fmt"

	"../def"
	"../glob"
)

func CmdChat(player *glob.PlayerData, args string) {
	if len(args) > 0 {
		out := fmt.Sprintf("%s chats: %s", player.Name, args)
		us := fmt.Sprintf("You chat: %s", args)

		WriteToOthers(player, out)
		WriteToPlayer(player, us)
	} else {
		WriteToPlayer(player, "But, what do you want to say?")
	}
}

func CmdSay(player *glob.PlayerData, args string) {
	if len(args) > 0 {
		out := fmt.Sprintf("%s says: %s", player.Name, args)
		us := fmt.Sprintf("You say: %s", args)

		WriteToRoom(player, out)
		WriteToPlayer(player, us)
	} else {
		WriteToPlayer(player, "But, what do you want to say?")
	}
}

func CmdEmote(player *glob.PlayerData, args string) {
	if len(args) > 0 {
		out := fmt.Sprintf("%s %s", player.Name, args)

		WriteToRoom(player, out)
		WriteToPlayer(player, out)
	} else {
		WriteToPlayer(player, "But, what do you want to say?")
	}
}

func CmdLook(player *glob.PlayerData, args string) {

	err := true
	sector := glob.SectorsList[player.Location.Sector]
	if sector.Valid {
		if sector.Rooms[player.Location.ID] != nil && sector.Rooms[player.Location.ID].Valid {
			room := sector.Rooms[player.Location.ID]
			roomName := room.Name
			roomDesc := room.Description
			buf := fmt.Sprintf("%s:\r\n%s", roomName, roomDesc)
			WriteToPlayer(player, buf)
			err = false
		}

		if player.Location.RoomLink != nil {
			exits := "["
			l := len(player.Location.RoomLink.Exits)
			x := 0
			for name, _ := range player.Location.RoomLink.Exits {
				x++
				if name != "" {
					exits = exits + name
					if x < l {
						exits = exits + ", "
					}
				}
			}
			if exits == "[" {
				exits = exits + " None... "
			}
			exits = exits + "]"

			WriteToPlayer(player, "\r\nExits: "+exits)

			names := ""
			unlinked := ""
			for _, target := range player.Location.RoomLink.Players {
				if target != nil && target != player {
					if target.Connection != nil && target.Connection.Valid == false {
						unlinked = " (lost connection)"
					}
					names = names + fmt.Sprintf("%s is here.%s", target.Name, unlinked)
				}
			}
			//Newline if there are players here.
			if names != "" {
				WriteToPlayer(player, "\r\n"+names)
			}
		} else {
			err = true
		}
	}
	if err {
		WriteToPlayer(player, "You are floating in the VOID...")
		PlayerToRoom(player, def.PLAYER_START_SECTOR, def.PLAYER_START_ROOM)
	}

}
