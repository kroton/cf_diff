package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
)

type Problem struct {
	ContestID int    `json:"contestId"`
	Index     string `json:"index"`
	Name      string `json:"name"`
}

type Problems []Problem

func (s Problems) Len() int           { return len(s) }
func (s Problems) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s Problems) Less(i, j int) bool { return s[i].ContestID > s[j].ContestID }

type Submission struct {
	Verdict string  `json:"verdict"`
	Problem Problem `json:"problem"`
}

type Contest struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func fetchAcceptedProblems(handle string) []Problem {
	resq, err := http.Get("http://codeforces.com/api/user.status?from=1&handle=" + handle)
	if err != nil {
		panic(err)
	}
	defer resq.Body.Close()
	var us struct {
		Submissions []Submission `json:"result"`
	}
	if err := json.NewDecoder(resq.Body).Decode(&us); err != nil {
		panic(err)
	}
	var res []Problem
	for _, s := range us.Submissions {
		if s.Verdict == "OK" {
			res = append(res, s.Problem)
		}
	}
	return res
}

func fetchContests() map[int]Contest {
	resq, err := http.Get("http://codeforces.com/api/contest.list")
	if err != nil {
		panic(err)
	}
	defer resq.Body.Close()
	var cl struct {
		Contests []Contest `json:"result"`
	}
	if err := json.NewDecoder(resq.Body).Decode(&cl); err != nil {
		panic(err)
	}
	res := map[int]Contest{}
	for _, c := range cl.Contests {
		res[c.ID] = c
	}
	return res
}

func main() {
	if len(os.Args) < 3 {
		return
	}
	handle1, handle2 := os.Args[1], os.Args[2]

	ps1 := fetchAcceptedProblems(handle1)
	ps2 := fetchAcceptedProblems(handle2)
	diff := map[string]Problem{}
	for _, p := range ps1 {
		diff[p.Name] = p
	}
	for _, p := range ps2 {
		delete(diff, p.Name)
	}
	var ps []Problem
	for _, p := range diff {
		ps = append(ps, p)
	}
	sort.Sort(Problems(ps))

	cm := fetchContests()
	for _, p := range ps {
		if c, ok := cm[p.ContestID]; ok {
			fmt.Printf("%s / %s / %s\n", c.Name, p.Index, p.Name)
		}
	}
}
