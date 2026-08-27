package dashboard

import "testing"

func TestBuildOverview(t *testing.T) {
	view := BuildOverview("local", nil)
	if view.ServerName != "local" || view.Usage != nil {
		t.Fatalf("unexpected overview: %#v", view)
	}
}
