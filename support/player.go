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
	"../mlog"
)

func SetupNewCharacter(player *glob.PlayerData) {
	if player == nil && !player.Valid {
		return
	}
	player.Location.Sector = def.PLAYER_START_SECTOR
	player.Location.ID = def.PLAYER_START_ROOM

	/*Default config options*/
	player.Config.Ansi = true
	player.Config.PostNewline = true
	player.Config.PreNewline = true

	player.Fingerprint = MakeFingerprint(player.Name)
	WriteToAll("A newcomer has arrived, their name is " + player.Name + "...")
}

func CreatePlayer() *glob.PlayerData {
	loc := glob.LocationData{Sector: def.PLAYER_START_SECTOR, ID: def.PLAYER_START_ROOM}

	player := glob.PlayerData{
		Name:        def.STRING_UNKNOWN,
		Password:    "",
		PlayerType:  def.PLAYER_TYPE_NEW,
		Level:       0,
		State:       def.PLAYER_ALIVE,
		Location:    loc,
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

	player.Aliases = make(map[string]string)
	player.Connections = make(map[string]int)
	player.BytesIn = make(map[string]int)
	player.BytesOut = make(map[string]int)

	return &player
}

func CreatePlayerFromDesc(conn *glob.ConnectionData) *glob.PlayerData {
	loc := glob.LocationData{Sector: def.PLAYER_START_SECTOR, ID: def.PLAYER_START_ROOM}
	player := glob.PlayerData{
		Name:        conn.Name,
		Password:    "",
		PlayerType:  def.PLAYER_TYPE_NEW,
		Level:       0,
		State:       def.PLAYER_ALIVE,
		Location:    loc,
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
		//mlog.Write("Player not found: " + name)
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

				mlog.Write("Player loaded: " + player.Name)
				return player, true
			} else {
				CheckError("ReadPlayer: ReadFile", err, def.ERROR_NONFATAL)
				return nil, false
			}
		} else {
			//If we are just checking if player exists,
			//don't bother to actually load the file.
			//mlog.Write("Player found: " + name)
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

	player.LastSeen = time.Now()

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
	mlog.Write(buf)
	player.Dirty = false
	return true
}

func LinkPlayerConnection(player *glob.PlayerData, con *glob.ConnectionData) {

	if player == nil || con == nil || player.Valid == false {
		return
	}

	/*If player is already in the world, re-use*/
	for x := 1; x <= glob.PlayerListEnd; x++ {
		if glob.PlayerList[x] != nil &&
			glob.PlayerList[x].Name == player.Name &&
			glob.PlayerList[x].Fingerprint == player.Fingerprint {

			/* Invalidate old connection */
			if glob.PlayerList[x].Connection != nil {
				glob.PlayerList[x].Connection.Valid = false
			}
			/*Get rid of previous character from login*/
			con.Player.Valid = false
			con.Player = player //Replace pfile data with live

			player.Connection = con

			/*Re-activate old body*/
			player.UnlinkedTime = time.Time{} //Reset unlinked timer
			player.Valid = true
			player.Connection.Valid = true

			/* MOTD message here */
			WriteToPlayer(player, "\r\n")

			PlayerToRoom(player, player.Location.Sector, player.Location.ID)
			buf := fmt.Sprintf("%s reconnects to their body.", player.Name)
			WriteToRoom(player, buf)
			CmdLook(player, "")
			WriteToPlayer(player, "You reconnect to your body.")
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

	/*Recycle players*/
	recycled := false
	if glob.PlayerListEnd > 1 {
		for x := 1; x <= glob.PlayerListEnd; x++ {
			if glob.PlayerList[x].Valid == false {
				glob.PlayerList[x] = player
				recycled = true
				buf := fmt.Sprintf("Recycling player #%v.", x)
				log.Println(buf)
			}
		}
	}
	/* Create new if needed */
	if recycled == false {
		glob.PlayerListEnd++
		glob.PlayerList[glob.PlayerListEnd] = player
		buf := fmt.Sprintf("Creating new player #%v.", glob.PlayerListEnd)
		log.Println(buf)
	}

	/* MOTD message here */

	PlayerToRoom(player, player.Location.Sector, player.Location.ID)

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
	if player.Location.RoomLink != nil {
		room := player.Location.RoomLink
		delete(room.Players, player.Fingerprint)
	}

	//Add player to room, add error handling
	//Automatically generate "players" map if it doesn't exist
	if glob.SectorsList[sectorID].Valid &&
		glob.SectorsList[sectorID].Rooms != nil {

		room := glob.SectorsList[sectorID].Rooms[roomID]
		if room != nil {
			room.Players[player.Fingerprint] = player
		}

		player.Location.RoomLink = glob.SectorsList[sectorID].Rooms[roomID]
		player.Location.Sector = sectorID
		player.Location.ID = roomID
		player.Dirty = true

		if player.OLCEdit.Active &&
			player.OLCEdit.Mode == def.OLC_ROOM &&
			player.OLCSettings.OlcRoomFollow {

			player.OLCEdit.Room = player.Location
		}

	} else {
		mlog.Write("PlayerToRoom: That sector or room is not valid.")
		mlog.Write(fmt.Sprintf("Sector: %v, Room: %v, Player: %v", sectorID, roomID, player.Name))
		if sectorID != def.PLAYER_START_SECTOR && roomID != def.PLAYER_START_ROOM {
			PlayerToRoom(player, def.PLAYER_START_SECTOR, def.PLAYER_START_ROOM)
		} else {
			log.Println("Default room/sector not found. Quitting.")
			os.Exit(1)
		}

	}
}

func RemovePlayer(player *glob.PlayerData) {
	/* Check if data is valid */
	if player == nil {
		fmt.Println("RemovePlayer: nil player")
		return
	}
	if player.Valid == false {
		fmt.Println("RemovePlayer: non-valid player")
	}

	/* Remove player from room */
	if player.Location.RoomLink != nil {
		room := player.Location.RoomLink
		delete(room.Players, player.Fingerprint)
	}

	/* Set player and connection as invalid, clear room pointer */
	player.Location.RoomLink = nil
	player.Valid = false
	if player.Connection != nil {
		player.Connection.State = def.CON_STATE_DISCONNECTED
		player.Connection.Valid = false
		player.Connection = nil
	}

	buf := fmt.Sprintf("%v invalidated, end: %v", player.Name, glob.PlayerListEnd)
	player.Name = def.STRING_UNKNOWN
	mlog.Write(buf)
}
