package glob

import (
	"net"
	"sync"
	"time"

	"../def"
	//"github.com/sasha-s/go-deadlock"
)

var ServerState = def.SERVER_RUNNING
var ServerListener *net.TCPListener

//Fixed size arrays are faster
var ConnectionListMax int
var ConnectionList [def.MAX_USERS + 1]ConnectionData
var ConnectionListLock sync.RWMutex

var PlayerListMax int
var PlayerList [def.MAX_USERS + 1]PlayerData
var PlayerListLock sync.RWMutex

type ConnectionData struct {
	Name    string
	Desc    *net.TCPConn `json:"-,"`
	Address string

	State        int
	ConnectedFor time.Time
	IdleTime     time.Time
	Id           int

	BytesOut int
	BytesIn  int

	TempPass string      `json:"-,"`
	Player   *PlayerData `json:"-,"`
	Valid    bool        `json:"-,"`
}

type PlayerData struct {
	Name     string
	Password string

	PlayerType int
	Level      int
	State      int
	Sector     int
	Vnum       int

	Created     time.Time
	LastSeen    time.Time
	Seconds     int
	IPs         []string
	Connections []int
	BytesIn     []int
	BytesOut    []int
	Email       string

	Description string
	Sex         string

	Connection *ConnectionData `json:"-,"`
	Valid      bool            `json:"-,"`
}
