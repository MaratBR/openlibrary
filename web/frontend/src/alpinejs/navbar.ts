import { Subject } from '@/common/rx'
import Alpine from 'alpinejs'

Alpine.data('siteNav', () => ({
  elements: {} as Record<string, HTMLElement>,
  activeSubmenu: new Subject(''),
  navLink: {
    'x-init'() {
      const submenu = this.$el.dataset.submenu
      if (!submenu) return

      this.elements[`link_${submenu}`] = this.$el
    },
  },
  navSubMenu: {
    'x-init'() {
      const submenu = this.$el.dataset.submenu
      if (!submenu) return

      this.$el.style.display = 'none'
      this.$el.setAttribute('data-open', 'false')

      this.activeSubmenu.subscribe((active) => {
        const isActive = this.$el.getAttribute('data-open') === 'true'
        const shouldBeActive = active === submenu

        if (shouldBeActive === isActive) return

        this.$el.setAttribute('data-open', shouldBeActive ? 'true' : 'false')
      })

      window.requestAnimationFrame(() => {
        const link = this.elements[`link_${submenu}`]
        this.$el.style.left = `${link.getBoundingClientRect().left}px`

        link.addEventListener('mouseenter', () => {
          this.activeSubmenu.set(submenu)
        })
      })
    },
  },
}))
