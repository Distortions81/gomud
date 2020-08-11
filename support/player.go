package support

import (
	"fmt"
	"log"
	"time"

	"../def"
	"../glob"
)

func GetPTypeString(ptype int) string {
	for _, a := range glob.PlayerTypes {
		if a.PType == ptype {
			return a.PName
		}
	}

	return ""
}

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

	player.OLCSettings.NoOLCPrefix = true
	player.OLCSettings.OLCRoomFollow = true
	player.OLCSettings.OLCShowCodes = true

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
