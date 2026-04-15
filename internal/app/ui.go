package app

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const envUIMode = "TEAMS_MIGRATOR_UI"

type uiMode string

const (
	uiModeAuto  uiMode = "auto"
	uiModeRich  uiMode = "rich"
	uiModePlain uiMode = "plain"
)

type uiTheme struct {
	mode        uiMode
	useColor    bool
	useUnicode  bool
	borderColor string
	titleColor  string
	labelColor  string
	hintColor   string
	errorColor  string
	reset       string
}

func currentUITheme() uiTheme {
	mode := uiMode(strings.ToLower(strings.TrimSpace(os.Getenv(envUIMode))))
	switch mode {
	case uiModeRich:
		return richTheme()
	case uiModePlain:
		return plainTheme()
	default:
		if supportsRichUI() {
			return richTheme()
		}
		return plainTheme()
	}
}

func supportsRichUI() bool {
	if !isInteractiveTerminal() {
		return false
	}
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	term := strings.ToLower(strings.TrimSpace(os.Getenv("TERM")))
	if term == "" || term == "dumb" {
		return false
	}
	if runtime.GOOS == "windows" {
		return os.Getenv("WT_SESSION") != "" || os.Getenv("TERM_PROGRAM") != "" || strings.EqualFold(os.Getenv("ConEmuANSI"), "ON")
	}
	return true
}

func richTheme() uiTheme {
	return uiTheme{
		mode:        uiModeRich,
		useColor:    true,
		useUnicode:  true,
		borderColor: "\x1b[38;5;246m",
		titleColor:  "\x1b[1;38;5;45m",
		labelColor:  "\x1b[1;38;5;252m",
		hintColor:   "\x1b[38;5;110m",
		errorColor:  "\x1b[1;38;5;203m",
		reset:       "\x1b[0m",
	}
}

func plainTheme() uiTheme {
	return uiTheme{mode: uiModePlain}
}

func (t uiTheme) style(text, code string) string {
	if !t.useColor || code == "" {
		return text
	}
	return code + text + t.reset
}

func (t uiTheme) borderLine(title, heading string) []string {
	if !t.useUnicode {
		return []string{
			"================================================================",
			title,
			heading,
			"================================================================",
		}
	}

	width := 70
	line := func(left, fill, right string) string {
		return left + strings.Repeat(fill, width) + right
	}

	titleLine := fmt.Sprintf("│ %s  │  %s%s", t.style("Teams Migrator", t.titleColor), t.style(titleSuffix(title), t.labelColor), strings.Repeat(" ", max(0, width-4-visibleLen("Teams Migrator")-visibleLen(titleSuffix(title)))))
	headingLine := fmt.Sprintf("│ %s%s │", t.style(heading, t.labelColor), strings.Repeat(" ", max(0, width-2-visibleLen(heading))))
	return []string{
		t.style(line("┌", "─", "┐"), t.borderColor),
		t.style(titleLine, t.borderColor),
		t.style(line("├", "─", "┤"), t.borderColor),
		t.style(headingLine, t.borderColor),
		t.style(line("└", "─", "┘"), t.borderColor),
	}
}

func titleSuffix(title string) string {
	parts := strings.Split(title, " | ")
	if len(parts) < 2 {
		return title
	}
	return parts[1]
}

func visibleLen(s string) int {
	return len([]rune(s))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
