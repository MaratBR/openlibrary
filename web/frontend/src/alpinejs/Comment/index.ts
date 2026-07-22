import { httpLikeComment } from '@/api/comments'
import Alpine from 'alpinejs'
import { CommentEditorController, initCommentEditor } from './CommentEditor'
import { CommentRepliesController, initCommentReplies } from './CommentReplies'

Alpine.data(
  'Comment',
  ({ id, liked = false, replies = 0 }: { id: string; replies?: number; liked?: boolean }) => ({
    liked,
    replies,
    replyController: null as CommentEditorController | null,
    repliesController: null as CommentRepliesController | null,

    init() {
      this.$watch(
        'liked',
        Alpine.throttle((newValue: boolean) => {
          httpLikeComment(id, newValue)
        }, 1000),
      )

      if (this.replies > 0) {
        const $slot = this.$refs.slotReplies
        if (!($slot instanceof HTMLElement))
          throw new Error('could not find slotReplies ref in Comment component')

        const initialReplies = JSON.parse($slot.dataset.initialReplies ?? '[]')
        const nextCursor = Number($slot.dataset.nextCursor ?? 0)
        this.repliesController = initCommentReplies($slot, id, initialReplies, nextCursor)
      }
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
        this.$el.setAttribute('aria-pressed', `${this.liked}`)
      },
    },

    reply: {
      '@click'() {
        if (this.$el.hasAttribute('data-set')) {
          this.replyController?.close()
          this.replyController = null
          this.$el.removeAttribute('data-set')
          return
        }

        const $slot = this.$refs.slotReply
        if (!($slot instanceof HTMLElement))
          throw new Error('could not find slotReply ref in Comment component')

        this.replyController = initCommentEditor($slot, id)
        this.$el.setAttribute('data-set', 'true')
      },
    },
  }),
)

function getFirstText(el: Element): Text | null {
  for (let i = 0; i < el.childNodes.length; i++) {
    const node = el.childNodes.item(i)
    if (node instanceof Text) return node
  }

  return null
}
