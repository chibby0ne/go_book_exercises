package word

import (
	"math/rand"
	"testing"
	"time"
	"unicode"
)

func randomPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25) // random length up to 24 [0 - 24]
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := rune(rng.Intn(0x1000))
		runes[i] = r
		runes[n-1-i] = r
	}
	return string(runes)
}

// Creates random strings made of only ASCII chars for simplicity
func randomNonPalindrome(rng *rand.Rand) string {
	n := rng.Intn(25-2+1) + 2 // random length from 2 to 25 chars, [2 - 25], because every 1 letter string made of ascii will be a palindrome
	runes := make([]rune, n)
	for i := 0; i < (n+1)/2; i++ {
		r := chooseChar(rng)
		runes[i] = r
		runes[n-1-i] = chooseAnotherChar(r, rng)
	}
	return string(runes)
}

const ASCIIChars string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func chooseChar(rng *rand.Rand) rune {
	return rune(ASCIIChars[rng.Intn(len(ASCIIChars))])
}

func chooseAnotherChar(r rune, rng *rand.Rand) rune {
	p := rune(ASCIIChars[rng.Intn(len(ASCIIChars))])
	for unicode.ToLower(r) == unicode.ToLower(p) {
		p = rune(ASCIIChars[rng.Intn(len(ASCIIChars))])
	}
	return p
}

func TestRandomPalindromes(t *testing.T) {
	// Initialize a pseudo-random number generator
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 1000; i++ {
		p := randomPalindrome(rng)
		if !IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = false", p)
		}
	}
}

func TestRandomNonPalindromes(t *testing.T) {
	// Initialize a pseudo-random number generator
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 1000; i++ {
		p := randomNonPalindrome(rng)
		if IsPalindrome(p) {
			t.Errorf("IsPalindrome(%q) = true", p)
		}
	}
}
