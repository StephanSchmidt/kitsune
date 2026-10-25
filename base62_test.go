package kitsune

import (
	"testing"

	"github.com/gofrs/uuid"
)

func TestBase62RoundTrip(t *testing.T) {
	for _, s := range []string{
		"0059c4be-dc64-43c2-acd3-1592f67b6426",
		"0000c4be-dc64-43c2-acd3-1592f67b6426",
		"00000000-0000-0000-0000-000000000001",
		"00000000-0000-0000-0000-000000000000",
		"ff59c4be-dc64-43c2-acd3-1592f67b6426",
		"ffffffff-ffff-ffff-ffff-ffffffffffff",
	} {
		id := uuid.Must(uuid.FromString(s))
		got, err := FromBase62(ToBase62(id))
		if err != nil || got != id {
			t.Errorf("round trip of %s gave %s, %v", s, got, err)
		}
	}
	for i := 0; i < 2000; i++ {
		id := uuid.Must(uuid.NewV4())
		if got, err := FromBase62(ToBase62(id)); err != nil || got != id {
			t.Fatalf("round trip of %s gave %s, %v", id, got, err)
		}
	}
}

func TestFromBase62RejectsWhatIsNoUUID(t *testing.T) {
	for _, s := range []string{"", "not base62!", "-1", "zzzzzzzzzzzzzzzzzzzzzzzzzz"} {
		if _, err := FromBase62(s); err == nil {
			t.Errorf("FromBase62(%q) accepted", s)
		}
	}
}
