package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"./def"
	"./glob"
	"./mlog"
	"./support"
)

func setupListenerSSL() {
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
	buf := fmt.Sprintf("SSL listener online at: %s", def.DEFAULT_PORT_SSL)
	mlog.Write(buf)
}

func setupListener() {
	/*Find Network*/
	addr, err := net.ResolveTCPAddr("tcp4", def.DEFAULT_PORT)
	support.CheckError("main: resolveTCP", err, def.ERROR_FATAL)

	/*Open Listener*/
	listener, err := net.ListenTCP("tcp4", addr)
	glob.ServerListener = listener
	support.CheckError("main: ListenTCP", err, def.ERROR_FATAL)

	/*Print Connection*/
	buf := fmt.Sprintf("TCP listener online at: %s", addr.String())
	mlog.Write(buf)
}

func WaitNewConnectionSSL() {

	for glob.ServerState == def.SERVER_RUNNING {

		time.Sleep(def.CONNECT_THROTTLE_MS * time.Millisecond)
		desc, err := glob.ServerListenerSSL.Accept()
		support.AddNetDesc()
		time.Sleep(def.CONNECT_THROTTLE_MS * time.Millisecond)

		/* If there is a connection flood, sleep listeners */
		if err != nil || support.CheckNetDesc() {
			time.Sleep(5 * time.Second)
			desc.Close()
			support.RemoveNetDesc()
		} else {

			_, err = desc.Write([]byte(
				"You have connected to GOMud: " + def.VERSION + ", port " + def.DEFAULT_PORT_SSL + " (With SSL!)\r\n" + glob.Greeting +
					"(Type NEW to create character) Name:"))
			time.Sleep(def.CONNECT_THROTTLE_MS * time.Millisecond)
			support.NewDescriptor(desc, true)
		}

	}

	glob.ServerListenerSSL.Close()
}

func WaitNewConnection() {

	for glob.ServerState == def.SERVER_RUNNING {

		time.Sleep(def.CONNECT_THROTTLE_MS * time.Millisecond)
		desc, err := glob.ServerListener.Accept()
		support.AddNetDesc()
		time.Sleep(def.CONNECT_THROTTLE_MS * time.Millisecond)

		/* If there is a connection flood, sleep listeners */
		if err != nil || support.CheckNetDesc() {
			time.Sleep(5 * time.Second)
			desc.Close()
			support.RemoveNetDesc()
		} else {

			_, err = desc.Write([]byte(
				"You have connected to GOMud: " + def.VERSION + ", port " + def.DEFAULT_PORT + " (insecure telnet)\r\n" +
					"If your client supports it, please use port " + def.DEFAULT_PORT_SSL + ", for AES-256 encryption!\r\n" +
					glob.Greeting +
					"(Type NEW to create character) Name:"))

			time.Sleep(def.CONNECT_THROTTLE_MS * time.Millisecond)
			support.NewDescriptor(desc, false)
		}
	}

	glob.ServerListener.Close()
}

func main() {

	var err error

	t := time.Now()

	logName := fmt.Sprintf("log/%v-%v-%v.log", t.Day(), t.Month(), t.Year())
	glob.MudLog, err = os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Unable to open log file!")
		os.Exit(1)
		return
	}

	support.CreateShortcuts()
	support.MakeQuickHelp()
	support.MakeWizHelp()
	support.ReadSectorList()
	support.ReadTextFiles()

	setupListener()
	setupListenerSSL()
	go WaitNewConnection()
	go WaitNewConnectionSSL()

	/*Process connections*/
	mainLoop()

	//After starting loops, wait here for process signals
	sc := make(chan os.Signal, 1)

	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	glob.ConnectionListLock.Lock()
	//glob.ServerState = def.SERVER_CLOSING
	support.WriteToAll("Server killed!")

	/*Save everything*/
	for x := 1; x <= glob.ConnectionListEnd; x++ {
		con := glob.ConnectionList[x]
		if con.Player != nil && con.Player.Valid {
			support.WritePlayer(con.Player)
			if con.Desc != nil {
				con.Desc.Write([]byte("Your character has been saved!\r\n"))
			}
		}
	}
	support.WriteSectorList()
	support.WriteToAll("All sectors saved!")
	support.WriteToAll("")
	support.WriteToAll(glob.AuRevoir)

	glob.ConnectionListLock.Unlock()
	time.Sleep(time.Second)
}

