package api

import "testing"

func TestAddSig(t *testing.T) {
	tests := []struct {
		q   *QueryOptions
		sig string
	}{
		{&QueryOptions{Key: "key", Secret: "secret", Expire: "1444155568"}, "f6c74dccc67c7aa52c64d2098d1569f8"},
		{&QueryOptions{Key: "x", Secret: "y", Expire: "1"}, "329e219571afb9926aa7a88a6c53cf87"},
		{&QueryOptions{Key: "eb0d3b6904107b76c14d0a307bf444c2", Expire: "1447085034", Format: "json", Page: "1", SessionID: "1447084436-trnesh"}, ""},
	}
	for _, tt := range tests {
		tt.q.AddSig()
		if tt.q.Sig != tt.sig {
			t.Errorf("Got sig %s, want %s", tt.q.Sig, tt.sig)
		}
	}
}
