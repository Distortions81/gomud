package support

import (
	"fmt"
	"log"
	"os"
)

func CheckError(source string, err error, fatal bool) {
	if err != nil {
		buf := fmt.Sprintf("error: %s: %s", source, err.Error())
		log.Println(buf)
		if fatal {
			os.Exit(1)
		}
	}
}
