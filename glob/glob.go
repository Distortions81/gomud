package glob

import (
	"net"
	"time"

	"../def"
	"github.com/sasha-s/go-deadlock"
)

var ServerState = def.SERVER_RUNNING
var ServerListener *net.TCPListener

var ConnectionListMax int
var ConnectionList []ConnectionData
var ConnectionListLock deadlock.RWMutex

type ConnectionData struct {
	Name          string
	Desc          *net.TCPConn
	Address       string
	State         int
	ConnectedTime time.Time
	IdleTime      time.Time
	Id            int
	BytesOut      int
	BytesIn       int

	Player *PlayerData
	Valid  bool
}

type PlayerData struct {
	Name          string
	Description   string
	State         int
	ConnectedTime time.Time
	IdleTime      time.Time
	Admin         bool

	Desc  *net.Conn
	Valid bool
}
