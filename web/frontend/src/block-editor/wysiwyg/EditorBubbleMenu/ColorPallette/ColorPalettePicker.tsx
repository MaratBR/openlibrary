import React from 'react'
import './ColorPalettePicker.scss'

const COLORS: [number, number, number][] = [
  [230, 57, 70],
  [244, 114, 73],
  [249, 168, 37],
  [230, 196, 58],
  [138, 201, 38],

  [52, 168, 83],
  [42, 157, 143],
  [38, 166, 154],
  [35, 154, 214],
  [67, 97, 238],

  [92, 72, 205],
  [126, 87, 194],
  [168, 85, 247],
  [209, 79, 193],
  [224, 71, 138],

  [146, 64, 14],
  [121, 85, 61],
  [67, 90, 111],
  [46, 64, 83],
  [42, 46, 54],

  [255, 183, 197],
  [255, 214, 165],
  [180, 228, 190],
  [168, 218, 220],
  [190, 192, 255],
]

export function ColorPallettePicker({
  onColorChosen,
}: {
  onColorChosen: (color: null | { r: number; g: number; b: number }) => void
}) {
  return (
    <div className="BeColorPalettePicker">
      <button
        className="BeColorPalettePicker-item"
        style={{ '--color': `Var(--foreground)` } as React.CSSProperties}
        onClick={() => onColorChosen(null)}
      />
      {COLORS.map(([r, g, b]) => (
        <button
          key={`${r}-${g}-${b}`}
          className="BeColorPalettePicker-item"
          style={{ '--color': `rgb(${r}, ${g}, ${b})` } as React.CSSProperties}
          onClick={() => onColorChosen({ r, g, b })}
        />
      ))}
    </div>
  )
}
