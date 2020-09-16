package loghttpltsv

import (
	"fmt"
	"testing"
)

func TestLv(t *testing.T) {
	tests := []struct{
		key string
		value interface{}
		wanted string
	}{
		{"a", 3, "a:3"},
		{"b\n\nc", "d\n\ne", "b  c:d  e"},
		{"a\t\ta", "a\t\ta", "a  a:a  a"},
		{"a::a", "a::a", "a;;a:a::a"},
		{"", "", ":-"},
	}
	
	for i, v := range tests {
		t.Run(fmt.Sprint("case", i), func(t *testing.T) {
			got := lv(v.key, v.value)
			if got != v.wanted {
				t.Errorf("got %q, want %q", got, v.wanted)
			}
		})
	}
}

func TestLtsv(t *testing.T) {
	tests := []struct{
		lvs []string
		wanted string
	}{
		{[]string{
			lv("a", "1"), lv("b", 2.5), lv("c", nil),
		}, "a:1\tb:2.5\tc:<nil>"},
		{[]string{lv("a", 1)}, "a:1"},
	}

	for i, v := range tests {
		t.Run(fmt.Sprint("case", i), func(t *testing.T) {
			got := ltsv(v.lvs)
			if got != v.wanted {
				t.Errorf("got %q, want %q", got, v.wanted)
			}
		})
	}
}