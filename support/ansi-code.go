package support

func ANSIColor(in string) string {
	colorized := false

	output := ""
	input := in

	length := len(input) - 1

	for i := 0; i < length; i++ {

		cur := input[i]
		length = len(input)
		if i+1 < length {
			next := input[i+1]

			if cur == '{' && next >= 'A' && next <= 'z' {
				color := getColor(next)

				colorized = true
				output = input[:i] + color + input[i+2:]
				input = output
			}
		}
	}
	if colorized {
		input = input + getColor('x')
	}
	return input

}

func getColor(i byte) string {

	if i == 'x' {
		return "\033[0m"
	} else if i == 'i' { //italic
		return "\033[0;3m"
	} else if i == 'u' { //underline
		return "\033[0;4m"
	} else if i == 'n' { //inverse
		return "\033[0;7m"
	} else if i == 's' { //strike
		return "\033[0;9m"

	} else if i == 'k' { //black
		return "\033[0;30m"
	} else if i == 'r' { //red
		return "\033[0;31m"
	} else if i == 'g' { //green
		return "\033[0;32m"
	} else if i == 'y' { //yellow
		return "\033[0;33m"
	} else if i == 'b' { //blue
		return "\033[0;34m"
	} else if i == 'm' { //magenta
		return "\033[0;35m"
	} else if i == 'c' { //cyan
		return "\033[0;36m"
	} else if i == 'w' { //gray
		return "\033[0;37m"

	} else if i == 'K' { //light gray
		return "\033[1;30m"
	} else if i == 'R' { //bright red
		return "\033[1;31m"
	} else if i == 'G' { //bright green
		return "\033[1;32m"
	} else if i == 'Y' { //bright yellow
		return "\033[1;33m"
	} else if i == 'B' { //bright blue
		return "\033[1;34m"
	} else if i == 'M' { //bright magenta
		return "\033[10;35m"
	} else if i == 'C' { //bright cyan
		return "\033[1;36m"
	} else if i == 'W' { //bright white
		return "\033[1;37m"
	} else {
		return "{\033[0m}" //error reset
	}
}
