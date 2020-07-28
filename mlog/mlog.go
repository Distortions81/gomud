package mlog

import (
	"fmt"
	"log"
	"time"

	"../def"
	"../glob"
)

func Write(line string) {
	t := time.Now()
	dss := fmt.Sprintf("%v:%v:%v: ", t.Hour(), t.Minute(), t.Second())

	log.Println(line)
	writeToMods(dss + line)
}

func writeToMods(text string) {
	if text == "" {
		return
	}

	for x := 1; x <= glob.PlayerListEnd; x++ {
		player := glob.PlayerList[x]

		if player != nil && player.Valid && player.Connection.Valid {
			if player.Connection.State == def.CON_STATE_PLAYING && player.PlayerType >= def.PLAYER_TYPE_BUILDER {
				message := fmt.Sprintf("[LOG] %s\r\n", text)
				player.Connection.Desc.Write([]byte(message))
			}
		}
	}

}
