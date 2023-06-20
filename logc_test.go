package main

import "testing"

func TestGetAttributes(t *testing.T) {
	for _, tc := range []struct {
		in   string
		want string
	}{
		{`{"msg":"foo"}`, ""},
		{`{"msg":"foo","attr1":"val\"with quote"}`, `: "attr1":"val\"with quote"`},
		{`{"msg":"foo\\"}`, ``},
		{`{"msg":"foo\\","a1":"v1"}`, `: "a1":"v1"`},
		{`{"msg":"hello \"world\"","a1":"v\n1,"a2":"v2"}`, `: "a1":"v\n1,"a2":"v2"`},
	} {
		got := GetAttributes(tc.in)
		if got != tc.want {
			t.Errorf("GetAttributes(%s) = '%s', want '%s'", tc.in, got, tc.want)
		}
	}
}
