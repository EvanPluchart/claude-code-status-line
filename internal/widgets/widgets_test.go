package widgets

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/EvanPluchart/claude-code-status-line/internal/config"
	"github.com/EvanPluchart/claude-code-status-line/internal/i18n"
	"github.com/EvanPluchart/claude-code-status-line/internal/parser"
	"github.com/EvanPluchart/claude-code-status-line/internal/themes"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(s string) string {
	return ansiPattern.ReplaceAllString(s, "")
}

// testContext builds a render context with a realistic session.
func testContext() *Context {
	cfg := config.Default()
	cfg.Widgets.Cost.Currency = "USD"

	input := &parser.Input{
		CWD:     "/home/user/projects/app",
		Version: "2.1.90",
		Model:   parser.Model{ID: "claude-opus-5", DisplayName: "Opus"},
		Workspace: parser.Workspace{
			ProjectDir: "/home/user/projects/app",
			Repo:       &parser.Repo{Owner: "acme", Name: "app"},
		},
		Cost: parser.Cost{
			TotalCostUSD:       1.5,
			TotalDurationMS:    3_600_000,
			TotalAPIDurationMS: 900_000,
			TotalLinesAdded:    42,
			TotalLinesRemoved:  7,
		},
		ContextWindow: parser.ContextWindow{
			ContextWindowSize: 200000,
			UsedPercentage:    45,
			TotalInputTokens:  75000,
			TotalOutputTokens: 15000,
			CurrentUsage: &parser.CurrentUsage{
				InputTokens:              60000,
				CacheCreationInputTokens: 5000,
				CacheReadInputTokens:     25000,
			},
		},
	}

	return &Context{Input: input, Config: cfg, Theme: themes.Get("default")}
}

func TestRegistryContainsAllWidgets(t *testing.T) {
	expected := []string{
		"model", "effort", "thinking", "fast-mode", "agent", "output-style", "session-name", "claude-version", "vim-mode",
		"directory", "repo", "worktree", "git-branch", "git-status", "git-changes", "git-ahead-behind", "pr", "nested-repos", "lines-changed",
		"cost", "burn-rate", "duration", "api-time", "timestamp",
		"token-bar", "context-percent", "token-count", "context-remaining", "total-tokens", "cache-ratio", "exceeds-200k",
		"session-usage", "weekly-usage", "os-info", "hostname", "separator", "spacer",
	}

	for _, id := range expected {
		if Get(id) == nil {
			t.Errorf("widget %q is not registered", id)
		}
	}

	if len(IDs()) != len(expected) {
		t.Errorf("registry has %d widgets, test expects %d: %v", len(IDs()), len(expected), IDs())
	}
}

