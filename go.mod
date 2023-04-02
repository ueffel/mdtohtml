module github.com/ueffel/mdtohtml

go 1.20

require (
	github.com/alecthomas/chroma/v2 v2.7.0
	github.com/spf13/pflag v1.0.5
	github.com/yuin/goldmark v1.5.4
	github.com/yuin/goldmark-highlighting/v2 v2.0.0-20220924101305-151362477c87
	go.abhg.dev/goldmark/mermaid v0.4.0
)

require github.com/dlclark/regexp2 v1.8.1 // indirect

replace go.abhg.dev/goldmark/mermaid => github.com/ueffel/goldmark-mermaid v0.4.1-0.20230402203024-f44fc8134c82
