import "time"

func createPlayer(desc *glob.ConnectionData) glob.PlayerData {
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
