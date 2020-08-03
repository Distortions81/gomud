package glob

import (
	"net"
	"time"

	"../def"
)

type HelpMain struct {
	Version string `json:",omitempty"`

	Preface string                 `json:",omitempty"`
	Topics  map[string]*HelpTopics `json:",omitempty"`

	//Keyword links?

	Dirty bool `json:"-"`
	Valid bool `json:"-"`
}

type HelpTopics struct {
	Name string `json:"-"`
	Desc string `json:",omitempty"`

	Author  string    `json:",omitempty"`
	Created time.Time `json:",omitempty"`

	//Time, name
	EditHistory    map[string]string `json:",omitempty"`
	Changes        map[string]string `json:",omitempty"`
	QuickReference string            `json:",omitempty"`

	Preface      string              `json:",omitempty"`
	Chapters     map[string]HelpPage `json:",omitempty"`
	TermAbbrUsed map[string]string   `json:",omitempty"`
	Footnotes    map[string]string   `json:",omitempty"`

	Valid bool `json:"-"`
}

type HelpPage struct {
	Keywords map[int]string `json:",omitempty"`
	Pages    map[int]string `json:",omitempty"`

	Valid bool `json:"-"`
}

type MleData struct {
	Active     bool `json:"-"`
	Lines      map[int]string
	ColorCodes bool `json:",omitempty"`

	NumLines  int     `json:",omitempty"`
	CurLine   int     `json:",omitempty"`
	CallBackP *string `json:"-"`
	CallBack  string  `json:",omitempty"`

	Valid bool `json:"-"`
}

type DoorData struct {
	Door bool `json:",omitempty"`

	OpenString  string `json:",omitempty"`
	CloseString string `json:",omitempty"`

	Closed    bool `json:",omitempty"`
	AutoOpen  bool `json:",omitempty"`
	AutoClose bool `json:",omitempty"`

	Hidden bool `json:",omitempty"`
	Keyed  bool `json:",omitempty"`

	Valid bool `json:"-"`
}

type ExitData struct {
	RoomP          *RoomData `json:"-"`
	ColorName      string    `json:",omitempty"`
	Description    string    `json:",omitempty"`
	ExitFromString string    `json:",omitempty"`
	EnterToString  string    `json:",omitempty"`

	ToRoom LocationData `json:",omitempty"`
	Door   *DoorData    `json:",omitempty"`
	Hidden bool         `json:",omitempty"`

	Valid bool `json:"-"`
}

type RoomData struct {
	SectorP     *SectorData            `json:"-"`
	Name        string                 `json:",omitempty"`
	ColorName   string                 `json:",omitempty"`
	Description string                 `json:",omitempty"`
	Players     map[string]*PlayerData `json:"-"`
	Objects     map[string]*ObjectData `json:"-"`
	PermObjects map[string]*ObjectData `json:",omitempty"`

	Exits map[string]*ExitData `json:",omitempty"`

	Valid bool `json:"-"`
}

type SectorData struct {
	Version string `json:",omitempty"`

	NumRooms    int    `json:",omitempty"`
	ID          int    `json:",omitempty"`
	Fingerprint string `json:",omitempty"`

	Name        string `json:",omitempty"`
	ColorName   string `json:",omitempty"`
	Area        string `json:",omitempty"`
	Description string `json:",omitempty"`

	Rooms   map[int]*RoomData      `json:",omitempty"`
	Objects map[string]*ObjectData `json:",omitempty"`

	Dirty bool `json:"-"`

	Valid bool `json:"-"`
}

type InputBuffer struct {
	BufferInPos   int `json:"-"`
	BufferInCount int `json:"-"`

	BufferOutPos   int                             `json:"-"`
	BufferOutCount int                             `json:"-"`
	InputBuffer    [def.MAX_INPUT_LINES + 1]string `json:"-"`

	Valid bool ` json:"-"`
}

type ConnectionData struct {
	Input InputBuffer `json:"-"`

	Name    string   `json:"-"`
	Desc    net.Conn `json:"-"`
	Address string   `json:"-"`
	SSL     bool     `json:"-"`

	State        int       `json:"-"`
	ConnectedFor time.Time `json:"-"`
	IdleTime     time.Time `json:"-"`

	BytesOut int `json:"-"`
	BytesIn  int `json:"-"`

	BytesOutRecorded int `json:"-"`
	BytesInRecorded  int `json:"-"`

	TempPass   string      `json:"-"`
	TempPlayer *PlayerData `json:"-"`
	Player     *PlayerData `json:"-"`

	Valid bool `json:"-"`
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

	Valid bool `json:"-"`
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

	Valid bool `json:"-"`
}

