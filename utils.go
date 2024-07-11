package main

import "strings"

func StringToBool(s string) bool {
	switch strings.ToLower(s) {
	case "true", "yes", "1":
		return true
	case "false", "no", "0":
		return false
	default:
		// Handle unrecognized input; this example treats it as false
		return false
	}
}
