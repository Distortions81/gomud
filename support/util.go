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

func ToHourMinute(t time.Duration) string {
	out := ""

	if int(t.Hours()) > 0 {
		out = out + fmt.Sprintf("%dh", int(t.Hours()))
	}
	if int(t.Minutes()) > 0 {
		out = out + fmt.Sprintf("%dm", int(t.Minutes()))
	}
	return out
}

func MakeFingerprint(prefix string) string {
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
