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
	WriteToConn(c, text, true, false)
}

func WriteToConn(c *glob.ConnectionData, text string, color bool, codes bool) {

	if c == nil || !c.Valid {
		return
	}
	text, overflow := TruncateString(text, def.MAX_OUTPUT_LENGTH)
	if overflow {
		cstring := " Name: " + c.Name + ", Addr:" + c.Address
		mlog.Write("WriteToDesc: string too large, Truncated!" + cstring)
	}

	bytes := 0
	var err error

	message := fmt.Sprintf("%s\r\n", text)
	if color {
		bytes, err = c.Desc.Write([]byte(ANSIColor(message)))
	} else if codes {
		bytes, err = c.Desc.Write([]byte(message))
	} else {
		bytes, err = c.Desc.Write([]byte(StripColorCodes(message)))
	}
	c.BytesOut += bytes
	trackBytesOut(c)

	DescWriteError(c, err)
}

func WriteToPlayer(player *glob.PlayerData, text string) {
	if player == nil || !player.Valid || player.Connection == nil || !player.Connection.Valid {
		return
	}
	WriteToConn(player.Connection, text, player.Config.Ansi, false)
}

func WriteToPlayerCodes(player *glob.PlayerData, text string) {

	if player == nil || !player.Valid || player.Connection == nil || !player.Connection.Valid {
		return
	}
	WriteToConn(player.Connection, text, false, true)
}

func WriteToAll(text string) {

	for x := 0; x <= glob.ConnectionListEnd; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con != nil && con.Valid && con.Player != nil && con.State == def.CON_STATE_PLAYING {
			if con != nil && con.Valid && con.Player != nil && con.Player.Valid {
				WriteToConn(con, text, con.Player.Config.Ansi, false)
			} else if con != nil && con.Valid {
				WriteToConn(con, text, true, false)
			}
		}
	}
}

func WriteToOthers(player *glob.PlayerData, text string) {
	if player == nil || !player.Valid || player.Connection == nil || !player.Connection.Valid || text == "" {
		return
	}
	for x := 0; x <= glob.ConnectionListEnd; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con != player.Connection {
			if con != nil && con.Valid && con.Player != nil && con.State == def.CON_STATE_PLAYING {
				if con != nil && con.Valid && con.Player != nil && con.Player.Valid {
					WriteToConn(con, text, con.Player.Config.Ansi, false)
				} else if con != nil && con.Valid {
					WriteToConn(con, text, true, false)
				}
			}
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
