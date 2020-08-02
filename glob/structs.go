package glob

import (
	"net"
	"time"

	"../def"
)

type HelpMain struct {
	Version string

	Preface string
	Topics  map[string]HelpTopics

	//Keyword links?

	Dirty bool `json:"-"`
}

type HelpTopics struct {
	Name string `json:"-"`
	Desc string

	Author  string
	Created time.Time

	//Time, name
	EditHistory    map[string]string
	Changes        map[string]string
	QuickReference string

	Preface      string
	Chapters     map[string]HelpPage
	TermAbbrUsed map[string]string
	Footnotes    map[string]string
}

type HelpPage struct {
	Keywords map[int]string
	Pages    map[int]string
}

type MleData struct {
	Active     bool `json:"-"`
	Lines      map[int]string
	ColorCodes bool `json:",omitempty"`

	NumLines  int     `json:",omitempty"`
	CurLine   int     `json:",omitempty"`
	CallBackP *string `json:"-"`
	CallBack  string  `json:",omitempty"`
}

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
	Name        string                 `json:",omitempty"`
	Description string                 `json:",omitempty"`
	Players     map[string]*PlayerData `json:"-"`

	//Convert to map?
	Exits map[string]*ExitData `json:",omitempty"`

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

	Name    string   `json:",omitempty"`
	Desc    net.Conn `json:"-"`
	Address string   `json:",omitempty"`
	SSL     bool     `json:",omitempty"`

	State        int       `json:",omitempty"`
	ConnectedFor time.Time `json:",omitempty"`
	IdleTime     time.Time `json:",omitempty"`

	BytesOut int `json:",omitempty"`
	BytesIn  int `json:",omitempty"`

	BytesOutRecorded int `json:",omitempty"`
	BytesInRecorded  int `json:",omitempty"`

	TempPass   string      `json:"-"`
	TempPlayer *PlayerData `json:"-"`
	Player     *PlayerData `json:"-"`
	Valid      bool
}

type OLCSettingData struct {
	//Forward all input to OLC
	NoOLCPrefix bool `json:",omitempty"`
	//Tell users how to exit olc
	NoHint bool `json:",omitempty"`
	//Automatically switch room editor to current room
	OLCRoomFollow bool `json:",omitempty"`
	//Show color codes in OLC
	OLCShowCodes bool `json:",omitempty"`
	//Show color codes for whOLC world
	OLCShowAllCodes bool `json:",omitempty"`
	//OLC Promt enable
	OLCPrompt bool `json:",omitempty"`
	//OLC prompt string
	OLCPromptString string `json:",omitempty"`
}

type SettingsData struct {
	//Color
	Ansi bool `json:",omitempty"`
	//Short direction names, no room desc
	Brief bool `json:",omitempty"`
	//PromptString
	PromptString string `json:",omitempty"`
	//Hide prompt except in battle
	PromptHide bool `json:",omitempty"`
	//Telnet backspace/clear line prompt out of history
	PromptDelete bool `json:",omitempty"`
	//page wait, if one "round" of text exceeds N lines
	Paging int `json:",omitempty"`
	//Global chat off
	Deafen bool `json:",omitempty"`
	//long short brief off
	Affects int `json:",omitempty"`
	//none, friends, clan, all
	WhoHide int `json:",omitempty"`
	//Newline before commands
	PreNewline bool `json:",omitempty"`
	//Newline after commands
	PostNewline bool `json:",omitempty"`
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
	OLCEdit      OLCEdit   `json:",omitempty"`
	CurEdit      MleData   `json:",omitempty"`

	Aliases     map[string]string
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
	Active bool `json:"-"`
	Mode   int  `json:"-"`

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
	PType int    `json:",omitempty"`
	PName string `json:",omitempty"`
}

type ConfigData struct {
	ID   int     `json:",omitempty"`
	Name string  `json:",omitempty"`
	Help string  `json:",omitempty"`
	Ref  *bool   `json:",omitempty"`
	RefS *string `json:",omitempty"`
	RefI *int    `json:",omitempty"`
}

type PConfigData struct {
	Ansi         bool   `json:",omitempty"`
	Brief        bool   `json:",omitempty"`
	PromptString string `json:",omitempty"`
	PromptHide   bool   `json:",omitempty"`
	PromptDelete bool   `json:",omitempty"`
	Paging       int    `json:",omitempty"`
	Deafen       bool   `json:",omitempty"`
	Affects      int    `json:",omitempty"`
	WhoHide      int    `json:",omitempty"`
	PreNewline   bool   `json:",omitempty"`
	PostNewline  bool   `json:",omitempty"`
}
