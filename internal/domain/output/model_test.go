package output

import (
	"testing"
	"time"
)

func TestOutputBuffer_Append(t *testing.T) {
	buf := &OutputBuffer{
		TaskID: "test-1",
		Events: make([]OutputEvent, 0),
	}

	event := OutputEvent{
		ID:        "1",
		TaskID:    "test-1",
		Timestamp: time.Now(),
		Source:    StdoutSource,
		Level:     InfoLevel,
		Message:   "hello",
	}
	buf.Append(event)

	if len(buf.Events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(buf.Events))
	}
	if buf.Events[0].Message != "hello" {
		t.Errorf("expected message 'hello', got '%s'", buf.Events[0].Message)
	}
}

func TestOutputBuffer_FIFOEviction(t *testing.T) {
	maxLen := 3
	buf := &OutputBuffer{
		TaskID: "test-1",
		Events: make([]OutputEvent, 0),
		MaxLen: maxLen,
	}

	// 追加 5 个事件，超出限制
	for i := 0; i < 5; i++ {
		buf.Append(OutputEvent{
			ID:      string(rune('0' + i)),
			Message: "msg",
		})
	}

	// 最多保留 maxLen 条
	if len(buf.Events) > maxLen {
		t.Errorf("expected at most %d events, got %d", maxLen, len(buf.Events))
	}
}

func TestOutputBuffer_NoEvictionWhenNoLimit(t *testing.T) {
	buf := &OutputBuffer{
		TaskID: "test-1",
		Events: make([]OutputEvent, 0),
		MaxLen: 0,
	}

	for i := 0; i < 10; i++ {
		buf.Append(OutputEvent{ID: string(rune('0' + i))})
	}

	if len(buf.Events) != 10 {
		t.Errorf("expected 10 events, got %d", len(buf.Events))
	}
}

func TestOutputBuffer_GetEventsReturnsCopy(t *testing.T) {
	buf := &OutputBuffer{
		TaskID: "test-1",
		Events: []OutputEvent{{ID: "1", Message: "original"}},
	}

	events := buf.GetEvents()
	events[0].Message = "modified"

	if buf.Events[0].Message != "original" {
		t.Error("GetEvents should return a copy, not a reference")
	}
}

func TestOutputSource_Values(t *testing.T) {
	if StdoutSource != "stdout" {
		t.Errorf("expected 'stdout', got '%s'", StdoutSource)
	}
	if StderrSource != "stderr" {
		t.Errorf("expected 'stderr', got '%s'", StderrSource)
	}
	if SystemSource != "system" {
		t.Errorf("expected 'system', got '%s'", SystemSource)
	}
}

func TestLogLevel_Values(t *testing.T) {
	if TraceLevel != "trace" {
		t.Errorf("expected 'trace', got '%s'", TraceLevel)
	}
	if DebugLevel != "debug" {
		t.Errorf("expected 'debug', got '%s'", DebugLevel)
	}
	if InfoLevel != "info" {
		t.Errorf("expected 'info', got '%s'", InfoLevel)
	}
	if WarnLevel != "warn" {
		t.Errorf("expected 'warn', got '%s'", WarnLevel)
	}
	if ErrorLevel != "error" {
		t.Errorf("expected 'error', got '%s'", ErrorLevel)
	}
}
