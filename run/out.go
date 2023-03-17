package run

import "fmt"

func Out(msg string) {
	fmt.Printf("[NMD-BOX] %s\n", msg)
}

func Error(msg string) {
	fmt.Printf("[NMD-BOX]    ERROR => %s\n", msg)
}

func Warn(msg string) {
	fmt.Printf("[NMD-BOX]    WARN => %s\n", msg)
}

func Header(msg string) {
	fmt.Printf("[NMD-BOX] ===== %s =====\n", msg)
}
