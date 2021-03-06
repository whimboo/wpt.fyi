// Copyright 2017 The WPT Dashboard Project. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package shared

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/deckarep/golang-set"
)

// FetchLatestRuns fetches the TestRun metadata for the latest runs, using the
// API on the given host.
func FetchLatestRuns(wptdHost string) TestRuns {
	return FetchRuns(wptdHost, "latest", nil, nil)
}

// FetchRuns fetches the TestRun metadata for the given sha / labels, using the
// API on the given host.
func FetchRuns(wptdHost, sha string, maxCount *int, labels mapset.Set) TestRuns {
	url := "https://" + wptdHost + "/api/runs"

	filters := TestRunFilter{
		Labels:   labels,
		SHA:      sha,
		MaxCount: maxCount,
	}
	url += "?" + filters.ToQuery(true).Encode()
	log.Printf("Fetching %s...", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		log.Fatal(errors.New("Bad response code from " + url + ": " +
			strconv.Itoa(resp.StatusCode)))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var runs TestRuns
	if err := json.Unmarshal(body, &runs); err != nil {
		log.Fatal(err)
	}
	return runs
}
