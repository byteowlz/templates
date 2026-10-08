package schemas

import _ "embed"

//go:embed envelope.schema.json
var Envelope []byte

//go:embed config.schema.json
var Config []byte
