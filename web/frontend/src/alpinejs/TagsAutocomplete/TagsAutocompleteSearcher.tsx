import { DefinedTagDto } from '@/api/search'
import TagsInput from '@/components/TagsInput'
import { Wrapper } from '@/preact'
import { useMemo, useRef, useState } from 'react'
import { hydrateRoot } from 'react-dom/client'

export function initTagsAutocomplete($root: HTMLElement): () => void {
  const initialTagsValue = parseTagsValue($root.dataset.value)
  if (!$root.parentElement) throw new Error('TagsAutocomplete root is missing')
  const root = hydrateRoot(
    $root.parentElement,
    <TagsAutocompleteSearcher initialTagsValue={initialTagsValue} />,
  )

  return () => root.unmount()
}

function TagsAutocompleteSearcher({ initialTagsValue }: { initialTagsValue: DefinedTagDto[] }) {
  const [tags, setTags] = useState<DefinedTagDto[]>(initialTagsValue)
  const inputValue = useMemo(() => {
    return tags.map((x) => x.id).join(',')
  }, [tags])
  const inputRef = useRef<HTMLInputElement | null>(null)

  return (
    <Wrapper>
      <TagsInput onInput={setTags} tags={tags} />
      <input ref={inputRef} name="test" type="hidden" aria-hidden="true" value={inputValue} />
    </Wrapper>
  )
}

function parseTagsValue(str: string | undefined): DefinedTagDto[] {
  if (!str) return []

  try {
    return JSON.parse(str) as DefinedTagDto[]
  } catch {
    return []
  }
}
