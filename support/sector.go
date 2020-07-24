package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"../def"
	"../glob"
)

func ReadSectorList() {

	files, err := ioutil.ReadDir(def.DATA_DIR + def.SECTOR_DIR)
	if err != nil {
		CheckError("ReadSectorList:", err, def.ERROR_NONFATAL)
	}

	for _, file := range files {
		sector := ReadSector(file.Name())
		if glob.SectorsList[sector.ID].Valid {
			buf := fmt.Sprint("%v has same sector ID as %v! Skipping!", sector.Name, glob.SectorsList[sector.ID].Name)
			log.Println(buf)
		}
	}
}

func ReloadSector() {

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

	return nil
}
