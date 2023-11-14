package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

var rFlag = flag.String("r", "", "Корневая папка всех проектов")
var pFlag = flag.String("p", "", "Имена папок проектов через запятую")
var dFlag = flag.Bool("d", false, "Собирать только прямые зависимости")

const (
	licenseChecker = "license-checker-rseidelsohn"
	git            = "git"
	npm            = "npm"
)

type ModuleInfo struct {
	Repository  string `json:"repository,omitempty"`  // Repository URL
	Publisher   string `json:"publisher,omitempty"`   // Publisher name
	Email       string `json:"email,omitempty"`       // Publisher e-mail
	Licenses    string `json:"licenses,omitempty"`    // Array of licenses
	LicenseFile string `json:"licenseFile,omitempty"` // Path to license file, if available
	Path        string `json:"path,omitempty"`        // Path to module
}

type RepoInfo struct {
	Type string `json:"type,omitempty"`
	Url  string `json:"url,omitempty"`
}

func main() {
	var root string
	var dirs []string
	var direct bool

	flag.Parse()

	if _, err := exec.LookPath(npm); err != nil {
		log.Fatalf("%s not found in $PATH", npm)
	}

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

	allDeps := make(map[string]RepoInfo)

	sort.Strings(dirs)

	for _, v := range dirs {
		deps := grabPackageNames(root, v, direct)

		fmt.Printf("Собрано %d зависимостей по проекту %s...\n", len(deps), v)

		for _, v := range deps {
			allDeps[trimVersion(v)] = RepoInfo{}
		}
	}

	fmt.Printf("Собрано %d зависимостей по всем проектам...\n", len(allDeps))

	allDepsLen := len(allDeps)
	counter := 0

	for k := range allDeps {
		counter++
		fmt.Printf("%d/%d: Получение репозитория зависимости %q...\n", counter, allDepsLen, k)
		allDeps[k] = getRepoInfo(k)
	}

	fmt.Printf("Собраны все репозитории...\n")

	data, err := json.MarshalIndent(allDeps, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}

	filename := path.Join(root, strings.Join(dirs, "_")+"_repos_info.json")

	os.WriteFile(filename, data, 0666)

	fmt.Printf("Данные записаны в файл %q\n", filename)

	fmt.Println("Готово!")
}

func getRepoInfo(name string) RepoInfo {
	args := []string{
		"info",
		name,
		"--json",
		"repository",
	}

	cmd := exec.Command(npm, args...)
	data, err := cmd.Output()

	if err != nil {
		log.Fatalln(err)
	}

	var info RepoInfo

	if err := json.Unmarshal(data, &info); err != nil {
		log.Fatalln(err)
	}

	return info
}

func trimVersion(name string) string {
	i := strings.LastIndex(name, "@")

	if i <= 0 {
		return name
	}

	return string([]rune(name)[:i])
}

func grabPackageNames(root, name string, direct bool) []string {
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

	deps := make(map[string]any)

	if err := json.Unmarshal(data, &deps); err != nil {
		log.Fatal(err)
	}

	names := make([]string, 0, len(deps))

	for k := range deps {
		names = append(names, k)
	}

	return names
}