type ObjectData struct {
	Sector      int             `json:",omitempty"`
	ID          int             `json:",omitempty"`
	Fingerprint string          `json:",omitempty"`
	Container   ObjectContainer `json:",omitempty"`

	Name        string `json:",omitempty"`
	Description string `json:",omitempty"`
	ColorName   string `json:",omitempty"`
	Owner       string `json:",omitempty"`

	Player     *PlayerData  `json:"-"`
	InRoom     *RoomData    `json:"-"`
	Location   LocationData `json:",omitempty"`
	Persistant bool         `json:",omitempty"`

	PlayerRestring string         `json:",omitempty"`
	Owners         map[int]string `json:",omitempty"`

	Type          int     `json:",omitempty"`
	PlayerUseable bool    `json:",omitempty"`
	PlayerTake    bool    `json:",omitempty"`
	Unique        bool    `json:",omitempty"`
	Weight        int     `json:",omitempty"`
	SlotsUsed     int     `json:",omitempty"`
	Bound         bool    `json:",omitempty"`
	Health        float64 `json:",omitempty"`
	HealthPerUse  float64 `json:",omitempty"`
	WearSlot      int     `json:",omitempty"`

	Valid bool `json:"-"`
}

type ObjectContainer struct {
	Name string `json:",omitempty"`

	Contents  map[string]ObjectData `json:",omitempty"`
	Slots     int                   `json:",omitempty"`
	MaxWeight int                   `json:",omitempty"`

	Closeable bool `json:",omitempty"`
	Closed    bool `json:",omitempty"`

	Valid bool `json:"-"`
}

type PlayerData struct {
	Version     string `json:",omitempty"`
	Fingerprint string `json:",omitempty"`

	Name      string `json:",omitempty"`
	ColorName string `json:",omitempty"`
	Password  string `json:",omitempty"`
	Dirty     bool   `json:"-"`

	PlayerType int          `json:",omitempty"`
	Level      int          `json:",omitempty"`
	State      int          `json:",omitempty"`
	Location   LocationData `json:",omitempty"`
	Recall     LocationData `json:",omitempty"`

	Created      time.Time `json:",omitempty"`
	TimePlayed   int       `json:",omitempty"`
	LastSeen     time.Time `json:",omitempty"`
	UnlinkedTime time.Time `json:"-"`
	OLCEdit      OLCEdit   `json:",omitempty"`
	CurEdit      MleData   `json:",omitempty"`

	Aliases     map[string]string `json:",omitempty"`
	Connections map[string]int    `json:",omitempty"`
	BytesIn     map[string]int    `json:",omitempty"`
	BytesOut    map[string]int    `json:",omitempty"`

	Config      PConfigData    `json:",omitempty"`
	OLCSettings OLCSettingData `json:",omitempty"`

	Email string `json:",omitempty"`

	Description string `json:",omitempty"`
	Sex         string `json:",omitempty"`

	Connection *ConnectionData `json:"-"`
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
	Valid       bool   `json:"-"`
}

type LocationData struct {
	Sector int `json:",omitempty"`
	ID     int `json:",omitempty"`

	RoomLink *RoomData `json:"-"`
	//Function to print
	//Function to parse
	Valid bool `json:"-"`
}

type Command struct {
	AS    bool                                  `json:",omitempty"`
	Short string                                `json:",omitempty"`
	Name  string                                `json:",omitempty"`
	Cmd   func(player *PlayerData, args string) `json:"-"`
	Type  int                                   `json:",omitempty"`
	Help  string                                `json:",omitempty"`
	Valid bool                                  `json:"-"`
}

type pTypeData struct {
	PType int    `json:",omitempty"`
	PName string `json:",omitempty"`
	Valid bool   `json:"-"`
}

type ConfigData struct {
	ID    int     `json:",omitempty"`
	Name  string  `json:",omitempty"`
	Help  string  `json:",omitempty"`
	Ref   *bool   `json:",omitempty"`
	RefS  *string `json:",omitempty"`
	RefI  *int    `json:",omitempty"`
	Valid bool    `json:"-"`
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
	Valid        bool   `json:"-"`
}
