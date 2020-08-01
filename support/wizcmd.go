package support

import (
	"fmt"

	"../glob"
)

func CmdWizHelp(player *glob.PlayerData, args string) {
	WriteToPlayer(player, glob.WizHelp)
}

func CmdStats(player *glob.PlayerData, args string) {
	output := ""

	header := fmt.Sprintf("%-5v %25v: %16v (%4vc) %v/%v%v\r\n", "#", "Name", "ip", "count", "in", "out", "SSL")
	for x := 1; x <= glob.ConnectionListEnd; x++ {
		buf := ""
		con := &glob.ConnectionList[x]
		if con.Valid {
			target := con.Player
			ssl := ""

			if con.SSL {
				ssl = " (SSL)"
			}

			if target != nil {
				for key, value := range target.Connections {
					buf = fmt.Sprintf("%-5v %25v: %16v (%5v) %v/%v%v\r\n", x, target.Name, key, value, ScaleBytes(target.BytesIn[key]), ScaleBytes(target.BytesOut[key]), ssl)
				}
			} else if con != nil {
				buf = fmt.Sprintf("%-5v %25v: %16v (%5v) %v/%v%v\r\n", x, con.Name, "", "", ScaleBytes(con.BytesIn), ScaleBytes(con.BytesOut), ssl)
			}
		} else {
			buf = fmt.Sprintf("%-5v %25v: %16v (%5v) %v/%v%v\r\n", x, "Closed", "none", "0", "0", "0", "")
		}
		output = output + buf
	}
	WriteToPlayer(player, header+output)
}

func CmdReloadText(player *glob.PlayerData, args string) {
	ReadTextFiles()
	WriteToPlayer(player, "Text files reloaded.")
}
