package support

import (
	"../glob"
)

func LocationDataFromID(sector int, id int) (glob.LocationData, bool) {

	room := glob.SectorsList[sector].Rooms[id]
	if glob.SectorsList[sector].Rooms[id] != nil {
		loc := glob.LocationData{Sector: sector, ID: id, RoomLink: room}
		return loc, true
	}
	return glob.LocationData{}, false
}

func CreateSector() *glob.SectorData {

	sector := glob.SectorData{
		ID:          glob.SectorsListEnd,
		Fingerprint: "",
		Name:        "",
		Area:        "",
		Description: "",
		Rooms:       make(map[int]*glob.RoomData),
		Objects:     make(map[int]*glob.ObjectData),
		Dirty:       false,

		Valid: true,
	}

	return &sector
}

func CreateRoom() *glob.RoomData {
	room := glob.RoomData{
		Name:        "new room",
		Description: "",
		Players:     make(map[string]*glob.PlayerData),
		Exits:       make(map[string]*glob.ExitData),
		PermObjects: make(map[string]*glob.ObjectData),
		Objects:     make(map[string]*glob.ObjectData),
		Valid:       true,
	}

	return &room
}

func CreateExit() *glob.ExitData {
	exit := glob.ExitData{
		Valid: true,
	}

	exit.Door = &glob.DoorData{Valid: true}
	return &exit
}

func CreateObject() *glob.ObjectData {
	obj := glob.ObjectData{
		Valid: true,
	}
	return &obj
}
