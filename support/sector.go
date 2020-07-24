package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"../def"
	"../glob"
)

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
					log.Println(buf)
				} else {
					glob.SectorsList[sector.ID] = *sector
				}
			} else {
				log.Println("Invalid sector file: " + file.Name())
			}
		}
	}
}

func ReloadSector() {

	//reload sector, handle future load handles, fix player pointers
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
	log.Println(buf)
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

			if sector.Rooms == nil {
				sector.Rooms = make(map[int]glob.RoomData)
			}
			for key, _ := range sector.Rooms {
				p := make(map[string]*glob.PlayerData)
				room := sector.Rooms[key]
				room.Players = p
			}

			prefix := ""
			if sector.Fingerprint == "" {
				if sector.Name != "" {
					prefix = sector.Name + "-"
				}
				log.Println(sector.Name + " assigned fingerprint.")
				sector.Fingerprint = MakeFingerprint(prefix)
			}
			log.Println("Sector loaded: " + sector.Name)
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
		sector.Rooms = make(map[int]glob.RoomData)
	}

	for _, value := range sector.Rooms {
		value.Players = make(map[string]*glob.PlayerData)
	}

	return &sector
}
