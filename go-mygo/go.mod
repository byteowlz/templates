module github.com/byteowlz/{{project_name}}

go 1.27.1

require (
	github.com/byteowlz/design-system/impls/mygo-go v0.0.0
	github.com/egoist/mygo v0.3.2
	github.com/pelletier/go-toml/v2 v2.2.4
	github.com/santhosh-tekuri/jsonschema/v6 v6.0.2
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/ebitengine/purego v0.11.1 // indirect
	github.com/go-text/typesetting v0.3.5 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	golang.org/x/image v0.46.0 // indirect
	golang.org/x/text v0.42.0 // indirect
)

tool github.com/egoist/mygo/cmd/mygo

// Pinned byte-identical source snapshot; portable, no local checkout dependency.
replace github.com/byteowlz/design-system/impls/mygo-go => ./third_party/mygo-go
