package commonutil

import "github.com/gofrs/uuid"

func ParseStringArrayToUUID(ids []string) []uuid.UUID {
	uuids := []uuid.UUID{}
	for _, id := range ids {
		id := uuid.FromStringOrNil(id)
		if id != uuid.Nil {
			uuids = append(uuids, id)
		}
	}
	return uuids
}

func StringifyUUIDArray(uuids []uuid.UUID) []string {
	ids := make([]string, len(uuids))
	for i, id := range uuids {
		ids[i] = id.String()
	}
	return ids
}
