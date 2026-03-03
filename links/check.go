package links

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dacharyc/skill-validator/types"
)

type linkResult struct {
	url    string
	result types.Result
}

// CheckLinks validates external (HTTP/HTTPS) links in the skill body.
func CheckLinks(dir, body string) []types.Result {
	ctx := types.ResultContext{Category: "Links", File: "SKILL.md"}
	allLinks := ExtractLinks(body)
	if len(allLinks) == 0 {
		return nil
	}

	var (
		results   []types.Result
		httpLinks []string
		mu        sync.Mutex
		wg        sync.WaitGroup
	)

	// Collect HTTP links only
	for _, link := range allLinks {
		// Skip template URLs containing {placeholder} variables (RFC 6570 URI Templates)
		if strings.Contains(link, "{") {
			continue
		}
		if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
			httpLinks = append(httpLinks, link)
		}
	}

	if len(httpLinks) == 0 {
		return nil
	}

	// Check HTTP links concurrently
	httpResults := make([]linkResult, len(httpLinks))
	for i, link := range httpLinks {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()
			r := checkHTTPLink(ctx, url)
			mu.Lock()
			httpResults[idx] = linkResult{url: url, result: r}
			mu.Unlock()
		}(i, link)
	}
	wg.Wait()

	for _, hr := range httpResults {
		results = append(results, hr.result)
	}

	return results
}

func checkHTTPLink(ctx types.ResultContext, url string) types.Result {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return ctx.Errorf("%s (invalid URL: %v)", url, err)
	}
	req.Header.Set("User-Agent", "skill-validator/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return ctx.Errorf("%s (request failed: %v)", url, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return ctx.Passf("%s (HTTP %d)", url, resp.StatusCode)
	}
	if resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return ctx.Passf("%s (HTTP %d redirect)", url, resp.StatusCode)
	}
	if resp.StatusCode == http.StatusForbidden {
		return ctx.Infof("%s (HTTP 403 — may block automated requests)", url)
	}
	return ctx.Errorf("%s (HTTP %d)", url, resp.StatusCode)
}
