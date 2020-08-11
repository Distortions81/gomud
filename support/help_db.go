package support

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"../def"
	"../glob"
	"../mlog"
)

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
	fileName := def.DATA_DIR + def.TEXTS_DIR + def.HELPS_FILE

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

	_, err := os.Stat(def.DATA_DIR + def.TEXTS_DIR + def.HELPS_FILE)
	notfound := os.IsNotExist(err)

	if notfound {
		CheckError("ReadHelps: os.Stat", err, def.ERROR_NONFATAL)
		mlog.Write("Help file not found!")
		return false

	} else {

		file, err := ioutil.ReadFile(def.DATA_DIR + def.TEXTS_DIR + def.HELPS_FILE)

		if file != nil && err == nil {
			helps := CreateHelps()
			helps.Valid = true

			err := json.Unmarshal([]byte(file), &helps)
			if err != nil {
				CheckError("ReadPlayer: Unmashal", err, def.ERROR_NONFATAL)
			}

			if helps.Topics == nil {
				helps.Topics = make(map[string]*glob.HelpTopics)
			}

			for x, _ := range helps.Topics {
				helps.Topics[x].Valid = true
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
	help.Topics = make(map[string]*glob.HelpTopics)

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
