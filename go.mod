module github.com/ueffel/mdtohtml

go 1.19

require (
	github.com/alecthomas/chroma v0.10.0
	github.com/spf13/pflag v1.0.5
	github.com/yuin/goldmark v1.5.3
	github.com/yuin/goldmark-highlighting v0.0.0-20220208100518-594be1970594
	go.abhg.dev/goldmark/mermaid v0.3.0
)

require github.com/dlclark/regexp2 v1.7.0 // indirect

replace go.abhg.dev/goldmark/mermaid v0.3.0 => ../goldmark-mermaid
