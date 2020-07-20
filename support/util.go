package support

import (
	"fmt"
	"log"
	"os"
)

func CheckError(err error, fatal bool) {
	if err != nil {
		buf := fmt.Sprintf("error: %s", err.Error())
		log.Println(buf)
		if fatal {
			os.Exit(1)
		}
	}
}
