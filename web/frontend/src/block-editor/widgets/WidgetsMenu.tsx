import { useEffect, useState } from 'react'
import { Widget } from './core'
import { WidgetsService } from './service'

export function WidgetsMenu({ service }: { service: WidgetsService }) {
  const [widgets, setWidgets] = useState<Widget[]>([])

  useEffect(() => {
    service.loadWidgets().then(setWidgets)
  }, [service])

  return (
    <section>
      <div className="grid grid-cols-2">
        {widgets.map((widget) => (
          <div key={widget.name} className="be-widget-card">
            {widget.name}
          </div>
        ))}
      </div>
    </section>
  )
}
