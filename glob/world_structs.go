package glob

type SectorData struct {
	Version string `json:",omitempty"`

	NumRooms    int    `json:",omitempty"`
	ID          int    `json:",omitempty"`
	Fingerprint string `json:",omitempty"`

	Name        string `json:",omitempty"`
	ColorName   string `json:",omitempty"`
	Area        string `json:",omitempty"`
	Description string `json:",omitempty"`

	Rooms   map[int]*RoomData   `json:",omitempty"`
	Objects map[int]*ObjectData `json:",omitempty"`
	Resets  map[int]*ResetsData `json:",omitempty"`

	Dirty bool `json:"-"`

	Valid bool `json:",omitempty"`
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

type LocationData struct {
	Sector int `json:",omitempty"`
	ID     int `json:",omitempty"`

	RoomLink *RoomData `json:"-"`
	//Function to print
	//Function to parse
	Valid bool `json:"-"`
}

type ResetsData struct {
	Name string `json:",omitempty"`

	Sector int `json:",omitempty"`
	ObjID  int `json:",omitempty"`
	//MobID int

	Quanity  int
	Interval string

	RoomLink *RoomData `json:"-"`
	Valid    bool
}
