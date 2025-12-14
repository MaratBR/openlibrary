import { DefinedTagDto, definedTagDtoSchema } from '@/api/search'
import TagsInput from '@/components/TagsInput'
import Wrapper from '@/preact/wrapper'
import { hydrate, render } from 'preact'
import { useMemo, useRef, useState } from 'preact/hooks'
import z from 'zod'

export function initTagsAutocomplete($root: HTMLElement): () => void {
  const initialTagsValue = parseTagsValue($root.dataset.value)
  hydrate(
    <TagsAutocompleteSearcher initialTagsValue={initialTagsValue} />,
    $root.parentElement as any,
  )
  return () => render(null, $root)
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

const definedTagDtoArraySchema = z.array(definedTagDtoSchema)

function parseTagsValue(str: string | undefined): DefinedTagDto[] {
  if (!str) return []

  try {
    const jsonValue = JSON.parse(str)
    return definedTagDtoArraySchema.parse(jsonValue)
  } catch {
    return []
  }
}
