import { useEffect, useState } from 'preact/hooks'
import { Widget } from './core'
import { WidgetsService } from './service'

export function WidgetsMenu({ service }: { service: WidgetsService }) {
  const [widgets, setWidgets] = useState<Widget[]>([])

  useEffect(() => {
    service.loadWidgets().then(setWidgets)
  }, [service])

  return (
    <section>
      <div class="grid grid-cols-2">
        {widgets.map((widget) => (
          <div key={widget.name} class="be-widget-card">
            {widget.name}
          </div>
        ))}
      </div>
    </section>
  )
}
