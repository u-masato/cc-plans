package plan

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func newTestRepo(t *testing.T) (*Repository, string) {
	t.Helper()
	dir := t.TempDir()
	return &Repository{plansDir: dir}, dir
}

func writeTestPlan(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name+".md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func TestRepository_List_Empty(t *testing.T) {
	repo, _ := newTestRepo(t)
	plans, err := repo.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(plans) != 0 {
		t.Errorf("List returned %d plans, want 0", len(plans))
	}
}

func TestRepository_List(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "plan-a", "# Plan A\nContent A\n")
	writeTestPlan(t, dir, "plan-b", "# Plan B\nContent B\n")

	plans, err := repo.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(plans) != 2 {
		t.Fatalf("List returned %d plans, want 2", len(plans))
	}

	names := map[string]bool{}
	for _, p := range plans {
		names[p.Name] = true
		if p.Path == "" {
			t.Error("plan Path is empty")
		}
		if p.Size == 0 {
			t.Error("plan Size is 0")
		}
	}
	if !names["plan-a"] || !names["plan-b"] {
		t.Errorf("unexpected plan names: %v", names)
	}
}

func TestRepository_List_SkipsNonMd(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "plan-a", "# Plan A\n")
	// Write a non-md file
	os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("not a plan"), 0644)

	plans, err := repo.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(plans) != 1 {
		t.Errorf("List returned %d plans, want 1 (should skip .txt)", len(plans))
	}
}

func TestRepository_List_SkipsDirectories(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "plan-a", "# Plan A\n")
	os.Mkdir(filepath.Join(dir, "subdir.md"), 0755)

	plans, err := repo.List()
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(plans) != 1 {
		t.Errorf("List returned %d plans, want 1 (should skip directories)", len(plans))
	}
}

func TestRepository_List_NonexistentDir(t *testing.T) {
	repo := &Repository{plansDir: "/nonexistent/dir"}
	_, err := repo.List()
	if err != ErrPlansNotExist {
		t.Errorf("List error = %v, want ErrPlansNotExist", err)
	}
}

func TestRepository_Get_ExactMatch(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "my-plan", "# My Plan\n")

	p, err := repo.Get("my-plan")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if p.Name != "my-plan" {
		t.Errorf("Get Name = %q, want 'my-plan'", p.Name)
	}
}

func TestRepository_Get_PartialMatch(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "my-long-plan-name", "# Plan\n")

	p, err := repo.Get("long-plan")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if p.Name != "my-long-plan-name" {
		t.Errorf("Get Name = %q, want 'my-long-plan-name'", p.Name)
	}
}

func TestRepository_Get_CaseInsensitive(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "MyPlan", "# Plan\n")

	p, err := repo.Get("myplan")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	if p.Name != "MyPlan" {
		t.Errorf("Get Name = %q, want 'MyPlan'", p.Name)
	}
}

func TestRepository_Get_NotFound(t *testing.T) {
	repo, _ := newTestRepo(t)
	_, err := repo.Get("nonexistent")
	if err != ErrNotFound {
		t.Errorf("Get error = %v, want ErrNotFound", err)
	}
}

func TestRepository_Get_Ambiguous(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "plan-alpha", "# Alpha\n")
	writeTestPlan(t, dir, "plan-beta", "# Beta\n")

	_, err := repo.Get("plan")
	if err != ErrAmbiguous {
		t.Errorf("Get error = %v, want ErrAmbiguous", err)
	}
}

func TestRepository_GetContent(t *testing.T) {
	repo, dir := newTestRepo(t)
	expected := "# Test Plan\n\nThis is content.\n"
	writeTestPlan(t, dir, "test", expected)

	content, err := repo.GetContent("test")
	if err != nil {
		t.Fatalf("GetContent returned error: %v", err)
	}
	if content != expected {
		t.Errorf("GetContent = %q, want %q", content, expected)
	}
}

func TestRepository_GetContent_NotFound(t *testing.T) {
	repo, _ := newTestRepo(t)
	_, err := repo.GetContent("nonexistent")
	if err != ErrNotFound {
		t.Errorf("GetContent error = %v, want ErrNotFound", err)
	}
}

func TestRepository_Delete(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "to-delete", "# Delete Me\n")

	err := repo.Delete("to-delete")
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}

	// Verify file is gone
	path := filepath.Join(dir, "to-delete.md")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("file still exists after Delete")
	}
}

func TestRepository_Delete_NotFound(t *testing.T) {
	repo, _ := newTestRepo(t)
	err := repo.Delete("nonexistent")
	if err != ErrNotFound {
		t.Errorf("Delete error = %v, want ErrNotFound", err)
	}
}

func TestRepository_Search_ByContent(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "plan-a", "# Plan A\nThis has keyword magic in it.\n")
	writeTestPlan(t, dir, "plan-b", "# Plan B\nNo special words here.\n")

	results, err := repo.Search("magic", false)
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Search returned %d results, want 1", len(results))
	}
	if results[0].Plan.Name != "plan-a" {
		t.Errorf("Search result = %q, want 'plan-a'", results[0].Plan.Name)
	}
	if results[0].LineNumber != 2 {
		t.Errorf("Search LineNumber = %d, want 2", results[0].LineNumber)
	}
}

func TestRepository_Search_ByName(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "feature-auth", "# Auth\n")
	writeTestPlan(t, dir, "bugfix-login", "# Login Fix\n")

	results, err := repo.Search("auth", true)
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Search returned %d results, want 1", len(results))
	}
	if results[0].Plan.Name != "feature-auth" {
		t.Errorf("Search result = %q, want 'feature-auth'", results[0].Plan.Name)
	}
}

func TestRepository_Search_CaseInsensitive(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "plan", "# Plan\nHas UPPERCASE keyword.\n")

	results, err := repo.Search("uppercase", false)
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Search returned %d results, want 1", len(results))
	}
}

func TestRepository_Search_NoResults(t *testing.T) {
	repo, dir := newTestRepo(t)
	writeTestPlan(t, dir, "plan", "# Plan\nSome content.\n")

	results, err := repo.Search("nonexistent-query", false)
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("Search returned %d results, want 0", len(results))
	}
}

func TestSortByModTime(t *testing.T) {
	now := time.Now()
	plans := []Plan{
		{Name: "old", ModTime: now.Add(-2 * time.Hour)},
		{Name: "newest", ModTime: now},
		{Name: "middle", ModTime: now.Add(-1 * time.Hour)},
	}

	SortByModTime(plans)

	if plans[0].Name != "newest" || plans[1].Name != "middle" || plans[2].Name != "old" {
		t.Errorf("SortByModTime order = [%s, %s, %s], want [newest, middle, old]",
			plans[0].Name, plans[1].Name, plans[2].Name)
	}
}

func TestSortByName(t *testing.T) {
	plans := []Plan{
		{Name: "charlie"},
		{Name: "alpha"},
		{Name: "bravo"},
	}

	SortByName(plans)

	if plans[0].Name != "alpha" || plans[1].Name != "bravo" || plans[2].Name != "charlie" {
		t.Errorf("SortByName order = [%s, %s, %s], want [alpha, bravo, charlie]",
			plans[0].Name, plans[1].Name, plans[2].Name)
	}
}
