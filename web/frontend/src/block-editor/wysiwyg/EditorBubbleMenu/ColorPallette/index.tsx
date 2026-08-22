import { useBubbleState } from '../state'
import { ColorPalletteSection } from './ColorPaletteSection'

export { ColorPalletteSection }

export function ColorPallette() {
  const toggle = useBubbleState((x) => x.toggleColorPicker)

  return (
    <button className="be-listitem be-listitem--btn" onClick={toggle}>
      {generateCircleSvg(7, 24)}
    </button>
  )
}

function generateCircleSvg(sectors: number, size: number) {
  if (sectors <= 0) return null

  const cx = 0.5
  const cy = 0.5
  const r = 0.5
  const step = (Math.PI * 2) / sectors

  const overlap = 0.06 // radians

  const paths = Array.from({ length: sectors }, (_, i) => {
    const startAngle = step * i - overlap
    const endAngle = step * (i + 1) + overlap

    const x1 = cx + r * Math.cos(startAngle)
    const y1 = cy + r * Math.sin(startAngle)

    const x2 = cx + r * Math.cos(endAngle)
    const y2 = cy + r * Math.sin(endAngle)

    return (
      <path
        key={i}
        d={`
        M ${cx},${cy}
        L ${x1},${y1}
        A ${r},${r} 0 0,1 ${x2},${y2}
        Z
      `}
        fill={`hsl(${(i / sectors) * 360}deg 100% 50%)`}
        opacity={0.85}
      />
    )
  })

  return (
    <svg width={size} height={size} viewBox="0 0 1 1" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <filter
          id="blur"
          x="-0.2"
          y="-0.2"
          width="1.4"
          height="1.4"
          filterUnits="objectBoundingBox"
        >
          <feGaussianBlur stdDeviation="0.09" />
        </filter>

        <clipPath id="circle-clip">
          <circle cx={cx} cy={cy} r={r} />
        </clipPath>
      </defs>

      <g clipPath="url(#circle-clip)">{paths}</g>
    </svg>
  )
}
