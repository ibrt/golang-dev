package assets

import (
	"embed"
)

// SQLCTemplatesFS embeds an asset.
//
//go:embed sqlc/templates
var SQLCTemplatesFS embed.FS

// SQLCTemplatesFSRoot points to the root directory of SQLCTemplatesFS.
const SQLCTemplatesFSRoot = "sqlc/templates"

// SQLCPluginWASMGZEmbed embeds an asset.
//
//go:embed sqlc/plugin.wasm.gz
var SQLCPluginWASMGZEmbed []byte
