import { PreactIsland, PreactIslandProps } from '../common/preact-island'
import { useEffect, useMemo, useState } from 'preact/hooks'
import Popper from '@/components/Popper'
import { useMutation, useQuery } from '@tanstack/react-query'
import {
  httpAddBookToCollections,
  httpCreateCollection,
  httpGetCollectionsContainingBook,
  httpGetRecentCollections,
} from './api/collections'
import { z } from 'zod'
import Checkbox from '@/components/Checkbox'
import { isSameCollection } from '@/common/util/collections'

function BookCollectionsPopup({ data }: PreactIslandProps) {
  const { bookId } = useMemo(() => z.object({ bookId: z.string() }).parse(data), [data])

  const target = document.getElementById('BookCollectionsPopup:target')

  const [open, setOpen] = useState(false)
  const [checked, setChecked] = useState<Record<string, boolean>>({})

  const {
    data: recentCollections,
    isLoading: recentIsLoading,
    refetch,
  } = useQuery({
    queryKey: ['recent-collections'],
    queryFn: httpGetRecentCollections,
    enabled: open,
  })

  const { isLoading: containingBookIsLoading, data: containingBookCollections } = useQuery({
    queryKey: ['collections', 'containingBook', bookId],
    queryFn: async () => {
      const response = await httpGetCollectionsContainingBook(bookId)
      setChecked(Object.fromEntries(response.map(({ id }) => [id, true])))
      return response
    },
    enabled: open,
    staleTime: 0,
  })

  const isLoading = containingBookIsLoading || recentIsLoading

  const collections = useMemo(() => {
    const list = recentCollections ? [...recentCollections] : []

    if (containingBookCollections) {
      for (const addedCol of containingBookCollections) {
        if (!list.some((x) => x.id === addedCol.id)) {
          list.push(addedCol)
        }
      }
    }

    return list
  }, [recentCollections, containingBookCollections])

  const hasChanges = useMemo(() => {
    if (isLoading) return false
    if (!collections || !containingBookCollections) return false

    const checkedCollectionIds = Object.entries(checked)
      .filter((x) => x[1])
      .map((x) => x[0])

    return !isSameCollection(
      checkedCollectionIds,
      containingBookCollections.map((x) => x.id),
    )
  }, [checked, collections, containingBookCollections, isLoading])

  useEffect(() => {
    if (!target) return
    const cb = () => {
      setOpen(true)
    }
    target.addEventListener('click', cb)
    return () => target.removeEventListener('click', cb)
  }, [target])

  const [showCollectionNameInput, setShowCollectionNameInput] = useState(false)
  const [collectionName, setCollectionName] = useState('')

  function handleAddCollectionClick() {
    setCollectionName('')
    setShowCollectionNameInput(true)
  }

  const createCollection = useMutation({
    mutationFn: async () => {
      await httpCreateCollection(collectionName)
    },
  })

  async function handleCreateCollection() {
    if (collectionName.trim().length === 0) return
    await createCollection.mutateAsync()
    setCollectionName('')
    setShowCollectionNameInput(false)
    refetch()
  }

  function handleClose() {
    setOpen(false)
    setShowCollectionNameInput(false)
    setCollectionName('')
    setChecked({})
  }

  const addToCollections = useMutation({
    mutationFn: async () => {
      const collectionIds = Object.entries(checked)
        .filter(([_id, added]) => added)
        .map(([id, _added]) => id)
        .filter((id) => collections?.some((x) => x.id === id))

      await httpAddBookToCollections(bookId, collectionIds)

      handleClose()
    },
  })

  if (!target || !open) return null

  return (
    <Popper anchorEl={target} class="z-20">
      <div class="card card--elevated p-0 min-h-32 w-64 relative">
        <button class="btn btn--ghost absolute right-1 top-1" onClick={handleClose}>
          <i class="fa-solid fa-xmark" />
        </button>

        <header class="p-2">
          <p class="font-medium">{window._('book.collection.collectionsListTitle')}</p>
        </header>

        <ul class="max-h-96 overflow-y-auto" style="scrollbar-width:thin">
          {isLoading && (
            <div class="flex items-center justify-center my-8">
              <span class="loader" />
            </div>
          )}
          {collections?.map((c) => (
            <li
              role="button"
              class="listitem listitem--lg"
              key={c.id}
              onClick={(e) => {
                if (e.target !== e.currentTarget) return
                setChecked((v) => ({
                  ...v,
                  [c.id]: !v[c.id],
                }))
              }}
            >
              <Checkbox
                checked={checked[c.id] ?? false}
                onChange={(e) => {
                  setChecked((v) => ({
                    ...v,
                    [c.id]: (e.target as HTMLInputElement).checked,
                  }))
                }}
              >
                <span class="ml-2 cursor-pointer select-none">{c.name}</span>
              </Checkbox>
            </li>
          ))}
          {showCollectionNameInput ? (
            <div class="p-1 flex gap-1">
              <input
                value={collectionName}
                onInput={(e) => setCollectionName((e.target as HTMLInputElement).value)}
                class="input"
                placeholder={window._('book.collection.namePlaceholder')}
              />
              <div class="btn-group">
                <button
                  disabled={collectionName.trim().length === 0}
                  class="btn btn--ghost btn--sm"
                  onClick={handleCreateCollection}
                >
                  <i class="fa-solid fa-check" />
                </button>
                <button
                  class="btn btn--ghost btn--sm"
                  onClick={() => setShowCollectionNameInput(false)}
                >
                  <i class="fa-solid fa-xmark" />
                </button>
              </div>
            </div>
          ) : (
            <li role="button" class="listitem listitem--lg" onClick={handleAddCollectionClick}>
              <i class="fa-solid fa-square-plus mr-2" />
              {window._('book.collection.addCollection')}
            </li>
          )}
        </ul>
        {hasChanges && (
          <button
            class="btn primary m-2"
            disabled={addToCollections.isPending}
            onClick={() => addToCollections.mutate()}
          >
            {window._('book.collection.applyChanges')}
          </button>
        )}
      </div>
    </Popper>
  )
}

export default new PreactIsland(BookCollectionsPopup)
