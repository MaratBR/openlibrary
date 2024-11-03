import { MinimalTiptapEditor } from '@/components/minimal-tiptap'
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
import { useMutation } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { ButtonSpinner } from '@/components/spinner'
import { Textarea } from '@/components/ui/textarea'
import { CreateBookChapterRequest, httpCreateBookChapter } from '../../api'

export type ChapterEditorProps = {
  bookId: string
  chapterId: string | null
}

const formSchema = z.object({
  name: z.string().max(50).trim().optional(),
  content: z.string().min(1).max(50000),
  isAdultOverride: z.boolean().default(false),
  summary: z.string().max(1000).optional(),
})

export default function ChapterEditor({ bookId, chapterId }: ChapterEditorProps) {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
  })

  const createChapterMutation = useMutation({
    mutationFn: (req: CreateBookChapterRequest) => httpCreateBookChapter(bookId, req),
  })

  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (createChapterMutation.isPending) return
    createChapterMutation.mutate({
      name: values.name || '',
      content: values.content,
      isAdultOverride: values.isAdultOverride,
      summary: values.summary || '',
    })
  }

  return (
    <main className="container-default">
      <header className="page-header">
        <h1 className="page-header-text">{chapterId ? '' : 'New chapter'}</h1>
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
                <FormControl>
                  <MinimalTiptapEditor
                    value={field.value}
                    onChange={field.onChange}
                    className="w-full"
                    output="html"
                    editorContentClassName="p-5"
                    editable
                  />
                </FormControl>
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
    </main>
  )
}
