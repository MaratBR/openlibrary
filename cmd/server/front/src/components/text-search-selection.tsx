import React, { useMemo } from 'react'

export type TextSearchSelectionProps = {
  label: string
  selection: string
}

export default function TextSearchSelection({ label, selection }: TextSearchSelectionProps) {
  const elements = useMemo(() => {
    const sel = selection.trim()
    if (sel === '') return [label]

    const parts = label.split(sel)
    const elements: React.JSX.Element[] = []

    for (let i = 0; i < parts.length; i++) {
      elements.push(<span key={i}>{parts[i]}</span>)
      if (i !== parts.length - 1) {
        elements.push(
          <span key={i + 0.5} className="bg-primary/20 text-primary-content">
            {sel}
          </span>,
        )
      }
    }

    return elements
  }, [label])

  return <span>{elements}</span>
}
