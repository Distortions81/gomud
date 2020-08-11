package support

import (
	"net"
	"strings"
	"time"

	"../def"
	"../glob"
)

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
