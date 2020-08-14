package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"../def"
	"../glob"
	"../mlog"
)

//TODO ASYNC READ
func ReadPlayer(name string, load bool) (*glob.PlayerData, bool) {

	_, err := os.Stat(def.DATA_DIR + def.PLAYER_DIR + strings.ToLower(name))
	notfound := os.IsNotExist(err)

	if notfound {
		//CheckError("ReadPlayer: os.Stat", err, def.ERROR_NONFATAL)
		//mlog.Write("Player not found: " + name)
		return nil, false

	} else {

		if load {

			glob.PlayerFileLock.Lock()
			defer glob.PlayerFileLock.Unlock()

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
				/*Re-link OLC pointer*/
				if player.OLCEdit.Active {
					loc, found := LocationDataFromID(player.OLCEdit.Room.Sector, player.OLCEdit.Room.ID)
					if found {
						player.OLCEdit.Room.RoomLink = loc.RoomLink
					}
					obj, found := GetObjectFromID(player.OLCEdit.Object.Sector, player.OLCEdit.Object.ID)
					if found {
						player.OLCEdit.Object.ObjectLink = obj
					}
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

func WritePlayer(player *glob.PlayerData, asyncSave bool) bool {
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

	//Async write
	if asyncSave {
		go writePlayerFile(outbuf, fileName)
	} else {
		writePlayerFile(outbuf, fileName)
	}

	player.Dirty = false
	return true
}

func writePlayerFile(outbuf *bytes.Buffer, fileName string) {
	glob.PlayerFileLock.Lock()
	defer glob.PlayerFileLock.Unlock()

	err := ioutil.WriteFile(fileName, []byte(outbuf.String()), 0644)

	if err != nil {
		CheckError("WritePlayer: WriteFile", err, def.ERROR_NONFATAL)
	}

	buf := fmt.Sprintf("Wrote %v, %v.", fileName, ScaleBytes(len(outbuf.String())))
	mlog.Write(buf)
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
