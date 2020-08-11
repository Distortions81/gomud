package support

import (
	"fmt"
	"log"
	"strings"

	"../def"
	"../glob"
)

func MakeQuickHelp() {
	buf := "Commands:\r\n"
	buf = buf + fmt.Sprintf("%-5v:%12v : %-48v%10v\r\n", "short", "command", "Help info", "Type")
	buf = buf + def.LINESEPB

	for _, cmd := range CommandList {
		ptype := ""
		if cmd.Type >= 700 {
			continue
			//ptype = " " + GetPTypeString(cmd.Type)
		}
		help, _ := TruncateString(cmd.Help, 48)
		short, _ := TruncateString(strings.ToLower(cmd.Short), 5)
		buf = buf + fmt.Sprintf("%-5v:%12v : %-48v%10v\r\n", short, strings.ToLower(cmd.Name), help, ptype)
	}
	buf = buf + "\r\nCommands that require arguments will show extended help, if run with no arguments."
	glob.QuickHelp = buf
	log.Println("MakeQuickHelp: QuickHelp loaded.")
}

func MakeWizHelp() {
	buf := "Commands:\r\n"
	buf = buf + fmt.Sprintf("%-5v:%12v : %-48v%10v\r\n", "short", "command", "Help info", "Type")
	buf = buf + def.LINESEPB

	for _, cmd := range CommandList {
		ptype := ""
		if cmd.Type >= 700 {
			ptype = " " + GetPTypeString(cmd.Type)
		} else {
			continue
		}
		help, _ := TruncateString(cmd.Help, 48)
		short, _ := TruncateString(strings.ToLower(cmd.Short), 5)
		buf = buf + fmt.Sprintf("%-5v:%12v : %-48v%10v\r\n", short, strings.ToLower(cmd.Name), help, ptype)
	}
	buf = buf + "\r\nCommands that require arguments will show extended help, if run with no arguments."
	glob.WizHelp = buf
	log.Println("MakeWizHelp: WizHelp loaded.")
}
