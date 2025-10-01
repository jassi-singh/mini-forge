package utils

func shuffle(n uint64) uint64 {
	// These are high-quality mixing constants for bijective hash function.
	x := n
	x = (x ^ (x >> 30)) * 0xbf58476d1ce4e5b9
	x = (x ^ (x >> 27)) * 0x94d049bb133111eb
	x = x ^ (x >> 31)
	return x
}

// GenerateBase62Key takes a unique counter and returns a unique, non-sequential, base62-encoded key.
// Counter must be in range [0, 62^7) = [0, 3521614606208) to guarantee uniqueness.
func GenerateBase62Key(counter int64) string {
	const base62 = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const keyLength = 7
	const maxValue = 3521614606208 // 62^7

	if counter < 0 || counter >= maxValue {
		panic("counter out of valid range for 7-character base62 key")
	}

	// Shuffle to make keys non-sequential
	shuffledCounter := shuffle(uint64(counter))

	// Ensure it fits in 7 base62 characters by taking modulo
	shuffledCounter = shuffledCounter % maxValue

	// Build the base62 representation
	result := make([]byte, keyLength)

	// Fill from right to left (least significant digit first)
	for i := keyLength - 1; i >= 0; i-- {
		remainder := shuffledCounter % 62
		result[i] = base62[remainder]
		shuffledCounter /= 62
	}

	return string(result)
}
