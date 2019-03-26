/*
Godl is a CLI tool used to download and install go binary releases on mac.

It downloads the go binary archive specified from https:golang.org/dl/,
saves it at specified path and unpacks it into /usr/local/.
The downloaded archive can be found at specified download path
or $HOME/Downloads by default.

Usage:
	godl [go_archive] [path_to_save_archive] [flags]

Examples:
	godl go1.11.4.darwin-amd64.tar.gz ~/Downloads -r

Flags:
	-h, --help      help for godl
	-r, --remove    Remove flag is optional and is used to remove the downloaded archive after installing go.
	--version   version for godl
*/
package main
