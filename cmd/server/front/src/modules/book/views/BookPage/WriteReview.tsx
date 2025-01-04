import StarRating from '@/components/star-rating'
import './WriteReview.css'
import { Button } from '@/components/ui/button'
import { useTranslation } from 'react-i18next'
import { ArrowRightIcon, X } from 'lucide-react'
import { z } from 'zod'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Form, FormControl, FormField, FormItem, FormMessage } from '@/components/ui/form'
import { useCallback, useRef, useState } from 'react'
import { useMinimalTiptapEditorComponent } from '@/components/minimal-tiptap'
import { cn } from '@/lib/utils'
import sanitize from 'sanitize-html'
import { CreateReviewRequest, httpUpdateReview, ratingSchema, RatingValue } from '../../api'
import { useMutation } from '@tanstack/react-query'

const formSchema = z.object({
  rating: ratingSchema.nullable(),
  content: z.string(),
})

export default function WriteReview({ bookId }: { bookId: string }) {
  const { t } = useTranslation()
  const [active, setActive] = useState(false)

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {},
  })

  const updateReview = useMutation({
    mutationFn: (request: CreateReviewRequest) => httpUpdateReview(bookId, request),
  })

  function onSubmit(values: z.infer<typeof formSchema>) {
    if (updateReview.isPending) {
      return
    }
    if (values.rating === null) {
      return
    }
    const content = htmlToText(values.content)
    if (content === '') {
      return
    }
    updateReview.mutate({
      content: values.content,
      rating: values.rating,
    })
  }

  return (
    <section
      className={cn('book-write-review', {
        'book-write-review__active': active,
      })}
    >
      {!active && (
        <Button
          onClick={() => setActive(true)}
          variant="outline3"
          className="rounded-full"
          size="lg"
        >
          {t('book.review.writeReview')}
        </Button>
      )}

      {active && (
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
            <FormField
              control={form.control}
              name="rating"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <StarRatingInput value={field.value} onChange={field.onChange} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              name="content"
              render={({ field }) => (
                <FormItem>
                  <FormControl>
                    <ReviewTextInput value={field.value} onChange={field.onChange} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div className="flex mt-3 gap-2">
              <Button className="rounded-full group pl-6" type="submit">
                {t('book.review.submit')}
                <ArrowRightIcon className="transition-all opacity-0 ml-[-12px] group-hover:opacity-100 group-hover:ml-0" />
              </Button>
            </div>

            <button
              onClick={() => {
                setActive(false)
                form.reset()
              }}
              className="book-write-review__close"
            >
              <X />
            </button>
          </form>
        </Form>
      )}
    </section>
  )
}

type StarRatingInputProps = {
  value: RatingValue | null
  onChange: (value: number) => void
}

function StarRatingInput({ value, onChange }: StarRatingInputProps) {
  const propsRef = useRef({ value, onChange })
  propsRef.current = { value, onChange }

  const handleClick = useCallback((e: React.MouseEvent<HTMLDivElement>) => {
    if (e.target instanceof HTMLDivElement) {
      const rect = e.target.getBoundingClientRect()
      const x = e.clientX - rect.left
      const width = rect.width
      const value = Math.max(Math.min(Math.ceil((x / width) * 10), 10), 1)
      if (value != propsRef.current.value) {
        propsRef.current.onChange(value)
      }
    }
  }, [])

  return (
    <StarRating
      className="[&>*]:pointer-events-none cursor-pointer"
      onClick={handleClick}
      value={value ? value / 2 : 0}
    />
  )
}

function ReviewTextInput({
  value,
  onChange,
}: {
  value: string
  onChange: (value: string) => void
}) {
  const { editorElement } = useMinimalTiptapEditorComponent({
    value,
    onChange: (content) => {
      if (typeof content === 'string') {
        onChange(content)
      }
    },
    editorContentClassName: 'm-2 [&>.ProseMirror]:min-h-[239px]',
    output: 'html',
    extensions: {
      disableImage: true,
      disabledColor: true,
      disableLink: true,
      disableHeadings: true,
    },
  })

  return editorElement
}

function htmlToText(html: string): string {
  if (html === '' || html === '<p></p>') {
    return ''
  }

  const safeHtml = sanitize(html)

  const doc = document.createElement('div')
  doc.innerHTML = safeHtml

  const brTags = doc.querySelectorAll('br')
  brTags.forEach((br) => br.replaceWith('\n'))

  const pTags = doc.querySelectorAll('p')
  pTags.forEach((p) => p.appendChild(document.createTextNode('\n')))

  let text = doc.textContent || doc.innerText || ''

  // Clean up whitespace
  text = text
    .replace(/\s+/g, ' ') // Replace multiple spaces with single space
    .replace(/^\s+|\s+$/g, '') // Trim start and end
    .replace(/\n\s*\n/g, '\n\n') // Normalize multiple newlines

  return text
}
