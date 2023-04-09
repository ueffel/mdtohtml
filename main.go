package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	chromaHtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/spf13/pflag"
	"github.com/ueffel/mdtohtml/embedimg"
	"github.com/ueffel/mdtohtml/tasklistitem"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	gmhtml "github.com/yuin/goldmark/renderer/html"
	mermaid "go.abhg.dev/goldmark/mermaid"
)

//go:embed github.css
var githubCSS []byte

//go:embed github-markdown.css
var githubMarkdownCSS []byte

func main() {
	help := pflag.BoolP("help", "h", false, "Show this help text")
	overwrite := pflag.BoolP("overwrite", "y", false, "Don't ask, overwrite all files")
	pflag.Usage = func() {
		fmt.Println("Converts Markdown to HTML with Images embedded.")
		fmt.Printf("Usage: %s [options] [files]\n", os.Args[0])
		fmt.Println("Options:")
		fmt.Print(pflag.CommandLine.FlagUsages())
	}
	pflag.Parse()

	if *help {
		pflag.Usage()
		os.Exit(0)
	}

	for _, path := range pflag.Args() {
		err := makeHTML(path, *overwrite)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR '%s': %s", path, err)
		}
	}
}

func makeHTML(path string, overwrite bool) error {
	buf.Reset()
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithFormatOptions(
					chromaHtml.WithClasses(true),
				),
			),
			extension.Linkify,
			extension.Strikethrough,
			extension.TaskList,
			embedimg.EmbedImg,
			tasklistitem.TaskListItemClass,
			&mermaid.Extender{MermaidJS: "<embed>"},
		),
		goldmark.WithParserOptions(
			parser.WithAttribute(),
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			gmhtml.WithUnsafe(),
			gmhtml.WithXHTML(),
		),
	)

	basePath := filepath.Dir(path)
	newName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) + ".html"
	resultPath := filepath.Join(basePath, newName)

	if !overwrite {
		if _, err := os.Stat(resultPath); err == nil || !os.IsNotExist(err) {
			invalidResponse := true
			var response string
			for invalidResponse {
				fmt.Printf("'%s' already exists. Overwrite? (y/N): ", resultPath)
				n, _ := fmt.Scanln(&response)
				if n > 0 && response == "n" || n == 0 {
					return nil
				}
				if n > 0 && response == "y" {
					invalidResponse = false
				}
			}
		}
	}

	oldCwd, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(basePath)
	if err != nil {
		return err
	}
	_, err = buf.WriteString(beforeTitle)
	if err != nil {
		return err
	}
	_, err = buf.WriteString(html.EscapeString(filepath.Base(path)))
	if err != nil {
		return err
	}
	_, err = buf.WriteString(afterTitle)
	if err != nil {
		return err
	}
	_, err = buf.Write(githubMarkdownCSS)
	if err != nil {
		return err
	}
	_, err = buf.Write(githubCSS)
	if err != nil {
		return err
	}
	_, err = buf.WriteString(beforeRender)
	if err != nil {
		return err
	}
	err = md.Convert(input, buf)
	if err != nil {
		return err
	}
	_, err = buf.WriteString(afterRender)
	if err != nil {
		return err
	}
	err = os.Chdir(oldCwd)
	if err != nil {
		return err
	}
	fmt.Printf("Writing '%s'\n", resultPath)
	err = os.WriteFile(resultPath, buf.Bytes(), 0o644)
	if err != nil {
		return err
	}

	return nil
}

var (
	buf         = new(bytes.Buffer)
	beforeTitle = `<!DOCTYPE html>
<html>
<head>
	<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>`
	afterTitle = `</title>
	<style>
`
	beforeRender = `
	.markdown-body {
		box-sizing: border-box;
		min-width: 200px;
		max-width: 980px;
		margin: 0 auto;
		padding: 45px;
	}
	@media (max-width: 767px) {
		.markdown-body {
			padding: 15px;
		}
	}
	.anchor-link {
		opacity: 0;
		font-size: .6em;
		margin-left: -1.4em;
	}
	h1:hover > .anchor-link,
	h2:hover > .anchor-link,
	h3:hover > .anchor-link,
	h4:hover > .anchor-link,
	h5:hover > .anchor-link,
	h6:hover > .anchor-link {
		opacity: 1;
		text-decoration: none;
	}
	.copy-button {
		position: absolute;
		display: none;
		right: .5em;
		top: .5em;
		font-size: 1.1em;
		padding: .2em;
		font-family: inherit;
	}
	.code-container {
		position: relative;
	}
	.code-container:hover .copy-button {
		display: block;
	}
	</style>
</head>
<body>
	<div class="markdown-body">
`
	afterRender = `	</div>
	<script>
		let headlines = []
		let tagNames = ["h1", "h2", "h3", "h4", "h5", "h6"];
		tagNames.forEach(function(tag, index)
		{
			let elements = document.getElementsByTagName(tag);
			for(let i = 0; i < elements.length; i++)
			{
				headlines.push(elements.item(i));
			}
		});

		headlines.forEach(function(item, index)
		{
			let anchor = document.createElement("a");
			anchor.setAttribute("href", "#" + item.id);
			anchor.classList.add("anchor-link");
			anchor.setAttribute("title", "Direct link");
			let text = document.createTextNode("🔗");
			anchor.append(text);
			item.prepend(anchor);
		});

		let codeBlocks = document.getElementsByTagName("pre");
		for(let i = 0; i < codeBlocks.length; i++)
		{
			codeBlock = codeBlocks.item(i);
			let container = document.createElement("div");
			container.classList.add("code-container");
			let btn = document.createElement("button");
			btn.classList.add("copy-button");
			btn.title = "Copy to clipboard";
			let text = document.createTextNode("✂");
			btn.append(text);
			btn.onclick = function (event)
			{
				copyToClipboard(codeBlocks[i], event);
			}
			codeBlock.parentNode.insertBefore(container, codeBlock);
			container.appendChild(codeBlock);
			container.appendChild(btn);
		}

		function copyToClipboard(sender, event)
		{
			let range;
			if (document.selection)
			{
				range = document.body.createTextRange();
				range.moveToElementText(sender);
				range.select();
			}
			else if (window.getSelection)
			{
				range = document.createRange();
				range.selectNodeContents(sender);
				window.getSelection().removeAllRanges();
				window.getSelection().addRange(range);
			}
			if (navigator.clipboard && window.isSecureContext)
			{
				navigator.clipboard.writeText(range.toString());
			} else
			{
				document.execCommand("copy");
			}
		}
	</script>
</body>
</html>
`
)
