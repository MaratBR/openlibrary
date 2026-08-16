import { searchModerationUsers } from '@/api/moderation'
import type { ModerationUserListEntry } from '@/api/moderation'
import { Checkbox, Popover, PopoverContent, PopoverTrigger } from '@/components'
import { useEffect, useRef, useState } from 'react'

export function UserMultiSelect({ value, onChange, max = 20 }: { value: ModerationUserListEntry[]; onChange: (users: ModerationUserListEntry[]) => void; max?: number }) {
  const [open, setOpen] = useState(false)
  const [search, setSearch] = useState('')
  const [results, setResults] = useState<ModerationUserListEntry[]>([])
  const [pending, setPending] = useState(false)
  const request = useRef(0)

  useEffect(() => {
    if (!open) return
    const current = ++request.current
    const timeout = window.setTimeout(async () => {
      setPending(true)
      try {
        const response = await searchModerationUsers(search.trim(), '', '', 1, 10)
        response.throwIfError()
        if (request.current === current) setResults(response.data.entries)
      } catch (error) {
        if (request.current === current) window.toast.error(error)
      } finally {
        if (request.current === current) setPending(false)
      }
    }, 200)
    return () => window.clearTimeout(timeout)
  }, [open, search])

  const selected = new Set(value.map((user) => user.id))
  const toggle = (user: ModerationUserListEntry) => {
    if (selected.has(user.id)) onChange(value.filter((entry) => entry.id !== user.id))
    else if (value.length < max) onChange([...value, user])
  }

  return <div>
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild><button type="button" className="select w-full justify-between"><span className={value.length ? '' : 'text-secondary-foreground'}>{value.length ? window._('moderationPortal.loginHistoryPage.usersSelected', { count: String(value.length) }) : window._('moderationPortal.loginHistoryPage.allUsers')}</span><i className="fa-solid fa-chevron-down text-xs opacity-60" /></button></PopoverTrigger>
      <PopoverContent className="w-[var(--radix-popover-trigger-width)] min-w-72 p-0">
        <div className="p-2 border-b border-border"><input className="input w-full" value={search} onChange={(event) => setSearch(event.target.value)} placeholder={window._('moderationPortal.loginHistoryPage.searchUsers')} autoFocus /></div>
        <div role="listbox" aria-multiselectable="true" className="max-h-72 overflow-y-auto p-1">
          {pending ? <div className="p-5 text-center"><span className="circle-loader" /></div> : results.length ? results.map((user) => <div key={user.id} role="option" aria-selected={selected.has(user.id)} tabIndex={0} className="w-full flex cursor-pointer items-center gap-3 rounded-sm px-2 py-2 hover:bg-foreground/5 focus-visible:ring-2 focus-visible:ring-primary/45" onClick={() => toggle(user)} onKeyDown={(event) => { if (event.key === 'Enter' || event.key === ' ') { event.preventDefault(); toggle(user) } }}><Checkbox checked={selected.has(user.id)} tabIndex={-1} aria-hidden="true" className="pointer-events-none" /><img className="size-8 shrink-0 rounded-full object-cover" src={user.avatar} alt="" /><span className="min-w-0 flex-1"><span className="block truncate font-medium">{user.name}</span><span className="block truncate font-mono text-xs text-secondary-foreground">{user.id}</span></span></div>) : <div className="p-5 text-center text-secondary-foreground">{window._('moderationPortal.loginHistoryPage.noUsers')}</div>}
        </div>
        <div className="px-3 py-2 border-t border-border text-xs text-secondary-foreground">{window._('moderationPortal.loginHistoryPage.selectionLimit', { count: String(max) })}</div>
      </PopoverContent>
    </Popover>
    {value.length > 0 && <div className="flex flex-wrap gap-2 mt-2">{value.map((user) => <button key={user.id} type="button" className="chip hover:bg-foreground/10" onClick={() => toggle(user)}>{user.name}<i className="fa-solid fa-xmark ml-2" /></button>)}</div>}
  </div>
}
