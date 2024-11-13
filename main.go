package main

import (
	_ "github.com/V2hiddify/V2hiddifyExtension/hiddify_extension"

	"github.com/hiddify/hiddify-core/extension/server"
)

func main() {
	server.StartTestExtensionServer()
}
