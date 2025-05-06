package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"sync"

	"github.com/MeteorsLiu/reporank/info"
)

var matchRepoRegex = regexp.MustCompile("github.com/(.*?)/(.*?)/")

type pkgWithScore struct {
	Info  *info.PkgInfo
	Score float64
}

func main() {
	var dir string
	var token string
	var repoURL string
	flag.StringVar(&dir, "dir", "", "Conan center index")
	flag.StringVar(&repoURL, "url", "", "Github Repo URL")
	flag.Parse()

	pkgs, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(token)

	var allPkgsMu sync.Mutex
	var allPkgs []pkgWithScore

	var wg sync.WaitGroup

	for _, pkg := range pkgs {
		pkgInfo, err := info.ReadPackageInfoWithReturn(pkg.Name(), dir)
		if err != nil {
			continue
		}
		results := matchRepoRegex.FindStringSubmatch(pkgInfo.URLs[0])
		if len(results) != 3 {
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			score := New(results[2], results[1]).Score()
			fmt.Println(pkgInfo, score)

			allPkgsMu.Lock()
			defer allPkgsMu.Unlock()

			thisPkg := pkgWithScore{
				Info:  pkgInfo,
				Score: score,
			}
			allPkgs = append(allPkgs, thisPkg)

			fmt.Println(thisPkg.Score, thisPkg.Info)
		}()
	}

	wg.Wait()

	sort.Slice(allPkgs, func(i, j int) bool {
		return allPkgs[i].Score > allPkgs[j].Score
	})

	b, _ := json.Marshal(&allPkgs)
	os.WriteFile("result.json", b, 0644)

	// // results := matchRepoRegex.FindStringSubmatch(repoURL)
	// // if len(results) != 3 {
	// // 	log.Fatal("not a github repo")
	// // }

	// // fmt.Println(New(results[2], results[1], github.NewClient(nil).WithAuthToken(token)).Score())
}
