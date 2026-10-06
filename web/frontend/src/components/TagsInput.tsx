import { useEffect, useId, useRef, useState } from 'react'
import { DefinedTagDto, useTagsSearch } from '@/features/search'

export type TagsInputProps = {
  tags?: DefinedTagDto[]
  onInput?: (tags: DefinedTagDto[]) => void
  id?: string
}

export default function TagsInput({ tags = [], onInput, id }: TagsInputProps) {
  const generatedId = useId()
  const inputId = id ?? generatedId
  const menuId = `${inputId}-options`
  const inputRef = useRef<HTMLInputElement>(null)
  const [searchQuery, setSearchQuery] = useState('')
  const [open, setOpen] = useState(false)
  const [activeId, setActiveId] = useState<string | null>(null)
  const hasSearchQuery = searchQuery.trim().length > 0
  const query = useTagsSearch({ query: searchQuery, enabled: hasSearchQuery })
  const results = (query.data ?? []).filter((tag) => !tags.some((x) => x.id === tag.id))
  const activeIndex = results.findIndex((tag) => tag.id === activeId)
  const pending = query.isLoading || (query.data === undefined && !query.error)
  const activeOptionId =
    open && !pending && !query.error && activeIndex >= 0 ? `${menuId}-${activeId}` : undefined

  useEffect(() => {
    if (activeOptionId)
      document.getElementById(activeOptionId)?.scrollIntoView({ block: 'nearest' })
  }, [activeOptionId])

  function add(tag: DefinedTagDto) {
    if (tags.some((x) => x.id === tag.id)) return
    onInput?.([...tags, tag])
    setSearchQuery('')
    setOpen(false)
    setActiveId(null)
    inputRef.current?.focus()
  }

  return (
    <div
      className="input-dropdown TagsAutocomplete"
      data-open={open}
      onBlur={(event) => {
        if (!event.currentTarget.contains(event.relatedTarget)) {
          setOpen(false)
          setActiveId(null)
        }
      }}
    >
      <div className="TagsAutocomplete-control" onClick={() => inputRef.current?.focus()}>
        {tags.map((tag) => (
          <span key={tag.id} className="Chip TagsAutocomplete-chip">
            <span className="TagsAutocomplete-name">{tag.name}</span>
            {tag.adult && <span className="Tag-adult">18+</span>}
            <button
              type="button"
              className="Chip-close TagsAutocomplete-remove"
              aria-label={`${window._('search.removeTag')}: ${tag.name}`}
              onClick={(event) => {
                event.stopPropagation()
                onInput?.(tags.filter((x) => x.id !== tag.id))
                inputRef.current?.focus()
              }}
            >
              <i className="fa-solid fa-xmark" aria-hidden="true" />
            </button>
          </span>
        ))}
        <input
          ref={inputRef}
          id={inputId}
          className="Dropdown-input"
          value={searchQuery}
          placeholder={window._('search.findTags')}
          aria-label={id ? undefined : window._('search.tags')}
          role="combobox"
          aria-autocomplete="list"
          aria-expanded={open}
          aria-controls={menuId}
          aria-activedescendant={activeOptionId}
          autoComplete="off"
          onFocus={() => setOpen(hasSearchQuery)}
          onChange={(event) => {
            setSearchQuery(event.target.value)
            setActiveId(null)
            setOpen(event.target.value.trim().length > 0)
          }}
          onKeyDown={(event) => {
            if (event.nativeEvent.isComposing) return
            if (event.key === 'Escape') {
              if (open) {
                event.preventDefault()
                event.stopPropagation()
                setOpen(false)
                setActiveId(null)
              }
            } else if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
              event.preventDefault()
              if (!hasSearchQuery) return
              setOpen(true)
              if (!pending && !query.error && results.length) {
                const next =
                  event.key === 'ArrowDown'
                    ? (activeIndex + 1) % results.length
                    : activeIndex < 0
                      ? results.length - 1
                      : (activeIndex - 1 + results.length) % results.length
                setActiveId(results[next].id)
              }
            } else if (event.key === 'Enter' && open) {
              event.preventDefault()
              if (activeOptionId) add(results[activeIndex])
            }
          }}
        />
      </div>
      <div className="Dropdown-menu TagsAutocomplete-menu" hidden={!open}>
        <div role="status" className="TagsAutocomplete-status" aria-live="polite">
          {pending
            ? window._('search.tagsLoading')
            : query.error
              ? window._('search.tagsError')
              : results.length === 0
                ? window._('search.noResults')
                : null}
        </div>
        <ul id={menuId} role="listbox" aria-label={window._('search.tags')} aria-busy={pending}>
          {!pending &&
            !query.error &&
            results.map((tag) => (
              <li
                key={tag.id}
                id={`${menuId}-${tag.id}`}
                role="option"
                aria-selected={activeId === tag.id}
                className="TagsAutocomplete-option"
                onMouseDown={(event) => event.preventDefault()}
                onMouseMove={() => setActiveId(tag.id)}
                onClick={() => add(tag)}
              >
                <span className="TagsAutocomplete-name">{tag.name}</span>
                {tag.adult && <span className="Tag-adult">18+</span>}
                <i
                  className="fa-solid fa-plus text-xs text-secondary-foreground"
                  aria-hidden="true"
                />
              </li>
            ))}
        </ul>
      </div>
    </div>
  )
}
