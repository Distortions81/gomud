package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"../def"
	"../glob"
)

func SetupNewCharacter(player *glob.PlayerData) {
	player.Sector = def.PLAYER_START_SECTOR
	player.Room = def.PLAYER_START_ROOM
	player.Fingerprint = MakeFingerprint(player.Name)
	WriteToPlayer(player, "Welcome! Type LOOK to see around you, and HELP to see more commands.")
	WriteToAll("A newcomer has arrived, their name is " + player.Name + "...")
}

func CreatePlayer() *glob.PlayerData {
	player := glob.PlayerData{
		Name:        def.STRING_UNKNOWN,
		Password:    "",
		PlayerType:  def.PLAYER_TYPE_NEW,
		Level:       0,
		State:       def.PLAYER_ALIVE,
		Sector:      0,
		Room:        0,
		Created:     time.Now(),
		LastSeen:    time.Now(),
		TimePlayed:  0,
		Connections: nil,
		BytesIn:     nil,
		BytesOut:    nil,
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
		Room:        0,
		Created:     time.Now(),
		LastSeen:    time.Now(),
		TimePlayed:  0,
		Connections: nil,
		BytesIn:     nil,
		BytesOut:    nil,
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
			//log.Println("Player found: " + name)
			return nil, true
		}
	}
}

func WritePlayer(player *glob.PlayerData) bool {
	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	player.Version = def.PFILE_VERSION

	if player == nil {
		log.Println("WritePlayer: nil player")
		return false
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
		for x := 0; x < def.MAX_USERS; x++ {
			if glob.PlayerList[x].Fingerprint == player.Fingerprint && glob.PlayerList[x].Name == player.Name {
				con.Player = player //Replace pfile data with live
				player.Connection = con

				buf := fmt.Sprintf("%s reconnects to their body.", player.Name)
				PlayerToRoom(player, player.Sector, player.Room)
				WriteToRoom(player, buf)
				return
			}
		}

		if player.Connections == nil {
			player.Connections = make(map[string]int)
		}
		player.Connections[con.Address]++
		player.Connection = con

		PlayerToRoom(player, player.Sector, player.Room)
		con.Player = player
		con.State = def.CON_STATE_PLAYING
	}
}

func PlayerToRoom(player *glob.PlayerData, sectorID int, roomID int) {

	//Remove player from room, if they are in one
	if player.RoomLink != nil {
		buf := fmt.Sprintf("%s left.", player.Name)
		WriteToRoom(player, buf)

		room := player.RoomLink
		delete(room.Players, player.Fingerprint)
	}

	//Add player to room, add error handling
	glob.SectorsList[sectorID].Rooms[roomID].Players[player.Fingerprint] = player
	room := glob.SectorsList[sectorID].Rooms[roomID]
	player.RoomLink = &room
	player.Sector = sectorID
	player.Room = roomID

	//Send to message handler.
	buf := fmt.Sprintf("%s has arrived.", player.Name)
	WriteToRoom(player, buf)
}

func TrackBytesPlayer(con *glob.ConnectionData, player *glob.PlayerData) {

	if player != nil && player.Valid && con != nil && con.Valid {
		player.BytesOut[con.Address] += con.BytesOut
		player.BytesIn[con.Address] += con.BytesIn
	}

}
