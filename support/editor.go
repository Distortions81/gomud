package support

import (
	"fmt"
	"strconv"
	"strings"

	"../def"
	"../glob"
)

func MleEditor(player *glob.PlayerData, input string) {
	command, argTwoThrough := SplitArgsTwo(input, " ")
	cmdB, _ := SplitArgsTwo(argTwoThrough, " ")
	cmdl := strings.ToLower(command)
	//cmdBl := strings.ToLower(cmdB)

	player.CurEdit.Active = true
	e := &player.CurEdit

	if player.CurEdit.Lines == nil {
		player.CurEdit.Lines = make(map[int]string)
	}

	if e.CurLine > e.NumLines {
		e.CurLine = e.NumLines
	}
	if cmdl == "deleteall" {
		e.Active = false
		e.CurLine = 0
		e.NumLines = 0
		e.Lines[0] = ""
		WriteToPlayer(player, "Note cleared.")
	} else if cmdl == "done" {
		e.Active = false
		newDesc := ""
		if e.CallBackP != nil {
			for x := 0; x <= player.CurEdit.NumLines; x++ {
				newDesc = newDesc + player.CurEdit.Lines[x] + "\r\n"
			}
			*e.CallBackP = newDesc
			WriteToPlayer(player, "Edit sent to "+e.CallBack+".")
			e.CallBackP = nil
			e.CallBack = ""
		}
		if player.CurEdit.CallBack != "" {
			WriteToPlayer(player, "Exiting editor, unable to send edit to "+player.CurEdit.CallBack+". Use the paste command instead.")
		}
		e.CallBackP = nil
		e.CallBack = ""
		return
	} else if cmdl == "line" {
		num, err := strconv.Atoi(cmdB)
		if err == nil {
			e.CurLine = num - 1
		}
	} else if cmdl == "add" {
		if e.NumLines < def.MAX_MLE {
			e.CurLine = e.NumLines
			e.NumLines++
			e.CurLine++
			e.Lines[e.NumLines] = argTwoThrough
		} else {
			WriteToPlayer(player, "Sorry, 100 lines max.")
		}
	} else if cmdl == "remove" {
		e.Lines[e.CurLine] = ""
		if e.CurLine > e.NumLines {
			for x := e.CurLine; x < e.NumLines; x++ {
				e.Lines[x] = e.Lines[x+1]
			}
		}
		e.NumLines--
		e.CurLine--
	} else if cmdl == "insert" {
		e.NumLines++
		for x := e.NumLines; x > e.NumLines; x-- {
			e.Lines[x] = e.Lines[x-1]
		}
		e.Lines[e.CurLine] = argTwoThrough

	} else if cmdl == "replace" {
		e.Lines[e.CurLine] = argTwoThrough
	} else if cmdl == "colorcodes" {
		if e.ColorCodes {
			e.ColorCodes = false
		} else {
			e.ColorCodes = true
		}
	}

	//Print
	WriteToPlayer(player, "Syntax: line <line num>\r\n        <command> <text>\r\nCommands: add, remove, insert, replace, deleteAll, colorCodes, done")
	cl := " "
	for i := 0; i <= e.NumLines; i++ {
		if i == e.CurLine {
			cl = "@"
		} else {
			cl = " "
		}
		if e.ColorCodes {
			WriteToPlayerCodes(player, fmt.Sprintf("%3v %v: %v", i+1, cl, e.Lines[i]))
		} else {
			WriteToPlayer(player, fmt.Sprintf("%3v %v: %v", i+1, cl, e.Lines[i]))
		}
	}
}
