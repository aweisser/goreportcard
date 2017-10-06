package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/boltdb/bolt"
	humanize "github.com/dustin/go-humanize"
	"github.com/gojp/goreportcard/check"
	"github.com/gojp/goreportcard/download"
	"github.com/gojp/goreportcard/report"
)

func dirName(repo string) string {
	return fmt.Sprintf("_repos/src/%s", repo)
}

func getFromCache(repo string) (report.Card, error) {
	// try and fetch from boltdb
	db, err := bolt.Open(DBPath, 0600, &bolt.Options{Timeout: 3 * time.Second})
	if err != nil {
		return report.Card{}, fmt.Errorf("failed to open bolt database during GET: %v", err)
	}
	defer db.Close()

	resp := report.Card{}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(RepoBucket))
		if b == nil {
			return errors.New("No repo bucket")
		}
		cached := b.Get([]byte(repo))
		if cached == nil {
			return fmt.Errorf("%q not found in cache", repo)
		}

		err = json.Unmarshal(cached, &resp)
		if err != nil {
			return fmt.Errorf("failed to parse JSON for %q in cache", repo)
		}
		return nil
	})

	if err != nil {
		return resp, err
	}

	resp.LastRefresh = resp.LastRefresh.UTC()
	resp.LastRefreshFormatted = resp.LastRefresh.Format(time.UnixDate)
	resp.LastRefreshHumanized = humanize.Time(resp.LastRefresh.UTC())

	return resp, nil
}

func newChecksResp(repo string, forceRefresh bool) (report.Card, error) {
	if !forceRefresh {
		resp, err := getFromCache(repo)
		if err != nil {
			// just log the error and continue
			log.Println(err)
		} else {
			resp.Grade = report.GradeOf(resp.Average * 100) // grade is not stored for some repos, yet
			return resp, nil
		}
	}

	// fetch the repo and grade it
	repoRoot, err := download.Download(repo, "_repos/src")
	if err != nil {
		return report.Card{}, fmt.Errorf("could not clone repo: %v", err)
	}

	repo = repoRoot.Root

	dir := dirName(repo)
	filenames, skipped, err := check.GoFiles(dir)
	if err != nil {
		return report.Card{}, fmt.Errorf("could not get filenames: %v", err)
	}
	if len(filenames) == 0 {
		return report.Card{}, fmt.Errorf("no .go files found")
	}

	err = check.RenameFiles(skipped)
	if err != nil {
		log.Println("Could not remove files:", err)
	}
	defer check.RevertFiles(skipped)

	checks := createChecks(dir, filenames)

	ch := make(chan report.Score)
	for _, c := range checks {
		go func(c report.Check) {
			p, summaries, err := c.Percentage()
			errMsg := ""
			if err != nil {
				log.Printf("ERROR: (%s) %v", c.Name(), err)
				errMsg = err.Error()
			}
			s := report.Score{
				Name:          c.Name(),
				Description:   c.Description(),
				FileSummaries: summaries,
				Weight:        c.Weight(),
				Percentage:    p,
				Error:         errMsg,
			}
			ch <- s
		}(c)
	}

	t := time.Now().UTC()
	resp := report.Card{
		Repo:                 repo,
		ResolvedRepo:         repoRoot.Repo,
		Files:                len(filenames),
		LastRefresh:          t,
		LastRefreshFormatted: t.Format(time.UnixDate),
		LastRefreshHumanized: humanize.Time(t),
	}

	var total, totalWeight float64
	var issues = make(map[string]bool)
	for i := 0; i < len(checks); i++ {
		s := <-ch
		resp.Checks = append(resp.Checks, s)
		total += s.Percentage * s.Weight
		totalWeight += s.Weight
		for _, fs := range s.FileSummaries {
			issues[fs.Filename] = true
		}
	}
	total /= totalWeight

	sort.Sort(report.ByWeight(resp.Checks))
	resp.Average = total
	resp.Issues = len(issues)
	resp.Grade = report.GradeOf(total * 100)

	return resp, nil
}
