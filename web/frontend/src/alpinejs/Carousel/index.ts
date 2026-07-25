import Alpine from 'alpinejs'

Alpine.data('Carousel', () => {
  let list: HTMLElement | null = null

  return {
    init() {
      list = getListElement(this.$root)
      if (!list) {
        console.error(
          'cannot find list element in carousel, must have a ul/ol element directly inside carousel component',
        )
        return
      }

      this.onScrollEnd = this.onScrollEnd.bind(this)
      list.addEventListener('scrollend', this.onScrollEnd)

      requestAnimationFrame(() => this.updateCurtains())
    },

    onScrollEnd() {
      this.updateCurtains()
    },

    destroy() {
      if (list) {
        list.removeEventListener('scrollend', this.onScrollEnd)
      }
    },

    updateCurtains() {
      if (!list) return
      const { scrollLeft, scrollWidth, clientWidth } = list

      const leftCurtainVisible = scrollLeft > 1

      const rightCurtainVisible = scrollWidth - (scrollLeft + clientWidth) > 1

      this.$root.dataset.canScrollL = `${leftCurtainVisible}`
      this.$root.dataset.canScrollR = `${rightCurtainVisible}`
    },

    scroll(back = false) {
      if (!list) return

      const items = Array.from(list.children) as HTMLElement[]
      if (items.length === 0) return

      const listRect = list.getBoundingClientRect()
      const visibleLeft = Math.max(listRect.left, 0)
      const visibleRight = Math.min(listRect.right, window.innerWidth)

      let target: HTMLElement

      if (back) {
        const firstFullyVisibleIndex = items.findIndex((item) => {
          const rect = item.getBoundingClientRect()

          return rect.left >= visibleLeft && rect.right <= visibleRight
        })

        if (firstFullyVisibleIndex === -1) return

        target = items[Math.max(0, firstFullyVisibleIndex - 1)]

        const targetRect = target.getBoundingClientRect()

        list.scrollBy({
          // Align the target's right edge with the viewport's right edge.
          left: targetRect.right - visibleRight,
          behavior: 'smooth',
        })
      } else {
        let lastFullyVisibleIndex = -1

        for (let index = 0; index < items.length; index++) {
          const rect = items[index].getBoundingClientRect()

          if (rect.left >= visibleLeft && rect.right <= visibleRight) {
            lastFullyVisibleIndex = index
          }
        }

        if (lastFullyVisibleIndex === -1) return

        target = items[Math.min(items.length - 1, lastFullyVisibleIndex + 1)]

        const targetRect = target.getBoundingClientRect()

        list.scrollBy({
          // Align the target's left edge with the viewport's left edge.
          left: targetRect.left - visibleLeft,
          behavior: 'smooth',
        })
      }

      requestAnimationFrame(() => this.updateCurtains())
    },
    btnL: {
      '@click'() {
        this.scroll(true)
      },
    },

    btnR: {
      '@click'() {
        this.scroll()
      },
    },
  }
})

function getListElement(el: HTMLElement): HTMLElement | null {
  if (el instanceof HTMLUListElement || el instanceof HTMLOListElement) {
    return el
  }

  const list = el.querySelector(':scope > ul, :scope > ol')
  if (list instanceof HTMLElement) {
    return list
  }

  return null
}
