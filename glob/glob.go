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
var ConnectionList [def.MAX_DESCRIPTORS + 1]ConnectionData
var ConnectionListLock deadlock.Mutex

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

	temp   string      `json:"-,"`
	Player *PlayerData `json:"-,"`
	Valid  bool        `json:"-,"`
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

	Desc  *net.TCPConn `json:"-,"`
	Valid bool         `json:"-,"`
}
