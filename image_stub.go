//go:build !image

package main

func runImageCommand(args []string) {
	panic("image tooling not available; rebuild with `-tags image`")
}
