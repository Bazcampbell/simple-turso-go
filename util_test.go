package simpleturso

import "testing"

func TestValidationLogic(t *testing.T) {
	t.Run("isValidJWTStructure", func(t *testing.T) {
		tests := []struct {
			token    string
			expected bool
		}{
			{"header.payload.signature", true},
			{"eyHdr.eyPay.sig", true},
			{"invalid-token", false},
			{"part1.part2", false},
			{"part1.part2.part3.part4", false},
			{"..", false},
			{"", false},
		}

		for _, tc := range tests {
			if res := isValidJWTStructure(tc.token); res != tc.expected {
				t.Errorf("isValidJWTStructure(%q) = %v; want %v", tc.token, res, tc.expected)
			}
		}
	})

	t.Run("isValidTursoUrl", func(t *testing.T) {
		tests := []struct {
			url      string
			expected bool
		}{
			{"https://my-db.turso.io", true},
			{"libsql://my-db.turso.io", true},
			{"http://my-db.turso.io", false},
			{"libsql://my-db.com", false},
			{"not-a-url", false},
			{"https://turso.io", false},
			{"", false},
		}

		for _, tc := range tests {
			if res := isValidTursoUrl(tc.url); res != tc.expected {
				t.Errorf("isValidTursoUrl(%q) = %v; want %v", tc.url, res, tc.expected)
			}
		}
	})
}
