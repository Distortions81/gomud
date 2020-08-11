package support

import (
	"../def"
	"../glob"
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
