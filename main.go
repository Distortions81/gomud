package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
	"time"

	"./def"
	"./glob"
	"./support"
)

func setupSSL() {
	//openssl ecparam -genkey -name prime256v1 -out server.key
	//openssl req -new -x509 -key server.key -out server.pem -days 3650
	cert, err := tls.LoadX509KeyPair(def.DATA_DIR+def.SSL_PEM, def.DATA_DIR+def.SSL_KEY)
	if err != nil {
		log.Fatal("Error loading certificate. ", err)
	}

	tlsCfg := &tls.Config{Certificates: []tls.Certificate{cert}}

	/*Open Listener*/
	listener, err := tls.Listen("tcp4", def.DEFAULT_PORT_SSL, tlsCfg)
	glob.ServerListenerSSL = listener
	support.CheckError("setupSSL: tls.listen", err, def.ERROR_FATAL)

	/*Print Connection*/
	buf := fmt.Sprintf("SSL Server online at: %s", def.DEFAULT_PORT_SSL)
	log.Println(buf)

	go listenSSL()
}

func listenSSL() {

	for glob.ServerState == def.SERVER_RUNNING {
		conn, err := glob.ServerListenerSSL.Accept()
		if err != nil {
			support.CheckError("listenSSL Accept():", err, def.ERROR_NONFATAL)
			conn.Close()
		} else {
			buf := fmt.Sprintf("New connection from %s.", conn.LocalAddr().String())
			log.Println(buf)

			//TODO Add full greeting/info
			/*Respond here, so we don't have to wait for lock*/
			_, err = conn.Write([]byte(
				"You have connected to GOMud: " + def.VERSION + ", port " + def.DEFAULT_PORT_SSL + " (With SSL!)\r\n" +
					"To create a new character, type: NEW\r\n\r\nName: \r\n"))
			support.CheckError("main: desc.Write", err, def.ERROR_NONFATAL)

			go support.NewDescriptor(conn, true)
		}
	}

	glob.ServerListenerSSL.Close()
}

func main() {

	support.MakeQuickHelp()
	support.ReadSectorList()
	setupSSL()

	//Disabled, for creating inital sector and room
	if 1 == 2 {
		defaultSector := glob.SectorData{

			ID:          1,
			Fingerprint: support.MakeFingerprint("Default-"),
			Area:        "Default",
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

		glob.SectorsListEnd++
	}

	/*Find Network*/
	addr, err := net.ResolveTCPAddr("tcp4", def.DEFAULT_PORT)
	support.CheckError("main: resolveTCP", err, def.ERROR_FATAL)

	/*Open Listener*/
	listener, err := net.ListenTCP("tcp4", addr)
	glob.ServerListener = listener
	support.CheckError("main: ListenTCP", err, def.ERROR_FATAL)

	/*Print Connection*/
	buf := fmt.Sprintf("Server online at: %s", addr.String())
	log.Println(buf)

	/*Listen for connections*/
	mainLoop()
	listener.Close()
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
				player := glob.PlayerList[x]

				if player != nil &&
					player.Valid &&
					player.Location.RoomLink != nil {

					if player.UnlinkedTime.IsZero() == false && time.Since(player.UnlinkedTime) > (2*time.Minute) {
						player.UnlinkedTime = time.Time{}
						support.WriteToRoom(player, fmt.Sprintf("%s fades into nothing...", player.Name))
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
		desc, err := glob.ServerListener.Accept()

		if err == nil {
			buf := fmt.Sprintf("New connection from %s.", desc.LocalAddr().String())
			log.Println(buf)

			//TODO Add full greeting/info
			/*Respond here, so we don't have to wait for lock*/
			_, err = desc.Write([]byte(
				"You have connected to GOMud: " + def.VERSION + ", port " + def.DEFAULT_PORT + " (insecure telnet)\r\n" +
					"If your client supports it, please use port " + def.DEFAULT_PORT_SSL + ", for AES-256 encryption!\r\n" +
					"To create a new character, type: NEW\r\n\r\nName: \r\n"))
			support.CheckError("main: desc.Write", err, def.ERROR_NONFATAL)

			go support.NewDescriptor(desc, false)
			/* netconn/netconn.go */
		}

		/* Limit to this speed.
		 * But don't sleep unless needed,
		 * so we stay responsive */
		sleepFor := time.Until(startTime.Add(def.CONNECT_THROTTLE * time.Millisecond))
		time.Sleep(sleepFor)
	}
}
