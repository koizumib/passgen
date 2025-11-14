// main.go
package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"
)

func main() {
	// Preprocess args to allow long options with `--`
	os.Args = normalizeLongOpts(os.Args, map[string]bool{
		"help":   true,
		"length": true, "l": true,
		"number": true, "n": true,
		"add": true, "a": true,
		"delete": true, "d": true,
	})

	// Custom FlagSet to control usage output
	fs := flag.NewFlagSet("passgen", flag.ContinueOnError)
	fs.SetOutput(io.Discard) // we'll handle errors ourselves

	var (
		length int
		number int
		add    bool
		del    bool
	)

	// Define flags with aliases
	fs.IntVar(&length, "l", 8, "パスワードの長さ (別名: --length)")
	fs.IntVar(&length, "length", 8, "パスワードの長さ")
	fs.IntVar(&number, "n", 1, "生成するパスワードの数 (別名: --number)")
	fs.IntVar(&number, "number", 1, "生成するパスワードの数")
	fs.BoolVar(&add, "a", false, "標準入力/引数の文字をデフォルト文字セットに追加 (別名: --add)")
	fs.BoolVar(&add, "add", false, "標準入力/引数の文字をデフォルト文字セットに追加")
	fs.BoolVar(&del, "d", false, "標準入力/引数の文字をデフォルト文字セットから削除 (別名: --delete)")
	fs.BoolVar(&del, "delete", false, "標準入力/引数の文字をデフォルト文字セットから削除")

	// Custom usage message (matches the style in your example)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of passgen:\n")
		fmt.Fprintf(os.Stderr, "  -a\t標準入力/引数の文字をデフォルト文字セットに追加 (別名: --add)\n")
		fmt.Fprintf(os.Stderr, "  -add\n\t\t標準入力/引数の文字をデフォルト文字セットに追加\n")
		fmt.Fprintf(os.Stderr, "  -d\t標準入力/引数の文字をデフォルト文字セットから削除 (別名: --delete)\n")
		fmt.Fprintf(os.Stderr, "  -delete\n\t\t標準入力/引数の文字をデフォルト文字セットから削除\n")
		fmt.Fprintf(os.Stderr, "  -l int\n\t\tパスワードの長さ (別名: --length) (default 8)\n")
		fmt.Fprintf(os.Stderr, "  -length int\n\t\tパスワードの長さ (default 8)\n")
		fmt.Fprintf(os.Stderr, "  -n int\n\t\t生成するパスワードの数 (別名: --number) (default 1)\n")
		fmt.Fprintf(os.Stderr, "  -number int\n\t\t生成するパスワードの数 (default 1)\n")
	}

	// Parse
	if err := fs.Parse(os.Args[1:]); err != nil {
		// Show our usage on parse error
		fs.Usage()
		os.Exit(2)
	}

	// Validate
	if add && del {
		fmt.Fprintln(os.Stderr, "error: -a/--add と -d/--delete は同時に指定できません")
		os.Exit(2)
	}
	if length < 1 {
		length = 1
	}
	if number < 1 {
		number = 1
	}

	// Default character set: A-Z, a-z, 0-9
	defaultCharset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Collect characters from STDIN (if piped) and from remaining args
	provided := readStdin()
	if leftover := strings.Join(fs.Args(), ""); leftover != "" {
		provided += leftover
	}
	provided = stripNewlines(provided)

	// Build the final charset according to flags
	var charset []rune
	switch {
	case add:
		charset = uniqueRunes(defaultCharset + provided)
	case del:
		charset = subtractRunes(uniqueRunes(defaultCharset), toSet(provided))
	case provided != "":
		charset = uniqueRunes(provided)
	default:
		charset = uniqueRunes(defaultCharset)
	}

	if len(charset) == 0 {
		fmt.Fprintln(os.Stderr, "error: 利用可能な文字セットが空です（指定の結果、全て除去されました）")
		os.Exit(2)
	}

	// Generate
	for i := 0; i < number; i++ {
		pw, err := genPassword(length, charset)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println(pw)
	}
}

// --- helpers ---

func normalizeLongOpts(args []string, known map[string]bool) []string {
	out := make([]string, 0, len(args))
	for _, a := range args {
		if strings.HasPrefix(a, "--") {
			trim := strings.TrimPrefix(a, "--")
			// Support --flag=value form too
			name, rest, hasEq := splitFlagEq(trim)
			if known[name] {
				if hasEq {
					out = append(out, "-"+name+rest)
				} else {
					out = append(out, "-"+name)
				}
				continue
			}
		}
		out = append(out, a)
	}
	return out
}

func splitFlagEq(s string) (name, rest string, hasEq bool) {
	if i := strings.IndexByte(s, '='); i >= 0 {
		return s[:i], s[i:], true
	}
	return s, "", false
}

func readStdin() string {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return ""
	}
	// If not a character device, something is piped in
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		b, _ := io.ReadAll(os.Stdin)
		return string(b)
	}
	return ""
}

func stripNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	return s
}

func uniqueRunes(s string) []rune {
	seen := map[rune]bool{}
	var out []rune
	for _, r := range s {
		if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	return out
}

func toSet(s string) map[rune]bool {
	m := map[rune]bool{}
	for _, r := range s {
		m[r] = true
	}
	return m
}

func subtractRunes(base []rune, remove map[rune]bool) []rune {
	var out []rune
	for _, r := range base {
		if !remove[r] {
			out = append(out, r)
		}
	}
	return out
}

func genPassword(length int, charset []rune) (string, error) {
	var b strings.Builder
	b.Grow(length)
	max := big.NewInt(int64(len(charset)))
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		b.WriteRune(charset[n.Int64()])
	}
	return b.String(), nil
}
