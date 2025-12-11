import clsx from 'clsx'
import { HTMLAttributes } from 'preact'
import {
  createContext,
  ForwardedRef,
  forwardRef,
  PropsWithChildren,
  useContext,
  useRef,
} from 'preact/compat'

export type TabsProps = PropsWithChildren<{
  value?: string
  onChange: (value: string) => void
}>

const TabValueContext = createContext('')

const TabCallbackContext = createContext<{ current: (value: string) => void }>({
  current: () => {},
})

function TabsRoot({ value, children, onChange }: TabsProps) {
  const onChangeRef = useRef(onChange)
  onChangeRef.current = onChange

  return (
    <TabCallbackContext.Provider value={onChangeRef}>
      <TabValueContext.Provider value={value || ''}>{children}</TabValueContext.Provider>
    </TabCallbackContext.Provider>
  )
}

export type TabProps = PropsWithChildren<{
  value: string
}>

function TabsTab({ value, children }: TabProps) {
  const activeValue = useContext(TabValueContext)
  const onChangeRef = useContext(TabCallbackContext)

  return (
    <li
      role="tab"
      class={clsx('tab', {
        'tab--active': value === activeValue,
      })}
      onClick={(e) => {
        e.preventDefault()
        onChangeRef.current(value)
      }}
    >
      <span class="tabs__tab__title">{children}</span>
    </li>
  )
}

const TabsBody = forwardRef(
  (
    { class: class_, className, ...props }: HTMLAttributes<HTMLDivElement>,
    ref: ForwardedRef<HTMLDivElement>,
  ) => {
    return <div ref={ref} class={clsx('tabs__body', className || class_)} {...props} />
  },
)

const TabsMenu = forwardRef(
  (
    { class: class_, className, children, ...props }: HTMLAttributes<HTMLUListElement>,
    ref: ForwardedRef<HTMLUListElement>,
  ) => {
    return (
      <ul
        ref={ref}
        role="tablist"
        class={clsx('tabs tabs--primary', className || class_)}
        {...props}
      >
        {children}
      </ul>
    )
  },
)

const Tabs = {
  Root: TabsRoot,
  List: TabsMenu,
  Tab: TabsTab,
  Body: TabsBody,
}

export default Tabs
