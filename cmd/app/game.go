package main

import (
	"math/rand"
	"strings"
	"time"
)

type game struct {
	Current  []string
	Attempts []string
	Answer   string
}

type guess struct {
	Letters []letter
}

type letter struct {
	Character string
	Status    status
}

type status string

const (
	Miss        status = "miss"
	InWord      status = "in-word"
	Hit         status = "hit"
	Unattempted status = "unattempted"
)

func (s *status) Value() int {
	values := map[status]int{
		Unattempted: 0,
		Miss:        1,
		InWord:      2,
		Hit:         3,
	}

	return values[*s]
}

func (g *game) Guess(attempt string) {
	g.Attempts = append(g.Attempts, attempt)
}

func (g *game) HasWon() bool {
	for _, attempt := range g.Attempts {
		if attempt == g.Answer {
			return true
		}
	}

	return false
}

func (g *game) Guesses() []guess {
	guesses := []guess{}
	for _, attempt := range g.Attempts {
		guess := guess{}
		for idx, ltr := range strings.Split(attempt, "") {
			guess.Letters = append(guess.Letters, letter{
				Character: ltr,
				Status:    g.LetterStatus(ltr, idx),
			})
		}
		guesses = append(guesses, guess)
	}
	return guesses
}

func (g *game) KeyboardLetterStatus(l string) status {
	status := Unattempted

	for _, guess := range g.Guesses() {
		for _, ltr := range guess.Letters {
			if ltr.Character == l && ltr.Status.Value() > status.Value() {
				status = ltr.Status
			}
		}
	}

	return status
}

func (g *game) LetterStatus(l string, pos int) status {
	if string(g.Answer[pos]) == l {
		return Hit
	}

	if strings.Contains(g.Answer, string(l)) {
		return InWord
	}

	return Miss
}

func (g *game) HasLost() bool {
	return len(g.Attempts) >= 6
}

func (g *game) IsOver() bool {
	return g.HasWon() || g.HasLost()
}

// Get a random word from the most common 1000. We want to be fair.
func (app *application) randomWord() string {
	rand.Seed(time.Now().UnixNano())
	randIdx := rand.Intn(1000)

	return app.words[randIdx]
}
