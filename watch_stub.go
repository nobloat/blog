//go:build !watch

package main

func watchFiles() {
	panic("watch mode not available; rebuild with `-tags watch`")
}
