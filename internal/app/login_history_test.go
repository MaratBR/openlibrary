package app

import (
	"context"
	"testing"
	"time"

	"github.com/gofrs/uuid"
)

type loginHistorySessionStub struct{ sessions []SessionInfo }

func (s loginHistorySessionStub) GetBySID(context.Context, string) (*SessionInfo, error) {
	return nil, nil
}
func (s loginHistorySessionStub) Create(context.Context, CreateSessionCommand) (*SessionInfo, error) {
	return nil, nil
}
func (s loginHistorySessionStub) GetByUserID(context.Context, uuid.UUID) ([]SessionInfo, error) {
	return s.sessions, nil
}
func (s loginHistorySessionStub) TerminateBySID(context.Context, string) error { return nil }
func (s loginHistorySessionStub) TerminateAllByUserID(context.Context, uuid.UUID) error {
	return nil
}
func (s loginHistorySessionStub) Renew(context.Context, RenewSessionCommand) (*SessionInfo, error) {
	return nil, nil
}

func TestRecentLoginLocationsReturnsThreeNewestDistinctLocations(t *testing.T) {
	now := time.Now()
	berlin := IPLocation{Country: "Germany", Region: "Berlin", City: "Berlin"}
	tokyo := IPLocation{Country: "Japan", Region: "Tokyo", City: "Tokyo"}
	sydney := IPLocation{Country: "Australia", Region: "New South Wales", City: "Sydney"}
	toronto := IPLocation{Country: "Canada", Region: "Ontario", City: "Toronto"}
	svc := NewLoginHistoryService(moderationAuthStub{}, loginHistorySessionStub{sessions: []SessionInfo{
		{CreatedAt: now.Add(-4 * time.Hour), Location: toronto},
		{CreatedAt: now.Add(-1 * time.Hour), Location: berlin},
		{CreatedAt: now.Add(-2 * time.Hour), Location: tokyo},
		{CreatedAt: now, Location: berlin},
		{CreatedAt: now.Add(-3 * time.Hour), Location: sydney},
		{CreatedAt: now.Add(time.Hour), Location: IPLocation{}},
	}})

	locations, err := svc.GetRecentLoginLocations(context.Background(), GetLoginHistoryQuery{})
	if err != nil {
		t.Fatal(err)
	}
	if len(locations) != 3 || locations[0].IPLocation != berlin || locations[1].IPLocation != tokyo || locations[2].IPLocation != sydney {
		t.Fatalf("unexpected recent locations: %#v", locations)
	}
}
