package themes

import "testing"

func TestGetFallsBackToDefault(t *testing.T) {
	if Get("unknown").Name != "default" {
		t.Error("unknown theme should fall back to default")
	}

	if Get("nord").Name != "nord" {
		t.Error("nord theme should be returned")
	}
}

func TestAllThemesHaveColors(t *testing.T) {
	for _, name := range Names() {
		th := Get(name)

		if th.Primary == "" || th.Success == "" || th.Warning == "" || th.Danger == "" || th.Muted == "" || th.Separator == "" || th.BarEmpty == "" {
			t.Errorf("theme %q has missing colors: %+v", name, th)
		}
	}
}
