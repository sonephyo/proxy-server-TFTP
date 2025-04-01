package helper

import "fmt"

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func ColorPrintln(color string, text string) {
	colorMap := map[string]string{
		"red":     Red,
		"green":   Green,
		"yellow":  Yellow,
		"blue":    Blue,
		"magenta": Magenta,
		"cyan":    Cyan,
		"gray":    Gray,
		"white":   White,
	}

	if c, exists := colorMap[color]; exists {
		fmt.Println(c + text + Reset)
	} else {
		fmt.Println(Red + text + Reset)
	}
}
