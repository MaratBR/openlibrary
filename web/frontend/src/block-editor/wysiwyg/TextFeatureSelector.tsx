import { ChapterContentEditor, useEditorToolbarState } from './editor'
import { DropdownMenu } from 'radix-ui'
import './TextFeatureSelector.scss'
import { JSX, Key } from 'preact'

type FeatureItem = {
  key: string
  label: () => string
  icon: JSX.Element
  onEnable: (editor: ChapterContentEditor) => void
}

const FEATURES: FeatureItem[] = [
  {
    key: 'text',
    label: () => window._('editor.text'),
    icon: <i class="fa-solid fa-t" />,
    onEnable(editor) {
      editor.chain().focus().clearNodes().unsetAllMarks().run()
    },
  },
  {
    key: 'h1',
    label: () => window._('editor.h1'),
    icon: <i class="fa-solid fa-heading" />,
    onEnable(editor) {
      editor.chain().focus().setHeading({ level: 1 }).run()
    },
  },
  {
    key: 'h2',
    label: () => window._('editor.h2'),
    icon: <i class="fa-solid fa-heading" />,
    onEnable(editor) {
      editor.chain().focus().setHeading({ level: 2 }).run()
    },
  },
  {
    key: 'h3',
    label: () => window._('editor.h3'),
    icon: <i class="fa-solid fa-heading" />,
    onEnable(editor) {
      editor.chain().focus().setHeading({ level: 3 }).run()
    },
  },
]

export default function TextFeatureSelector({ editor }: { editor: ChapterContentEditor }) {
  const { textType } = useEditorToolbarState(editor)

  return (
    <DropdownFacade
      items={FEATURES}
      getItemKey={(item) => item.key}
      value={FEATURES.find((x) => x.key === textType)}
      renderItem={(item) => (
        <>
          <div class="be-listitem__icon">{item.icon}</div>
          {item.label()}
        </>
      )}
      onItemSelected={(item) => {
        item.onEnable(editor)
      }}
    />
  )
}

/**
 *    <DropdownMenu.Root modal={true}>
      <DropdownMenu.Trigger asChild>
        <button class="be-listitem be-dropdown">
          Header
          <span class="text-xs text-muted-foreground">
            <i class="fa-solid fa-chevron-down" />
          </span>
        </button>
      </DropdownMenu.Trigger>

      <DropdownMenu.Content align="start">
        <ul class="be-dropdown__menu">
          <li role="button" class="be-listitem">
            <div class="be-listitem__icon">
              <i class="fa-solid fa-heading" />
            </div>
            {window._('editor.h1')}
          </li>
          <li role="button" class="be-listitem">
            <div class="be-listitem__icon">
              <i class="fa-solid fa-heading" />
            </div>
            {window._('editor.h2')}
          </li>
          <li role="button" class="be-listitem">
            <div class="be-listitem__icon">
              <i class="fa-solid fa-heading" />
            </div>
            {window._('editor.h3')}
          </li>
          <li role="button" class="be-listitem">
            <div class="be-listitem__icon">
              <i class="fa-solid fa-list-ul" />
            </div>
            {window._('editor.bulletList')}
          </li>
          <li role="button" class="be-listitem">
            <div class="be-listitem__icon">
              <i class="fa-solid fa-list-ol" />
            </div>
            {window._('editor.orderedList')}
          </li>
        </ul>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
 */

type DropdownFacadeProps<T> = {
  items: T[]
  getItemKey: (item: T) => Key
  renderItem: (item: T) => JSX.Element
  value?: T
  fallbackValue?: JSX.Element
  onItemSelected: (item: T) => void
}

function DropdownFacade<T>({
  items,
  value,
  renderItem,
  getItemKey,
  fallbackValue,
  onItemSelected,
}: DropdownFacadeProps<T>) {
  return (
    <DropdownMenu.Root modal={true}>
      <DropdownMenu.Trigger asChild>
        <button class="be-listitem be-dropdown">
          <div class="inline-block">{value ? renderItem(value) : fallbackValue}</div>

          <span class="text-xs text-muted-foreground">
            <i class="fa-solid fa-chevron-down" />
          </span>
        </button>
      </DropdownMenu.Trigger>

      <DropdownMenu.Content align="start">
        <ul class="be-dropdown__menu">
          {items.map((item) => (
            <li
              key={getItemKey(item)}
              role="button"
              class="be-listitem"
              onClick={() => onItemSelected(item)}
            >
              {renderItem(item)}
            </li>
          ))}
        </ul>
      </DropdownMenu.Content>
    </DropdownMenu.Root>
  )
}
