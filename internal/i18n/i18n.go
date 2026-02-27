package i18n

// Translations holds all translatable strings.
type Translations struct {
	GitStatusClean string
	GitStatusDirty string
	RepoSingular   string
	RepoPlural     string
	CacheLabel     string
}

var locales = map[string]Translations{
	"en": {
		GitStatusClean: "Clean",
		GitStatusDirty: "Dirty",
		RepoSingular:   "repo",
		RepoPlural:     "repos",
		CacheLabel:     "Cache",
	},
	"fr": {
		GitStatusClean: "Propre",
		GitStatusDirty: "Modifie",
		RepoSingular:   "depot",
		RepoPlural:     "depots",
		CacheLabel:     "Cache",
	},
}

// Get returns translations for a locale. Falls back to English.
func Get(locale string) Translations {
	if t, ok := locales[locale]; ok {
		return t
	}

	return locales["en"]
}
