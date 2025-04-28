package main

import (
	"context"
	"math"
	"time"

	"github.com/google/go-github/v69/github"
)

type Repo struct {
	repo   string
	owner  string
	client *github.Client
}

func New(repo, owner string, client *github.Client) *Repo {
	return &Repo{repo: repo, owner: owner, client: client}
}

func (r *Repo) stats() (ret []Data, err error) {
	repo, _, err := r.client.Repositories.Get(context.TODO(), r.owner, r.repo)
	if err != nil {
		return nil, err
	}
	contributors, _, err := r.client.Repositories.ListContributorsStats(context.TODO(), r.owner, r.repo)
	if err != nil {
		return nil, err
	}
	commits, _, err := r.client.Repositories.ListCommits(context.TODO(), r.owner, r.repo, &github.CommitsListOptions{})
	if err != nil {
		return nil, err
	}

	ret = []Data{
		{
			Key:   "stars",
			Value: math.Log(float64(repo.GetStargazersCount() + 1)),
		},
		{
			Key:   "forks",
			Value: math.Log(float64(repo.GetForksCount() + 1)),
		},
		{
			Key:   "watchers",
			Value: math.Log(float64(repo.GetSubscribersCount() + 1)),
		},
		{
			Key:   "issues",
			Value: math.Log(float64(repo.GetOpenIssues() + 1)),
		},
		{
			Key:   "networks",
			Value: math.Log(float64(repo.GetNetworkCount() + 1)),
		},
		{
			Key:   "contributors",
			Value: math.Log(float64(len(contributors) + 1)),
		},
		{
			Key:   "commit_rate",
			Value: math.Log(float64(len(commits))/float64(repo.UpdatedAt.Unix()-repo.CreatedAt.Unix()) + 1),
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
		time.Sleep((1 << i) * time.Second)
	}

	return sum
}
