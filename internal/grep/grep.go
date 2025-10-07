package grep

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Options struct {
	StringCountAfter              int  // -A
	StringCountBefore             int  // -B
	StringCountBeforeAndAfter     int  // -C
	PrintCount                    bool // -c
	IgnoreRegister                bool // -i
	InvertFilter                  bool // -v
	SampleIsNotARegularExpression bool // -F
	PrintStringNumberBeforeString bool // -n
}

type ringBuffer struct {
	data  []string
	index int
	full  bool
}

func newRingBuffer(size int) *ringBuffer {
	return &ringBuffer{data: make([]string, size)}
}

func (r *ringBuffer) add(s string) {
	if len(r.data) == 0 {
		return
	}
	r.data[r.index] = s
	r.index = (r.index + 1) % len(r.data)
	if r.index == 0 {
		r.full = true
	}
}

func (r *ringBuffer) contents() []string {
	if len(r.data) == 0 {
		return nil
	}
	if !r.full {
		return r.data[:r.index]
	}
	out := make([]string, len(r.data))
	copy(out, r.data[r.index:])
	copy(out[len(r.data)-r.index:], r.data[:r.index])
	return out
}

func formatLine(num int, s string, opt Options) string {
	if opt.PrintStringNumberBeforeString {
		return fmt.Sprintf("%d %s", num, s)
	}
	return s
}

func FilterRows(pattern string, options Options) (<-chan string, error) {
	out := make(chan string)

	before := options.StringCountBefore + options.StringCountBeforeAndAfter
	after := options.StringCountAfter + options.StringCountBeforeAndAfter
	var re *regexp.Regexp
	var err error
	if !options.SampleIsNotARegularExpression {
		rePattern := pattern
		if options.IgnoreRegister {
			rePattern = "(?i)" + pattern
		}
		re, err = regexp.Compile(rePattern)
		if err != nil {
			return nil, err
		}
	}

	patternLower := pattern
	if options.SampleIsNotARegularExpression && options.IgnoreRegister {
		patternLower = strings.ToLower(pattern)
	}

	go func() {
		defer close(out)

		scanner := bufio.NewScanner(os.Stdin)
		ring := newRingBuffer(before)
		lineNum := 0

		printUntil := 0
		lastPrinted := 0
		matchCount := 0

		for scanner.Scan() {
			lineNum++
			rawLine := strings.TrimRight(scanner.Text(), "\r")
			checkLine := rawLine

			matched := false
			if options.SampleIsNotARegularExpression {
				if options.IgnoreRegister {
					checkLine = strings.ToLower(rawLine)
					if strings.Contains(checkLine, patternLower) {
						matched = true
					}
				} else {
					if strings.Contains(checkLine, pattern) {
						matched = true
					}
				}
			} else {
				if re.MatchString(checkLine) {
					matched = true
				}
			}

			if options.InvertFilter {
				matched = !matched
			}

			if options.PrintCount {
				if matched {
					matchCount++
				}
				if !matched {
					if before > 0 {
						ring.add(rawLine)
					}
				}
				continue
			}

			if matched {
				matchCount++

				targetUntil := lineNum + after
				if targetUntil > printUntil {
					printUntil = targetUntil
				}

				contents := ring.contents()
				startIdx := lineNum - len(contents)
				for j, s := range contents {
					idx := startIdx + j
					if idx > lastPrinted {
						out <- formatLine(idx, s, options)
						lastPrinted = idx
					}
				}

				if lineNum > lastPrinted {
					out <- formatLine(lineNum, rawLine, options)
					lastPrinted = lineNum
				}
			} else if lineNum <= printUntil {
				if lineNum > lastPrinted {
					out <- formatLine(lineNum, rawLine, options)
					lastPrinted = lineNum
				}
			} else {
				if before > 0 {
					ring.add(rawLine)
				}
			}
		}

		if options.PrintCount {
			out <- fmt.Sprint(matchCount)
			return
		}

	}()

	return out, nil
}
