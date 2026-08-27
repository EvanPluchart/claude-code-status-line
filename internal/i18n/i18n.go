package i18n

// Translations holds all translatable strings.
type Translations struct {
	GitStatusClean string
	GitStatusDirty string
	RepoSingular   string
	RepoPlural     string
	CacheLabel     string
	RemainingLabel string
	// Placeholders when info is unavailable
	NoGitRepo      string
	NoNestedRepos  string
	NoLinesChanged string
	// Rate limit labels
	SessionLabel string
	WeeklyLabel  string
	ResetsIn     string
	// Session metadata
	APILabel      string
	PerHourSuffix string
	ThinkingLabel string
	FastLabel     string
	// Pull request review states
	PRApproved         string
	PRPending          string
	PRChangesRequested string
	PRDraft            string
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
		GitStatusClean:     "Clean",
		GitStatusDirty:     "Dirty",
		RepoSingular:       "repo",
		RepoPlural:         "repos",
		CacheLabel:         "Cache",
		RemainingLabel:     "left",
		NoGitRepo:          "no git repo",
		NoNestedRepos:      "0 nested repos",
		NoLinesChanged:     "+0 -0",
		SessionLabel:       "5h",
		WeeklyLabel:        "7d",
		ResetsIn:           "reset",
		APILabel:           "API",
		PerHourSuffix:      "/h",
		ThinkingLabel:      "think",
		FastLabel:          "fast",
		PRApproved:         "approved",
		PRPending:          "pending",
		PRChangesRequested: "changes",
		PRDraft:            "draft",
		DurationMonths:     "mo",
		DurationWeeks:      "w",
		DurationDays:       "d",
		DurationHours:      "h",
		DurationMinutes:    "m",
		DurationSeconds:    "s",
	},
	"fr": {
		GitStatusClean:     "Propre",
		GitStatusDirty:     "Modifié",
		RepoSingular:       "dépôt",
		RepoPlural:         "dépôts",
		CacheLabel:         "Cache",
		RemainingLabel:     "restants",
		NoGitRepo:          "pas de dépôt git",
		NoNestedRepos:      "0 dépôt imbriqué",
		NoLinesChanged:     "+0 -0",
		SessionLabel:       "5h",
		WeeklyLabel:        "7j",
		ResetsIn:           "reset",
		APILabel:           "API",
		PerHourSuffix:      "/h",
		ThinkingLabel:      "réflexion",
		FastLabel:          "rapide",
		PRApproved:         "approuvée",
		PRPending:          "en attente",
		PRChangesRequested: "modifs",
		PRDraft:            "brouillon",
		DurationMonths:     "mo",
		DurationWeeks:      "sem",
		DurationDays:       "j",
		DurationHours:      "h",
		DurationMinutes:    "min",
		DurationSeconds:    "s",
	},
}

// Get returns translations for a locale. Falls back to English.
func Get(locale string) Translations {
	if t, ok := locales[locale]; ok {
		return t
	}

	return locales["en"]
}

// Locales returns the supported locale codes.
func Locales() []string {
	return []string{"en", "fr"}
}
