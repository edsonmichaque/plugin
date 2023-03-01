package main

import (
	"io"
	"os"
	"strings"

	"github.com/edsonmichaque/plugin"
)

func main() {
	_, err := plugin.Search("openapi-gen", '-')
	if err != nil {
		panic(err)
	}

	opts := plugin.ExecuteOptions{
		Plugin: "bash",
		Prefix: "openapi-gen",
		Stdin:  strings.NewReader(os.Args[1]),
		Env: []string{
			"MIME_TYPE=application/json",
		},
		Sep: '-',
	}

	res, err := plugin.Execute(opts)
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, res.Out)
}
