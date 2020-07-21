package support

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"../def"
	"../glob"
)

func CreatePlayer(desc *glob.ConnectionData) *glob.PlayerData {
	player := glob.PlayerData{
		Name:        desc.Name,
		Password:    "",
		PlayerType:  def.PLAYER_TYPE_NEW,
		Level:       0,
		State:       def.PLAYER_ALIVE,
		Sector:      0,
		Vnum:        0,
		Created:     time.Now(),
		LastSeen:    time.Now(),
		Seconds:     0,
		IPs:         []string{},
		Connections: []int{},
		BytesIn:     []int{},
		BytesOut:    []int{},
		Email:       "",

		Description: "",
		Sex:         "",

		Desc:  desc.Desc,
		Valid: true,
	}
	return &player
}

func ReadPlayer(name string, load bool) (*glob.PlayerData, bool) {

	_, err := ioutil.ReadFile(def.DATA_DIR + def.PLAYER_DIR + name)
	if err != nil {
		return nil, false
	} else {
		if load == true {
			//
			return nil, true
		} else {
			return nil, true
		}
	}
}

func WritePlayer(player *glob.PlayerData) bool {
	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	if err := enc.Encode(&player); err != nil {
		CheckError("WritePlayer: enc.Encode", err, def.ERROR_NONFATAL)
		return false
	}
	_, err := os.Create(def.DATA_DIR + def.PLAYER_DIR + player.Name)

	if err != nil {
		CheckError("WritePlayer: os.Create", err, def.ERROR_NONFATAL)
		return false
	}

	err = ioutil.WriteFile(def.DATA_DIR+def.PLAYER_DIR+player.Name, []byte(outbuf.String()), 0644)

	if err != nil {
		CheckError("WritePlayer: WriteFile", err, def.ERROR_NONFATAL)
		return false
	}
	return true
}

func WriteDesc(desc *glob.ConnectionData) bool {
	writer := os.Stdout
	enc := json.NewEncoder(writer)
	enc.SetIndent("", "\t")

	if err := enc.Encode(&desc); err != nil {
		log.Println(err)
		return false
	}
	return true
}
