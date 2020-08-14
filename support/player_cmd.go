package support

import (
	"fmt"
	"strings"
	"time"

	"../def"
	"../glob"
)

func CmdNews(player *glob.PlayerData, input string) {
	WriteToPlayer(player, glob.News)
}

func CmdRelog(player *glob.PlayerData, input string) {
	c := player.Connection
	CmdSave(player, "")
	RemovePlayer(player)
	c.Input.BufferInPos = 0
	c.Input.BufferInCount = 0
	c.Input.BufferOutPos = 0
	c.Input.BufferOutCount = 0

	for p := 0; p < def.MAX_INPUT_LINES; p++ {
		c.Input.InputBuffer[p] = ""
	}

	c.State = def.CON_STATE_RELOG
	c.Valid = true

	WriteToDesc(c, "")
	WriteToDesc(c, "")
	WriteToDesc(c, "")
	WriteToDesc(c, glob.Greeting)

	WriteToDesc(c, "(Type NEW to create character) Name:")
	return
}

func CmdAlias(player *glob.PlayerData, input string) {
	command, longArgs := SplitArgsTwo(input, " ")
	firstArg, lastArgs := SplitArgsTwo(longArgs, " ")

	if command == "list" {
		aliases := ""
		for key, value := range player.Aliases {
			aliases = aliases + fmt.Sprintf("%s: %s\r\n", key, value)
		}
		if aliases == "" {
			aliases = "None"
		}
		WriteToPlayer(player, "Aliases:\r\n"+aliases)
	} else if command == "add" {

		/*Prevent problems*/
		if firstArg == "" {
			WriteToPlayer(player, "The alias needs a name.")
			return
		} else if firstArg == "alias" {
			WriteToPlayer(player, "That would be a bad idea.")
			return
		}
		if len(lastArgs) > 80 {
			WriteToPlayer(player, "That output is too long.")
			return
		}

		/*Write data to player*/
		player.Aliases[firstArg] = lastArgs
		WriteToPlayer(player, "Alias added.")
		player.Dirty = true
	} else if command == "delete" {
		found := false
		for key, _ := range player.Aliases {
			if strings.EqualFold(key, firstArg) {
				found = true
				break
			}
		}
		if found {
			delete(player.Aliases, firstArg)
			WriteToPlayer(player, "Alias deleted")
			player.Dirty = true
		}
	} else {
		WriteToPlayer(player, "Aliases can be the same name as commands,")
		WriteToPlayer(player, "but you can still alias the orignal command to something else.")
		WriteToPlayer(player, "Aliases can not call other aliases, or the alias command,")
		WriteToPlayer(player, "and must use full-length command names (no shorthand).")
		WriteToPlayer(player, "")
		WriteToPlayer(player, "alias add <shortcut> <output> (max 80 characters)")
		WriteToPlayer(player, "alias delete <shortcut>")
		WriteToPlayer(player, "alias list")
	}
}

func CmdCommands(player *glob.PlayerData, args string) {
	WriteToPlayer(player, glob.QuickHelp)
}

func CmdQuit(player *glob.PlayerData, args string) {
	okay := WritePlayer(player, true)
	if okay == false {
		WriteToPlayer(player, "Saving character failed!!!")
		return //Don't quit if we couldn't save
	} else {
		WriteToPlayer(player, "Character saved.")
	}
	buf := fmt.Sprintf("%s has quit.", player.Name)
	WriteToAll(buf)
	player.Connection.State = def.CON_STATE_DISCONNECTING
}

func CmdWho(player *glob.PlayerData, args string) {
	buf := "Players online:"
	WriteToPlayer(player, buf)

	pos := 0
	for x := 1; x <= glob.ConnectionListEnd; x++ {
		var p *glob.ConnectionData = &glob.ConnectionList[x]
		if p.Valid == false {
			continue
		}

		if p.State == def.CON_STATE_PLAYING {
			idleString := ""
			connectedString := ""

			if time.Since(p.IdleTime) > time.Minute {
				idleString = " (idle " + RoundSinceTime("m", p.IdleTime) + ")"
			}
			if time.Since(p.ConnectedFor) > time.Minute {
				connectedString = " (on " + RoundSinceTime("m", p.ConnectedFor) + ")"
			}
			pos++
			buf = fmt.Sprintf("%d: %s%s%s", pos, p.Name, connectedString, idleString)
			WriteToPlayer(player, buf)
		} else {
			pos++
			buf = fmt.Sprintf("%d: %s", pos, "(Connecting)")
			WriteToPlayer(player, buf)
		}
	}

	buf = fmt.Sprintf("Server uptime: %v", time.Since(glob.BootTime).Round(time.Second).String())
	WriteToPlayer(player, "")
	WriteToPlayer(player, buf)
}

func CmdSave(player *glob.PlayerData, args string) {
	if player.ReqSave == false {
		WriteToPlayer(player, "Saving character...")
	}
	player.Dirty = true
	player.ReqSave = true
}
