package utils

import "strings"

func shuffle(n int64) int64 {
	// Cast the input to uint64 to perform unsigned arithmetic.
	x := uint64(n)

	// These are the same high-quality mixing constants, now used correctly.
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)

	// Cast the final result back to int64.
	return int64(x)
}

// GenerateBase62Key takes a unique counter and returns a unique, non-sequential, base62-encoded key.
func GenerateBase62Key(counter int64) string {
	const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	shuffledCounter := shuffle(counter)

	if shuffledCounter < 0 {
		shuffledCounter = -shuffledCounter
	}

	if shuffledCounter == 0 {
		return string(base62[0])
	}

	var sb strings.Builder

	for shuffledCounter > 0 {
		remainder := shuffledCounter % 62
		sb.WriteByte(base62[remainder])
		shuffledCounter /= 62
	}

	return sb.String()
}
