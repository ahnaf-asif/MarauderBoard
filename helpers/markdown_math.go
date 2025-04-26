package helpers

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
)

func RenderMathToHTML(math string, displayMode bool) (string, error) {
	cmd := exec.Command("katex", "--no-throw-on-error")
	if displayMode {
		cmd.Args = append(cmd.Args, "--display-mode")
	}
	cmd.Stdin = strings.NewReader(math)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func MarkdownWithMath(input string) (string, error) {
	reInline := regexp.MustCompile(`\$(.+?)\$`)
	input = reInline.ReplaceAllStringFunc(input, func(match string) string {
		math := match[1 : len(match)-1]
		rendered, err := RenderMathToHTML(math, false)
		if err != nil {
			return match
		}
		return rendered
	})

	reBlock := regexp.MustCompile(`\$\$(.+?)\$\$`)
	input = reBlock.ReplaceAllStringFunc(input, func(match string) string {
		math := match[2 : len(match)-2]
		rendered, err := RenderMathToHTML(math, true)
		if err != nil {
			return match
		}
		return rendered
	})

	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(false),
				),
			),
		),
	)

	if err := md.Convert([]byte(input), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
} //
// 	markdown := `
// This is some markdown with inline math: $e^{i\pi} + 1 = 0$
//
// And a block:
//
// $$
// \int_0^\infty e^{-x^2} dx = \frac{\sqrt{\pi}}{2}
// $$
// `
// 	html, err := markdownWithMath(markdown)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(html)
