package support

import (
	"fmt"
	"log"
	"os"
	"time"
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

func ToDayHourMinute(time time.Duration) string {
	out := ""

	if int(time.Hours()) > 0 {
		out = out + fmt.Sprintf("%dh", int(time.Hours()))
	}
	if int(time.Minutes()) > 0 {
		out = out + fmt.Sprintf("%dm", int(time.Minutes()))
	}
	return out
}
