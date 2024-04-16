/*
Command 7zplugin compiles a plugin DLL with the specified plugin packages named
by their import paths.

This exists to allow arbitrary plugin choices to be easily built together since
there can only be one Go runtime in a process, so only one plugin DLL built
using 7zplugin can be used and must contain all desired plugins.

The arch argument must be set to one of:

  - 64   64-bit x86      (CGO_ENABLED=1 GOOS=windows GOARCH=amd64)
  - 32   32-bit x86_64   (CGO_ENABLED=1 GOOS=windows GOARCH=386)
  - A64  64-bit arm      (CGO_ENABLED=1 GOOS=windows GOARCH=arm64)

Specify environment variables for go build (CC, CXX, etc) as arguments before
the flags.

Flags are passed as-is to go build (see go help build). The only flag explicitly
set by this package is -buildmode=c-shared. It is highly recommended to also set
-ldflags '-s -w -extldflags=-static'.

Version information is added to the built binary:

  - FileDescription is set to a list of packages/files used to build the plugin.
  - FileVersion is set to YYYY.MM.DD.0 based on the current date, or the unix
    timestamp in SOURCE_DATE_EPOCH. If GITHUB_ACTION is set, GITHUB_RUN_NUMBER is
    used as the last component.
  - InternalName is set to the default output dll basename.
  - OriginalFilename is set based on InternalName and the architecture.
  - ProductName describes this library.
  - ProductVersion describes the target 7-Zip version and architecture.

An example invocation of this command to build all default plugins (from within
the 7zplugin source dir):

	CC=x86_64-w64-mingw32-gcc go run . 64 -ldflags '-s -w -extldflags=-static' -trimpath -v -x ./plugins/...

If using 7zplugin as a library, you can build the default plugins combined with
your plugin like:

	CC=x86_64-w64-mingw32-gcc go run github.com/pg9182/7zplugin 64 -ldflags '-s -w -extldflags=-static' -trimpath -v -x github.com/pg9182/7zplugin/plugins/...  ./pluginpkg1 ./pluginpkg2

Each plugin package or file should contain an init() function which calls
z7plugin.RegisterArc (or imports packages which do).

The built DLL file should be installed to the Formats directory beside the 7-Zip
executable. The architecture must match that of 7-Zip.
*/
package main

import (
	"bufio"
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/josephspurrier/goversioninfo"
	"github.com/pg9182/7zplugin/z7"
)

//go:embed build.go
var self []byte
var doc string

func init() {
	var ok bool
	if self, ok = bytes.CutPrefix(self, []byte("/*")); !ok {
		panic("failed to find start of package doc")
	}
	if self, _, ok = bytes.Cut(self, []byte("*/")); !ok {
		panic("failed to find end of package doc")
	}
	if crlf := []byte{'\r', '\n'}; bytes.Contains(self, crlf) {
		self = bytes.ReplaceAll(self, crlf, []byte{'\n'})
	} else {
		self = bytes.ReplaceAll(self, []byte{'\r'}, []byte{'\n'})
	}
	doc = string(self)
}

func help() {
	fmt.Printf("usage: %s arch [ENV=VALUE...] [flags] [package...]\n%s", os.Args[0], doc)
	os.Exit(0)
}

var dllname = "go7zPlugin"

