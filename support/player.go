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
	if player == nil && !player.Valid {
		return
	}
	player.Sector = def.PLAYER_START_SECTOR
	player.Room = def.PLAYER_START_ROOM
	player.Fingerprint = MakeFingerprint(player.Name)
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

	player.Connections = make(map[string]int)
	player.BytesIn = make(map[string]int)
	player.BytesOut = make(map[string]int)

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

	player.Connections = make(map[string]int)
	player.BytesIn = make(map[string]int)
	player.BytesOut = make(map[string]int)

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

				if player.Connections == nil {
					player.Connections = make(map[string]int)
				}
				if player.BytesIn == nil {
					player.BytesIn = make(map[string]int)
				}
				if player.BytesOut == nil {
					player.BytesOut = make(map[string]int)
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
	fileName := def.DATA_DIR + def.PLAYER_DIR + strings.ToLower(player.Name)

	if player == nil && !player.Valid {
		return false
	}

	if err := enc.Encode(&player); err != nil {
		CheckError("WritePlayer: enc.Encode", err, def.ERROR_NONFATAL)
		return false
	}
	_, err := os.Create(fileName)

	if err != nil {
		CheckError("WritePlayer: os.Create", err, def.ERROR_NONFATAL)
		return false
	}

	err = ioutil.WriteFile(fileName, []byte(outbuf.String()), 0644)

	if err != nil {
		CheckError("WritePlayer: WriteFile", err, def.ERROR_NONFATAL)
		return false
	}

	buf := fmt.Sprintf("Wrote %v, %v.", fileName, ScaleBytes(len(outbuf.String())))
	log.Println(buf)
	return true
}

func LinkPlayerConnection(player *glob.PlayerData, con *glob.ConnectionData) {

	if player == nil || !player.Valid || con == nil || !con.Valid {
		return
	}

	/*If player is already in the world, re-use*/
	for x := 0; x <= glob.PlayerListEnd; x++ {
		if glob.PlayerList[x] != nil && glob.PlayerList[x].Valid == false &&
			glob.PlayerList[x].Name == player.Name &&
			glob.PlayerList[x].Fingerprint == player.Fingerprint {

			/*Get rid of previous character from login*/
			con.Player.Valid = false
			con.Player = nil

			con.Player = player //Replace pfile data with live
			player.Connection = con

			/*Re-activate old body*/
			player.UnlinkedTime = time.Time{}
			player.Valid = true

			/* MOTD message here */
			WriteToPlayer(player, "\r\n")

			PlayerToRoom(player, player.Sector, player.Room)
			buf := fmt.Sprintf("%s reconnects to their body.", player.Name)
			WriteToRoom(player, buf)
			WriteToPlayer(player, "You reconnect to your body.")

			CmdWho(player, "")
			CmdLook(player, "")
			return
		}
	}

	if player.Connections == nil {
		player.Connections = make(map[string]int)
	}
	player.Connections[con.Address]++

	/*Link to each other*/
	player.Connection = con
	con.Player = player

	/*Add to global player list*/
	glob.PlayerListEnd++
	glob.PlayerList[glob.PlayerListEnd] = player

	/* MOTD message here */
	WriteToPlayer(player, "\r\n")

	PlayerToRoom(player, player.Sector, player.Room)

	buf := fmt.Sprintf("%s suddenly appears.", player.Name)
	WriteToRoom(player, buf)

	CmdWho(player, "")
	CmdLook(player, "")
}

func PlayerToRoom(player *glob.PlayerData, sectorID int, roomID int) {

	if player == nil && !player.Valid {
		return
	}
	//Remove player from room, if they are in one
	if player.RoomLink != nil {
		room := player.RoomLink
		delete(room.Players, player.Fingerprint)
	}

	//Add player to room, add error handling
	//Automatically generate "players" map if it doesn't exist
	if glob.SectorsList[sectorID].Valid &&
		glob.SectorsList[sectorID].Rooms != nil {

		room := glob.SectorsList[sectorID].Rooms[roomID]
		if room.Players == nil {
			room.Players = make(map[string]*glob.PlayerData)
		}
		room.Players[player.Name] = player

		glob.SectorsList[sectorID].Rooms[roomID] = room

		player.RoomLink = &room
		player.Sector = sectorID
		player.Room = roomID

	} else {
		log.Println("PlayerToRoom: That sector or room is not valid.")
		log.Println(fmt.Sprintf("Sector: %v, Room: %v, Player: %v", sectorID, roomID, player.Name))

	}
}

func RemovePlayerWorld(player *glob.PlayerData) {
	if player == nil && !player.Valid {
		return
	}
	PlayerToRoom(player, 0, 0)
	player.Valid = false

	buf := fmt.Sprintf("%v invalidated, end: %v", player.Name, glob.PlayerListEnd)
	log.Println(buf)
}
