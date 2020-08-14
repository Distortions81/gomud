package support

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"../glob"
)

func CmdShutdown(player *glob.PlayerData, args string) {
	glob.SignalHandle <- syscall.SIGINT
}

func CmdPerfStat(player *glob.PlayerData, args string) {
	glob.PerLock.Lock() //PERF-LOCK
	WriteToPlayer(player, glob.PerfStats)
	glob.PerLock.Unlock() //PERF-UNLOCK
}

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

func CmdReloadHelpst(player *glob.PlayerData, args string) {
	ReadHelps()
	WriteToPlayer(player, "Help file reloaded.")
}

func CmdReloadText(player *glob.PlayerData, args string) {
	ReadTextFiles()
	WriteToPlayer(player, "Text files reloaded.")
}

func CmdReloadPlayer(player *glob.PlayerData, args string) {
	for i := 1; i <= glob.PlayerListEnd; i++ {
		target := glob.PlayerList[i]

		if strings.EqualFold(target.Name, args) {
			rtarget, found := ReadPlayer(strings.ToLower(args), true)
			if found {
				glob.PlayerList[i] = rtarget
				LinkPlayerConnection(rtarget, target.Connection)
				WriteToPlayer(target, "Your character file has been re-loaded.")
				WriteToPlayer(player, "Player reloaded.")
				player.Dirty = true
				return
			}
		}
	}
	WriteToPlayer(player, "I don't see that player online.")
}

func CmdPlayerType(player *glob.PlayerData, args string) {
	pname, level := SplitArgsTwo(args, " ")

	plevel, err := strconv.Atoi(level)

	if err != nil {
		WriteToPlayer(player, "Syntax: playerType <playerName> <typeNumber>")
		return
	}

	//TODO: NAMED TYPES
	pname = strings.ToLower(pname)
	for x := 1; x <= glob.PlayerListEnd; x++ {
		target := glob.PlayerList[x]
		if strings.EqualFold(target.Name, pname) {
			target.PlayerType = plevel
			WriteToPlayer(player, "Player type set.")
			WriteToPlayer(target, "Your player-type has been changed.")
			player.Dirty = true
			return
		}
	}
	WriteToPlayer(player, "I couldn't find anyone online by that name.")

}

func CmdSavePlayers(player *glob.PlayerData, args string) {
	for x := 1; x <= glob.PlayerListEnd; x++ {
		target := glob.PlayerList[x]
		if target != nil && target.Valid {
			WritePlayer(target, true)
			WriteToPlayer(player, player.Name+" was saved.")
		}
	}
}
