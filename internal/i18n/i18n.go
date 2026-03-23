package i18n

// Translations holds all translatable strings.
type Translations struct {
	GitStatusClean string
	GitStatusDirty string
	RepoSingular   string
	RepoPlural     string
	CacheLabel     string
	// Placeholders when info is unavailable
	NoGitRepo      string
	NoNestedRepos  string
	NoLinesChanged string
	// Rate limit labels
	SessionLabel string
	WeeklyLabel  string
	ResetsIn     string
	// Duration units
	DurationMonths  string
	DurationWeeks   string
	DurationDays    string
	DurationHours   string
	DurationMinutes string
	DurationSeconds string
}

var locales = map[string]Translations{
	"en": {
		GitStatusClean:  "Clean",
		GitStatusDirty:  "Dirty",
		RepoSingular:    "repo",
		RepoPlural:      "repos",
		CacheLabel:      "Cache",
		NoGitRepo:       "no git repo",
		NoNestedRepos:   "0 nested repos",
		NoLinesChanged:  "+0 -0",
		SessionLabel:    "5h",
		WeeklyLabel:     "7d",
		ResetsIn:        "reset",
		DurationMonths:  "mo",
		DurationWeeks:   "w",
		DurationDays:    "d",
		DurationHours:   "h",
		DurationMinutes: "m",
		DurationSeconds: "s",
	},
	"fr": {
		GitStatusClean:  "Propre",
		GitStatusDirty:  "Modifie",
		RepoSingular:    "depot",
		RepoPlural:      "depots",
		CacheLabel:      "Cache",
		NoGitRepo:       "pas de repo git",
		NoNestedRepos:   "0 depots imbriques",
		NoLinesChanged:  "+0 -0",
		SessionLabel:    "5h",
		WeeklyLabel:     "7j",
		ResetsIn:        "reset",
		DurationMonths:  "mo",
		DurationWeeks:   "sem",
		DurationDays:    "j",
		DurationHours:   "h",
		DurationMinutes: "min",
		DurationSeconds: "s",
	},
}

// Get returns translations for a locale. Falls back to English.
func Get(locale string) Translations {
	if t, ok := locales[locale]; ok {
		return t
	}

	return locales["en"]
}
