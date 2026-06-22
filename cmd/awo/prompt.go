package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

var stdinReader = bufio.NewReader(os.Stdin)

// ask prints a prompt and reads one line. Returns def if input is empty.
func ask(label, def string) (string, error) {
	if def != "" {
		fmt.Printf("  %s [%s]: ", label, def)
	} else {
		fmt.Printf("  %s: ", label)
	}
	line, err := stdinReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return def, nil
	}
	return line, nil
}

// askRequired keeps prompting until user provides a non-empty value.
func askRequired(label string) (string, error) {
	for {
		v, err := ask(label, "")
		if err != nil {
			return "", err
		}
		if v != "" {
			return v, nil
		}
		fmt.Println("  (required — please enter a value)")
	}
}

// askPassword reads a password without echo. Falls back to plain read if not a TTY.
func askPassword(label string) (string, error) {
	fmt.Printf("  %s: ", label)
	fd := int(os.Stdin.Fd())
	if term.IsTerminal(fd) {
		b, err := term.ReadPassword(fd)
		fmt.Println()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(b)), nil
	}
	line, err := stdinReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(line), nil
}

// askSelect prompts for one of the allowed choices, looping until valid.
func askSelect(label string, choices []string, def string) (string, error) {
	choiceStr := strings.Join(choices, "|")
	for {
		v, err := ask(fmt.Sprintf("%s (%s)", label, choiceStr), def)
		if err != nil {
			return "", err
		}
		for _, c := range choices {
			if strings.EqualFold(v, c) {
				return c, nil
			}
		}
		fmt.Printf("  (choose one of: %s)\n", choiceStr)
	}
}
