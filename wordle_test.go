package main

import (
	"testing"
)

func TestWordle(t *testing.T) {
	for _, tc := range []struct {
		Word, Guess, Expect string
	}{
		{"amuse", "amaze", "🟩🟩⬜⬜🟩"},
		{"amuse", "amuse", "🟩🟩🟩🟩🟩"},
		{"amuse", "amula", "🟩🟩🟩⬜⬜"},
		{"amuse", "amuck", "🟩🟩🟩⬜⬜"},
		{"amuse", "ameed", "🟩🟩🟨⬜⬜"},
		{"amuse", "amend", "🟩🟩🟨⬜⬜"},
		{"amuse", "level", "⬜🟨⬜⬜⬜"},
		{"dodge", "boats", "⬜🟩⬜⬜⬜"},
		{"dodge", "muddy", "⬜⬜🟩🟨⬜"},
		{"dodge", "daddy", "🟩⬜🟩⬜⬜"},
		{"dodge", "unadd", "⬜⬜⬜🟨🟨"},
	} {
		res := wordle(tc.Guess, tc.Word)
		if res != tc.Expect {
			t.Errorf("wordle(%q, %q): got %v, want %v", tc.Guess, tc.Word, res, tc.Expect)
		}
	}
}