func mainLoop() {
	rand.Seed(time.Now().UTC().UnixNano())

	/* Player rounds */
	go func() {
		numPlayerLast := glob.ConnectionListEnd
		sleepFor := time.Duration(def.ROUND_LENGTH_uS)
		for glob.ServerState == def.SERVER_RUNNING {

			glob.ConnectionListLock.Lock() /*--- LOCK ---*/

			/*Handle 0 players*/
			if numPlayerLast <= 0 {
				sleepFor = time.Duration(def.ROUND_LENGTH_uS) * time.Microsecond
				time.Sleep(sleepFor)
			} else {
				sleepFor = time.Duration(def.ROUND_LENGTH_uS/numPlayerLast) * time.Microsecond
			}
			cEnd := glob.ConnectionListEnd
			glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/

			tempCount := 0
			for x := 0; x <= cEnd; x++ {

				glob.ConnectionListLock.Lock() /*--- LOCK ---*/
				if glob.ConnectionList[x].Valid {
					start := time.Now()
					tempCount++

					/*Check for stale connections*/
					if time.Since(glob.ConnectionList[x].ConnectedFor) > (def.WELCOME_TIMEOUT_S*time.Second) &&
						glob.ConnectionList[x].State <= def.CON_STATE_WELCOME {
						glob.ConnectionList[x].Valid = false
						glob.ConnectionList[x].Desc.Close()
						support.RemoveNetDesc()
					}

					support.ReadPlayerInputBuffer(&glob.ConnectionList[x])

					glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/
					end := time.Now()
					spent := end.Sub(start) /*Round sleep, slice per player*/
					time.Sleep(sleepFor - spent)
					glob.ConnectionListLock.Lock() /*--- LOCK ---*/
				}
				glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/
			}

			numPlayerLast = tempCount
			time.Sleep(def.ROUND_REST_MS) /*Limit max CPU, and allow background to run*/

		}
	}()

	/*Background tasks*/
	go func() {
		for glob.ServerState == def.SERVER_RUNNING {
			/* ONCE A MINUTE */
			time.Sleep(time.Minute)

			/*--- LOCK ---*/
			glob.ConnectionListLock.Lock()
			/*--- LOCK ---*/

			/*Autosave sectors*/
			for x := 1; x <= glob.SectorsListEnd; x++ {
				if glob.SectorsList[x].Dirty {
					support.WriteSector(&glob.SectorsList[x])
				}
			}

			/*Autosave players*/
			for x := 1; x <= glob.PlayerListEnd; x++ {
				if glob.PlayerList[x].Dirty {
					support.WritePlayer(glob.PlayerList[x])
				}
			}

			/*Remove disconnected players after a while*/
			for x := 1; x <= glob.PlayerListEnd; x++ {
				player := glob.PlayerList[x]

				if player != nil &&
					player.Valid &&
					player.Location.RoomLink != nil {

					if player.Connection.Valid == false &&
						player.UnlinkedTime.IsZero() == false &&
						time.Since(player.UnlinkedTime) > (2*time.Minute) {

						player.UnlinkedTime = time.Time{}
						support.WritePlayer(player)
						support.WriteToRoom(player, fmt.Sprintf("%s fades into nothing...", player.Name))
						support.RemovePlayer(player)
					}
				}
			}

			/*--- UNLOCK ---*/
			glob.ConnectionListLock.Unlock()
			/*--- UNLOCK ---*/
		}
	}()
}
