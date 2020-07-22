package support

import (
	"fmt"
	"log"
	"math/rand"
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

func ToHourMinute(time time.Duration) string {
	out := ""

	if int(time.Hours()) > 0 {
		out = out + fmt.Sprintf("%dh", int(time.Hours()))
	}
	if int(time.Minutes()) > 0 {
		out = out + fmt.Sprintf("%dm", int(time.Minutes()))
	}
	return out
}

func MakeFingerprint(prefix string) string {
	fingerprint := fmt.Sprintf("%v%v-%v", prefix, time.Now().UnixNano(), rand.Uint64())
	return fingerprint
}
