import { httpLikeComment } from '@/api/comments'
import Alpine from 'alpinejs'

Alpine.data('Comment', ({ id, liked = false }: { id: string; liked?: boolean }) => ({
  liked,
  openEditor: false,

  init() {
    this.$watch(
      'liked',
      Alpine.throttle((newValue: boolean) => {
        httpLikeComment(id, newValue)
      }, 1000),
    )
  },

  like: {
    '@click'() {
      this.liked = !this.liked

      let text = getFirstText(this.$el)
      if (!text) {
        text = document.createTextNode('')
        this.$el.appendChild(text)
      }

      const newCount = +text.textContent + (this.liked ? 1 : -1)
      if (newCount === 0) {
        text.data = ''
      } else {
        text.data = `${newCount}`
      }
      this.$el.setAttribute('data-set', `${this.liked}`)
    },
  },

  reply: {},
}))

function getFirstText(el: Element): Text | null {
  for (let i = 0; i < el.childNodes.length; i++) {
    const node = el.childNodes.item(i)
    if (node instanceof Text) return node
  }

  return null
}
