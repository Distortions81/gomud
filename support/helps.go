package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"../def"
	"../glob"
	"../mlog"
)

func CmdGetHelps(player *glob.PlayerData, input string) {

	input = strings.ToLower(input)

	argOne, argTwoAndOn := SplitArgsTwo(input, " ")
	argTwo, _ := SplitArgsTwo(argTwoAndOn, " ")

	pageNum, err := strconv.Atoi(argTwo)
	if err != nil {
		pageNum = 1
	}

	found := false
	totalPages := 1

	//TODO
	//Show closest match, allow topic:chapter searching

	for topicName, topicData := range glob.HelpSystem.Topics {
		for chapterName, chapterData := range topicData.Chapters {
			if strings.HasPrefix(strings.ToLower(chapterName), argOne) {
				totalPages = len(chapterData.Pages)
				found = true

				buf := fmt.Sprintf("Topic: %v\r\nChapter: %v\r\nPage: %v of %v\r\n", topicName, chapterName, pageNum, totalPages)
				WriteToPlayer(player, buf)
				WriteToPlayer(player, chapterData.Pages[pageNum]+"\r\n")
				break
			}
		}
	}
	if found == false {
		WriteToPlayer(player, "No chapter found with that search term, Topics found:")

		for topicName, topicData := range glob.HelpSystem.Topics {
			if strings.HasPrefix(strings.ToLower(topicName), argOne) {
				WriteToPlayer(player, "Topic: "+topicName)
				for chapterName, _ := range topicData.Chapters {
					buf := fmt.Sprintf("Chapter: %v:%v", topicName, chapterName)
					WriteToPlayer(player, buf)
					found = true
				}
				WriteToPlayer(player, "") //Topic spacer
			}
		}
	}

	if found && totalPages > 1 {
		WriteToPlayer(player, "Help: <chapter> <page> to see more.")
	}
	if found == false {
		WriteToPlayer(player, "I couldn't find anything!")
	}
}

func CmdAddHelps(player *glob.PlayerData, input string) {

}

func CmdWriteHelps(player *glob.PlayerData, input string) {
	if WriteHelps() {
		WriteToPlayer(player, "Helps saved!")
	} else {
		WriteToPlayer(player, "Writing helps failed!")
	}
}

func WriteHelps() bool {
	outbuf := new(bytes.Buffer)
	enc := json.NewEncoder(outbuf)
	enc.SetIndent("", "\t")

	glob.HelpSystem.Version = def.HELPS_VERSION
	fileName := def.DATA_DIR + def.HELPS_FILE

	if err := enc.Encode(&glob.HelpSystem); err != nil {
		CheckError("WriteHelps: enc.Encode", err, def.ERROR_NONFATAL)
		return false
	}

	_, err := os.Create(fileName)

	if err != nil {
		CheckError("WriteHelps: os.Create", err, def.ERROR_NONFATAL)
		return false
	}

	//Async write
	go func(outbuf bytes.Buffer) {
		err = ioutil.WriteFile(fileName, []byte(outbuf.String()), 0644)

		if err != nil {
			CheckError("WriteHelps: WriteFile", err, def.ERROR_NONFATAL)
		}

		buf := fmt.Sprintf("Wrote %v, %v.", fileName, ScaleBytes(len(outbuf.String())))
		mlog.Write(buf)
	}(*outbuf)

	glob.HelpSystem.Dirty = false
	return true
}

func ReadHelps() bool {

	_, err := os.Stat(def.DATA_DIR + def.HELPS_FILE)
	notfound := os.IsNotExist(err)

	if notfound {
		CheckError("ReadHelps: os.Stat", err, def.ERROR_NONFATAL)
		mlog.Write("Help file not found!")
		return false

	} else {

		file, err := ioutil.ReadFile(def.DATA_DIR + def.HELPS_FILE)

		if file != nil && err == nil {
			helps := CreateHelps()

			err := json.Unmarshal([]byte(file), &helps)
			if err != nil {
				CheckError("ReadPlayer: Unmashal", err, def.ERROR_NONFATAL)
			}

			if helps.Topics == nil {
				helps.Topics = make(map[string]glob.HelpTopics)
			}

			glob.HelpSystem = helps

			mlog.Write("Helps loaded.")
			return true
		} else {
			CheckError("ReadHelps: ReadFile", err, def.ERROR_NONFATAL)
			return false
		}
	}
}

func CreateHelps() glob.HelpMain {
	help := glob.HelpMain{Version: def.HELPS_VERSION, Dirty: false}
	help.Topics = make(map[string]glob.HelpTopics)

	return help
}

func CreateHelpTopic() glob.HelpTopics {
	topic := glob.HelpTopics{Name: "New"}

	topic.Preface = "This is a new topic."

	topic.EditHistory = make(map[string]string)
	topic.Changes = make(map[string]string)
	topic.Chapters = make(map[string]glob.HelpPage)
	topic.TermAbbrUsed = make(map[string]string)
	topic.Footnotes = make(map[string]string)

	return topic
}

func CreateChapter() glob.HelpPage {
	helpPage := glob.HelpPage{}

	helpPage.Pages = make(map[int]string)
	helpPage.Keywords = make(map[int]string)

	return helpPage
}
