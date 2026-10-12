// Instantiate verifies a real replaced scaffold. It needs only the Go stdlib.
package main

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("usage: instantiate SOURCE")
	}
	source, err := filepath.Abs(os.Args[1])
	if err != nil {
		return err
	}
	target, err := os.MkdirTemp("", "mygo-native-scaffold-")
	if err != nil {
		return err
	}
	fmt.Println("scaffold retained at", target)
	err = filepath.WalkDir(source, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case "build", "artifacts", ".git", ".mygo", "node_modules":
				return filepath.SkipDir
			}
			return os.MkdirAll(filepath.Join(target, relative), 0755)
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("template contains nonregular file %s", relative)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		data = bytes.ReplaceAll(data, []byte("{{project_name}}"), []byte("mygo-native-proof"))
		if bytes.Contains(data, []byte("{{project_name}}")) {
			return fmt.Errorf("unreplaced placeholder %s", relative)
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(target, relative), data, info.Mode().Perm())
	})
	if err != nil {
		return err
	}
	env := []string{}
	for _, value := range os.Environ() {
		if strings.HasPrefix(value, "GO111MODULE=") || strings.HasPrefix(value, "XDG_CONFIG_HOME=") || strings.HasPrefix(value, "MYGO_APP_") {
			continue
		}
		env = append(env, value)
	}
	env = append(env, "XDG_CONFIG_HOME="+filepath.Join(target, "isolated-config"))
	for _, args := range [][]string{{"just", "install"}, {"just", "check"}, {"just", "build"}, {"just", "agent-check"}} {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = target
		cmd.Env = env
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%v: %w (scaffold %s retained)", args, err, target)
		}
	}
	fmt.Printf("verified native scaffold: %s\n", target)
	return nil
}
