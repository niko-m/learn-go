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
	"time"
)

var rFlag = flag.String("r", "", "Корневая папка всех проектов")
var pFlag = flag.String("p", "", "Имена папок проектов через запятую (относительно корня)")
var oFlag = flag.String("o", "deprepos", "Имя выходного файла (относительно корня)")
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
	Type  string `json:"type,omitempty"`
	Url   string `json:"url,omitempty"`
	Error *error `json:"error,omitempty"`
}

func main() {
	var root string
	var dirs []string
	var direct bool
	var out string

	flag.Parse()

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
	out = *oFlag

	allDeps := make(map[string]RepoInfo)

	sort.Strings(dirs)

	for _, v := range dirs {
		deps := grabPackageNames(root, v, direct)

		fmt.Printf("Собрано %d зависимостей по проекту %s\n", len(deps), v)

		for _, v := range deps {
			allDeps[trimVersion(v)] = RepoInfo{}
		}
	}

	fmt.Printf("Собрано %d зависимостей по всем проектам\n", len(allDeps))

	allDepsLen := len(allDeps)
	counter := 0

	for k := range allDeps {
		counter++
		fmt.Printf("%d/%d: Получение репозитория зависимости %q...\n", counter, allDepsLen, k)
		i := getRepoInfo(k)
		allDeps[k] = i
		if i.Error != nil {
			fmt.Println("Ошибка:", i.Error)
		}
	}

	fmt.Printf("Собраны все репозитории\n")

	data, err := json.MarshalIndent(allDeps, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}

	filename := path.Join(root, fmt.Sprintf("%s_%d.json", out, time.Now().Unix()))

	fmt.Printf("Запись данных в файл %q...\n", filename)

	os.WriteFile(filename, data, 0666)

	fmt.Println("Готово!")
}

// func cloneRepo(repo string) {
// 	if _, err := exec.LookPath(git); err != nil {
// 		log.Fatalf("%s not found in $PATH", git)
// 	}
// }

func getRepoInfo(name string) RepoInfo {
	if _, err := exec.LookPath(npm); err != nil {
		log.Fatalf("%s not found in $PATH", npm)
	}

	args := []string{
		"info",
		name,
		"--json",
		"repository",
	}

	var info RepoInfo
	cmd := exec.Command(npm, args...)
	data, err := cmd.Output()

	if err != nil {
		info.Error = &err
		return info
	}

	if err := json.Unmarshal(data, &info); err != nil {
		info.Error = &err
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
	if _, err := exec.LookPath(licenseChecker); err != nil {
		log.Fatalf("%s not found in $PATH", licenseChecker)
	}

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

	cmd := exec.Command(licenseChecker, args...)

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
