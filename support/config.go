package support

import (
	"fmt"
	"strconv"
	"strings"

	"../glob"
)

func CmdConfig(player *glob.PlayerData, input string) {
	command, longArg := SplitArgsTwo(input, " ")

	var PlayerConfig = []glob.ConfigData{
		{ID: 100, Name: "Color", Help: "ANSI color", Ref: &player.Config.Ansi},
		{ID: 200, Name: "Brief", Help: "Compact direction names, no room descriptions while walking.", Ref: &player.Config.Brief},
		{ID: 300, Name: "PromptString", Help: "Customize your prompt", RefS: &player.Config.PromptString},
		{ID: 400, Name: "PromptHide", Help: "Hide the prompt", Ref: &player.Config.PromptHide},
		{ID: 500, Name: "PromptDelete", Help: "Try to erase prompt from scroll.", Ref: &player.Config.Ansi},
		{ID: 600, Name: "Paging", Help: "If more than (X) lines in a command/round, prompt to continue.", RefI: &player.Config.Paging},
		{ID: 700, Name: "Deafen", Help: "Mute global messages", Ref: &player.Config.Deafen},
		{ID: 800, Name: "Affects", Help: "Full, Normal or Compact", RefI: &player.Config.Affects},
		{ID: 900, Name: "WhoVis", Help: "All, Friends & Clan, Friends, Clan, None", RefI: &player.Config.WhoHide},
		{ID: 1000, Name: "PreNewline", Help: "Blank line before commands", Ref: &player.Config.PreNewline},
		{ID: 1100, Name: "PostNewline", Help: "Blank line after commands", Ref: &player.Config.PostNewline},
	}
	if PlayerConfig != nil {

		for p, cfg := range PlayerConfig {
			n, err := strconv.Atoi(input)
			if input != "" && (strings.EqualFold(command, cfg.Name) || (p+1 == n && err == nil)) {

				setCfg(player, PlayerConfig, p, longArg)
				WriteToPlayer(player, fmt.Sprintf("%v is now %v", cfg.Name, printCfgType(PlayerConfig, p)))
				break

			}
		}
		for p, cfg := range PlayerConfig {
			WriteToPlayer(player, fmt.Sprintf("%15v (%-3v)[ %v] -- %v", cfg.Name, p+1, printCfgType(PlayerConfig, p), cfg.Help))
		}
		WriteToPlayer(player, "config <name/number> to toggle, config <name/number> <number/text>")

	}

}

func printCfgType(pc []glob.ConfigData, p int) string {
	isOn := "Error"
	if pc[p].Ref != nil {
		isOn = boolToOnOff(*pc[p].Ref)
	} else if pc[p].RefI != nil {
		isOn = fmt.Sprintf("{Y%-4v{x", *pc[p].RefI)
	} else if pc[p].RefS != nil && *pc[p].RefS != "" {
		isOn = "{GSET{x "
	} else {
		isOn = "{rOFF{x "
	}

	return isOn
}

func setCfg(player *glob.PlayerData, pc []glob.ConfigData, p int, longArg string) {
	if pc[p].Ref != nil {
		if *pc[p].Ref {
			*pc[p].Ref = false
		} else {
			*pc[p].Ref = true
		}
	} else if pc[p].RefI != nil {
		i, err := strconv.Atoi(longArg)
		if err == nil {
			*pc[p].RefI = i
		} else {
			WriteToPlayer(player, "Error: Syntax: config <option> <number>")
		}
	} else if pc[p].RefS != nil && *pc[p].RefS != "" {
		i, err := strconv.Atoi(longArg)
		if err == nil {
			*pc[p].RefI = i
		} else {
			WriteToPlayer(player, "Error: Syntax: config <option> <number>")
		}
	}
}
