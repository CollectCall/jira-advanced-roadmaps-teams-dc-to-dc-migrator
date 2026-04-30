package app

import (
	"strings"
	"testing"
)

func TestFormatTaskLineShowsCompletedTask(t *testing.T) {
	line := formatTaskLine(&progressTaskState{
		label:   "Applying issue team rewrites",
		current: 3,
		total:   3,
		done:    true,
	}, "-")

	if !strings.Contains(line, "100% Applying issue team rewrites") {
		t.Fatalf("expected completed task line to remain visible, got %q", line)
	}
	if !strings.Contains(line, strings.Repeat("#", 16)) {
		t.Fatalf("expected completed task line to show a full bar, got %q", line)
	}
}

func TestFormatTaskLineShowsRequestDetail(t *testing.T) {
	line := formatTaskLine(&progressTaskState{
		label:   "Applying filter rewrites",
		detail:  "updating filter 9001",
		current: 1,
		total:   2,
	}, "|")

	if !strings.Contains(line, "Applying filter rewrites - updating filter 9001") {
		t.Fatalf("expected task detail in progress line, got %q", line)
	}
}
