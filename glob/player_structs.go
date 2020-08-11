package glob

import (
	"net"
	"time"

	"../def"
)

type PlayerData struct {
	Version     string `json:",omitempty"`
	Fingerprint string `json:",omitempty"`

	Name      string `json:",omitempty"`
	ColorName string `json:",omitempty"`
	Password  string `json:",omitempty"`
	Dirty     bool   `json:"-"`
	ReqSave   bool   `json:"-"`

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

type InputBuffer struct {
	BufferInPos   int `json:"-"`
	BufferInCount int `json:"-"`

	BufferOutPos   int                             `json:"-"`
	BufferOutCount int                             `json:"-"`
	InputBuffer    [def.MAX_INPUT_LINES + 1]string `json:"-"`

	Valid bool ` json:"-"`
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
	Ref   *bool   `json:"-"`
	RefS  *string `json:"-"`
	RefI  *int    `json:"-"`
	Valid bool    `json:"-"`
}

type MleData struct {
	Active     bool           `json:",omitempty"`
	Lines      map[int]string `json:",omitempty"`
	ColorCodes bool           `json:",omitempty"`

	NumLines  int     `json:",omitempty"`
	CurLine   int     `json:",omitempty"`
	CallBackP *string `json:"-"`
	CallBack  string  `json:",omitempty"`

	Valid bool `json:"-"`
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

	Valid bool `json:"-"`
}
