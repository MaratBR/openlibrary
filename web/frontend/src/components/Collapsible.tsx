import { animate, linear } from 'popmotion'
import { Component, ComponentChild, ComponentChildren, RenderableProps } from 'preact'

export class Collapsible extends Component<{
  in: boolean
  duration: number
  children: ComponentChild
}> {
  private $el: HTMLElement | null = null
  private $inner: HTMLElement | null = null
  private animation?: { stop: () => void }
  private animationProgress = 0
  private expectedHeight = 0
  private resizeObserver: ResizeObserver | null = null

  componentDidMount() {
    const $el = this.$el
    if (!$el) return

    $el.style.willChange = 'height'
    $el.style.transition = 'height 300ms'

    if (this.props.in) {
      this.animationProgress = 1
    } else {
      this.animationProgress = 0
      $el.style.height = '0px'
    }

    this._recalculateHeight()
    this.resizeObserver = new ResizeObserver(() => {
      this._recalculateHeight()
    })
    if (this.$inner) this.resizeObserver.observe(this.$inner)
  }

  componentWillUnmount(): void {
    this.animation?.stop()
    this.animation = undefined
    this.resizeObserver?.disconnect()
  }

  componentDidUpdate(previousProps: Readonly<{ in: boolean; children: ComponentChild }>): void {
    if (previousProps.in !== this.props.in) {
      this.animation?.stop()
      if (this.props.in) {
        this.animation = animate({
          from: this.animationProgress,
          to: 1,
          ease: linear,
          duration: this.props.duration * (1 - this.animationProgress),
          onUpdate: this._onUpdate.bind(this),
        })
      } else {
        this.animation = animate({
          from: this.animationProgress,
          to: 0,
          ease: linear,
          duration: this.props.duration * this.animationProgress,
          onUpdate: this._onUpdate.bind(this),
        })
      }
    }
  }

  private _onUpdate(latest: number) {
    this.animationProgress = latest
    const $el = this.$el
    if (!$el || this.expectedHeight === 0) return
    $el.style.height = `${Math.ceil(latest * this.expectedHeight)}px`
  }

  private _recalculateHeight() {
    const $inner = this.$inner
    if (!$inner) return
    this.expectedHeight = $inner.getBoundingClientRect().height
  }

  render(
    props?: RenderableProps<{ in: boolean; children: ComponentChild }, unknown> | undefined,
  ): ComponentChildren {
    return (
      <div
        ref={(el) => {
          this.$el = el
        }}
        class="overflow-y-hidden"
      >
        <div
          ref={(el) => {
            this.$inner = el
          }}
        >
          {props?.children}
        </div>
      </div>
    )
  }
}
