package stream_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"

	llmerrors "go.osspkg.com/llm-client/pkg/errors"
	"go.osspkg.com/llm-client/pkg/stream"
)

type fuzzValue struct {
	Value string `json:"value"`
}

func TestNDJSON(t *testing.T) {
	tests := []struct {
		name string
		data string
		want []int
	}{
		{name: "multiple lines", data: "1\n2\n", want: []int{1, 2}},
		{name: "final line without newline", data: "3", want: []int{3}},
		{name: "blank lines", data: "\n4\n\n5", want: []int{4, 5}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			iterator := stream.NewNDJSON(io.NopCloser(strings.NewReader(test.data)), decodeInt, 64)
			var got []int
			for iterator.Next(context.Background()) {
				got = append(got, iterator.Value())
			}
			if err := iterator.Err(); err != nil {
				t.Fatalf("stream error: %v", err)
			}
			if len(got) != len(test.want) {
				t.Fatalf("got %v, want %v", got, test.want)
			}
			for index := range got {
				if got[index] != test.want[index] {
					t.Errorf("value %d: got %d, want %d", index, got[index], test.want[index])
				}
			}
		})
	}
}

func TestSSE(t *testing.T) {
	data := "event: message\ndata: 7\n\ndata: [DONE]\n\n"
	iterator := stream.NewSSE(io.NopCloser(strings.NewReader(data)), decodeInt, 64)
	if !iterator.Next(context.Background()) || iterator.Value() != 7 {
		t.Fatalf("unexpected first SSE value: %v", iterator.Value())
	}
	if iterator.Next(context.Background()) {
		t.Fatal("expected [DONE] to end SSE")
	}
	if err := iterator.Err(); err != nil {
		t.Fatalf("stream error: %v", err)
	}
}

func TestStreamContextAndLimit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	iterator := stream.NewNDJSON(io.NopCloser(strings.NewReader("1\n")), decodeInt, 64)
	if iterator.Next(ctx) {
		t.Fatal("expected canceled context")
	}
	if !errors.Is(iterator.Err(), context.Canceled) {
		t.Fatalf("got %v, want context canceled", iterator.Err())
	}

	large := stream.NewNDJSON(io.NopCloser(strings.NewReader("12345\n")), decodeInt, 2)
	if large.Next(context.Background()) {
		t.Fatal("expected oversized event failure")
	}
	if !errors.Is(large.Err(), llmerrors.ErrBodyTooLarge) {
		t.Fatalf("got %v, want body-too-large", large.Err())
	}
}

func decodeInt(data []byte) (int, error) {
	return strconv.Atoi(string(data))
}

func FuzzNDJSONDoesNotPanic(f *testing.F) {
	f.Add(`{"value":"one"}` + "\n" + `{"value":"two"}`)
	f.Fuzz(func(t *testing.T, data string) {
		iterator := stream.NewNDJSON(io.NopCloser(strings.NewReader(data)), func(payload []byte) (fuzzValue, error) {
			var value fuzzValue
			err := json.Unmarshal(payload, &value)
			return value, err
		}, 1024)
		for iterator.Next(context.Background()) {
			_ = iterator.Value()
		}
		_ = iterator.Close()
	})
}

func FuzzSSEDoesNotPanic(f *testing.F) {
	f.Add("data: {\"value\":\"one\"}\n\n")
	f.Fuzz(func(t *testing.T, data string) {
		iterator := stream.NewSSE(io.NopCloser(strings.NewReader(data)), func(payload []byte) (fuzzValue, error) {
			var value fuzzValue
			err := json.Unmarshal(payload, &value)
			return value, err
		}, 1024)
		for iterator.Next(context.Background()) {
			_ = iterator.Value()
		}
		_ = iterator.Close()
	})
}

func BenchmarkNDJSONDecode(b *testing.B) {
	data := `{"value":"benchmark"}` + "\n"
	for b.Loop() {
		iterator := stream.NewNDJSON(io.NopCloser(strings.NewReader(data)), func(payload []byte) (fuzzValue, error) {
			var value fuzzValue
			err := json.Unmarshal(payload, &value)
			return value, err
		}, 1024)
		if !iterator.Next(context.Background()) {
			b.Fatal(iterator.Err())
		}
		_ = iterator.Close()
	}
}

func BenchmarkSSEDecode(b *testing.B) {
	data := "data: {\"value\":\"benchmark\"}\n\n"
	for b.Loop() {
		iterator := stream.NewSSE(io.NopCloser(strings.NewReader(data)), func(payload []byte) (fuzzValue, error) {
			var value fuzzValue
			err := json.Unmarshal(payload, &value)
			return value, err
		}, 1024)
		if !iterator.Next(context.Background()) {
			b.Fatal(iterator.Err())
		}
		_ = iterator.Close()
	}
}