func TestFormatTokens(t *testing.T) {
	cases := map[int]string{
		0: "0", 999: "999", 1000: "1k", 1500: "1.5k", 45200: "45.2k",
		200000: "200k", 1_000_000: "1M", 1_250_000: "1.2M",
	}

	for in, want := range cases {
		if got := formatTokens(in); got != want {
			t.Errorf("formatTokens(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestTokenWidgets(t *testing.T) {
	ctx := testContext()

	if got := stripANSI((&TokenCountWidget{}).Render(ctx)); got != "(90k/200k)" {
		t.Errorf("token-count = %q", got)
	}

	if got := stripANSI((&ContextPercentWidget{}).Render(ctx)); got != "45%" {
		t.Errorf("context-percent = %q", got)
	}

	if got := stripANSI((&ContextRemainingWidget{}).Render(ctx)); got != "110k left" {
		t.Errorf("context-remaining = %q", got)
	}

	if got := stripANSI((&TotalTokensWidget{}).Render(ctx)); got != "↑75k ↓15k" {
		t.Errorf("total-tokens = %q", got)
	}

	if got := stripANSI((&CacheRatioWidget{}).Render(ctx)); got != "Cache: 28%" {
		t.Errorf("cache-ratio = %q", got)
	}

	bar := stripANSI((&TokenBarWidget{}).Render(ctx))

	if strings.Count(bar, "━") != 7 || strings.Count(bar, "─") != 9 {
		t.Errorf("token-bar = %q, want 7 filled / 9 empty", bar)
	}

	if got := (&Exceeds200KWidget{}).Render(ctx); got != "" {
		t.Errorf("exceeds-200k should be empty, got %q", got)
	}

	ctx.Input.Exceeds200K = true

	if got := stripANSI((&Exceeds200KWidget{}).Render(ctx)); got != ">200k" {
		t.Errorf("exceeds-200k = %q", got)
	}
}

func TestThresholdColor(t *testing.T) {
	ctx := testContext()
	th := config.ThresholdGroup{Yellow: 60, Orange: 80, Red: 90}

	if got := thresholdColor(10, th, ctx); got != ctx.Theme.Success {
		t.Errorf("10 should be success")
	}

	if got := thresholdColor(60, th, ctx); got != ctx.Theme.Warning {
		t.Errorf("60 should be warning")
	}

	if got := thresholdColor(85, th, ctx); got != ctx.Theme.Danger {
		t.Errorf("85 should be danger")
	}

	if got := thresholdColor(95, th, ctx); !strings.Contains(got, ctx.Theme.Danger) || !strings.Contains(got, "\x1b[1m") {
		t.Errorf("95 should be bold danger")
	}
}

func TestCostWidget(t *testing.T) {
	t.Setenv("CLAUDE_STATUSLINE_DIR", t.TempDir())
	t.Setenv("CLAUDE_STATUSLINE_OFFLINE", "1")

	ctx := testContext()

	if got := stripANSI((&CostWidget{}).Render(ctx)); got != "$1.50" {
		t.Errorf("cost USD = %q", got)
	}

	ctx.Config.Widgets.Cost.Currency = "EUR"
	got := stripANSI((&CostWidget{}).Render(ctx))

	if !strings.HasSuffix(got, "€") || !strings.Contains(got, ",") {
		t.Errorf("cost EUR should use comma and trailing symbol, got %q", got)
	}

	ctx.Config.Widgets.Cost.Currency = "GBP"
	ctx.Config.Widgets.Cost.Decimals = 3

	if got := stripANSI((&CostWidget{}).Render(ctx)); !strings.HasPrefix(got, "£") || len(strings.SplitN(got, ".", 2)[1]) != 3 {
		t.Errorf("cost GBP with 3 decimals = %q", got)
	}

	ctx.Config.Widgets.Cost.Currency = "XXX"

	if got := stripANSI((&CostWidget{}).Render(ctx)); got != "XXX1.500" {
		t.Errorf("unknown currency should fall back to code and rate 1, got %q", got)
	}
}

func TestBurnRateWidget(t *testing.T) {
	t.Setenv("CLAUDE_STATUSLINE_DIR", t.TempDir())
	t.Setenv("CLAUDE_STATUSLINE_OFFLINE", "1")

	ctx := testContext()

	if got := stripANSI((&BurnRateWidget{}).Render(ctx)); got != "$1.50/h" {
		t.Errorf("burn-rate = %q", got)
	}

	ctx.Input.Cost.TotalDurationMS = 30_000

	if got := (&BurnRateWidget{}).Render(ctx); got != "" {
		t.Errorf("burn-rate should be empty during the first minute, got %q", got)
	}
}

func TestDurationFormatting(t *testing.T) {
	en := i18n.Get("en")
	fr := i18n.Get("fr")

	cases := []struct {
		d    time.Duration
		t    i18n.Translations
		want string
	}{
		{45 * time.Second, en, "45s"},
		{12*time.Minute + 34*time.Second, en, "12m34s"},
		{time.Hour + 5*time.Minute, en, "1h 05m"},
		{time.Hour + 5*time.Minute, fr, "1h 05min"},
		{26 * time.Hour, en, "1d 2h 00m"},
		{8 * 24 * time.Hour, en, "1w 1d 0h 00m"},
		{-5 * time.Second, en, "0s"},
	}

	for _, c := range cases {
		if got := formatSessionDuration(c.d, c.t); got != c.want {
			t.Errorf("formatSessionDuration(%v) = %q, want %q", c.d, got, c.want)
		}
	}

	compact := []struct {
		d    time.Duration
		want string
	}{
		{12 * time.Minute, "12m"},
		{3*time.Hour + 42*time.Minute, "3h42m"},
		{4*24*time.Hour + 8*time.Hour + 5*time.Minute, "4d8h05m"},
	}

	for _, c := range compact {
		if got := formatCompactDuration(c.d, en); got != c.want {
			t.Errorf("formatCompactDuration(%v) = %q, want %q", c.d, got, c.want)
		}
	}
}

func TestDurationAndAPITimeWidgets(t *testing.T) {
	ctx := testContext()

	if got := stripANSI((&DurationWidget{}).Render(ctx)); got != "1h 00m" {
		t.Errorf("duration = %q", got)
	}

	if got := stripANSI((&APITimeWidget{}).Render(ctx)); got != "API 15m00s (25%)" {
		t.Errorf("api-time = %q", got)
	}

	ctx.Input.Cost.TotalAPIDurationMS = 0

	if got := (&APITimeWidget{}).Render(ctx); got != "" {
		t.Errorf("api-time should be empty without API time, got %q", got)
	}
}

func TestRateLimitWidgets(t *testing.T) {
	ctx := testContext()

	if got := (&SessionUsageWidget{}).Render(ctx); got != "" {
		t.Errorf("session-usage should be empty without rate limits, got %q", got)
	}

	ctx.Input.RateLimits = &parser.RateLimits{
		FiveHour: &parser.RateLimit{UsedPercentage: 23.5, ResetsAt: time.Now().Add(2*time.Hour + 10*time.Minute).Unix()},
		SevenDay: &parser.RateLimit{UsedPercentage: 91, ResetsAt: 0},
	}

	session := stripANSI((&SessionUsageWidget{}).Render(ctx))

	if !strings.HasPrefix(session, "5h ") || !strings.Contains(session, " 24% ") || !strings.Contains(session, "(reset 2h") {
		t.Errorf("session-usage = %q", session)
	}

	weekly := stripANSI((&WeeklyUsageWidget{}).Render(ctx))

	if !strings.HasPrefix(weekly, "7d ") || !strings.HasSuffix(weekly, " 91%") {
		t.Errorf("weekly-usage = %q", weekly)
	}
}

func TestSessionWidgets(t *testing.T) {
	ctx := testContext()

	empties := []Widget{
		&EffortWidget{}, &ThinkingWidget{}, &FastModeWidget{}, &AgentWidget{}, &OutputStyleWidget{},
		&SessionNameWidget{}, &VimModeWidget{}, &PRWidget{}, &WorktreeWidget{},
	}

	for _, w := range empties {
		if got := w.Render(ctx); got != "" {
			t.Errorf("%s should be empty when data is absent, got %q", w.ID(), got)
		}
	}

	ctx.Input.Effort = &parser.Effort{Level: "xhigh"}
	ctx.Input.Thinking = &parser.Thinking{Enabled: true}
	ctx.Input.FastMode = true
	ctx.Input.Agent = &parser.Agent{Name: "reviewer"}
	ctx.Input.OutputStyle = &parser.OutputStyle{Name: "Explanatory"}
	ctx.Input.SessionName = "Fix login bug"
	ctx.Input.Vim = &parser.Vim{Mode: "INSERT"}
	ctx.Input.PR = &parser.PullRequest{Number: 42, ReviewState: "changes_requested"}
	ctx.Input.Worktree = &parser.Worktree{Name: "feat-x"}

	cases := map[Widget]string{
		&EffortWidget{}:        "⚙ xhigh",
		&ThinkingWidget{}:      "✦ think",
		&FastModeWidget{}:      "⚡ fast",
		&AgentWidget{}:         "@reviewer",
		&OutputStyleWidget{}:   "Explanatory",
		&SessionNameWidget{}:   "Fix login bug",
		&VimModeWidget{}:       "INSERT",
		&PRWidget{}:            "PR #42 changes",
		&WorktreeWidget{}:      "⌥ feat-x",
		&RepoWidget{}:          "acme/app",
		&ClaudeVersionWidget{}: "v2.1.90",
	}

	for w, want := range cases {
		if got := stripANSI(w.Render(ctx)); got != want {
			t.Errorf("%s = %q, want %q", w.ID(), got, want)
		}
	}

	ctx.Input.OutputStyle.Name = "default"

	if got := (&OutputStyleWidget{}).Render(ctx); got != "" {
		t.Errorf("output-style should hide the default style, got %q", got)
	}

	ctx.Input.PR = &parser.PullRequest{Number: 7, ReviewState: "approved", Kind: "mr"}

	if got := stripANSI((&PRWidget{}).Render(ctx)); got != "MR #7 ✓" {
		t.Errorf("pr (gitlab) = %q", got)
	}
}

func TestLinesChangedSessionSource(t *testing.T) {
	ctx := testContext()
	ctx.Config.Widgets.LinesChanged.Source = "session"

	if got := stripANSI((&LinesChangedWidget{}).Render(ctx)); got != "+42 -7" {
		t.Errorf("lines-changed (session) = %q", got)
	}
}

func TestParseNumstat(t *testing.T) {
	out := "10\t2\tmain.go\n-\t-\timage.png\n3\t0\tREADME.md\n"

	added, removed := parseNumstat(out)

	if added != 13 || removed != 2 {
		t.Errorf("parseNumstat = +%d -%d, want +13 -2", added, removed)
	}
}

func TestCountPorcelain(t *testing.T) {
	status := "M  staged.go\n M modified.go\nMM both.go\n?? new.txt\nA  added.go\n"

	staged, modified, untracked := countPorcelain(status)

	if staged != 3 || modified != 2 || untracked != 1 {
		t.Errorf("countPorcelain = %d/%d/%d, want 3/2/1", staged, modified, untracked)
	}
}

func TestTruncateBranch(t *testing.T) {
	if got := truncateBranch("main", 40); got != "main" {
		t.Errorf("short branch should be untouched, got %q", got)
	}

	got := truncateBranch("feature/very-long-branch-name-with-ticket-ABC-1234", 20)

	if len([]rune(got)) != 20 || !strings.Contains(got, "…") {
		t.Errorf("truncated branch = %q", got)
	}
}

func TestShortenPath(t *testing.T) {
	cases := []struct {
		dir   string
		depth int
		want  string
	}{
		{"/home/user/projects/app", 0, "app"},
		{"/home/user/projects/app", 1, "app"},
		{"/home/user/projects/app", 2, "projects/app"},
		{"/home/user/projects/app", 10, "home/user/projects/app"},
		{"/", 1, "~"},
		{"", 1, "~"},
	}

	for _, c := range cases {
		if got := shortenPath(c.dir, c.depth); got != c.want {
			t.Errorf("shortenPath(%q, %d) = %q, want %q", c.dir, c.depth, got, c.want)
		}
	}
}

func TestFindNestedRepos(t *testing.T) {
	root := t.TempDir()

	mustMkdir(t, root, ".git")
	mustMkdir(t, root, "libs/one/.git")
	mustMkdir(t, root, "libs/two/.git")
	mustMkdir(t, root, "node_modules/dep/.git")
	mustMkdir(t, root, "a/b/c/d/.git") // too deep

	repos := findNestedRepos(root)

	if len(repos) != 2 {
		t.Errorf("found %d nested repos, want 2: %v", len(repos), repos)
	}
}

func TestSeparatorAndSpacer(t *testing.T) {
	ctx := testContext()

	if got := stripANSI((&SeparatorWidget{}).Render(ctx)); got != " │ " {
		t.Errorf("separator = %q", got)
	}

	ctx.Config.Widgets.Separator.Char = "•"

	if got := stripANSI((&SeparatorWidget{}).Render(ctx)); got != " • " {
		t.Errorf("custom separator = %q", got)
	}

	if got := (&SpacerWidget{}).Render(ctx); got != " " {
		t.Errorf("spacer = %q", got)
	}
}
