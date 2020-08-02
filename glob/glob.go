package glob

import (
	"net"
	"os"
	"sync"

	"../def"
)

/*The big dataset*/
var OpenDesc int
var OpenDescLock sync.Mutex

var MudLog *os.File

var ServerState = def.SERVER_RUNNING
var ServerListener *net.TCPListener
var ServerListenerSSL net.Listener
var Round <-chan struct{}

var ConnectionListEnd int
var ConnectionList [def.MAX_USERS + 1]ConnectionData
var ConnectionListLock sync.Mutex

var PlayerListEnd int
var PlayerList [def.MAX_USERS + 1]*PlayerData

var SectorsListEnd int
var SectorsList [def.MAX_SECTORS]SectorData

var QuickHelp string
var WizHelp string

//Texts
var Greeting string
var AuRevoir string
var News string
