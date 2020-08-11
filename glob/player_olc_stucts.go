package glob

type WearLocations struct {
	Name          string `json:",omitempty"`
	ID            int    `json:",omitempty"`
	WearMessage   string `json:",omitempty"`
	RemoveMessage string `json:",omitempty"`
	LookDesc      string `json:",omitempty"`

	ConflictLocationA int `json:",omitempty"`
	ConflictLocationB int `json:",omitempty"`

	Valid bool `json:"-"`
}

type EditLink struct {
	Name   string `json:",omitempty"`
	Sector int    `json:",omitempty"`
	ID     int    `json:",omitempty"`

	RoomLink   *RoomData   `json:"-"`
	ObjectLink *ObjectData `json:"-"`
	//TriggerLink *TriggerData `json:"-"`
	//MobileLink  *MobileData `json:"-"`
	ExitLink *ExitData `json:"-"`
}

type OLCEdit struct {
	Active bool `json:",omitempty"`
	Mode   int  `json:",omitempty"`

	Room   EditLink `json:",omitempty"`
	Object EditLink `json:",omitempty"`
	//Trigger EditLink `json:",omitempty"`
	//Mobile  EditLink `json:",omitempty"`
	//Quest   EditLink `json:",omitempty"`
	Exit   EditLink `json:",omitempty"`
	Sector EditLink `json:",omitempty"`
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
