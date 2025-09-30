package utils

func GenerateBase62Key(counter int64) string {
	const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const keyLength = 7

	// Shuffle the bits of the counter to add randomness
	counter = ((counter >> 16) | (counter << 48)) & ((1 << 63) - 1)

	// Convert to base62
	result := make([]byte, keyLength)
	for i := keyLength - 1; i >= 0; i-- {
		result[i] = base62[counter%62]
		counter /= 62
	}

	return string(result)
}
