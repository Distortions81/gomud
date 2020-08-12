package support

import (
	"bufio"
	"fmt"
	"strings"

	"../def"
	"../glob"
)

func ReadConnection(con *glob.ConnectionData) {

	glob.ConnectionListLock.Lock() /*--- LOCK ---*/
	if con == nil {
		return
	}
	if con.Valid == false {
		return
	}
	if glob.ConnectionListEnd >= def.MAX_USERS-1 {
		return
	}

	reader := bufio.NewReader(con.Desc)

	glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/

	for con.Valid && con.Desc != nil {

		input, err := reader.ReadString('\n')
		glob.ConnectionListLock.Lock() /*--- LOCK ---*/

		/*Connection died*/
		if err != nil {
			DescWriteError(con, err)
			glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/
			return
		}

		filter := StripControl(input)
		limit, _ := TruncateString(filter, def.MAX_INPUT_LENGTH)

		if con.Input.BufferInCount-con.Input.BufferOutCount >= def.MAX_INPUT_LINES-1 {
			for x := 0; x <= 3; x++ {
				WriteToDesc(con, "Too many lines, stop spamming!")
				CloseConnection(con)
				glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/
				return
			}
		}

		lines := strings.Split(limit, ";")
		for i, line := range lines {
			if i < def.MAX_COMMANDS_PER_LINE {
				con.Input.BufferInPos++
				con.Input.BufferInCount++
				if con.Input.BufferInPos >= def.MAX_INPUT_LINES {
					con.Input.BufferInPos = 0
				}
				con.Input.InputBuffer[con.Input.BufferInPos] = line
			} else {
				buf := fmt.Sprintf("Too many commands on one line, stopped at #%v", def.MAX_COMMANDS_PER_LINE)
				WriteToDesc(con, buf)
				break
			}
		}
		glob.ConnectionListLock.Unlock() /*--- UNLOCK ---*/

	}
}

func ReadPlayerInputBuffer(con *glob.ConnectionData) {
	/* Only run if we have something */
	if con.Input.BufferInCount > con.Input.BufferOutCount {

		con.Input.BufferOutCount++
		con.Input.BufferOutPos++

		if con.Input.BufferOutPos >= def.MAX_INPUT_LINES {
			con.Input.BufferOutPos = 0
		}
		input := con.Input.InputBuffer[con.Input.BufferOutPos]

		bIn := len(input)

		con.BytesIn += bIn
		trackBytesIn(con)

		HandleReadConnection(con, input)
	}
}

func HandleReadConnection(con *glob.ConnectionData, input string) {

	//Newline before commands
	WriteToDesc(con, "")

	/*Player aliases*/
	if con.Player != nil && con.Player.Valid {
		if con.Player.Aliases != nil {

			if input != "" {
				for key, value := range con.Player.Aliases {

					if strings.EqualFold(key, input) {
						//add ; newline support
						interpretInput(con, value, true)
						return
					}

				}
			}
		}
	}

	/*Handles all user input*/
	interpretInput(con, input, false)

	/*Handle players marked for disconnection*/
	/*Doing this at the end is cleaner*/
	if con.State == def.CON_STATE_DISCONNECTING {
		CloseConnection(con)
		if con.Player != nil {
			RemovePlayer(con.Player)
		}
	}

}
