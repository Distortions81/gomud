package glob

import (
	"net"
	"time"

	"../def"
)

type DoorData struct {
	Door bool `json:",omitempty"`

	Closed    bool `json:",omitempty"`
	AutoOpen  bool `json:",omitempty"`
	AutoClose bool `json:",omitempty"`

	Hidden bool `json:",omitempty"`
	Keyed  bool `json:",omitempty"`
}

type ExitData struct {
	ToRoom LocationData

	Door DoorData `json:",omitempty"`

	//function to print
	//fucntion to parse
	Valid bool `json:",omitempty"`
}

type RoomData struct {
	//Location    LocationData           `json:"-"`
	Name        string                 `json:",omitempty"`
	Description string                 `json:",omitempty"`
	Players     map[string]*PlayerData `json:"-"`

	//Convert to map?
	Exits map[string]*ExitData

	Valid bool
}

type SectorData struct {
	Version string

	NumRooms    int
	ID          int
	Fingerprint string `json:",omitempty"`

	Name        string `json:",omitempty"`
	Area        string `json:",omitempty"`
	Description string `json:",omitempty"`

	Rooms map[int]*RoomData `json:",omitempty"`
	Dirty bool              `json:"-"`

	Valid bool
}

type InputBuffer struct {
	BufferInPos   int `json:"-"`
	BufferInCount int `json:"-"`

	BufferOutPos   int                             `json:"-"`
	BufferOutCount int                             `json:"-"`
	InputBuffer    [def.MAX_INPUT_LINES + 1]string `json:"-"`
}

type ConnectionData struct {
	Input InputBuffer `json:"-"`

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

type OLCSettingData struct {
	//Forward all input to OLC
	NoOLCPrefix bool
	//Tell users how to exit olc
	NoHint bool
	//Automatically switch room editor to current room
	OLCRoomFollow bool
	//Show color codes in OLC
	OLCShowCodes bool
	//Show color codes for whOLC world
	OLCShowAllCodes bool
	//OLC Promt enable
	OLCPrompt bool
	//OLC prompt string
	OLCPromptString string
}

type SettingsData struct {
	//Color
	Ansi bool
	//Short direction names, no room desc
	Brief bool
	//PromptString
	PromptString string
	//Hide prompt except in battle
	PromptHide bool
	//Telnet backspace/clear line prompt out of history
	PromptDelete bool
	//page wait, if one "round" of text exceeds N lines
	Paging int
	//Global chat off
	Deafen bool
	//long short brief off
	Affects int
	//none, friends, clan, all
	WhoHide int
	//Newline before commands
	PreNewline bool
	//Newline after commands
	PostNewline bool
}

type PlayerData struct {
	Version     string
	Fingerprint string
	Name        string
	Password    string
	Dirty       bool `json:"-"`

	PlayerType int `json:",omitempty"`
	Level      int
	State      int
	Location   LocationData
	Recall     LocationData `json:",omitempty"`

	Created      time.Time
	LastSeen     time.Time
	TimePlayed   int
	UnlinkedTime time.Time `json:"-"`
	OLCEdit      OLCEdit   `json:"-"`

	Aliases     map[string]string `json:",omitempty"`
	Connections map[string]int
	BytesIn     map[string]int
	BytesOut    map[string]int

	Config      PConfigData    `json:",omitempty"`
	OLCSettings OLCSettingData `json:",omitempty"`

	Email string `json:",omitempty"`

	Description string `json:",omitempty"`
	Sex         string `json:",omitempty"`

	Connection *ConnectionData `json:"-"`
	Banned     bool            `json:",omitempty"`
	Valid      bool            `json:"-"`
}

type OLCEdit struct {
	Active   bool `json:",omitempty"`
	Mode     int  `json:",omitempty"`
	EditDesc bool `json:",omitempty"`

	/*Current selection & past selections*/
	Sector int `json:",omitempty"`

	Room     LocationData `json:",omitempty"`
	Object   LocationData `json:",omitempty"`
	Trigger  LocationData `json:",omitempty"`
	Mobile   LocationData `json:",omitempty"`
	Quest    LocationData `json:",omitempty"`
	Exit     *ExitData    `json:",omitempty"`
	ExitName string       `json:",omitempty"`

	Description string `json:",omitempty"`
}

type LocationData struct {
	Sector int `json:",omitempty"`
	ID     int `json:",omitempty"`

	RoomLink *RoomData `json:"-"`
	//Function to print
	//Function to parse
}

type Command struct {
	AS    bool
	Short string
	Name  string
	Cmd   func(player *PlayerData, args string)
	Type  int
	Help  string
}

type pTypeData struct {
	PType int
	PName string
}

type ConfigData struct {
	ID   int
	Name string
	Help string
	Ref  *bool
	RefS *string
	RefI *int
}

type PConfigData struct {
	Ansi         bool
	Brief        bool
	PromptString string
	PromptHide   bool
	PromptDelete bool
	Paging       int
	Deafen       bool
	Affects      int
	WhoHide      int
	PreNewline   bool
	PostNewline  bool
}
