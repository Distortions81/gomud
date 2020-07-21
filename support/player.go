package support

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"../def"
	"../glob"
)

func CreatePlayer() *glob.PlayerData {
	player := glob.PlayerData{
		Name:        def.STRING_UNKNOWN,
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

		Connection: nil,
		Valid:      true,
	}
	return &player
}

func CreatePlayerFromDesc(conn *glob.ConnectionData) *glob.PlayerData {
	player := glob.PlayerData{
		Name:        conn.Name,
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

		Connection: conn,
		Valid:      true,
	}
	return &player
}

func ReadPlayer(name string, load bool) (*glob.PlayerData, bool) {

	_, err := os.Stat(def.DATA_DIR + def.PLAYER_DIR + strings.ToLower(name))
	notfound := os.IsNotExist(err)

	if notfound {
		//CheckError("ReadPlayer: os.Stat", err, def.ERROR_NONFATAL)
		log.Println("Player not found: " + name)
		return nil, false

	} else {

		if load {
			file, err := ioutil.ReadFile(def.DATA_DIR + def.PLAYER_DIR + strings.ToLower(name))

			if file != nil && err == nil {
				player := CreatePlayer()

				err := json.Unmarshal([]byte(file), &player)
				if err != nil {
					CheckError("ReadPlayer: Unmashal", err, def.ERROR_NONFATAL)
				}

				log.Println("Player loaded: " + player.Name)
				return player, true
			} else {
				CheckError("ReadPlayer: ReadFile", err, def.ERROR_NONFATAL)
				return nil, false
			}
		} else {
			//If we are just checking if player exists,
			//don't bother to actually load the file.
			log.Println("Player found: " + name)
			return nil, true
		}
	}
}

func WritePlayer(player *glob.PlayerData, report bool) bool {
	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	if player == nil {
		log.Println("WritePlayer: nil player")
		return false
	}

	if report {
		WriteToPlayer(player, "Character saved.")
	}

	if err := enc.Encode(&player); err != nil {
		CheckError("WritePlayer: enc.Encode", err, def.ERROR_NONFATAL)
		return false
	}
	_, err := os.Create(def.DATA_DIR + def.PLAYER_DIR + strings.ToLower(player.Name))

	if err != nil {
		CheckError("WritePlayer: os.Create", err, def.ERROR_NONFATAL)
		return false
	}

	err = ioutil.WriteFile(def.DATA_DIR+def.PLAYER_DIR+strings.ToLower(player.Name), []byte(outbuf.String()), 0644)

	if err != nil {
		CheckError("WritePlayer: WriteFile", err, def.ERROR_NONFATAL)
		return false
	}
	return true
}

func LinkPlayerConnection(player *glob.PlayerData, con *glob.ConnectionData) {

	if player != nil && player.Valid && con != nil && con.Valid {
		player.Connection = con
		con.Player = player
	}
}
