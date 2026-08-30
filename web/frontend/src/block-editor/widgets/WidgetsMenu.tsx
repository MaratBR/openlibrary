import { useEffect, useState } from 'react'
import { useAtomValue } from 'jotai'
import { Widget } from './core'
import { WidgetsService } from './service'
import { wysiwygEditorAtom } from '../wysiwyg/state'

export function WidgetsMenu({ service }: { service: WidgetsService }) {
  const [widgets, setWidgets] = useState<Widget[]>([])
  const editor = useAtomValue(wysiwygEditorAtom)

  useEffect(() => {
    service.getWidgets().then(setWidgets)
  }, [service])

  return (
    <section className="p-3">
      <h2 className="text-sm font-semibold mb-3">{window._('editor.widgets')}</h2>
      <div className="grid grid-cols-2 gap-2">
        {widgets.map((widget) => (
          <button
            type="button"
            key={widget.name}
            disabled={!editor}
            className="be-widget-card"
            onClick={() => editor && widget.apply(editor)}
          >
            <span className="be-widget-card__icon">{widget.icon}</span>
            <span>{widget.name}</span>
          </button>
        ))}
      </div>
    </section>
  )
}
