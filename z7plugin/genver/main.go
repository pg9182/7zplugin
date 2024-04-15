// Command genver is an optional helper to generate version info for a plugin.
package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/josephspurrier/goversioninfo"
	"github.com/pg9182/7zplugin/z7"
)

var (
	Owner       = flag.String("owner", "", "only use github actions info if this is the repo owner")
	CompanyName = flag.String("company", "", "set the company name")
	Copyright   = flag.String("copyright", "", "set the copyright")
	Description = flag.String("description", "", "set the file description (set based on the dir name if not provided)")
	Major       = flag.Int("major", 0, "set the file major version")
	Minor       = flag.Int("minor", 0, "set the file minor version")
	Patch       = flag.Int("patch", 0, "set the file patch version")
)

func main() {
	flag.Parse()

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	pkg := filepath.Base(dir)

	var ver goversioninfo.VersionInfo

	ver.FixedFileInfo.FileVersion = goversioninfo.FileVersion{
		Major: *Major,
		Minor: *Minor,
		Patch: *Patch,
		Build: 0,
	}
	ver.FixedFileInfo.ProductVersion = goversioninfo.FileVersion{
		Major: 0,
		Minor: 0,
		Patch: 0,
		Build: 0,
	}
	ver.StringFileInfo.CompanyName = *CompanyName
	if *Description != "" {
		ver.StringFileInfo.FileDescription = *Description
	} else {
		ver.StringFileInfo.FileDescription = strings.ToUpper(pkg) + " Plugin for 7-Zip"
	}
	ver.StringFileInfo.FileVersion = ver.FixedFileInfo.FileVersion.GetVersionString() + " (devel)"
	ver.StringFileInfo.InternalName = strings.ToUpper(pkg)
	ver.StringFileInfo.LegalCopyright = *Copyright
	ver.StringFileInfo.ProductName = "github.com/pg9182/7zplugin"
	ver.StringFileInfo.ProductVersion = ver.FixedFileInfo.ProductVersion.GetVersionString() + " (devel), for 7-Zip " + z7.MY_VERSION_NUMBERS

	if _, ok := os.LookupEnv("CI"); ok {
		if _, ok := os.LookupEnv("GITHUB_ACTION"); ok {
			if v, ok := os.LookupEnv("GITHUB_REPOSITORY_OWNER"); (ok && v == *Owner) || *Owner == "" {
				if v, ok := os.LookupEnv("GITHUB_RUN_NUMBER"); !ok {
					panic("missing GITHUB_RUN_NUMBER")
				} else if v, err := strconv.ParseInt(v, 10, 16); err != nil {
					panic("invalid GITHUB_RUN_NUMBER: " + err.Error())
				} else {
					ver.FixedFileInfo.FileVersion.Build = int(v)
				}
				if v, ok := os.LookupEnv("GITHUB_SHA"); !ok {
					panic("missing GITHUB_SHA")
				} else {
					if len(v) > 7 {
						v = v[:7]
					}
					ver.StringFileInfo.FileVersion = ver.FixedFileInfo.FileVersion.GetVersionString() + " (" + v + ")"
					ver.StringFileInfo.ProductVersion = ver.FixedFileInfo.ProductVersion.GetVersionString() + " (" + v + ")"
				}
			}
			// so the update doesn't make the git repo dirty
			if err := os.WriteFile(".gitignore", []byte(".gitignore\n*.syso\n"), 0644); err != nil {
				panic(err)
			}
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(ver.StringFileInfo)

	ver.StringFileInfo.OriginalFilename = ver.StringFileInfo.InternalName + "64.DLL"
	ver.Build()
	ver.Walk()
	if err := ver.WriteSyso("rsrc_windows_amd64.syso", "amd64"); err != nil {
		panic(err)
	}

	ver.StringFileInfo.OriginalFilename = ver.StringFileInfo.InternalName + "A64.DLL"
	ver.Build()
	ver.Walk()
	if err := ver.WriteSyso("rsrc_windows_arm64.syso", "arm64"); err != nil {
		panic(err)
	}

	ver.StringFileInfo.OriginalFilename = ver.StringFileInfo.InternalName + "32.DLL"
	ver.Build()
	ver.Walk()
	if err := ver.WriteSyso("rsrc_windows_386.syso", "386"); err != nil {
		panic(err)
	}
}
