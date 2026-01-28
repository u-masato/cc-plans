package plan

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/masato-uno/cc-plans/internal/config"
)

var (
	ErrNotFound      = errors.New("plan not found")
	ErrAmbiguous     = errors.New("ambiguous plan name: multiple matches found")
	ErrPlansNotExist = errors.New("plans directory does not exist")
)

// Repository handles plan file operations.
type Repository struct {
	plansDir string
}

// NewRepository creates a new Repository.
func NewRepository() *Repository {
	return &Repository{
		plansDir: config.PlansDir(),
	}
}

// List returns all plans.
func (r *Repository) List() ([]Plan, error) {
	if _, err := os.Stat(r.plansDir); os.IsNotExist(err) {
		return nil, ErrPlansNotExist
	}

	entries, err := os.ReadDir(r.plansDir)
	if err != nil {
		return nil, err
	}

	var plans []Plan
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".md") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		path := filepath.Join(r.plansDir, name)
		plans = append(plans, Plan{
			Name:    strings.TrimSuffix(name, ".md"),
			Path:    path,
			ModTime: info.ModTime(),
			Size:    info.Size(),
			Title:   extractTitle(path),
		})
	}

	return plans, nil
}

// Get returns a plan by name (supports partial matching).
func (r *Repository) Get(name string) (*Plan, error) {
	plans, err := r.List()
	if err != nil {
		return nil, err
	}

	// Exact match first
	for _, p := range plans {
		if p.Name == name {
			return &p, nil
		}
	}

	// Partial match
	var matches []Plan
	for _, p := range plans {
		if strings.Contains(strings.ToLower(p.Name), strings.ToLower(name)) {
			matches = append(matches, p)
		}
	}

	if len(matches) == 0 {
		return nil, ErrNotFound
	}
	if len(matches) > 1 {
		return nil, ErrAmbiguous
	}

	return &matches[0], nil
}

// GetContent returns the content of a plan.
func (r *Repository) GetContent(name string) (string, error) {
	plan, err := r.Get(name)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(plan.Path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Search searches plans by query.
func (r *Repository) Search(query string, nameOnly bool) ([]SearchResult, error) {
	plans, err := r.List()
	if err != nil {
		return nil, err
	}

	var results []SearchResult
	queryLower := strings.ToLower(query)

	for _, p := range plans {
		if nameOnly {
			if strings.Contains(strings.ToLower(p.Name), queryLower) {
				results = append(results, SearchResult{Plan: p})
			}
			continue
		}

		// Search in content
		matches, err := searchInFile(p, queryLower)
		if err != nil {
			continue
		}
		results = append(results, matches...)
	}

	return results, nil
}

func searchInFile(p Plan, query string) ([]SearchResult, error) {
	f, err := os.Open(p.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var results []SearchResult
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), query) {
			results = append(results, SearchResult{
				Plan:       p,
				MatchLine:  line,
				LineNumber: lineNum,
			})
		}
	}

	return results, scanner.Err()
}

// SortByModTime sorts plans by modification time (newest first).
func SortByModTime(plans []Plan) {
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].ModTime.After(plans[j].ModTime)
	})
}

// SortByName sorts plans by name.
func SortByName(plans []Plan) {
	sort.Slice(plans, func(i, j int) bool {
		return plans[i].Name < plans[j].Name
	})
}
