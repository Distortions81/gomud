package support

import (
	"strconv"
	"strings"

	"../glob"
)

func getLocationFromString(player *glob.PlayerData, input string) (glob.LocationData, bool) {
	loc := strings.Split(input, ":")
	locLen := len(loc)

	if locLen >= 1 {
		sector := 0
		id := 0
		var erra error
		var errb error
		if locLen == 1 {
			sector = player.Location.Sector
			id, erra = strconv.Atoi((loc[0]))
		} else if locLen == 2 {
			sector, errb = strconv.Atoi(loc[0])
			id, erra = strconv.Atoi((loc[1]))
		}

		if erra == nil && errb == nil {
			room, roomFound := LocationDataFromID(sector, id)
			if roomFound {
				return room, true
			}
		}
	}
	return glob.LocationData{}, false
}

func doDig(player *glob.PlayerData, rooms map[int]*glob.RoomData, found int, dir string, sector int) {
	room := rooms[found]

	room.Valid = true
	room.Name = "new room"

	room.Exits[GetStandardDirectionMirror(dir)] = CreateExit()
	room.Exits[GetStandardDirectionMirror(dir)].ToRoom = player.Location

	player.Location.RoomLink.Exits[strings.Title(dir)] = CreateExit()
	player.Location.RoomLink.Exits[strings.Title(dir)].ToRoom.ID = found
	player.Location.RoomLink.Exits[strings.Title(dir)].ToRoom.Sector = sector

	CmdGo(player, dir)
	glob.SectorsList[sector].Dirty = true //Autosave
}

func doDigCustom(player *glob.PlayerData, rooms map[int]*glob.RoomData, found int, dirOne string, dirTwo string, sector int) {
	room := rooms[found]

	room.Valid = true
	room.Name = "new room"

	room.Exits[dirTwo] = CreateExit()
	room.Exits[dirTwo].ToRoom = player.Location

	player.Location.RoomLink.Exits[dirOne] = CreateExit()
	player.Location.RoomLink.Exits[dirOne].ToRoom.ID = found
	player.Location.RoomLink.Exits[dirOne].ToRoom.Sector = sector

	CmdGo(player, dirOne)
	glob.SectorsList[sector].Dirty = true //Autosave
}

func WriteToBuilder(player *glob.PlayerData, text string) {

	var bytes int
	var err error

	if player == nil || !player.Valid || player.Connection == nil || !player.Connection.Valid {
		return
	}

	if player.OLCSettings.OLCShowCodes {
		bytes, err = player.Connection.Desc.Write([]byte(text + "\r\n"))
	} else {
		bytes, err = player.Connection.Desc.Write([]byte(ANSIColor(text) + "\r\n"))
	}
	player.Connection.BytesOut += bytes
	trackBytesOut(player.Connection)

	DescWriteError(player.Connection, err)
}
