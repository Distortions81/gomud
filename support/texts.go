package support

import (
	"io/ioutil"
	"log"

	"../def"
	"../glob"
)

func ReadTextFiles() {
	ReadGreet()
	ReadAuRevoir()
	ReadNews()
}
func ReadGreet() {
	message, err := ioutil.ReadFile(def.DATA_DIR + def.TEXTS_DIR + def.GREET_FILE)
	if err != nil {
		log.Println("Unable to read greeting file.")
		glob.Greeting = "Welcome!"
		return
	} else {
		glob.Greeting = ANSIColor(string(message))
	}
	log.Println("ReadGreet: Greeting loaded.")
}

func ReadAuRevoir() {
	message, err := ioutil.ReadFile(def.DATA_DIR + def.TEXTS_DIR + def.AUREVOIR_FILE)
	if err != nil {
		log.Println("Unable to read AuRevoir file.")
		glob.AuRevoir = "Farewell!"
		return
	} else {
		glob.AuRevoir = ANSIColor(string(message))
	}
	log.Println("ReadAuRevoir: AuRevoir loaded.")
}

func ReadNews() {
	message, err := ioutil.ReadFile(def.DATA_DIR + def.TEXTS_DIR + def.NEWS_FILE)
	if err != nil {
		log.Println("Unable to read news file.")
		glob.News = "No news is good news."
		return
	} else {
		glob.News = string(message)
	}
	log.Println("ReadNews: news loaded.")
}
