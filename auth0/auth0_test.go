package auth0

import (
	"testing"
)

var pkceChallenge = []struct {
	v    string
	want string
}{
	{"lWuwhEs5s3Gjd0HBf0VniNhkiNl8VIRQWcfw5GfIc1I", "bsDef56Rqb5ztE3AK4PdKHNFsjU5G1Ix4GTscTOwgTA"},
	{"fb_JhIX6T6KmjsMzTN6Nxe0wNStGkldt0CAnFcJ9egg", "LXL4nkynVRwZCah_7sltoMMmoGFM65lAcBmGZjVDSRA"},
	{"JpQo4q_xf-j9ZSpRtYrEahfHrEeVJQGS9U1ng8jXcCs", "BT_mFbq0nUoQG1WMd9BAsP7PzanQIZ09w4XIPilSnA4"},
}

// TestPKCEChallenge tests CreatePKCEChallenge creation.
func TestPKCEChallenge256(t *testing.T) {
	for _, tt := range pkceChallenge {
		got := CreatePKCEChallengeS256(tt.v)
		if got != tt.want {
			t.Errorf("auth0.CreatePKCEChallengeS256(\"%v\") Mismatch: want %v, got %v", tt.v, tt.want, got)
		}
	}
}
