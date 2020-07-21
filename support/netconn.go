package support

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"../def"
	"../glob"
)

func NewDescriptor(desc *net.TCPConn) {

	/*--- LOCK ---*/
	glob.ConnectionListLock.Lock()
	defer glob.ConnectionListLock.Unlock()
	/*--- LOCK ---*/

	for x := 1; x <= glob.ConnectionListMax; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con.Valid == true {
			continue
		} else {
			/*Replace*/
			con.Name = def.STRING_UNKNOWN
			con.Desc = desc
			con.Address = desc.LocalAddr().String()
			con.State = def.CON_STATE_WELCOME
			con.ConnectedFor = time.Now()
			con.IdleTime = time.Now()
			con.Id = x
			con.BytesOut = 0
			con.BytesIn = 0
			con.Player = nil
			con.Valid = true

			buf := fmt.Sprintf("Recycling connection #%d.", x)
			log.Println(buf)

			go ReadConnection(&glob.ConnectionList[x])
			return
		}
	}

	/*Generate new descriptor data*/
	if glob.ConnectionListMax >= def.MAX_USERS-1 {
		log.Println("Create ConnectionData: MAX_USERS REACHED!")
		desc.Write([]byte("Sorry, something has gone wrong (MAX_DESCRIPTORS)!\r\nGoodbye!\r\n"))
		return
	}

	/*Create*/
	glob.ConnectionListMax++
	newConnection := glob.ConnectionData{
		Name:         def.STRING_UNKNOWN,
		Desc:         desc,
		Address:      desc.LocalAddr().String(),
		State:        def.CON_STATE_WELCOME,
		ConnectedFor: time.Now(),
		IdleTime:     time.Now(),
		Id:           glob.ConnectionListMax,
		BytesOut:     0,
		BytesIn:      0,
		Player:       nil,
		Valid:        true}
	glob.ConnectionList[glob.ConnectionListMax] = newConnection
	buf := fmt.Sprintf("Created new connection #%d.", glob.ConnectionListMax)
	log.Println(buf)

	go ReadConnection(&glob.ConnectionList[glob.ConnectionListMax])
	return

}

func ReadConnection(con *glob.ConnectionData) {

	/*--- LOCK ---*/
	glob.ConnectionListLock.Lock()
	/*Create reader*/
	reader := bufio.NewReader(con.Desc)
	glob.ConnectionListLock.Unlock()
	/*--- UNLOCK ---*/

	for {

		input, err := reader.ReadString('\n')

		if err != nil {
			glob.ConnectionListLock.Lock()
			DescWriteError(con, err)
			glob.ConnectionListLock.Unlock()
			return
		}

		//Move all this to handlers, to get rid of if/elseif mess,
		//and to enable autocomplete and shortcuts/aliases.
		/*--- LOCK ---*/
		glob.ConnectionListLock.Lock()
		/*--- LOCK ---*/

		interpretInput(con, input)

		if con.State == def.CON_STATE_DISCONNECTING {
			con.Desc.Close()
			con.Valid = false
			if con.Player != nil && con.Player.Valid {
				con.Player.Valid = false
				con.Player.Connection = nil
			}
			/*--- UNLOCK ---*/
			glob.ConnectionListLock.Unlock()
			/*--- UNLOCK ---*/
			return
		}
		/*--- UNLOCK ---*/
		glob.ConnectionListLock.Unlock()
		/*--- UNLOCK ---*/
	}
}
func DescWriteError(c *glob.ConnectionData, err error) {
	if err != nil {
		CheckError("DescWriteError", err, def.ERROR_NONFATAL)

		if c.Name != def.STRING_UNKNOWN && c.State == def.CON_STATE_PLAYING {
			buf := fmt.Sprintf("%s lost their connection.", c.Name)
			WriteToOthers(c, buf)
		} else {
			buf := fmt.Sprintf("%s disconnected.", c.Address)
			log.Println(buf)
		}

		c.State = def.CON_STATE_DISCONNECTED
		c.Valid = false
		c.Desc.Close()
	}
}

func WriteToDesc(c *glob.ConnectionData, text string) {

	if c == nil {

		log.Println("Attempted to write to invalid ConnectionData.")
		return
	}
	message := fmt.Sprintf("%s\r\n", text)
	bytes, err := c.Desc.Write([]byte(message))
	c.BytesOut += bytes

	DescWriteError(c, err)
}

func WriteToPlayer(player *glob.PlayerData, text string) {

	if player != nil && player.Valid &&
		player.Connection != nil && player.Connection.Valid &&
		player.Connection.Desc != nil {
		message := fmt.Sprintf("%s\r\n", text)
		bytes, err := player.Connection.Desc.Write([]byte(message))
		player.Connection.BytesOut += bytes

		DescWriteError(player.Connection, err)
	} else {
		log.Println("Attempted to write to invalid or disconnected player.")
		log.Println(text)
		return
	}
}

func WriteToAll(text string) {

	for x := 0; x <= glob.ConnectionListMax; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con.Valid == false {
			continue
		}
		if con.State == def.CON_STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			bytes, err := con.Desc.Write([]byte(message))
			con.BytesOut += bytes

			DescWriteError(con, err)
		}
	}
	log.Println(text)
}

func WriteToOthers(c *glob.ConnectionData, text string) {

	for x := 0; x <= glob.ConnectionListMax; x++ {
		var con *glob.ConnectionData
		con = &glob.ConnectionList[x]
		if con.Valid == false {
			continue
		}
		if con.Desc != c.Desc && con.State == def.CON_STATE_PLAYING {
			message := fmt.Sprintf("%s\r\n", text)
			bytes, err := con.Desc.Write([]byte(message))
			con.BytesOut += bytes

			DescWriteError(c, err)
		}
	}
	log.Println(text)
}
