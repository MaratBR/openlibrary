import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Switch } from '@/components/ui/switch'
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation, useQuery } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { ButtonSpinner } from '@/components/spinner'
import { Textarea } from '@/components/ui/textarea'
import {
  CreateBookChapterRequest,
  createBookChapterRequestSchema,
  httpCreateBookChapter,
  httpManagerGetBookChapter,
  httpUpdateBookChapter,
  UpdateBookChapterRequest,
  updateBookChapterRequestSchema,
} from '@/modules/book-manager/api'
import { useBookManager } from '../book-manager-context'
import { useChapterName } from '@/modules/book/utils'
import React from 'react'
import BackToBookButton from '../BackToBookButton'
import { useMinimalTiptapEditorComponent } from '@/components/minimal-tiptap'

export type ChapterEditorProps = {
  chapterId: string | null
}

const formSchema = z.object({
  name: z.string().max(50).trim().optional(),
  content: z.string().min(1).max(50000),
  isAdultOverride: z.boolean().default(false),
  summary: z.string().max(1000).optional(),
})

export default function ChapterEditor({ chapterId }: ChapterEditorProps) {
  const { book } = useBookManager()

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  })

  useQuery({
    queryKey: ['manager', 'book', book.id, 'chapter', chapterId],
    queryFn: async () => {
      const data = await httpManagerGetBookChapter(book.id, chapterId!)

      form.setValue('name', data.name)
      form.setValue('content', data.content)
      editor?.commands.setContent(data.content)
      form.setValue('isAdultOverride', data.isAdultOverride)
      form.setValue('summary', data.summary)
      setExistingChapterName({ order: data.order, name: data.name })
      return data
    },
    enabled: !!chapterId,
    gcTime: 0,
    staleTime: 0,
  })

  const [existingChapterName, setExistingChapterName] = React.useState({ name: '', order: 0 })
  const chapterName = useChapterName(existingChapterName.name, existingChapterName.order)

  const { editorElement, editor } = useMinimalTiptapEditorComponent({
    editorContentClassName: 'px-4 py-2',
  })
  React.useEffect(() => {
    if (!editor) return

    editor.commands.setContent(form.getValues('content'))
  }, [editor, form])

  const createChapterMutation = useMutation({
    mutationFn: (req: CreateBookChapterRequest) => httpCreateBookChapter(book.id, req),
  })

  const updateChapterMutation = useMutation({
    mutationFn: (req: UpdateBookChapterRequest) => httpUpdateBookChapter(book.id, chapterId!, req),
  })

  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (!editor) {
      throw new Error('WYSIWYG editor is not initialized yet')
    }

    // a hack because the editor does not work well with react-hook-form
    form.setValue('content', editor.getHTML() || '')
    values = form.getValues()

    if (chapterId) {
      updateChapterMutation.mutate(
        updateBookChapterRequestSchema.parse({
          name: values.name || '',
          content: values.content,
          isAdultOverride: values.isAdultOverride,
          summary: values.summary || '',
        } satisfies UpdateBookChapterRequest),
      )
    } else {
      if (createChapterMutation.isPending) return

      createChapterMutation.mutate(
        createBookChapterRequestSchema.parse({
          name: values.name || '',
          content: values.content,
          isAdultOverride: values.isAdultOverride,
          summary: values.summary || '',
        } satisfies CreateBookChapterRequest),
      )
    }
  }

  return (
    <section className="page-section">
      <header className="section-header">
        <BackToBookButton />
        <h1 className="section-header-text">{chapterId ? chapterName : 'New chapter'}</h1>
      </header>

      <Form {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Chapter name </FormLabel>
                <FormControl>
                  <Input {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="summary"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Chapter summary</FormLabel>
                <FormControl>
                  <Textarea
                    placeholder="This will be added at the beginning of the chapter"
                    {...field}
                  />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="isAdultOverride"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between ">
                <div className="space-y-0.5">
                  <FormLabel>Contains adult content</FormLabel>
                  <FormDescription className="max-w-[600px] mr-10">
                    If checked, this chapter will be marked as adults-only. If your book is NOT
                    marked as adults-only, your readers will be informed that only a portion of the
                    content is marked as adult.
                  </FormDescription>
                </div>
                <FormControl>
                  <Switch
                    checked={field.value}
                    onCheckedChange={field.onChange}
                    disabled={field.disabled}
                    ref={field.ref}
                    name={field.name}
                  />
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
                <FormLabel>Contents</FormLabel>
                <FormControl>{editorElement}</FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <Button type="submit" disabled={createChapterMutation.isPending}>
            {createChapterMutation.isPending && <ButtonSpinner />}
            Create
          </Button>
        </form>
      </Form>
    </section>
  )
}
