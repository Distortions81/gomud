package support

import (
	"fmt"
	"log"
	"os"
	"strings"

	"../def"
	"../glob"
	"../mlog"
)

//Hard coded aliases
func CmdNorth(player *glob.PlayerData, input string) {
	CmdGo(player, "north")
}
func CmdSouth(player *glob.PlayerData, input string) {
	CmdGo(player, "south")
}
func CmdEast(player *glob.PlayerData, input string) {
	CmdGo(player, "east")
}
func CmdWest(player *glob.PlayerData, input string) {
	CmdGo(player, "west")
}
func CmdUp(player *glob.PlayerData, input string) {
	CmdGo(player, "up")
}
func CmdDown(player *glob.PlayerData, input string) {
	CmdGo(player, "down")
}

func CmdRecall(player *glob.PlayerData, input string) {

	if input == "set" {
		//I love how elegant this is
		player.Recall = player.Location
		WriteToPlayer(player, "Recall set!")
		return
	}

	if player.Location.ID == player.Recall.ID && player.Location.Sector == player.Recall.Sector {
		WriteToPlayer(player, "You try to recall, but strain to remember... wait, does this place look familiar?")
		return
	} else {
		WriteToPlayer(player, "You recall, and are suddenly transported, in a bright blue {Cflash{x of {Ylight.")
	}

	WriteToRoom(player, fmt.Sprintf("%v {Kvanishes{x with a bright blue {Cflash{x of {Ylight.", player.Name))
	if player.Recall.Sector != 0 || player.Recall.ID != 0 {
		PlayerToRoom(player, player.Recall.Sector, player.Recall.ID)
	} else {
		PlayerToRoom(player, def.PLAYER_START_SECTOR, def.PLAYER_START_ROOM)
	}
	WriteToRoom(player, fmt.Sprintf("%v suddenly {mappears{x, with a bright blue {Cflash{x of {Ylight.", player.Name))
	player.Dirty = true
	CmdLook(player, "")
}

func PlayerToRoom(player *glob.PlayerData, sectorID int, roomID int) {

	if player == nil && !player.Valid {
		return
	}
	//Remove player from room, if they are in one
	if player.Location.RoomLink != nil {
		room := player.Location.RoomLink
		delete(room.Players, player.Fingerprint)
	}

	//Add player to room, add error handling
	//Automatically generate "players" map if it doesn't exist
	if glob.SectorsList[sectorID].Valid &&
		glob.SectorsList[sectorID].Rooms != nil {

		room := glob.SectorsList[sectorID].Rooms[roomID]
		if room != nil {
			room.Players[player.Fingerprint] = player
		}

		player.Location.RoomLink = glob.SectorsList[sectorID].Rooms[roomID]
		player.Location.Sector = sectorID
		player.Location.ID = roomID
		player.Dirty = true

		if player.OLCEdit.Active &&
			player.OLCEdit.Mode == def.OLC_ROOM &&
			player.OLCSettings.OLCRoomFollow {

			player.OLCEdit.Room.ID = player.Location.ID
			player.OLCEdit.Room.Sector = player.Location.Sector
			player.OLCEdit.Room.RoomLink = player.Location.RoomLink
		}

	} else {
		mlog.Write("PlayerToRoom: That sector or room is not valid.")
		mlog.Write(fmt.Sprintf("Sector: %v, Room: %v, Player: %v", sectorID, roomID, player.Name))
		if sectorID != def.PLAYER_START_SECTOR && roomID != def.PLAYER_START_ROOM {
			PlayerToRoom(player, def.PLAYER_START_SECTOR, def.PLAYER_START_ROOM)
		} else {
			log.Println("Default room/sector not found. Quitting.")
			os.Exit(1)
		}

	}
}

func movePlayerExit(player *glob.PlayerData, arg string, exit *glob.ExitData) {
	WriteToPlayer(player, "You go "+arg+".")
	WriteToRoom(player, player.Name+" goes "+arg+".")
	PlayerToRoom(player, exit.ToRoom.Sector, exit.ToRoom.ID)

	WriteToRoom(player, player.Name+" arrives.")
	player.Dirty = true

	CmdLook(player, "")
}

func CmdGo(player *glob.PlayerData, args string) {
	found := false

	if args == "" {
		WriteToPlayer(player, "Go where?")
		return
	}
	for exitName, exit := range player.Location.RoomLink.Exits {
		if strings.HasPrefix(strings.ToLower(exitName), strings.ToLower(args)) {
			found = true

			if IsStandardDirection(exitName) {
				WriteToPlayer(player, "You go "+exitName+".")
				WriteToRoom(player, player.Name+" goes "+exitName+".")
			} else {
				WriteToPlayer(player, "You '"+exitName+"'.")
				WriteToRoom(player, player.Name+" went '"+exitName+"'.")
			}

			PlayerToRoom(player, exit.ToRoom.Sector, exit.ToRoom.ID)

			if IsStandardDirection(exitName) {
				WriteToRoom(player, player.Name+" arrives from '"+GetStandardDirectionMirror(exitName)+"'.")
			} else {
				WriteToRoom(player, player.Name+" went '"+exitName+"'.")
			}

			CmdLook(player, "")
			player.Dirty = true
			return
		}
	}
	if !found {

		var exitsList []string
		for exitName, _ := range player.Location.RoomLink.Exits {
			exitsList = append(exitsList, exitName)
		}
		result, _ := FindClosestMatch(exitsList, args)
		if result != "" {

		} else {
			WriteToPlayer(player, "Go where?")
		}
	}
}
