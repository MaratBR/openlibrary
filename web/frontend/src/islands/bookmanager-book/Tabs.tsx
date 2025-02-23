import clsx from 'clsx'
import { createContext, PropsWithChildren, useContext, useRef } from 'preact/compat'

export type TabsProps = PropsWithChildren<{
  value?: string
  type: string
  // eslint-disable-next-line no-unused-vars
  onChange: (value: string) => void
}>

const TabValueContext = createContext('')
// eslint-disable-next-line no-unused-vars
const TabCallbackContext = createContext<{ current: (value: string) => void }>({
  current: () => {},
})

export function Tabs({ value, type, children, onChange }: TabsProps) {
  const onChangeRef = useRef(onChange)
  onChangeRef.current = onChange

  return (
    <div role="tablist" class={`tabs tabs--${type}`}>
      <TabCallbackContext.Provider value={onChangeRef}>
        <TabValueContext.Provider value={value || ''}>{children}</TabValueContext.Provider>
      </TabCallbackContext.Provider>
    </div>
  )
}

export type TabProps = PropsWithChildren<{
  value: string
}>

export function Tab({ value, children }: TabProps) {
  const activeValue = useContext(TabValueContext)
  const onChangeRef = useContext(TabCallbackContext)

  return (
    <div
      role="tab"
      onClick={() => onChangeRef.current(value)}
      class={clsx('tab', {
        'tab--active': value === activeValue,
      })}
    >
      {children}
    </div>
  )
}
