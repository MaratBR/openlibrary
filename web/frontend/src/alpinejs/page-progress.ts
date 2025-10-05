import { throttleUntilNextFrame } from '@/common/dom'
import Alpine from 'alpinejs'

Alpine.data('pageProgress', () => ({
  init() {
    const selector = this.$root.dataset.contentElement
    if (!selector) return
    const pageContent = document.querySelector(selector)
    if (!pageContent || !(pageContent instanceof HTMLElement)) return
    this.$root.style.transform = 'scaleX(0)'

    const updateTransform = throttleUntilNextFrame((progress: number) => {
      this.$root.style.transform = `scaleX(${progress})`
    })

    const updateFn = () => {
      const rect = pageContent.getBoundingClientRect()
      const { scrollY, innerHeight } = window
      const containerBottomPosition = rect.bottom + scrollY
      const progress = Math.min(1, scrollY / (containerBottomPosition - innerHeight))
      updateTransform(progress)
    }

    window.requestAnimationFrame(updateFn)

    window.addEventListener('scroll', updateFn)
  },
}))
