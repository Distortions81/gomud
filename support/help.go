package support

import (
	"fmt"
	"strconv"
	"strings"

	"../glob"
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

	if input == "" {
		WriteToPlayer(player, "Help topics:")
		for topicName, _ := range glob.HelpSystem.Topics {
			WriteToPlayer(player, topicName)
		}
		return
	}

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
					buf := fmt.Sprintf("Chapter: %v", chapterName)
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
