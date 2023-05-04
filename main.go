package main

import (
	"ociswrapper/cmd"
	"ociswrapper/common"
)

func main() {
	common.Wg.Add(2)

	cmd.Execute()

	common.Wg.Wait()
}
