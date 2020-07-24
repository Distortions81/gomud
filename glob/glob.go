package glob

import (
	"net"
	"sync"
	"time"

	"../def"
)

/*The big dataset*/
var ServerState = def.SERVER_RUNNING
var ServerListener *net.TCPListener
var ServerListenerSSL net.Listener
var Round <-chan struct{}

var ConnectionListEnd int
var ConnectionList [def.MAX_USERS]ConnectionData
var ConnectionListLock sync.Mutex

var PlayerListEnd int
var PlayerList [def.MAX_USERS]*PlayerData
var PlayerListLock sync.Mutex

var SectorsListEnd int
var SectorsList [def.MAX_SECTORS]SectorData

var QuickHelp string

type DirectionData struct {
	Name         string
	ToRoom       *RoomData `json:"-"`
	ToRoomID     int       `json:",omitempty"`
	ToRoomSector int       `json:",omitempty"`

	Closed bool `json:",omitempty"`
	Hidden bool `json:",omitempty"`
	Keyed  bool `json:",omitempty"`

	Builders map[string]time.Time `json:",omitempty"`

	Valid bool
}

type RoomData struct {
	Name        string                 `json:",omitempty"`
	Description string                 `json:",omitempty"`
	Players     map[string]*PlayerData `json:"-"`

	//Convert to map?
	Exits map[string]DirectionData `json:",omitempty"`

	Builders map[string]time.Time `json:",omitempty"`

	Valid bool
}

type SectorData struct {
	Version string

	ID          int
	Fingerprint string `json:",omitempty"`

	Name        string `json:",omitempty"`
	Area        string `json:",omitempty"`
	Description string `json:",omitempty"`

	Rooms map[int]RoomData `json:",omitempty"`

	Valid bool
}

type ConnectionData struct {
	Name    string
	Desc    net.Conn `json:"-"`
	Address string
	SSL     bool

	State        int
	ConnectedFor time.Time
	IdleTime     time.Time

	BytesOut int
	BytesIn  int

	BytesOutRecorded int
	BytesInRecorded  int

	TempPass string      `json:"-"`
	Player   *PlayerData `json:"-"`
	Valid    bool
}

type PlayerData struct {
	Version     string
	Fingerprint string
	Name        string
	Password    string

	PlayerType int
	Level      int
	State      int
	Sector     int
	Room       int
	RoomLink   *RoomData `json:"-"`

	Created      time.Time
	LastSeen     time.Time
	TimePlayed   int
	UnlinkedTime time.Time `json:"-"`

	Connections map[string]int
	BytesIn     map[string]int
	BytesOut    map[string]int
	Email       string `json:",omitempty"`

	Description string `json:",omitempty"`
	Sex         string `json:",omitempty"`

	Connection *ConnectionData `json:"-"`
	Valid      bool
}
