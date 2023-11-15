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
	Name  string `json:"name,omitempty"`
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

	allDepsMap := make(map[string]bool)
	dirsCh := make(chan []string)

	for _, v := range dirs {
		go grabPackageNames(root, v, direct, dirsCh)
	}

	for range dirs {
		deps := <-dirsCh
		for _, v := range deps {
			allDepsMap[trimVersion(v)] = true
		}
	}

	allDeps := make([]string, 0, len(allDepsMap))

	for k := range allDepsMap {
		allDeps = append(allDeps, k)
	}

	fmt.Printf("Собрано %d зависимостей по всем проектам\n", len(allDeps))

	allRepos := make([]RepoInfo, 0, len(allDeps))
	repoCh := make(chan RepoInfo)
	count := 0
	l := len(allDeps)

	deps := allDeps

	for len(deps) > 0 {
		var t []string

		if len(deps) < 10 {
			t = deps
			deps = []string{}
		} else {
			t = deps[:10]
			deps = deps[10:]
		}

		for _, v := range t {
			go getRepoInfo(v, repoCh)
		}

		for range t {
			i := <-repoCh
			count++
			fmt.Printf("(%d/%d) Получен репозиторий пакета %s\n", count, l, i.Name)
			allRepos = append(allRepos, i)
		}
	}

	// for _, v := range allDeps {
	// 	go getRepoInfo(v, repoCh)
	// }

	// for range allDeps {
	// 	i := <-repoCh
	// 	count++
	// 	fmt.Printf("(%d/%d) Получен репозиторий пакета %s\n", count, l, i.Name)
	// 	allRepos = append(allRepos, i)
	// }

	fmt.Printf("Собраны все репозитории\n")

	sort.SliceStable(allRepos, func(i, j int) bool {
		return allRepos[i].Name < allRepos[j].Name
	})

	data, err := json.MarshalIndent(allRepos, "", "  ")

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

func getRepoInfo(name string, ch chan<- RepoInfo) {
	if _, err := exec.LookPath(npm); err != nil {
		log.Fatalf("%s not found in $PATH", npm)
	}

	args := []string{
		"info",
		name,
		"--json",
		"repository",
	}

	info := RepoInfo{Name: name}
	cmd := exec.Command(npm, args...)
	data, err := cmd.Output()

	if err != nil {
		log.Fatalln(err)
	}

	if err := json.Unmarshal(data, &info); err != nil {
		info.Error = &err
	}

	ch <- info
}

func trimVersion(name string) string {
	i := strings.LastIndex(name, "@")

	if i <= 0 {
		return name
	}

	return string([]rune(name)[:i])
}

func grabPackageNames(root, name string, direct bool, ch chan<- []string) {
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

	fmt.Printf("Собрано %d зависимостей по проекту %s\n", len(deps), name)

	ch <- names
}
