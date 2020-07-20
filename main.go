package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"./def"
	"./glob"
	"./support"
)

func main() {

	/*Find Network*/
	addr, err := net.ResolveTCPAddr("tcp", def.DEFAULT_PORT)
	support.CheckError(err, def.ERROR_FATAL)

	/*Open Listener*/
	ServerListener, err := net.ListenTCP("tcp", addr)
	glob.ServerListener = ServerListener
	support.CheckError(err, def.ERROR_FATAL)

	/*Print Connection*/
	buf := fmt.Sprintf("Server online at: %s", addr.String())
	log.Println(buf)

	/*Listen for connections*/
	mainLoop()
	ServerListener.Close()
}

func mainLoop() {

	/*separate thread, wait for new connections*/
	for glob.ServerState == def.SERVER_RUNNING {
		/*Check for new connections*/
		desc, err := glob.ServerListener.AcceptTCP()

		if err == nil {
			buf := fmt.Sprintf("New connection from %s.", desc.LocalAddr().String())
			log.Println(buf)

			/*Change connections settings*/
			err := desc.SetLinger(-1)
			support.CheckError(err, def.ERROR_NONFATAL)
			err = desc.SetNoDelay(true)
			support.CheckError(err, def.ERROR_NONFATAL)
			err = desc.SetReadBuffer(10000) //10k, 10 seconds of insanely-fast typing
			support.CheckError(err, def.ERROR_NONFATAL)
			err = desc.SetWriteBuffer(12500000) //12.5MB, 10 second buffer at 10mbit
			support.CheckError(err, def.ERROR_NONFATAL)

			//TODO Add full greeting/info
			/*Respond here, so we don't have to wait for lock*/
			_, err = desc.Write([]byte(def.VERSION + "\r\nTo create a new character, type: NEW\r\nName: \r\n"))
			support.CheckError(err, def.ERROR_NONFATAL)

			go desc.NewDescriptor(desc)
			/* desc/desc.go */
		}

	}
}

func createPlayer() glob.PlayerData {
	player := glob.PlayerData{
		Name:        def.STRING_UNKNOWN,
		Password:    "",
		PlayerType:  def.PLAYER_TYPE_NEW,
		Level:       0,
		State:       def.PLAYER_ALIVE,
		Sector:      0,
		Vnum:        0,
		Created:     time.Now(),
		LastSeen:    time.Now(),
		Seconds:     0,
		IPs:         []string{},
		Connections: []int{},
		BytesIn:     []int{},
		BytesOut:    []int{},
		Email:       "",

		Description: "",
		Sex:         "",

		Desc:  nil,
		Valid: true,
	}
	return player
}

func showCommands(c *glob.ConnectionData) {
	us := fmt.Sprintf("commands: say, who, quit")
	desc.WriteToDesc(c, us)
}
