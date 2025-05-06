package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

type repoInfo struct {
	Info struct {
		Branch    string `json:"defaultBranch"`
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		PushedAt  string `json:"pushedAt"`
		Stars     int64  `json:"stars"`
		Watchers  int64  `json:"watchers"`
		Forks     int64  `json:"forks"`
	} `json:"repo"`
}

type Repo struct {
	repo  string
	owner string
}

func New(repo, owner string) *Repo {
	return &Repo{repo: repo, owner: owner}
}

func (r *Repo) info() (ret *repoInfo, err error) {
	resp, err := http.Get(fmt.Sprintf("https://ungh.cc/repos/%s/%s", r.owner, r.repo))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result repoInfo
	err = json.NewDecoder(resp.Body).Decode(&result)

	ret = &result

	return
}
func (r *Repo) contributors() (ret []any, err error) {
	resp, err := http.Get(fmt.Sprintf("https://ungh.cc/repos/%s/%s/contributors", r.owner, r.repo))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var result struct {
		Contributors []any `json:"contributors"`
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	ret = result.Contributors
	return
}

func (r *Repo) stats() (ret []Data, err error) {
	repo, err := r.info()
	if err != nil {
		return nil, err
	}
	contributors, err := r.contributors()
	if err != nil {
		return nil, err
	}

	ret = []Data{
		{
			Key:   "stars",
			Value: math.Log(float64(repo.Info.Stars + 1)),
		},
		{
			Key:   "forks",
			Value: math.Log(float64(repo.Info.Forks + 1)),
		},
		{
			Key:   "watchers",
			Value: math.Log(float64(repo.Info.Watchers + 1)),
		},
		{
			Key:   "contributors",
			Value: math.Log(float64(len(contributors) + 1)),
		},
	}

	return
}

func (r *Repo) Score() float64 {
	var sum float64

	for i := range 10 {
		datas, err := r.stats()
		if err == nil {
			return sumData(datas)
		}
		log.Println("get stats fail: ", err)
		time.Sleep((1 << i) * time.Second)
	}

	return sum
}
