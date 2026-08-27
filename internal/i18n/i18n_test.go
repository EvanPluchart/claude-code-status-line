package i18n

import (
	"reflect"
	"testing"
)

func TestGetFallsBackToEnglish(t *testing.T) {
	if Get("de").CacheLabel != Get("en").CacheLabel {
		t.Error("unknown locale should fall back to English")
	}

	if Get("fr").DurationDays != "j" {
		t.Error("French translations not returned")
	}
}

func TestAllLocalesAreComplete(t *testing.T) {
	for _, locale := range Locales() {
		tr := Get(locale)
		v := reflect.ValueOf(tr)

		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).String() == "" {
				t.Errorf("locale %q: field %s is empty", locale, v.Type().Field(i).Name)
			}
		}
	}
}
