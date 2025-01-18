import Alpine from 'alpinejs'

Alpine.data('collapseContent', () => ({
  can: false,
  expand: false,
  resizeObserver: null as ResizeObserver | null,

  init() {
    this.can = this.$root.hasAttribute('data-collapsible-init')

    window.requestAnimationFrame(() => {
      this.resizeObserver = new ResizeObserver(() => {
        this.recalculate()
      })
      this.resizeObserver.observe(this.$refs.content)
    })
  },

  destroy() {
    this.resizeObserver?.disconnect()
    this.resizeObserver = null
  },

  recalculate() {
    const el = this.$refs.content
    const width = el.clientWidth
    const textSize = 18 // just hard code it what can possible go wrong
    this.can = approximateLines(width, textSize, el.innerHTML) >= 8
  },

  content: {
    'x-ref': 'content',
  },

  button: {
    '@click'() {
      this.expand = !this.expand
    },

    'x-show'() {
      return this.can
    },
  },

  buttonLabel: {
    'x-text'() {
      return this.expand ? window.i18n!['common.less'] : window.i18n!['common.more']
    },
  },

  buttonIcon: {
    'x-text'() {
      return this.expand ? 'collapse_all' : 'expand_all'
    },
  },
}))

/**
 * Calculates the approximate number of lines a given HTML text will take in an element.
 * @param width - The width of the element in pixels.
 * @param textSize - The text size in pixels.
 * @param html - The HTML string containing text and supported tags.
 * @returns The approximate number of lines.
 */
function approximateLines(width: number, textSize: number, html: string): number {
  // Remove HTML tags to get the raw text
  const text: string = stripHTMLTags(html)

  // Estimate the average character width based on the text size (pixels)
  // Assuming 0.6 of text size per character
  const avgCharWidth: number = textSize * 0.6

  // Calculate the total width of the text in pixels
  const totalTextWidth: number = Math.ceil(text.length * avgCharWidth)

  // Calculate the approximate number of lines
  const numLines: number = Math.ceil(totalTextWidth / width)

  return numLines
}

/**
 * Removes supported HTML tags from the input string.
 * @param input - The HTML string to process.
 * @returns The plain text string.
 */
function stripHTMLTags(input: string): string {
  const tagRegex: RegExp = /<\/?(p|b|strong|em|i|span|br)[^>]*?>/gi
  let output: string = input.replace(tagRegex, '')

  // Replace <br> tags with newlines
  output = output.replace(/<br\s*\/?>/gi, '\n')

  return output
}
