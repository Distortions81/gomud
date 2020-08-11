package support

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"../def"
	"../glob"
	"../mlog"
)

func GetObjectFromID(sector int, id int) (*glob.ObjectData, bool) {

	sData := glob.SectorsList[sector]
	if sData.Valid {
		oData := sData.Objects[id]
		if oData != nil && oData.Valid {
			return oData, true
		}
	}

	return nil, false
}

func ParseVnum(player *glob.PlayerData, input string) (int, int, bool) {
	a, b := SplitArgsTwo(input, ":")
	sector := 0
	id := 0
	var erra error
	var errb error

	if len(input) < 1 {
		return 0, 0, true
	}

	if b == "" {
		sector = player.Location.Sector
		id, erra = strconv.Atoi(a)
	} else {
		sector, _ = strconv.Atoi(a)
		id, errb = strconv.Atoi(b)
	}

	err := false
	if erra != nil || errb != nil {
		err = true
	}
	return sector, id, err
}

func MovingExpAvg(value, oldValue, fdtime, ftime float64) float64 {
	alpha := 1.0 - math.Exp(-fdtime/ftime)
	r := alpha*value + (1.0-alpha)*oldValue
	return r
}

func CheckError(source string, err error, fatal bool) {
	if err != nil {
		buf := fmt.Sprintf("error: %s: %s", source, err.Error())
		mlog.Write(buf)
		if fatal {
			os.Exit(1)
		}
	}
}

func RoundSinceTime(roundTo string, value time.Time) string {
	since := time.Since(value)
	if roundTo == "h" {
		since -= since % time.Hour
	}
	if roundTo == "m" {
		since -= since % time.Minute
	}
	if roundTo == "s" {
		since -= since % time.Second
	}
	return since.String()
}

func MakeFingerprint(prefix string) string {
	if prefix != "" {
		prefix = prefix + "-"
	}
	fingerprint := fmt.Sprintf("%v%v-%v", prefix, time.Now().UnixNano(), rand.Uint64())
	return fingerprint
}

func ScaleBytes(b int) string {

	output := "Error"

	if b >= 1024 { //kb
		output = fmt.Sprintf("%vkb", b/1024)
	} else if b >= 1024*1024 { //mb
		output = fmt.Sprintf("%vmb", b/1024/1024)
	} else if b >= 1024*1024*1024 { //gb
		output = fmt.Sprintf("%vgb", b/1024/1024/1024)
	} else if b >= 1024*1024*1024*1024 { //pb
		output = fmt.Sprintf("%vpb", b/1024/1024/1024/1024)
	} else { //b
		output = fmt.Sprintf("%vb", b)
	}

	return output
}

func IsStandardDirection(input string) bool {
	command := strings.ToLower(input)

	if command == "north" || command == "south" ||
		command == "east" || command == "west" ||
		command == "up" || command == "down" {
		return true
	} else {
		return false
	}
}

func GetStandardDirectionMirror(input string) string {
	command := strings.ToLower(input)
	output := "Error"

	if command == "north" {
		output = "South"
	} else if command == "south" {
		output = "North"
	} else if command == "east" {
		output = "West"
	} else if command == "west" {
		output = "East"
	} else if command == "up" {
		output = "Down"
	} else if command == "down" {
		output = "Up"
	}

	return output
}

func FindClosestMatch(CommandList []string, command string) (string, int) {
	var output [def.MAX_CMATCH_SEARCH + 1]string

	command = strings.ToLower(command)
	if command == "" || CommandList == nil || len(CommandList) <= 1 {
		return "", -1
	}

	//Find shortest unique name for a command from a list of commands
	for pos, aCmd := range CommandList {
		aName := aCmd
		aLen := len(aName)
		maxMatch := 1

		for x := 0; x < aLen; x++ { //Search up to full length of name
			for _, bCmd := range CommandList { //Search all commands except ourself
				bName := bCmd
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
			if pos >= def.MAX_CMATCH_SEARCH-1 {
				break
			}
		}
		output[pos] = (aName[0 : maxMatch+1])

		for y, cmd := range CommandList {
			if strings.HasPrefix(command, output[y]) && strings.HasPrefix(cmd, command) {
				return cmd, y
			}
		}
	}

	//No results
	return "", -1

}

func boolToOnOff(toggle bool) string {
	if toggle {
		return "{GON{x  "
	} else {
		return "{rOFF{x "
	}
}

func boolToYesNo(toggle bool) string {
	if toggle {
		return "{GYES{x "
	} else {
		return "{RNO{x  "
	}
}

func boolToTrueFalse(toggle bool) string {
	if toggle == true {
		return "{YTRUE{x"
	} else {
		return "{RFALSE{x"
	}
}

func NewBool(a bool) *bool {
	b := a
	return &b
}
