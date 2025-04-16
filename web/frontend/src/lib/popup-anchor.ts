export type Origin = {
  vertical: 'top' | 'bottom' | 'center'
  horizontal: 'left' | 'right' | 'center'
}

export type PopupControllerOptions = {
  anchorOrigin: Origin
  targetOrigin: Origin
}

export class PopupController {
  private readonly _anchor: HTMLElement
  private readonly _element: HTMLElement
  private readonly _options: PopupControllerOptions

  constructor(anchor: HTMLElement, element: HTMLElement, options: PopupControllerOptions) {
    this._anchor = anchor
    this._element = element
    this._options = options
  }

  /**
   * Updates the position of the target element based on provided options and anchor
   */
  update() {
    const anchorRect = this._anchor.getBoundingClientRect()
    const elementRect = this._element.getBoundingClientRect()

    let left = 0,
      top = 0

    left = anchorRect.left
    top = anchorRect.top

    const anchorVertical = this._options.anchorOrigin.vertical

    if (anchorVertical === 'center') {
      top += anchorRect.height / 2
    } else if (anchorVertical === 'bottom') {
      top += anchorRect.height
    }

    const anchorHorizontal = this._options.anchorOrigin.horizontal

    if (anchorHorizontal === 'center') {
      left += anchorRect.width / 2
    } else if (anchorHorizontal === 'right') {
      left += anchorRect.width
    }

    const targetVertical = this._options.targetOrigin.vertical

    if (targetVertical === 'center') {
      top -= elementRect.height / 2
    } else if (targetVertical === 'bottom') {
      top -= elementRect.height
    }

    const targetHorizontal = this._options.targetOrigin.horizontal

    if (targetHorizontal === 'center') {
      left -= elementRect.width / 2
    } else if (targetHorizontal === 'right') {
      left -= elementRect.width
    }

    this._element.style.position = 'fixed'
    this._element.style.left = `${left}px`
    this._element.style.top = `${top}px`
  }

  dispose() {}
}
