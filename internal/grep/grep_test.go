package grep

import (
	"os"
	"strings"
	"testing"
)

type testCase struct {
	name      string
	pattern   string
	inputFile string
	wantFile  string
	opts      Options
}

func TestSortFromFiles(t *testing.T) {
	tests := []testCase{
		{
			name:      "simple cmp test",
			pattern:   "Hello",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/simple_cmp_test.txt",
			opts:      Options{},
		},
		{
			name:      "ignore register test",
			pattern:   "hello",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/ignore_register_test.txt",
			opts:      Options{IgnoreRegister: true},
		},
		{
			name:      "counter test",
			pattern:   "error",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/counter_test.txt",
			opts:      Options{PrintCount: true},
		},
		{
			name:      "ignore register count test",
			pattern:   "error",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/ignore_register_counter_test.txt",
			opts: Options{IgnoreRegister: true,
				PrintCount: true},
		},
		{
			name:      "inverted filter test",
			pattern:   "error",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/inverted_filter_test.txt",
			opts:      Options{InvertFilter: true},
		},
		{
			name:      "n lines after",
			pattern:   "error",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/n_lines_after_test.txt",
			opts:      Options{StringCountAfter: 2, IgnoreRegister: true},
		},
		{
			name:      "n lines before",
			pattern:   "error",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/n_lines_before_test.txt",
			opts:      Options{StringCountBefore: 2, IgnoreRegister: true},
		},
		{
			name:      "regex test",
			pattern:   "^(error|fatal):",
			inputFile: "test_cases/regex_input.txt",
			wantFile:  "test_cases/regex_test.txt",
			opts:      Options{},
		},
		{
			name:      "with string number test",
			pattern:   "Info",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/with_string_number_test.txt",
			opts:      Options{PrintStringNumberBeforeString: true},
		},
		{
			name:      "with string number and ignore register test",
			pattern:   "hello",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/with_string_number_and_ignore_register_test.txt",
			opts:      Options{IgnoreRegister: true, PrintStringNumberBeforeString: true},
		},
		{
			name:      "ignore register inverted test",
			pattern:   "hello",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/ignore_register_inverted_test.txt",
			opts:      Options{IgnoreRegister: true, InvertFilter: true},
		},
		{
			name:      "count inverted test",
			pattern:   "hello",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/count_inverted_test.txt",
			opts:      Options{PrintCount: true, InvertFilter: true},
		},
		{
			name:      "3 lines after and 2 before",
			pattern:   "foo",
			inputFile: "test_cases/input.txt",
			wantFile:  "test_cases/3_lines_after_and_2_before_test.txt",
			opts:      Options{StringCountAfter: 2, StringCountBeforeAndAfter: 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			inBytes, err := os.ReadFile(tc.inputFile)
			if err != nil {
				t.Fatalf("can't read input: %v", err)
			}

			lines := string(inBytes)
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}

			oldStdin := os.Stdin
			defer func() { os.Stdin = oldStdin }()
			os.Stdin = r

			_, err = w.Write([]byte(lines))
			if err != nil {
				t.Fatal(err)
			}
			_ = w.Close()

			resultChannel, err := FilterRows(tc.pattern, tc.opts)
			if err != nil {
				t.Errorf("wrong regexp: %q", tc.pattern)
			}

			wantBytes, err := os.ReadFile(tc.wantFile)
			if err != nil {
				t.Fatalf("can't read want: %v", err)
			}
			wantLines := strings.Split(string(wantBytes), "\n")
			// fmt.Println("----------------------------")
			for i := range wantLines {
				got := <-resultChannel
				// fmt.Println(got)
				if got != wantLines[i] {
					t.Errorf("[%d] got %q, want %q", i, got, wantLines[i])
				}
			}
			select {
			case got, ok := <-resultChannel:
				if ok {
					t.Errorf("got %q, want \"\"", got)
				}
			default:
			}
		})
	}
}
