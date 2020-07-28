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
		if glob.SectorsList[x].Valid {
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

	err = ioutil.WriteFile(fileName, []byte(outbuf.String()), 0644)

	if err != nil {
		CheckError("WriteSector: WriteFile", err, def.ERROR_NONFATAL)
		return false
	}

	buf := fmt.Sprintf("Wrote %v, %v.", fileName, ScaleBytes(len(outbuf.String())))
	mlog.Write(buf)
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

			for x, _ := range sector.Rooms {
				room := sector.Rooms[x]
				room.Players = make(map[string]*glob.PlayerData)
			}

			prefix := ""
			if sector.Fingerprint == "" {
				if sector.Name != "" {
					prefix = sector.Name + "-"
				}
				mlog.Write(sector.Name + " assigned fingerprint.")
				sector.Fingerprint = MakeFingerprint(prefix)
			}
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
		Fingerprint: "",
		Name:        "",
		Area:        "",
		Description: "",
		Rooms:       nil,

		Valid: true,
	}

	if sector.Rooms == nil {
		sector.Rooms = make(map[int]*glob.RoomData)
	}

	for _, value := range sector.Rooms {
		value.Players = make(map[string]*glob.PlayerData)
	}

	return &sector
}

func CreateRoom() *glob.RoomData {
	room := glob.RoomData{
		Name:        "new room",
		Description: "",

		Valid: true,
	}
	room.Players = make(map[string]*glob.PlayerData)
	room.Exits = make(map[string]*glob.ExitData)

	return &room
}

func CreateExit() *glob.ExitData {
	exit := glob.ExitData{
		Valid: true,
	}

	return &exit
}
