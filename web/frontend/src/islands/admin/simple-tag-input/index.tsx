import { DefinedTagDto, useTagsSearch } from '@/api/search'
import Modal from '@/components/Modal'
import { ReactIslandProps } from '@/lib/island'
import { useMemo, useState } from 'react'
import { z } from 'zod'

const dataSchema = z.object({
  open: z.boolean(),
})

export function SimpleTagInputModal({ data, rootElement }: ReactIslandProps) {
  const { open } = useMemo(() => dataSchema.parse(data), [data])

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
              className="p-2 cursor-pointer hover:bg-muted hover:text-primary"
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
