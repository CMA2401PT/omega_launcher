//go:build darwin
// +build darwin

package embed_binary

import (
	_ "embed"
)

//go:embed cqhttp-macos
var embedding_cqhttp []byte
var PLANTFORM = MACOS_x86_64
