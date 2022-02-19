package main

import (
	"bufio"
	"os"
	"strings"
)

var validWords = map[string]struct{}{}

const (
	Right   rune = '🟩'
	Known        = '🟨'
	Unknown      = '⬜'
)

func checkGuessValid(guess string) bool {
	_, ok := validWords[strings.ToLower(guess)]
	return ok
}

func wordle(guess, inword string) string {
	if len(guess) != len(inword) {
		return ""
	}

	word := []rune(inword)
	out := []rune(strings.Repeat(string(Unknown), len(word)))
	seen := []rune(strings.Repeat(string(Unknown), len(word)))
	for i, r := range guess {
		if word[i] == r {
			out[i] = Right
			seen[i] = Right
		}
	}
	for i, r := range guess {
		if out[i] != Right {
			wordPos := 0
			for {
				if f := strings.IndexRune(inword[wordPos:], r); f != -1 {
					f = f + wordPos
					if seen[f] == Unknown {
						seen[f] = Known
						out[i] = Known
						break
					}
					if f < len(word)-1 {
						wordPos = f + 1
					} else {
						break
					}
				} else {
					break
				}
			}
		}
	}
	return string(out)
}

func dictLoad(length int) {
	words, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}
	defer words.Close()

	scanner := bufio.NewScanner(words)
	for scanner.Scan() {
		w := scanner.Text()
		if len(w) == length {
			if strings.ToLower(w) == w {
				validWords[w] = struct{}{}
			}
		}
	}
}
