package glob

import (
	"net"
	"time"
)

/*Connections*/
var ConnectionList []ConnectionData
var ServerListener *TCPListener

type ConnectionData struct {
	connection    net.Conn
	address       string
	state         int
	connectedTime time.Time
	idleTime      time.Time

	player *PlayerData
	valid  bool
}

type PlayerData struct {
	name          string
	description   string
	state         int
	connectedTime time.Time
	idleTime      time.Time
	admin         bool

	connection *net.Conn
	valid      bool
}
