package support

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"../def"
	"../glob"
)

func CreatePlayer(desc *glob.ConnectionData) glob.PlayerData {
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

		Desc:  desc,
		Valid: true,
	}
	return player
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
	writer := os.Stdout
	enc := json.NewEncoder(writer)
	enc.SetIndent("", "\t")

	if err := enc.Encode(&player); err != nil {
		log.Println(err)
		return false
	}
	fo, err := os.Create(def.DATA_DIR + def.PLAYER_DIR + name)

	if err != nil {
		CheckError("WritePlayer: os.Create", err, def.ERROR_NONFATAL)
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
