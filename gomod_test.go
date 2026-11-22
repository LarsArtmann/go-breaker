package breaker

import (
	"os"
	"regexp"
	"testing"
)

// TestModuleGoDirectiveStaysPinned guards the module's Go directive against
// silent re-pinning by external tooling. The repo deliberately declares the
// relaxed `go 1.27` (never a patch pin), and the README installation floor
// states the same version, so a writer bumping go.mod without updating the
// README would otherwise drift the two apart unnoticed.
func TestModuleGoDirectiveStaysPinned(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile("go.mod")
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}

	directive := regexp.MustCompile(`(?m)^go (\S+)$`).FindStringSubmatch(string(data))
	if directive == nil {
		t.Fatal("go.mod has no `go` directive line")
	}

	const want = "1.27"
	if directive[1] != want {
		t.Fatalf(
			"go.mod declares `go %s`, want `go %s`: re-pin intentionally and update every doc that states it",
			directive[1],
			want,
		)
	}
}
