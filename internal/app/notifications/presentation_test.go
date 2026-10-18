package notifications

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"uuid"

	"github.com/MaratBR/openlibrary/internal/store"
	"github.com/stretchr/testify/require"
)

func TestSanitizeNotificationContent(t *testing.T) {
	input := `<b>bold</b><strong>strong</strong><i>italic</i><em>em</em><u>u</u><s>s</s><del>d</del><sub>2</sub><sup>3</sup><br><code>x</code><a href="/book/test" title="Read">link</a>`
	safe := SanitizeContent(input)
	for _, part := range []string{"<b>bold</b>", "<i>italic</i>", "<sub>2</sub>", "<sup>3</sup>", "<br>", "<code>x</code>", `href="/book/test"`} {
		require.Contains(t, safe, part)
	}
	for _, bad := range []string{
		`<script>alert(1)</script>`, `<img src=x onerror="alert(1)">`,
		`<a href="javascript:alert(1)" onclick="alert(1)">link</a>`,
		`<a href="data:text/html,alert(1)">link</a>`,
		`<b style="color:red" onmouseover="alert(1)">safe</b>`,
		`<iframe src="https://example.com"></iframe>`,
	} {
		safe := SanitizeContent(bad)
		for _, forbidden := range []string{"<script", "<img", "javascript:", "onclick", "data:", "style=", "onmouseover", "<iframe"} {
			require.NotContains(t, safe, forbidden)
		}
	}
	for _, link := range []string{"https://example.com", "http://example.com", "mailto:reader@example.com"} {
		require.Contains(t, SanitizeContent(`<a href="`+link+`">link</a>`), `href="`+link+`"`)
	}
}

func TestNotificationPresentationDefaultsAndOverrides(t *testing.T) {
	for _, tc := range []struct{ category, icon string }{
		{"", ""}, {CategorySystem, ""}, {CategoryModeration, "fa-solid fa-shield"}, {CategoryAdmin, "fa-solid fa-gauge-high"},
	} {
		n, err := NewBuilder("Title", uuid.NewV7()).Category(tc.category).Urgent(true).Type("reply / suggestions 🔖").Build()
		require.NoError(t, err)
		category, urgent, icon := presentation(n.Data)
		require.True(t, urgent)
		require.Equal(t, category, n.Category())
		require.Equal(t, urgent, n.Urgent())
		require.Equal(t, icon, n.Icon())
		require.NotEmpty(t, category)
		require.Equal(t, "reply / suggestions 🔖", n.Type())
		if tc.icon == "" {
			require.Nil(t, icon)
		} else {
			require.Equal(t, tc.icon, icon.Payload)
		}
	}
	n, err := NewBuilder("Title", uuid.NewV7()).Category(CategoryAdmin).IconEmoji("📚").Build()
	require.NoError(t, err)
	_, _, icon := presentation(n.Data)
	require.Equal(t, &Icon{Type: IconEmoji, Payload: "📚"}, icon)
	for _, builder := range []*builder{
		NewBuilder("x", uuid.NewV7()).Category("unknown"),
		NewBuilder("x", uuid.NewV7()).Type(strings.Repeat("x", 101)),
		NewBuilder("x", uuid.NewV7()).IconImg("javascript:alert(1)"),
		NewBuilder("x", uuid.NewV7()).IconImg("//evil.example/icon"),
		NewBuilder("x", uuid.NewV7()).IconImg("https:missing-host"),
		NewBuilder("x", uuid.NewV7()).IconImg(`/\evil.example/icon`),
		NewBuilder("x", uuid.NewV7()).IconFontAwesome("fa-solid fa-bell arbitrary-class"),
	} {
		_, err := builder.Build()
		require.Error(t, err)
	}
}

func TestSendAndBuilderPersistPresentation(t *testing.T) {
	db := &recordingDB{}
	id := uuid.NewV7()
	n, err := NewBuilder("Title", id).Content(`<b>ok</b><script>bad</script>`).
		Category(CategoryModeration).Type("suggestions").Urgent(true).IconEmoji("📚").Build()
	require.NoError(t, err)
	require.NoError(t, NewSQLStore(db).Add(context.Background(), n))
	record := eventFromRow(store.Notification{Metadata: db.args[6].([]byte), Content: n.Content, Type: n.Type()})
	require.Equal(t, CategoryModeration, record.Category)
	require.True(t, record.Urgent)
	require.Equal(t, "📚", record.Icon.Payload)
	require.Equal(t, "suggestions", record.Type)
	require.Equal(t, "<b>ok</b>", record.Content)
	_, err = NewStreamReader(db).Send(context.Background(), id.String(), "", "Title", `<b>ok</b><script>bad</script>`, Options{
		Category: CategoryAdmin, Type: "reply", Urgent: true, Icon: &Icon{Type: IconEmoji, Payload: "📌"},
	})
	require.NoError(t, err)
	require.Equal(t, "reply", db.args[1])
	require.Equal(t, "<b>ok</b>", db.args[3])
	var data map[string]string
	require.NoError(t, json.Unmarshal(db.args[4].([]byte), &data))
	category, urgent, icon := presentation(data)
	require.Equal(t, CategoryAdmin, category)
	require.True(t, urgent)
	require.Equal(t, "📌", icon.Payload)
}

func TestLegacyNotificationsAreSanitizedAndInvalidIconsFallBack(t *testing.T) {
	event := eventFromRow(store.Notification{Content: `<b>hello</b><img src=x onerror="alert(1)">`, Metadata: []byte(`{"category":"moderation","urgent":"true","icon":"{\"type\":0,\"payload\":\"javascript:bad\"}","old":123}`)})
	require.Equal(t, "<b>hello</b>", event.Content)
	require.Equal(t, CategoryModeration, event.Category)
	require.True(t, event.Urgent)
	require.Equal(t, "fa-solid fa-shield", event.Icon.Payload)
	require.Nil(t, eventFromRow(store.Notification{}).Icon)
}
