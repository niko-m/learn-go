package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var rFlag = flag.String("r", "", "Корневая папка всех проектов")
var pFlag = flag.String("p", "", "Имена папок проектов через запятую")
var dFlag = flag.Bool("d", false, "Собирать только прямые зависимости")

const (
	licenseChecker = "license-checker-rseidelsohn"
	git            = "git"
)

type ModuleInfo struct {
	Repository  *string `json:"repository,omitempty"`  // Repository URL
	Publisher   *string `json:"publisher,omitempty"`   // Publisher name
	Email       *string `json:"email,omitempty"`       // Publisher e-mail
	Licenses    *string `json:"licenses,omitempty"`    // Array of licenses
	LicenseFile *string `json:"licenseFile,omitempty"` // Path to license file, if available
	Path        *string `json:"path,omitempty"`        // Path to module
}

func main() {
	var root string
	var dirs []string
	var direct bool

	flag.Parse()

	if _, err := exec.LookPath(git); err != nil {
		log.Fatalf("%s not found in $PATH", git)
	}

	if _, err := exec.LookPath(licenseChecker); err != nil {
		log.Fatalf("%s not found in $PATH", licenseChecker)
	}

	if r, err := filepath.Abs(*rFlag); err != nil {
		log.Fatalf("Invalid root path:\n%q", err)
	} else {
		root = r
	}

	if len(*pFlag) == 0 {
		log.Fatalf("projects paths are empty\n")
	} else {
		dirs = strings.Split(*pFlag, ",")
	}

	direct = *dFlag

	for _, v := range dirs {
		grabLicenses(root, v, direct)
	}
}

func grabLicenses(root, name string, direct bool) {
	args := []string{
		"--start",
		path.Join(root, name),
		"--excludePackages",
		"app",
		"--json",
	}

	if direct {
		args = append(args, "--direct", "0")
	}

	cmd := exec.Command("license-checker-rseidelsohn", args...)

	data, err := cmd.Output()

	if err != nil {
		log.Fatalf("%q\n", err)
	}

	deps := make(map[string]ModuleInfo)

	json.Unmarshal(data, &deps)

	for k, v := range deps {
		if v.Repository == nil {
			fmt.Println(k, "has no repo")
		}
	}
}
