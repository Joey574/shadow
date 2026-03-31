package shell

import (
	"strings"
)

const ALPHA_NUMERIC = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// Creates an increasing length variable name based on the idx
func toBase62(n uint64) string {
	if n == 0 {
		return string(ALPHA_NUMERIC[0])
	}

	var sb strings.Builder
	base := uint64(len(ALPHA_NUMERIC))

	for n > 0 {
		remainder := n % base
		sb.WriteByte(ALPHA_NUMERIC[remainder])
		n = n / base
	}

	// The digits are collected in reverse order, so we flip the string
	return reverse(sb.String())
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func varName(idx int) string {
	v := toBase62(uint64(idx))
	if strings.ContainsAny(v[:1], "1234567890") {
		return "_" + v
	}

	return v
}
