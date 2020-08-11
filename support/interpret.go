package support

import (
	"fmt"
	"log"
	"strings"

	"../glob"
)

var ShortCutList []glob.Command

func CreateShortcuts() {

	//Find shortest unique name for a command from a list of commands
	wa := 0
	wb := 0
	wc := 0
	wd := 0

	for pos, aCmd := range CommandList {
		wa++
		if aCmd.Short != "" || aCmd.AS == false {
			continue
		}
		aName := strings.ToLower(aCmd.Name)
		aLen := len(aName)
		maxMatch := 1

		for x := 0; x < aLen; x++ { //Search up to full length of name
			wb++
			for _, bCmd := range CommandList { //Search all commands except ourself
				wc++
				if bCmd.AS == false {
					continue
				}
				bName := strings.ToLower(bCmd.Name)
				bLen := len(bName)
				if x > bLen { //If we have reached max length of B, stop
					continue
				}
				if bName != aName {
					if aName[0:x] == bName[0:x] {
						maxMatch = x
					}
				}
			}
		}
		wd++
		CommandList[pos].Short = (aName[0 : maxMatch+1])
	}
	log.Println(fmt.Sprintf("CreateShortcuts: %v:%v:%v-%v", wa, wb, wc, wd))
}

func PlayerCommand(player *glob.PlayerData, command string, args string, isAlias bool) {

	if player != nil && player.Valid {
		for _, cmd := range CommandList {

			command = strings.ToLower(command)
			short := strings.ToLower(cmd.Short)

			//Don't allow alias loop
			if cmd.Name == "alias" && isAlias {
				continue
			}

			//Check if they are allowed to use the command
			if cmd.Type > player.PlayerType {
				//continue
			}

			if cmd.AS == false && strings.EqualFold(command, cmd.Name) {
				cmd.Cmd(player, args)
				return
			} else if strings.HasPrefix(command, short) && strings.HasPrefix(cmd.Name, command) && isAlias == false && cmd.AS == true {
				//aliases don't get shortcuts
				//If begining of the input matches with record for shortest unique name for the command,
				//If input is longer, if that also matches the full command name.
				cmd.Cmd(player, args)
				return
			}
		}

		WriteToPlayer(player, "Invalid command.")
		CmdCommands(player, "")
	}
}
