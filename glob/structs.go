package glob

import (
	"time"
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

type LocationData struct {
	Sector int `json:",omitempty"`
	ID     int `json:",omitempty"`

	RoomLink *RoomData `json:"-"`
	//Function to print
	//Function to parse
	Valid bool `json:"-"`
}

type Command struct {
	AS    bool                                  `json:"-"`
	Short string                                `json:"-"`
	Name  string                                `json:"-"`
	Cmd   func(player *PlayerData, args string) `json:"-"`
	Type  int                                   `json:"-"`
	Help  string                                `json:"-"`
	Valid bool                                  `json:"-"`
}