func main() {
	var err error
	if len(os.Args) <= 1 {
		help()
	}

	// override the help flag
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if arg == "-h" || arg == "--help" || arg == "-help" {
			help()
		}
	}

	// expand the arch arg
	var arch, z7arch string
	switch arch = os.Args[1]; arch {
	case "64":
		os.Args[1] = "GOARCH=amd64"
		z7arch = "x64"
	case "32":
		os.Args[1] = "GOARCH=386"
		z7arch = "x86"
	case "ARM64":
		os.Args[1] = "GOARCH=arm64"
		z7arch = "arm64"
	default:
		fmt.Fprintf(os.Stderr, "7zplugin: error: unknown arch %q\n", arch)
		os.Exit(2)
	}
	os.Args = slices.Insert(os.Args, 1, "CGO_ENABLED=1", "GOOS=windows")

	// get the current dir
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "7zplugin: error: get working directory: %v\n", err)
		os.Exit(2)
	}

	// get the current go env
	goEnv, err := goEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "7zplugin: error: get go env: %v\n", err)
		os.Exit(2)
	}

	// show the current go env
	fmt.Println("go env")
	for _, x := range goEnv {
		key, value, _ := strings.Cut(x, "=")
		if strings.ContainsFunc(value, unicode.IsSpace) {
			fmt.Printf("  %s=%q\n", key, value)
		} else {
			fmt.Printf("  %s=%s\n", key, value)
		}
	}
	fmt.Println()

	// extract and set env vars from args
	fmt.Println("go env override")
	var env int
	for env = 1; env < len(os.Args); env++ {
		key, value, ok := strings.Cut(os.Args[env], "=")
		if !ok {
			break
		}
		if strings.ContainsFunc(value, unicode.IsSpace) {
			fmt.Printf("  %s=%q\n", key, value)
		} else {
			fmt.Printf("  %s=%s\n", key, value)
		}
		if err := os.Setenv(key, value); err != nil {
			fmt.Fprintf(os.Stderr, "7zplugin: error: setenv(%q, %q): %v\n", key, value, err)
			os.Exit(1)
		}
	}
	os.Args = slices.Delete(os.Args, 1, env)
	fmt.Println()

	// extract the output filename from args if present, or set the default
	var out string
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}
		if arg, ok := strings.CutPrefix(arg, "-o="); ok {
			out = arg
			os.Args = slices.Delete(os.Args, i, i+1)
		} else if arg == "-o" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "error: error: expected output file after -o flag, got nothing\n")
				os.Exit(1)
			}
			out = os.Args[i+1]
			os.Args = slices.Delete(os.Args, i, i+2)
		}
	}
	if out == "" {
		out = dllname + arch + ".dll"
	}

	// extract package names from args
	var (
		pkg   []string
		files []string
	)
	for i := len(os.Args) - 1; i >= 1; i-- {
		x := os.Args[i]
		if strings.HasPrefix(x, "-") {
			break // found a flag
		}
		if strings.HasSuffix(x, ".go") {
			if _, err := os.Stat(x); err != nil {
				fmt.Fprintf(os.Stderr, "7zplugin: error: expand package list: file %q is not accessible: %v\n", x, err)
				if i > 1 && strings.HasPrefix(os.Args[i-1], "-") {
					fmt.Fprintf(os.Stderr, "7zplugin: note: if you are trying to specify a go build flag with a value, you must use the syntax '-flag=value rather than -flag value'\n")
				}
				os.Exit(1)
			}
			files = append(files, x)
		} else {
			pkg = append(pkg, x)
		}
		os.Args = slices.Delete(os.Args, i, i+1)
	}
	if pkg, err = goList(pkg...); err != nil {
		fmt.Fprintf(os.Stderr, "7zplugin: error: expand package list: %v\n", err)
		os.Exit(1)
	}

	// add the default build flags
	os.Args = slices.Insert(os.Args, 1, "-buildmode=c-shared")

	// get the build date
	built := time.Now().UTC()
	if x, ok := os.LookupEnv("SOURCE_DATE_EPOCH"); ok {
		n, err := strconv.ParseInt(x, 0, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "7zplugin: error: parse SOURCE_DATE_EPOCH: %v\n", err)
			os.Exit(1)
		}
		built = time.Unix(n, 0)
	}

	// set the version info and default filename
	var pluginDesc []string
	for _, x := range pkg {
		pluginDesc = append(pluginDesc, strings.TrimPrefix(x, "github.com/pg9182/7zplugin/plugins/"))
	}
	for _, x := range files {
		pluginDesc = append(pluginDesc, strings.TrimSuffix(x, ".go"))
	}
	var ver = goversioninfo.VersionInfo{
		FixedFileInfo: goversioninfo.FixedFileInfo{
			FileVersion: goversioninfo.FileVersion{
				Major: built.Year(),
				Minor: int(built.Month()),
				Patch: built.Day(),
				Build: 0,
			},
		},
		StringFileInfo: goversioninfo.StringFileInfo{
			ProductName:      "Go Plugins for 7-Zip (github.com/pg9182/7zplugin)",
			ProductVersion:   "7-Zip " + z7.MY_VERSION_NUMBERS + "+ (" + z7arch + ")",
			FileDescription:  strings.Join(pluginDesc, ", "),
			InternalName:     strings.ToUpper(dllname),
			OriginalFilename: strings.ToUpper(dllname + arch + ".dll"),
		},
	}
	if _, ok := os.LookupEnv("GITHUB_ACTION"); ok {
		if v, ok := os.LookupEnv("GITHUB_REPOSITORY_OWNER"); ok {
			ver.StringFileInfo.CompanyName = v + " (github.com/" + v + ")"
		}
		if v, ok := os.LookupEnv("GITHUB_RUN_NUMBER"); ok {
			if v, err := strconv.ParseInt(v, 10, 16); err == nil {
				ver.FixedFileInfo.FileVersion.Build = int(v)
			}
		}
		if v, ok := os.LookupEnv("GITHUB_SHA"); ok {
			if len(v) > 7 {
				v = v[:7]
			}
			ver.StringFileInfo.FileVersion = " (" + v + ")" + ver.StringFileInfo.FileVersion
		}
	}
	ver.StringFileInfo.FileVersion = ver.FixedFileInfo.FileVersion.GetVersionString() + ver.StringFileInfo.FileVersion

	// catch ctrl+c
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// run the build
	if err = build(ctx, dir, os.Args[1:], files, pkg, out, ver); err != nil {

		// got ctrl+c
		if errors.Is(err, context.Canceled) {
			fmt.Fprintf(os.Stderr, "interrupted\n")
		}

		// go build exit status non-zero
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			status := ee.ExitCode()
			fmt.Fprintf(os.Stderr, "go build exited with status %d\n", status)
			os.Exit(status)
		}

		// other error
		fmt.Fprintf(os.Stderr, "7zplugin: error: %v\n", err)
		os.Exit(1)
	}
}

