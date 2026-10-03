package board

import (
	"slices"
	"strings"

	"fuda/internal/taskfiles"
)

func Delivers(prNumber int, before *taskfiles.Task, after taskfiles.Task) bool {
	becomesMerged := normalizeStatus(after.Status) == "merged" &&
		(before == nil || normalizeStatus(before.Status) != "merged")
	gainsNumber := slices.Contains(after.PRs, prNumber) &&
		(before == nil || !slices.Contains(before.PRs, prNumber))
	return becomesMerged || gainsNumber
}

func Search(tasks []taskfiles.Task, query string) []string {
	query = strings.ToLower(strings.TrimSpace(query))
	out := []string{}
	if query == "" {
		return out
	}
	for _, t := range tasks {
		if strings.Contains(strings.ToLower(t.Body), query) || strings.Contains(strings.ToLower(t.Title), query) {
			out = append(out, t.ID)
		}
	}
	return out
}
