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
}
func ReadGreet() {
	message, err := ioutil.ReadFile(def.DATA_DIR + def.TEXTS_DIR + def.GREET_FILE)
	if err != nil {
		log.Println("Unable to read greeting file.")
		glob.Greeting = "Welcome!"
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
	} else {
		glob.AuRevoir = ANSIColor(string(message))
	}
	log.Println("ReadAuRevoir: AuRevoir loaded.")
}
