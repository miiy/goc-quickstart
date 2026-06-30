package media

import "testing"

func TestUploadsURL(t *testing.T) {
	if got, want := UploadsURL("post-covers/2026/06/a.png"), "/uploads/post-covers/2026/06/a.png"; got != want {
		t.Fatalf("UploadsURL() = %q, want %q", got, want)
	}
}
