package support

import (
	"fmt"
	"time"

	"../def"
	"../glob"
	"../mlog"
)

func DescWriteError(c *glob.ConnectionData, err error) {
	if err != nil {

		if c != nil {
			if c.Valid && c.State == def.CON_STATE_PLAYING && c.Name != def.STRING_UNKNOWN {
				if c.Player != nil && c.Player.Valid {
					buf := fmt.Sprintf("%s lost their network connection.", c.Player.Name)
					c.Player.UnlinkedTime = time.Now()
					c.Valid = false
					c.Name = def.STRING_UNKNOWN

					WriteToRoom(c.Player, buf)
				}
			}

		}
		RemoveNetDesc()
	}
}

func WriteToDesc(c *glob.ConnectionData, text string) {

	if c == nil || !c.Valid {
		return
	}
	text, overflow := TruncateString(text, def.MAX_OUTPUT_LENGTH)
	if overflow {
		cstring := " Name: " + c.Name + ", Addr:" + c.Address
		mlog.Write("WriteToDesc: string too large, Truncated!" + cstring)
	}

	message := fmt.Sprintf("%s\r\n", text)
	bytes, err := c.Desc.Write([]byte(ANSIColor(message)))
	c.BytesOut += bytes
	trackBytesOut(c)

	DescWriteError(c, err)
}

func WriteToPlayer(player *glob.PlayerData, text string) {

	if player == nil || !player.Valid || player.Connection == nil || !player.Connection.Valid {
		return
	}

	bytes, err := player.Connection.Desc.Write([]byte(ANSIColor(text) + "\r\n"))
	player.Connection.BytesOut += bytes
	trackBytesOut(player.Connection)

	DescWriteError(player.Connection, err)
}

func WriteToPlayerCodes(player *glob.PlayerData, text string) {

	if player == nil || !player.Valid || player.Connection == nil || !player.Connection.Valid {
		return
	}

	bytes, err := player.Connection.Desc.Write([]byte(text + "\r\n"))
	player.Connection.BytesOut += bytes
	trackBytesOut(player.Connection)

	DescWriteError(player.Connection, err)
}

func WriteToAll(text string) {
	if text == "" {
		return
	}

	for x := 1; x <= glob.ConnectionListEnd; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con.Valid == false {
			continue
		}
		if con.State == def.CON_STATE_PLAYING {
			bytes, err := con.Desc.Write([]byte(ANSIColor(text) + "\r\n"))
			con.BytesOut += bytes
			trackBytesOut(con)

			DescWriteError(con, err)
		}
	}
	//mlog.Write(text)
}

func WriteToOthers(player *glob.PlayerData, text string) {
	if player == nil || !player.Valid || text == "" {
		return
	}

	pc := player.Connection

	for x := 1; x <= glob.ConnectionListEnd; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con.Valid == false {
			continue
		}
		if con.Desc != pc.Desc && con.State == def.CON_STATE_PLAYING {
			bytes, err := con.Desc.Write([]byte(ANSIColor(text) + "\r\n"))
			con.BytesOut += bytes
			trackBytesOut(con)

			DescWriteError(con, err)
		}
	}
	mlog.Write(text)
}

func WriteToRoom(player *glob.PlayerData, text string) {
	if player == nil || !player.Valid || text == "" {
		return
	}

	if player.Location.RoomLink != nil {
		for _, target := range player.Location.RoomLink.Players {
			if target != nil && target != player {
				WriteToPlayer(target, "[Room] "+text)
			}
		}
	} else {
		buf := fmt.Sprintf("WriteToRoom: %v is not in a room.", player.Name)
		mlog.Write(buf)
	}
}
