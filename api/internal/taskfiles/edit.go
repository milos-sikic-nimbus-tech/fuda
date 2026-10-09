package taskfiles

import (
	"bytes"
	"errors"
	"strings"

	"go.yaml.in/yaml/v3"
)

var ErrNoStatusLine = errors.New("the frontmatter has no status line")

type Move struct {
	Status  string
	Claimed string
}

func ParseOne(p string, content []byte) (Task, error) {
	t, reason := parseTask(p, content, false)
	if reason != "" {
		return Task{}, errors.New(reason)
	}
	return t, nil
}

func ApplyMove(content []byte, m Move) ([]byte, error) {
	lines := bytes.SplitAfter(content, []byte("\n"))
	front := frontmatterLines(lines)
	statusAt, claimedAt := -1, -1
	for i := front.start; i < front.end; i++ {
		switch {
		case bytes.HasPrefix(lines[i], []byte("status:")):
			statusAt = i
		case bytes.HasPrefix(lines[i], []byte("claimed:")):
			claimedAt = i
		}
	}
	if statusAt < 0 {
		return nil, ErrNoStatusLine
	}

	eol := lineEnding(lines[statusAt])
	lines[statusAt] = []byte("status: " + yamlScalar(m.Status) + eol)
	if m.Claimed != "" {
		claimed := []byte("claimed: " + m.Claimed + eol)
		switch {
		case claimedAt < 0:
			lines = append(lines[:statusAt+1], append([][]byte{claimed}, lines[statusAt+1:]...)...)
		case strings.TrimSpace(string(lines[claimedAt][len("claimed:"):])) == "":
			lines[claimedAt] = claimed
		}
	}
	return bytes.Join(lines, nil), nil
}

type frontmatterRange struct{ start, end int }

func frontmatterLines(lines [][]byte) frontmatterRange {
	if len(lines) == 0 || !isFence(bytes.TrimPrefix(lines[0], []byte("\xef\xbb\xbf"))) {
		return frontmatterRange{}
	}
	for i := 1; i < len(lines); i++ {
		if isFence(lines[i]) {
			return frontmatterRange{start: 1, end: i}
		}
	}
	return frontmatterRange{}
}

func isFence(line []byte) bool {
	return string(bytes.TrimRight(line, "\r\n")) == "---"
}

func lineEnding(line []byte) string {
	if bytes.HasSuffix(line, []byte("\r\n")) {
		return "\r\n"
	}
	return "\n"
}

func yamlScalar(s string) string {
	out, err := yaml.Marshal(s)
	if err != nil {
		return s
	}
	return strings.TrimSpace(string(out))
}
