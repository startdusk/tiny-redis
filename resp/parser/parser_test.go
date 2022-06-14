package parser

import (
	"bufio"
	"bytes"
	"testing"
)

func TestReadLine(t *testing.T) {
	cases := []struct {
		name     string
		line     string
		expected string
		wantErr  bool
	}{
		{
			name:     "starts with *",
			line:     "*3\r\n",
			expected: "*3",
		},
		{
			name:     "starts with $",
			line:     "$3\r\n",
			expected: "$3",
		},
	}

	for _, c := range cases {
		var buf bytes.Buffer
		buf.Write([]byte(c.line))
		bufReader := bufio.NewReader(&buf)
		var state readState
		t.Run(c.name, func(t *testing.T) {
			data, ok, err := readLine(bufReader, &state)
			if c.wantErr && err == nil {
				t.Fatalf("%s expect error but got nil", c.name)
			}
			if !c.wantErr && err != nil {
				t.Fatalf("%s expect nil but got %+v", c.name, err)
			}
			if c.wantErr || ok {
				return
			}
			if string(data) != c.expected {
				t.Fatalf("%s expect msg %s but got %s", c.name, c.expected, data)
			}
		})
	}
}

// TODO: add aother test
// TODO: add fuzz test
