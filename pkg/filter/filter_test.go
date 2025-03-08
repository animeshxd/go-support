package filter

import (
	"reflect"
	"testing"

	"github.com/gobwas/glob"
)

func TestFilterType(t *testing.T) {
	t.Run("", func(t *testing.T) {
		//[],

		user := struct {
			Disabled bool
			Verified bool
			Name     string
			Role     []string
		}{
			Disabled: false,
			Verified: false,
			Name:     "hello",
			Role:     []string{"user", "host", "unknown"},
		}
		pattern := CreateGlobPattern([]string{"!*", "Name", "Role"})
		matcher := glob.MustCompile(pattern, '.')
		result, err := FilterType(user, make(map[uintptr]bool), "", matcher)
		if err != nil {
			t.Error(err)
		}
		if result == nil {
			t.Error("result is nil")
		}

		if reflect.DeepEqual(user, map[string]any{
			"Name": "hello",
			"Role": []string{"user", "host", "unknown"},
		}) {
			t.Errorf("expected %v, got %v", user, result)
		}

	})
}
