import clsx from 'clsx'
import {
  createContext,
  ForwardedRef,
  forwardRef,
  HTMLAttributes,
  PropsWithChildren,
  useContext,
  useRef,
} from 'preact/compat'

export type TabsProps = PropsWithChildren<{
  value?: string
  class?: string

  onChange: (value: string) => void
}>

const TabValueContext = createContext('')

const TabCallbackContext = createContext<{ current: (value: string) => void }>({
  current: () => {},
})

function SideTabsRoot({ value, class: className, children, onChange }: TabsProps) {
  const onChangeRef = useRef(onChange)
  onChangeRef.current = onChange

  return (
    <div role="tablist" class={clsx('side-tabs', className)}>
      <TabCallbackContext.Provider value={onChangeRef}>
        <TabValueContext.Provider value={value || ''}>{children}</TabValueContext.Provider>
      </TabCallbackContext.Provider>
    </div>
  )
}

export type TabProps = PropsWithChildren<{
  value: string
}>

function SideTabsTab({ value, children }: TabProps) {
  const activeValue = useContext(TabValueContext)
  const onChangeRef = useContext(TabCallbackContext)

  return (
    <li
      role="tab"
      class={clsx('side-tabs__tab', {
        'side-tabs__tab--active': value === activeValue,
      })}
      onClick={(e) => {
        e.preventDefault()
        onChangeRef.current(value)
      }}
    >
      <span class="side-tabs__tab__title">{children}</span>
    </li>
  )
}

const SideTabsBody = forwardRef(
  (
    { class: class_, className, ...props }: HTMLAttributes<HTMLDivElement>,
    ref: ForwardedRef<HTMLDivElement>,
  ) => {
    return <div ref={ref} class={clsx('side-tabs__body', className || class_)} {...props} />
  },
)

const SideTabsMenu = forwardRef(
  (
    { class: class_, className, children, ...props }: HTMLAttributes<HTMLDivElement>,
    ref: ForwardedRef<HTMLDivElement>,
  ) => {
    return (
      <div ref={ref} class={clsx('side-tabs__menu', className || class_)} {...props}>
        <ul>{children}</ul>
      </div>
    )
  },
)

const SideTabs = {
  Root: SideTabsRoot,
  List: SideTabsMenu,
  Tab: SideTabsTab,
  Body: SideTabsBody,
}

export default SideTabs
