package di

import "testing"

func TestSessionCookieSecure(t *testing.T) {
	tests := []struct {
		name        string
		env         string
		sessionName string
		want        bool
	}{
		{name: "local env", env: "local", sessionName: "blog_session", want: false},
		{name: "development env", env: "development", sessionName: "blog_session", want: false},
		{name: "test env", env: "test", sessionName: "blog_session", want: false},
		{name: "docker env", env: "docker", sessionName: "blog_session", want: false},
		{name: "production env", env: "production", sessionName: "blog_session", want: true},
		{name: "empty env defaults secure", env: "", sessionName: "blog_session", want: true},
		{name: "__Host prefix always secure", env: "local", sessionName: "__Host-blog_session", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sessionCookieSecure(tt.env, tt.sessionName); got != tt.want {
				t.Fatalf("sessionCookieSecure(%q, %q) = %v, want %v", tt.env, tt.sessionName, got, tt.want)
			}
		})
	}
}
