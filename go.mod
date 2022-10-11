module github.com/danielcmessias/sawsy

go 1.18

require (
	github.com/alecthomas/chroma v0.10.0
	github.com/aws/aws-sdk-go-v2 v1.16.16
	github.com/aws/aws-sdk-go-v2/config v1.17.8
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.34
	github.com/aws/aws-sdk-go-v2/service/glue v1.27.0
	github.com/aws/aws-sdk-go-v2/service/lakeformation v1.16.7
	github.com/aws/aws-sdk-go-v2/service/s3 v1.27.11
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.19
	github.com/charmbracelet/bubbles v0.13.0
	github.com/charmbracelet/bubbletea v0.22.0
	github.com/charmbracelet/glamour v0.5.0
	github.com/charmbracelet/lipgloss v0.5.0
	github.com/muesli/termenv v0.11.1-0.20220212125758-44cd13922739
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.21.6
	github.com/aws/aws-sdk-go-v2/service/iam v1.18.20
	github.com/aws/aws-sdk-go-v2/service/lambda v1.24.6
	github.com/aws/aws-sdk-go-v2/service/rds v1.26.1
)

require (
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.8 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.12.21 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.6 // indirect
	github.com/aws/smithy-go v1.13.3 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/containerd/console v1.0.3 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/microcosm-cc/bluemonday v1.0.17 // indirect
	github.com/muesli/ansi v0.0.0-20211018074035-2e021307bc4b // indirect
	github.com/muesli/cancelreader v0.2.1 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	github.com/yuin/goldmark v1.4.4 // indirect
	github.com/yuin/goldmark-emoji v1.0.1 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/sys v0.0.0-20220209214540-3681064d5158 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
)

replace github.com/charmbracelet/lipgloss v0.5.0 => github.com/danielcmessias/eugener-lipgloss v0.0.0-20220922195730-9e516bc662be
