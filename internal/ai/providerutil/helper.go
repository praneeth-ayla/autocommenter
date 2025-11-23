package providerutil

import (
	"fmt"
	"regexp"
	"strings"
)

// StripCodeFences removes markdown fences like ``` or ```go without damaging code layout.
func StripCodeFences(s string) string {
	// remove lines that only contain fences (e.g. ``` or ```go)
	reFenceBlock := regexp.MustCompile("(?m)^\\s*```(?:[a-zA-Z0-9]*)?\\s*$")
	s = reFenceBlock.ReplaceAllString(s, "")

	// remove remaining inline triple-backticks if any
	reInlineFence := regexp.MustCompile("```")
	s = reInlineFence.ReplaceAllString(s, "")

	// collapse excessive blank lines introduced by removals
	reMultiBlank := regexp.MustCompile("(?m)\n{3,}")
	s = reMultiBlank.ReplaceAllString(s, "\n\n")

	return strings.TrimSpace(s)
}

// EnsurePackageLine prepends the original package line if the output is missing it.
func EnsurePackageLine(out, original string) string {
	rePkg := regexp.MustCompile(`(?m)^\s*package\s+[a-zA-Z_]\w*`)
	if rePkg.MatchString(out) {
		return out
	}
	if m := rePkg.FindString(original); m != "" {
		return m + "\n\n" + out
	}
	return out
}

// NonCommentCodeChanged returns true if non-comment, non-whitespace code differs between orig and out.
// It also returns a short diff snippet for debugging.
func NonCommentCodeChanged(orig, out string) (bool, string) {
	strip := func(s string) string {
		// remove block comments
		reBlock := regexp.MustCompile(`(?s)/\*.*?\*/`)
		s = reBlock.ReplaceAllString(s, "")
		// remove line comments
		reLine := regexp.MustCompile(`//.*`)
		s = reLine.ReplaceAllString(s, "")
		// normalize whitespace and remove blank lines
		lines := strings.Split(s, "\n")
		var kept []string
		for _, l := range lines {
			t := strings.TrimSpace(l)
			if t == "" {
				continue
			}
			kept = append(kept, t)
		}
		return strings.Join(kept, "\n")
	}

	a := strip(orig)
	b := strip(out)

	if a == b {
		return false, ""
	}

	// tiny diff snippet: find first differing index
	min := len(a)
	if len(b) < min {
		min = len(b)
	}
	idx := 0
	for idx = 0; idx < min; idx++ {
		if a[idx] != b[idx] {
			break
		}
	}

	start := idx - 40
	if start < 0 {
		start = 0
	}
	endA := idx + 40
	if endA > len(a) {
		endA = len(a)
	}
	endB := idx + 40
	if endB > len(b) {
		endB = len(b)
	}
	snippet := fmt.Sprintf("orig...(%s)\n\nnew...(%s)", a[start:endA], b[start:endB])
	return true, snippet
}

// PruneExcessiveComments keeps comment blocks that are useful (precede declarations).
// It will drop many per-line comments to keep output concise.
// maxBlocks controls how many comment blocks are preserved overall.
func PruneExcessiveComments(src string, maxBlocks int) string {
	lines := strings.Split(src, "\n")

	type cblock struct {
		start int
		end   int
		keep  bool
	}

	var blocks []cblock
	i := 0
	for i < len(lines) {
		line := strings.TrimSpace(lines[i])
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			start := i
			inBlock := strings.HasPrefix(line, "/*")
			for j := i; j < len(lines); j++ {
				l := strings.TrimSpace(lines[j])
				if inBlock {
					if strings.Contains(l, "*/") {
						i = j + 1
						break
					}
				} else {
					if j > start && !strings.HasPrefix(strings.TrimSpace(lines[j]), "//") {
						i = j
						break
					}
					if j == len(lines)-1 {
						i = j + 1
						break
					}
				}
				if j == len(lines)-1 {
					i = j + 1
					break
				}
			}
			end := i - 1
			blocks = append(blocks, cblock{start: start, end: end, keep: false})
			continue
		}
		i++
	}

	declRe := regexp.MustCompile(`^\s*(func|type|var|const|package|import)\b`)

	for idx, b := range blocks {
		nextLineIdx := b.end + 1
		for nextLineIdx < len(lines) && strings.TrimSpace(lines[nextLineIdx]) == "" {
			nextLineIdx++
		}
		if nextLineIdx < len(lines) && declRe.MatchString(lines[nextLineIdx]) {
			blocks[idx].keep = true
		}
	}

	kept := 0
	for _, b := range blocks {
		if b.keep {
			kept++
		}
	}

	if kept < maxBlocks {
		for i := 0; i < len(blocks) && kept < maxBlocks; i++ {
			if blocks[i].keep {
				continue
			}
			blocks[i].keep = true
			kept++
		}
	}

	skip := map[int]bool{}
	for _, b := range blocks {
		if !b.keep {
			for k := b.start; k <= b.end; k++ {
				skip[k] = true
			}
		}
	}

	out := make([]string, 0, len(lines))
	for idx, l := range lines {
		if skip[idx] {
			continue
		}
		out = append(out, l)
	}
	return strings.Join(out, "\n")
}
