package main

import (
	"fmt"
	"log"
	"net"

	"./def"
	"./glob"
	"./support"
)

func main() {

	var defaultSector SectorsData = 
		ID: 0,
		Group: "Default",

		Name: "Default",
		Description: "Default sector"

		Rooms: nil,
		Valid: true,
	}

	defaultSector.Rooms[0] RoomData = {
		VNUM: nil,
		Name: "Default room",
		Description: "This is the default room.",
		Valid: true,
	}

	defaultSector.Rooms[0].VNUM

	/*Find Network*/
	addr, err := net.ResolveTCPAddr("tcp", def.DEFAULT_PORT)
	support.CheckError("main: resolveTCP", err, def.ERROR_FATAL)

	/*Open Listener*/
	ServerListener, err := net.ListenTCP("tcp", addr)
	glob.ServerListener = ServerListener
	support.CheckError("main: ListenTCP", err, def.ERROR_FATAL)

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
			support.CheckError("main: SetLinger", err, def.ERROR_NONFATAL)
			err = desc.SetNoDelay(true)
			support.CheckError("main: SetNoDelay", err, def.ERROR_NONFATAL)
			err = desc.SetReadBuffer(10000) //10k, 10 seconds of insanely-fast typing
			support.CheckError("main: SetReadBuffer", err, def.ERROR_NONFATAL)
			err = desc.SetWriteBuffer(12500000) //12.5MB, 10 second buffer at 10mbit
			support.CheckError("main: SetWriteBuffer", err, def.ERROR_NONFATAL)

			//TODO Add full greeting/info
			/*Respond here, so we don't have to wait for lock*/
			_, err = desc.Write([]byte(def.VERSION + "\r\nTo create a new character, type: NEW\r\nName: \r\n"))
			support.CheckError("main: desc.Write", err, def.ERROR_NONFATAL)

			go support.NewDescriptor(desc)
			/* netconn/netconn.go */
		}

	}
}
