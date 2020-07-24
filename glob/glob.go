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

type DoorData struct {
	Door bool

	Closed    bool
	AutoOpen  bool
	AutoClose bool

	Hidden bool
	Keyed  bool
}

type ExitData struct {
	Name   string
	ToRoom LocationData

	Door DoorData `json:",omitempty"`

	Builders map[string]time.Time `json:",omitempty"`

	//function to print
	//fucntion to parse
	Valid bool
}

type RoomData struct {
	Location    LocationData
	Name        string                 `json:",omitempty"`
	Description string                 `json:",omitempty"`
	Players     map[string]*PlayerData `json:"-"`

	//Convert to map?
	Exits    map[string]ExitData  `json:",omitempty"`
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

	TempPass   string      `json:"-"`
	TempPlayer *PlayerData `json:"-"`
	Player     *PlayerData `json:"-"`
	Valid      bool
}

type PlayerData struct {
	Version     string
	Fingerprint string
	Name        string
	Password    string

	PlayerType int
	Level      int
	State      int
	Location   LocationData

	OLCEdit OLCEdit `json:",omitempty"`

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

type OLCEdit struct {
	Active   bool
	Mode     int
	EditDesc bool
	AutoDig  bool

	/*Current selection & past selections*/
	Sector int

	Room    LocationData
	Object  LocationData
	Trigger LocationData
	Mobile  LocationData
	Quest   LocationData

	Description string

	Exit string
}

type LocationData struct {
	Sector int
	ID     int

	RoomLink *RoomData `json:"-"`
	Valid    bool
	//Function to print
	//Function to parse
}
