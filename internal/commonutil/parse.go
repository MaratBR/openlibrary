package commonutil

import "uuid"

func ParseStringArrayToUUID(ids []string) []uuid.UUID {
	uuids := []uuid.UUID{}
	for _, id := range ids {
		parsed, err := uuid.Parse(id)
		if err == nil && parsed != uuid.Nil() {
			uuids = append(uuids, parsed)
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
