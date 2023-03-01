package plugin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ExitCode int

const (
	ExitCodeSuccess = ExitCode(0)
)

type Env string

const (
	EnvVersion  = Env("VERSION")
	EnvMimeType = Env("MIME_TYPE")
)

type Result struct {
	ExitCode int
	Out      io.Reader
	Err      io.Reader
}

type ExecuteOptions struct {
	Plugin string
	Prefix string
	Stdin  io.Reader
	Env    []string
	Sep    rune
}

func Execute(in ExecuteOptions) (*Result, error) {
	plugins, err := Search(in.Prefix, in.Sep)
	if err != nil {
		return nil, err
	}

	p, ok := plugins[in.Plugin]
	if !ok {
		return nil, errors.New("no plugin found")
	}

	path, err := exec.LookPath(p.Bin)
	if err != nil {
		return nil, err
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)

	cmd := exec.Command(path)

	cmd.Stdin = in.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	cmd.Env = in.Env

	if err := cmd.Run(); err != nil {
		return nil, err
	}

	exitCode := ExitCode(cmd.ProcessState.ExitCode())

	switch exitCode {
	case ExitCodeSuccess:
	}

	return &Result{
		Out: stdout,
	}, nil
}

func Read(r io.Reader) ([]byte, error) {
	return io.ReadAll(r)
}

func ParseJSON(v interface{}, r io.Reader) error {
	data, err := Read(r)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}

func Search(prefix string, sep rune) (PluginRegistry, error) {
	plugins := make(PluginRegistry, 0)

	paths := make([]string, 0)
	path := os.Getenv("PATH")
	if path != "" {
		paths = strings.Split(path, ":")
	}

	prefix = fmt.Sprintf("%v%c", prefix, sep)

	for _, p := range paths {
		err := filepath.Walk(os.ExpandEnv(p), func(path string, info fs.FileInfo, err error) error {
			bin := filepath.Base(path)

			if strings.HasPrefix(bin, prefix) {
				pluginName := strings.TrimPrefix(bin, prefix)

				plugins[pluginName] = Plugin{
					Name: strings.TrimPrefix(bin, prefix),
					Path: path,
					Bin:  bin,
				}
			}

			return nil
		})

		if err != nil {
			return nil, err
		}

	}

	return plugins, nil
}

type PluginRegistry map[string]Plugin

type Plugin struct {
	Name string
	Path string
	Bin  string
}
