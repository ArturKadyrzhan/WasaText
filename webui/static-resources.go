//go:build webui

package webui

import "embed"

//go:embed "dist/*"
var Dist embed.FS
