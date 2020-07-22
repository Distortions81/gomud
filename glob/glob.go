package glob

import (
	"net"
	"sync"
	"time"

	"../def"
)

var ServerState = def.SERVER_RUNNING
var ServerListener *net.TCPListener

//Fixed size arrays are faster
var ConnectionListMax int
var ConnectionList [def.MAX_USERS + 1]ConnectionData
var ConnectionListLock sync.Mutex

var PlayerListMax int
var PlayerList [def.MAX_USERS + 1]PlayerData
var PlayerListLock sync.Mutex

var SectorsListMax int
var SectorsList [def.MAX_SECTORS + 1]SectorsData

type BuilderData struct {
	Builders []string
	Modified []time.Time

	CreatedBy string
	Created   Time.time

	Valid bool
}

type DirectionData struct {
	ToRoom       *RoomData
	ToRoomID     int
	ToRoomSector int

	Closed bool
	Hidden bool
	Keyed  bool

	Builders BuilderData

	Valid bool
}

type RoomData struct {
	RoomID   int
	SectorID int

	Name        string
	Description string

	North DirectionData
	South DirectionData
	East  DirectionData
	West  DirectionData
	Up    DirectionData
	Down  DirectionData

	Builders BuilderData

	Valid bool
}

type SectorsData struct {
	ID    string
	Group string

	Name        string
	Description string

	NumRooms int
	Rooms    [def.MAX_ROOMS_PER_SECTOR]RoomData

	Valid bool
}

type ConnectionData struct {
	Name    string
	Desc    *net.TCPConn `json:"-,"`
	Address string

	State        int
	ConnectedFor time.Time
	IdleTime     time.Time
	ID           int

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
