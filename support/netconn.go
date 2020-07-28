package support

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	"../def"
	"../glob"
	"../mlog"
)

func AddNetDesc() {
	glob.OpenDescLock.Lock()
	glob.OpenDesc++
	glob.OpenDescLock.Unlock()
}

func RemoveNetDesc() {
	glob.OpenDescLock.Lock()
	if glob.OpenDesc > 0 {
		glob.OpenDesc--
	}
	glob.OpenDescLock.Unlock()
}

func GetNetDesc() int {
	glob.OpenDescLock.Lock()
	count := glob.OpenDesc
	glob.OpenDescLock.Unlock()

	return count
}

func CheckNetDesc() bool {
	glob.OpenDescLock.Lock()
	if glob.OpenDesc >= def.MAX_DESC {
		return true
	}
	glob.OpenDescLock.Unlock()
	return false
}

func AutoResolveAddress(con *glob.ConnectionData) {

	addr := con.Desc.RemoteAddr().String()
	addrp := strings.Split(addr, ":")
	addrLen := len(addrp)
	if addrLen > 0 {
		addr = addrp[0]
	}
	con.Address = addr
}

func NewDescriptor(desc net.Conn, ssl bool) {

	if desc == nil || CheckNetDesc() {
		return
	}
	/* Recycle connection if we can */
	glob.ConnectionListLock.Lock()         /*--- LOCK ---*/
	defer glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/

	for x := 1; x <= glob.ConnectionListEnd; x++ {

		if glob.ConnectionList[x].Valid == true {
			continue
		} else {
			newConnection := glob.ConnectionData{
				Name:         def.STRING_UNKNOWN,
				Desc:         desc,
				Address:      "",
				SSL:          ssl,
				State:        def.CON_STATE_WELCOME,
				ConnectedFor: time.Now(),
				IdleTime:     time.Now(),
				BytesOut:     0,
				BytesIn:      0,
				Player:       nil,
				Valid:        true}

			AutoResolveAddress(&newConnection)
			glob.ConnectionList[x] = newConnection

			if glob.ConnectionListEnd >= def.MAX_USERS-1 {
				return
			}
			go ReadConnection(&glob.ConnectionList[x])
			return
		}
	}

	/*Create*/
	newConnection := glob.ConnectionData{
		Name:         def.STRING_UNKNOWN,
		Desc:         desc,
		Address:      "",
		SSL:          ssl,
		State:        def.CON_STATE_WELCOME,
		ConnectedFor: time.Now(),
		IdleTime:     time.Now(),
		BytesOut:     0,
		BytesIn:      0,
		Player:       nil,
		Valid:        true}
	AutoResolveAddress(&newConnection)

	glob.ConnectionListEnd++
	glob.ConnectionList[glob.ConnectionListEnd] = newConnection

	if glob.ConnectionListEnd >= def.MAX_USERS-1 {
		return
	}

	if !CheckNetDesc() {
		go ReadConnection(&glob.ConnectionList[glob.ConnectionListEnd])
	}
	return
}

func ReadConnection(con *glob.ConnectionData) {

	glob.ConnectionListLock.Lock() /*--- LOCK ---*/
	if con == nil {
		return
	}
	if con.Valid == false {
		return
	}
	if glob.ConnectionListEnd >= def.MAX_USERS-1 {
		return
	}

	reader := bufio.NewReader(con.Desc)

	glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/

	for con.Valid && con.Desc != nil {

		input, err := reader.ReadString('\n')
		glob.ConnectionListLock.Lock() /*--- LOCK ---*/

		/*Connection died*/
		if err != nil {
			DescWriteError(con, err)
			glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/
			return
		}

		filter := StripControl(input)
		limit, _ := TruncateString(filter, def.MAX_INPUT_LENGTH)

		if con.Input.BufferInCount-con.Input.BufferOutCount >= def.MAX_INPUT_LINES-1 {
			for x := 0; x <= 3; x++ {
				WriteToDesc(con, "Too many lines, stop spamming!")
				CloseConnection(con)
				glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/
				return
			}
		}

		con.Input.BufferInPos++
		con.Input.BufferInCount++
		if con.Input.BufferInPos >= def.MAX_INPUT_LINES {
			con.Input.BufferInPos = 0
		}
		con.Input.InputBuffer[con.Input.BufferInPos] = limit
		glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/

	}
}

func ReadPlayerInputBuffer(con *glob.ConnectionData) {
	/* Only run if we have something */
	if con.Input.BufferInCount > con.Input.BufferOutCount {

		con.Input.BufferOutCount++
		con.Input.BufferOutPos++

		if con.Input.BufferOutPos >= def.MAX_INPUT_LINES {
			con.Input.BufferOutPos = 0
		}
		input := con.Input.InputBuffer[con.Input.BufferOutPos]

		bIn := len(input)

		con.BytesIn += bIn
		trackBytesIn(con)

		HandleReadConnection(con, input)
	}
}

func HandleReadConnection(con *glob.ConnectionData, input string) {

	//Newline before commands
	WriteToDesc(con, "")

	/*Player aliases*/
	if con.Player != nil && con.Player.Valid {
		if con.Player.Aliases != nil {

			if input != "" {
				for key, value := range con.Player.Aliases {

					if strings.EqualFold(key, input) {
						//add ; newline support
						interpretInput(con, value, true)
						return
					}

				}
			}
		}
	}

	/*Handles all user input*/
	interpretInput(con, input, false)

	/*Handle players marked for disconnection*/
	/*Doing this at the end is cleaner*/
	if con.State == def.CON_STATE_DISCONNECTING {
		CloseConnection(con)
		if con.Player != nil {
			RemovePlayer(con.Player)
		}
	}

}

func trackBytesOut(con *glob.ConnectionData) {

	player := con.Player

	if player == nil || !player.Valid || con == nil || !con.Valid {
		return
	}
	player.BytesOut[con.Address] += (con.BytesOut - con.BytesOutRecorded)
	con.BytesOutRecorded = con.BytesOut
}

func trackBytesIn(con *glob.ConnectionData) {

	player := con.Player

	if player == nil || !player.Valid || con == nil || !con.Valid {
		return
	}
	player.BytesIn[con.Address] += (con.BytesIn - con.BytesInRecorded)
	con.BytesInRecorded = con.BytesIn
}

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
	bytes, err := c.Desc.Write([]byte(message))
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

func CloseConnection(con *glob.ConnectionData) {
	if con.Desc != nil {
		WriteToDesc(con, glob.AuRevoir)

		desc := con.Desc
		go func(desc net.Conn) {
			time.Sleep(2 * time.Second)
			desc.Close()
			RemoveNetDesc()
		}(desc)
	}
	con.Name = def.STRING_UNKNOWN
	con.Valid = false
	con.State = def.CON_STATE_DISCONNECTED
	con.Player = nil
}
