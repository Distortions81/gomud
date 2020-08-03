package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"../def"
	"../glob"
	"../mlog"
)

func LocationDataFromID(sector int, id int) (glob.LocationData, bool) {

	room := glob.SectorsList[sector].Rooms[id]
	if glob.SectorsList[sector].Rooms[id] != nil {
		loc := glob.LocationData{Sector: sector, ID: id, RoomLink: room}
		return loc, true
	}
	return glob.LocationData{}, false
}

func ReadSectorList() {

	files, err := ioutil.ReadDir(def.DATA_DIR + def.SECTOR_DIR)
	if err != nil {
		CheckError("ReadSectorList:", err, def.ERROR_NONFATAL)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), def.FILE_SUFFIX) {
			sector := ReadSector(file.Name())
			sector.Valid = true
			if sector != nil {
				if glob.SectorsList[sector.ID].Valid {
					buf := fmt.Sprintf("%v has same sector ID as %v! Skipping!", sector.Name, glob.SectorsList[sector.ID].Name)
					mlog.Write(buf)
				} else {
					glob.SectorsList[sector.ID] = *sector
					glob.SectorsListEnd++
				}
			} else {
				mlog.Write("Invalid sector file: " + file.Name())
			}
		}
	}
}

func ReloadSector() {

	//reload sector, handle future load handles, regen player pointers
}

func WriteSectorList() {
	for x := 1; x <= glob.SectorsListEnd; x++ {
		if glob.SectorsList[x].Valid && glob.SectorsList[x].Name != "" {
			if glob.SectorsList[x].Fingerprint == "" {
				glob.SectorsList[x].Fingerprint = MakeFingerprint(glob.SectorsList[x].Name)
			}
			glob.SectorsList[x].ID = x
			WriteSector(&glob.SectorsList[x])
		}
	}
}

func WriteSector(sector *glob.SectorData) bool {

	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	sector.Version = def.SECTOR_VERSION

	if sector == nil && !sector.Valid {
		return false
	}

	/*Write room count*/
	numRooms := 0
	for _, _ = range sector.Rooms {
		numRooms++
	}
	sector.NumRooms = numRooms

	fileName := def.DATA_DIR + def.SECTOR_DIR + def.SECTOR_PREFIX + fmt.Sprintf("%v", sector.ID) + def.FILE_SUFFIX

	if err := enc.Encode(&sector); err != nil {
		CheckError("WriteSector: enc.Encode", err, def.ERROR_NONFATAL)
		return false
	}
	_, err := os.Create(fileName)

	if err != nil {
		CheckError("WriteSector: os.Create", err, def.ERROR_NONFATAL)
		return false
	}

	//Async write
	go func(outbuf bytes.Buffer) {
		glob.WriteSectorsLock.Lock()
		defer glob.WriteSectorsLock.Unlock()

		err = ioutil.WriteFile(fileName, []byte(outbuf.String()), 0644)

		if err != nil {
			CheckError("WriteSector: WriteFile", err, def.ERROR_NONFATAL)
		}

		buf := fmt.Sprintf("Wrote %v, %v.", fileName, ScaleBytes(len(outbuf.String())))
		mlog.Write(buf)
	}(*outbuf)

	sector.Dirty = false
	return true
}

func ReadSector(name string) *glob.SectorData {

	_, err := os.Stat(def.DATA_DIR + def.SECTOR_DIR + name)
	notfound := os.IsNotExist(err)

	if notfound {
		CheckError("ReadSector: os.Stat", err, def.ERROR_NONFATAL)
		return nil

	} else {

		file, err := ioutil.ReadFile(def.DATA_DIR + def.SECTOR_DIR + name)

		if file != nil && err == nil {
			sector := CreateSector()

			err := json.Unmarshal([]byte(file), &sector)
			if err != nil {
				CheckError("ReadSector: Unmashal", err, def.ERROR_NONFATAL)
				return nil
			}
			if sector.ID > glob.SectorsListEnd {
				glob.SectorsListEnd = sector.ID
			}
			numRooms := 0
			for x, _ := range sector.Rooms {
				numRooms++
				room := sector.Rooms[x]

				if room.Players == nil {
					room.Players = make(map[string]*glob.PlayerData)
				}
				if room.Exits == nil {
					room.Exits = make(map[string]*glob.ExitData)
				}
				if room.PermObjects == nil {
					room.PermObjects = make(map[string]*glob.ObjectData)
				}
				if room.Objects == nil {
					room.Objects = make(map[string]*glob.ObjectData)
				}

				for x, _ := range room.Exits {
					exit := room.Exits[x]
					if exit.Door != nil {
						exit.Door.Valid = true
					}
					exit.RoomP = room
					exit.Valid = true
				}
				for x, _ := range room.PermObjects {
					pObj := room.PermObjects[x]
					pObj.Sector = sector.ID
					pObj.InRoom = room
					pObj.Valid = true
				}
				for x, _ := range room.Objects {
					obj := room.Objects[x]
					obj.Sector = sector.ID
					obj.InRoom = room
					obj.Valid = true
				}
				room.SectorP = sector
				room.Valid = true
			}
			sector.NumRooms = numRooms

			for x, _ := range sector.Objects {
				obj := sector.Objects[x]
				obj.Valid = true
			}

			prefix := ""
			if sector.Fingerprint == "" {
				if sector.Name != "" {
					prefix = sector.Name
				}
				mlog.Write(sector.Name + " assigned fingerprint.")
				sector.Fingerprint = MakeFingerprint(prefix)
			}
			sector.Valid = true
			mlog.Write("ReadSector: " + sector.Name)
			return sector
		} else {
			CheckError("ReadSector: ReadFile", err, def.ERROR_NONFATAL)
			return nil
		}

	}
}

func CreateSector() *glob.SectorData {

	sector := glob.SectorData{
		ID:          glob.SectorsListEnd,
		Fingerprint: "",
		Name:        "",
		Area:        "",
		Description: "",
		Rooms:       make(map[int]*glob.RoomData),
		Objects:     make(map[string]*glob.ObjectData),
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
