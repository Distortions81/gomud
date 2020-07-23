package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"./def"
	"./glob"
	"./support"
)

func main() {

	defaultSector := glob.SectorData{
		Area: "Default",

		Name:        "Default",
		Description: "Default sector",
		Valid:       true,
	}

	defaultRoom := glob.RoomData{
		Name:        "Default room",
		Description: "This is the default room.",
		Valid:       true,
	}

	defaultRoom.Players = make(map[string]*glob.PlayerData)
	defaultSector.Rooms = make(map[int]glob.RoomData)

	glob.SectorsList[1] = defaultSector
	glob.SectorsList[1].Rooms[1] = defaultRoom

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

func NewRound() (wait <-chan struct{}) {
	ch := make(chan struct{})
	go func() {
		time.Sleep(def.INPUT_THROTTLE * time.Millisecond)
		close(ch) // Broadcast to all receivers.
	}()
	return ch
}

func mainLoop() {
	rand.Seed(time.Now().UTC().UnixNano())

	/* Player rounds */
	go func() {
		for {
			glob.Round = NewRound()
			<-glob.Round
		}
	}()

	/*Background tasks*/
	go func() {
		for {
			time.Sleep(time.Minute)

			/*--- LOCK ---*/
			glob.ConnectionListLock.Lock()
			/*--- LOCK ---*/

			for x := 0; x <= glob.PlayerListEnd; x++ {
				if glob.PlayerList[x] != nil && glob.PlayerList[x].Valid == false {
					player := glob.PlayerList[x]
					if player.UnlinkedTime.IsZero() == false && time.Since(player.UnlinkedTime) > (2*time.Minute) {
						support.WriteToRoom(player, "fades into nothing...")
						support.RemovePlayerWorld(player)
					}
				}
			}

			/*--- UNLOCK ---*/
			glob.ConnectionListLock.Unlock()
			/*--- UNLOCK ---*/
		}
	}()

	/*Wait for new connections*/
	for glob.ServerState == def.SERVER_RUNNING {
		/*Create throttle*/
		startTime := time.Now()

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

		/* Limit to this speed.
		 * But don't sleep unless needed,
		 * so we stay responsive */
		sleepFor := time.Until(startTime.Add(def.CONNECT_THROTTLE * time.Millisecond))
		ranFor := time.Since(startTime)

		buf := fmt.Sprintf("connect: ran for %v", ranFor.String())
		log.Println(buf)

		buf = fmt.Sprintf("connect: sleeping for %v", sleepFor.String())
		log.Println(buf)

		time.Sleep(sleepFor)
	}
}
