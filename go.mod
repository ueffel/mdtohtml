module github.com/ueffel/mdtohtml

go 1.20

require (
	github.com/alecthomas/chroma/v2 v2.5.0
	github.com/spf13/pflag v1.0.5
	github.com/yuin/goldmark v1.5.4
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20220924101305-151362477c87
	go.abhg.dev/goldmark/mermaid v0.3.0
)

require github.com/dlclark/regexp2 v1.8.0 // indirect

replace go.abhg.dev/goldmark/mermaid v0.3.0 => ../goldmark-mermaid