func build(ctx context.Context, dir string, flags []string, files []string, pkg []string, out string, ver goversioninfo.VersionInfo) error {
	var err error

	// make a temp dir
	td := filepath.Join(dir, "z7plugin-build") // fixed name for reproducibility
	if err := os.Mkdir(td, 0777); err != nil {
		return fmt.Errorf("create build dir: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(td); err != nil {
			fmt.Printf("warning: failed to remove temp dir %s: %v\n", td, err)
			return
		}
	}()

	// show the temp dir
	fmt.Printf("in %s\n\n", td)

	// gitignore everything in it
	if err := os.WriteFile(filepath.Join(td, ".gitignore"), []byte("*\n"), 0666); err != nil {
		return fmt.Errorf("write gitgnore: %w", err)
	}

	// copy and show the standalone go files
	for _, x := range files {
		if buf, err := os.ReadFile(x); err != nil {
			return fmt.Errorf("copy standalone go file %s: read: %w", x, err)
		} else if err := os.WriteFile(filepath.Join(td, filepath.Base(x)), buf, 0666); err != nil {
			return fmt.Errorf("copy standalone go file %s: write: %w", x, err)
		}
		fmt.Printf("%s\n  < %s\n\n", filepath.Base(x), filepath.Join(dir, x))
	}

	// generate the go source
	var src bytes.Buffer
	fmt.Fprintln(&src, "package main")
	fmt.Fprintln(&src)
	fmt.Fprintf(&src, "import _ %q\n", "github.com/pg9182/7zplugin/z7plugin")
	for _, x := range pkg {
		fmt.Fprintf(&src, "import _ %q\n", x)
	}
	fmt.Fprintln(&src)
	fmt.Fprintln(&src, "func main() {}")

	// save the go source
	if err := writeFileExcl(filepath.Join(td, "z7plugin.go"), src.Bytes(), 0666); err != nil {
		return fmt.Errorf("save generated go source: %w", err)
	}

	// show the generated source file
	fmt.Println("z7plugin.go")
	sc, line := bufio.NewScanner(bytes.NewReader(src.Bytes())), 0
	for sc.Scan() {
		line++
		fmt.Printf("%3d | %s\n", line, sc.Text())
	}
	fmt.Println()

	// generate the resource files
	ver.Build()
	ver.Walk()

	// save the resource file syso
	if err := ver.WriteSyso(filepath.Join(td, "rsrc.syso"), os.Getenv("GOARCH")); err != nil {
		return fmt.Errorf("generate rsrc syso: %w", err)
	}

	// show the generated resource file
	fmt.Println("rsrc.syso")
	for v, i := reflect.ValueOf(ver.StringFileInfo), 0; i < v.NumField(); i++ {
		if f := v.Type().Field(i); f.Type.Kind() == reflect.String {
			if x := v.Field(i).String(); x != "" {
				fmt.Printf("%3d | %-16s   %q\n", i+1, f.Name, x)
			}
		}
	}
	fmt.Println()

	// resolve the output path
	if out, err = filepath.Abs(out); err != nil {
		return fmt.Errorf("resolve output path: %w", err)
	}

	// actually run the build
	cmd := exec.CommandContext(ctx, "go", "build")
	cmd.Args = append(cmd.Args, "-o", out)
	cmd.Args = append(cmd.Args, flags...)
	cmd.Dir = td
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = nil

	// show the build command
	for i, arg := range cmd.Args {
		if strings.ContainsFunc(arg, unicode.IsSpace) {
			arg = strconv.Quote(arg)
		}
		if i != 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg)
	}
	fmt.Println()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build: %w", err)
	}
	return nil
}

func goEnv() ([]string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("go", "env", "-json")
	cmd.Stdin = nil
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("run command %q: %w", cmd.Args, err)
	}

	var obj map[string]string
	if err := json.Unmarshal(buf.Bytes(), &obj); err != nil {
		return nil, fmt.Errorf("parse output of command %q: %w", cmd.Args, err)
	}

	var env []string
	for k, v := range obj {
		env = append(env, k+"="+v)
	}
	sort.Strings(env)

	return env, nil
}

func goList(pkg ...string) ([]string, error) {
	var buf bytes.Buffer

	cmd := exec.Command("go", "list", "-json", "--")
	cmd.Args = append(cmd.Args, pkg...)
	cmd.Stdin = nil
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("run command %q: %w", cmd.Args, err)
	}

	var imp []string
	for dec := json.NewDecoder(&buf); ; {
		var obj struct {
			ImportPath string
		}
		if err := dec.Decode(&obj); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("parse output of command %q: %w", cmd.Args, err)
		}
		if obj.ImportPath == "" {
			return nil, fmt.Errorf("parse output of command %q: missing ImportPath", cmd.Args)
		}
		imp = append(imp, obj.ImportPath)
	}
	return imp, nil
}

func writeFileExcl(name string, data []byte, perm os.FileMode) error {
	f, err := os.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, perm)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	if err != nil {
		os.Remove(name)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(name)
		return err
	}
	return nil
}
