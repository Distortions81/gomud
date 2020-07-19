package glob

import (
	"bufio"
	"net"
	"sync"
	"time"

	"../def"
)

var ServerState = def.SERVER_RUNNING
var ServerListener *net.TCPListener
var ConnectionList []ConnectionData
var ConnectionListLock sync.Mutex

var LastConnectionID int

type ConnectionData struct {
	Name          string
	Desc          *net.TCPConn
	Address       string
	State         int
	ConnectedTime time.Time
	IdleTime      time.Time
	Reader        *bufio.Reader
	Writer        *bufio.Writer
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
