package content

import "testing"

func TestProcessMinifiesHTML(t *testing.T) {
	result, err := Process(`
		<p class="intro">
			Hello <strong>world</strong>.
		</p>
		<p>Next paragraph.</p>
	`)
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	const want = `<p>Hello <strong>world</strong>.<p>Next paragraph.`
	if result.Sanitized != want {
		t.Errorf("Process().Sanitized = %q, want %q", result.Sanitized, want)
	}
	if result.Words != 4 {
		t.Errorf("Process().Words = %d, want 4", result.Words)
	}
}

func TestProcessPreservesPreformattedWhitespace(t *testing.T) {
	result, err := Process("<pre>first\n  second</pre>")
	if err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	const want = "<pre>first\n  second</pre>"
	if result.Sanitized != want {
		t.Errorf("Process().Sanitized = %q, want %q", result.Sanitized, want)
	}
}
