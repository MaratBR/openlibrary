import { Schema } from 'effect'
import { atomFamily } from 'jotai-family'
import {
  atomWithStorage,
  createJSONStorage,
  unstable_withStorageValidator as withStorageValidator,
} from 'jotai/utils'

const recentItemSchema = Schema.Struct({
  to: Schema.NonEmptyString,
  label: Schema.NonEmptyString,
})
type RecentItem = typeof recentItemSchema.Type

export const recentItemsAtom = atomFamily((key: string) => {
  const section = key.slice(key.indexOf(':') + 1)
  const schema = Schema.mutable(Schema.Array(recentItemSchema)).check(
    Schema.isMaxLength(2),
    Schema.makeFilter((items) =>
      items.every(
        (item) =>
          item.to.startsWith(`${section}/`) &&
          item.to.slice(section.length + 1).length > 0 &&
          !item.to.slice(section.length + 1).includes('/'),
      ),
    ),
  )
  const storage = withStorageValidator(Schema.is(schema))(createJSONStorage<unknown>())
  return atomWithStorage<RecentItem[]>(`openlibrary:dashboard-recent:${key}`, [], storage, {
    getOnInit: true,
  })
})

export function rememberDashboardItem(items: RecentItem[], item: RecentItem): RecentItem[] {
  if (items[0]?.to === item.to && items[0].label === item.label) return items
  return [item, ...items.filter((recent) => recent.to !== item.to)].slice(0, 2)
}
