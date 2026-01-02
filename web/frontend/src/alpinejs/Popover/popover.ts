import { AnimationController, ModalAnimation } from '@/lib/animate'
import { computePosition, ComputePositionConfig } from '@floating-ui/react'
import Alpine from 'alpinejs'

Alpine.data(
  'Popover',
  (params: { duration?: number; options?: Partial<ComputePositionConfig> }) => ({
    $anchorEl: null as null | HTMLElement,
    animation: null as AnimationController | null,
    isOpen: false,

    open(element: HTMLElement) {
      this.$anchorEl = element
    },

    close() {
      this.$anchorEl = null
    },

    init() {
      const { duration = 150, options } = params ?? {}
      this.animation = new ModalAnimation(this.$el, duration)
      this.animation.setShow(false, 0)

      this.$watch('$anchorEl', (newAnchor, oldValue) => {
        if (!!newAnchor === false) {
          this.animation?.setShow(false)
          requestAnimationFrame(() => {
            this.isOpen = false
          })
          return
        }

        if (newAnchor) {
          computePosition(newAnchor, this.$el, options).then((pos) => {
            this.$el.style.position = 'fixed'
            this.$el.style.left = `${pos.x}px`
            this.$el.style.top = `${pos.y}px`

            if (!!oldValue === false) {
              this.animation?.setShow(true)
            }

            requestAnimationFrame(() => {
              this.isOpen = true
            })
          })
        }
      })

      const onClickOutside = () => {
        this.close()
      }

      this.$watch('isOpen', (isOpen) => {
        if (isOpen) {
          window.addEventListener('click', onClickOutside)
        } else {
          window.removeEventListener('click', onClickOutside)
        }
      })
    },

    destroy() {
      this.animation?.dispose()
      this.animation = null
    },
  }),
)
