package lox_error

import "fmt"

var hadError = false;

func Error(line int, message string) {
	Report(line, "", message);
}

func Report(line int, where, message string) {
	fmt.Printf("[line %d] Error%s: %s\n", line, where, message);
	hadError = true;
}

func HadError() bool {
	return hadError;
}

func ResetError() {
	hadError = false;
}
