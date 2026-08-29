import { DefinedTagDto, useTagsSearch } from '@/features/search'
import Modal from '@/components/Modal'
import { ReactIslandProps } from '@/lib/island'
import { useMemo, useState } from 'react'
import { Schema } from 'effect'

const dataSchema = Schema.Struct({
  open: Schema.Boolean,
})

export function SimpleTagInputModal({ data, rootElement }: ReactIslandProps) {
  const { open } = useMemo(() => Schema.decodeUnknownSync(dataSchema)(data), [data])

  const [searchValue, setSearchValue] = useState('')

  const query = useTagsSearch({
    query: searchValue,
    fetchDefault: true,
    enabled: open,
  })

  const tags = query.data || []

  return (
    <Modal open={open} onClose={handleClose}>
      <div className="admin-card">
        <div className="p-2">
          <input
            className="input text-xl h-12"
            value={searchValue}
            onInput={(e) => setSearchValue((e.target as HTMLInputElement).value)}
          />
        </div>

        <ul className="space-y-1 h-96 overflow-auto p-4">
          {tags.map((tag) => (
            <li
              key={tag.id}
              onClick={() => handleSelect(tag)}
              className="p-2 cursor-pointer hover:bg-foreground/5 hover:text-primary"
            >
              {tag.name}
            </li>
          ))}
        </ul>
      </div>
    </Modal>
  )

  function handleClose() {
    rootElement.dispatchEvent(new CustomEvent('close'))
  }

  function handleSelect(tag: DefinedTagDto) {
    rootElement.dispatchEvent(new CustomEvent('selected', { detail: tag }))
  }
}
