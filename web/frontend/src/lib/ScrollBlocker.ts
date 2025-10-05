class ScrollBlocker extends HTMLElement {
  private static instances: Set<ScrollBlocker> = new Set()
  private static isScrolling = false
  private static scrollTimeout: number | null = null

  // Configuration
  private static readonly SCROLL_DELAY = 50 // ms after scroll ends before re-enabling pointer events

  constructor() {
    super()
    ScrollBlocker.instances.add(this)
  }

  connectedCallback() {
    this.applyStyles()
    this.setupEventListeners()
    this.updateBodyPointerEvents()
  }

  disconnectedCallback() {
    this.removeEventListeners()
    ScrollBlocker.instances.delete(this)
    this.updateBodyPointerEvents()
  }

  private applyStyles() {
    // Make element completely invisible and excluded from layout
    this.style.display = 'none'
    this.style.position = 'absolute'
    this.style.visibility = 'hidden'
    this.style.pointerEvents = 'none'
    this.style.userSelect = 'none'
    this.style.opacity = '0'
    this.style.width = '0'
    this.style.height = '0'
    this.style.margin = '0'
    this.style.padding = '0'
    this.style.border = '0'
    this.style.overflow = 'hidden'
    this.style.clip = 'rect(0, 0, 0, 0)'
    this.style.clipPath = 'inset(50%)'

    // Ensure it doesn't affect any layout
    this.style.flex = '0 0 0'
    this.style.flexBasis = '0'
    this.style.flexGrow = '0'
    this.style.flexShrink = '0'
    this.style.alignSelf = 'auto'
    this.style.order = '0'
    this.style.gridArea = 'auto'
    this.style.zIndex = '-1'
  }

  private setupEventListeners() {
    window.addEventListener('scroll', this.handleScroll, { passive: true })
  }

  private removeEventListeners() {
    window.removeEventListener('scroll', this.handleScroll)
  }

  private handleScroll = () => {
    if (!ScrollBlocker.isScrolling) {
      ScrollBlocker.isScrolling = true
      this.updateBodyPointerEvents()
    }

    // Clear existing timeout
    if (ScrollBlocker.scrollTimeout !== null) {
      window.clearTimeout(ScrollBlocker.scrollTimeout)
    }

    // Set new timeout to detect when scrolling stops
    ScrollBlocker.scrollTimeout = window.setTimeout(() => {
      ScrollBlocker.isScrolling = false
      this.updateBodyPointerEvents()
    }, ScrollBlocker.SCROLL_DELAY)
  }

  private updateBodyPointerEvents() {
    // Only disable pointer events if there are active instances AND scrolling is happening
    const shouldDisablePointerEvents = ScrollBlocker.instances.size > 0 && ScrollBlocker.isScrolling

    document.body.style.pointerEvents = shouldDisablePointerEvents ? 'none' : ''

    // Optional: Add/remove a class for more complex styling
    if (shouldDisablePointerEvents) {
      document.body.classList.add('scroll-blocker-active')
    } else {
      document.body.classList.remove('scroll-blocker-active')
    }
  }

  // Public method to manually trigger scroll end (useful for testing)
  public forceScrollEnd() {
    ScrollBlocker.isScrolling = false
    this.updateBodyPointerEvents()
  }

  // Static method to check if any scroll blockers are active
  public static get isActive(): boolean {
    return this.instances.size > 0
  }

  // Static method to get number of active instances
  public static get activeInstanceCount(): number {
    return this.instances.size
  }
}

// Define the custom element
customElements.define('ol-scroll-blocker', ScrollBlocker)

export default ScrollBlocker
