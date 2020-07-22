package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"../def"
	"../glob"
)

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
	return true
}
