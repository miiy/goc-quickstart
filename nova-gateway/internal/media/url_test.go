package media

import "testing"

func TestUploadsURL(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "object key", in: "post-covers/2026/06/a.png", want: "/uploads/post-covers/2026/06/a.png"},
		{name: "leading slash", in: "/avatars/a.png", want: "/uploads/avatars/a.png"},
		{name: "already uploads URL", in: "/uploads/avatars/a.png", want: "/uploads/avatars/a.png"},
		{name: "absolute URL", in: "https://cdn.test/avatars/a.png", want: "https://cdn.test/avatars/a.png"},
		{name: "empty", in: " ", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UploadsURL(tt.in); got != tt.want {
				t.Fatalf("UploadsURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestFileURLFallsBackToObjectKey(t *testing.T) {
	if got, want := FileURL("", "avatars/a.png"), "/uploads/avatars/a.png"; got != want {
		t.Fatalf("FileURL() = %q, want %q", got, want)
	}
	if got, want := FileURL("https://cdn.test/a.png", "avatars/a.png"), "https://cdn.test/a.png"; got != want {
		t.Fatalf("FileURL() = %q, want %q", got, want)
	}
}
