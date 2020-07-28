package support

import "../glob"

func ConfigSettings(player *glob.PlayerData, input string) {

	var PlayerSettings = []glob.CommandArgData{
		{Name: "Color", Help: "ANSI color text on/off", Ref: player.Settings.Ansi},
		{Name: "Brief", Help: "Compact direction names, and don't show room descriptions while walking.", Ref: player.Settings.Brief},
		{Name: "PromptString", Help: "Customize your prompt", RefS: player.Settings.PromptString},
		{Name: "PromptHide", Help: "Hide the standard prompt", Ref: player.Settings.PromptHide},
		{Name: "PromptDelete", Help: "Try to erase prompt from chat history.", Ref: player.Settings.Ansi},
		{Name: "Paging", Help: "If more than X lines happen in one command/round, prompt to continue, instead of scrolling off-screen", RefI: player.Settings.Paging},
		{Name: "Deafen", Help: "Mute global messages", Ref: player.Settings.Deafen},
		{Name: "Affects", Help: "full names, 1 letter names in row of dots, or just letters.", RefI: player.Settings.Affects},
		{Name: "WhoHide", Help: "Hide in 'who' command: Normal, friends/clan only, clan-only, everyone", RefI: player.Settings.WhoHide},
		{Name: "PreNewline", Help: "Place a bank line before commands", Ref: player.Settings.PreNewline},
		{Name: "PostNewline", Help: "Place a bank line after commands", Ref: player.Settings.PostNewline},
	}
	if PlayerSettings != nil {
		//
	}

}
