package support

import (
	"fmt"
	"strings"

	"../glob"
)

func OLCConfig(player *glob.PlayerData,
	input string, command string, cmdB string, cmdl string, cmdBl string,
	argTwoThrough string, argThreeThrough string) {
	OLCSettings := []glob.ConfigData{
		{ID: 1, Name: "follow", Help: "If on: you are always editing the room you are standing in.",
			Ref: &player.OLCSettings.OLCRoomFollow},
		{ID: 2, Name: "showCodes", Help: "If on: Show color codes in names / descriptions / etc",
			Ref: &player.OLCSettings.OLCShowCodes},
		//{ID: 3, Name: "showAllCodes", Help: "If on: Show color codes, instead of color for the whOLC mud.",
		//Ref: player.OLCSettings.OLCShowAllCodes},
		{ID: 4, Name: "prompt", Help: "If on: Change your prompt to OLC information while in editor.",
			Ref: &player.OLCSettings.OLCPrompt},
		//{ID: 5, Name: "promptString", Help: "Customize OLC prompt.",
		//Ref: player.OLCSettings.OLCPromptString},
		{ID: 6, Name: "noOLCPrefix", Help: "If on: When in editor, all input goes to olc by default.",
			Ref: &player.OLCSettings.NoOLCPrefix},
		{ID: 7, Name: "noHint", Help: "If on: Turn off message explaining STOP/CMD for NoOLCPrefix",
			Ref: &player.OLCSettings.NoHint},
	}

	cmdNames := []string{}
	for _, c := range OLCSettings {
		cmdNames = append(cmdNames, strings.ToLower(c.Name))
	}
	match, _ := FindClosestMatch(cmdNames, argTwoThrough)
	player.Dirty = true

	if match == "follow" {
		if player.OLCSettings.OLCRoomFollow {
			player.OLCSettings.OLCRoomFollow = false
			WriteToPlayer(player, "OLC will no longer change the room you are editing when you move.")
		} else {
			player.OLCSettings.OLCRoomFollow = true
			WriteToPlayer(player, "OLC will automatically edit whatever room you move to.")
		}
	} else if match == "showcodes" {
		if player.OLCSettings.OLCShowCodes {
			player.OLCSettings.OLCShowCodes = false
			WriteToPlayer(player, "OLC will now just show normal color.")
		} else {
			player.OLCSettings.OLCShowCodes = true
			WriteToPlayer(player, "OLC will show color codes in names and descriptions.")
		}
	} else if match == "prompt" {
		if player.OLCSettings.OLCPrompt {
			player.OLCSettings.OLCPrompt = false
			WriteToPlayer(player, "Your prompt will no longer change to OLC prompt while editing.")
		} else {
			player.OLCSettings.OLCPrompt = true
			WriteToPlayer(player, "Your prompt will now be OLC information.")
		}
	} else if match == "noolcprefix" {
		if player.OLCSettings.NoOLCPrefix {
			player.OLCSettings.NoOLCPrefix = false
			WriteToPlayer(player, "Your input will NOT be sent directly to olc. Prefix all OLC commands with: olc <command>.")
		} else {
			player.OLCSettings.NoOLCPrefix = true
			WriteToPlayer(player, "All your input will be directed to OLC, until you exit it.\r\ncmd <command> will pass-through commands.")
		}
	} else if match == "nohint" {
		if player.OLCSettings.NoHint {
			player.OLCSettings.NoHint = false
			WriteToPlayer(player, "After every line, remind you of STOP/CMD commands")
		} else {
			player.OLCSettings.NoHint = true
			WriteToPlayer(player, "Do not show reminder for STOP/CMD commands")
		}
	}

	//Show settings avaialble
	for _, cmd := range OLCSettings {
		WriteToPlayer(player, fmt.Sprintf("%10v:%5v --  %v", cmd.Name, boolToOnOff(*cmd.Ref), cmd.Help))
	}
	return
}
