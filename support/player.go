package support

import (
	"io/ioutil"
	"strings"
	"time"

	"../def"
	"../glob"
)

func CreatePlayer(desc *glob.ConnectionData) glob.PlayerData {
	player := glob.PlayerData{
		Name:        def.STRING_UNKNOWN,
		Password:    "",
		PlayerType:  def.PLAYER_TYPE_NEW,
		Level:       0,
		State:       def.PLAYER_ALIVE,
		Sector:      0,
		Vnum:        0,
		Created:     time.Now(),
		LastSeen:    time.Now(),
		Seconds:     0,
		IPs:         []string{},
		Connections: []int{},
		BytesIn:     []int{},
		BytesOut:    []int{},
		Email:       "",

		Description: "",
		Sex:         "",

		Desc:  nil,
		Valid: true,
	}
	return player
}

func ReadPlayer(name string, load bool) (*glob.PlayerData, bool) {

	pdata, err := ioutil.ReadFile(def.PLAYER_DIR + name)
	if err != nil {
		return nil, false
	} else {
		if load == true {
			lines := strings.Split(string(pdata), ";")
			numlines := len(lines)

			var args [def.PFILE_MAXARGS]string
			for x := 0; x <= numlines; x++ {
				args[x] = strings.Split(string(lines[x]), ",")
			}

			for x := 0; x <= numlines; x++ {
				cline := strings.ReplaceAll(line[x], "\n\r", "")
				data := line[x]
			}

		} else {
			return nil, true
		}
	}

