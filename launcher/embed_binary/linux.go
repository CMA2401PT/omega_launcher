//go:build linux
// +build linux

package embed_binary

import (
	_ "embed"
)

//go:embed cqhttp-linux
var embedding_cqhttp []byte
var PLANTFORM = Linux_x86_64
