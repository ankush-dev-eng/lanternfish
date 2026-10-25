package lanternfish

import "embed"

// WebFS embeds the static web user interface assets.
//go:embed all:web
var WebFS embed.FS
