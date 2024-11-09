import { ButtonSpinner } from '@/components/spinner'
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
import { zodResolver } from '@hookform/resolvers/zod'
import { useMutation } from '@tanstack/react-query'
import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { AgeRating } from '../../book/api'
import { httpCreateBook } from '../api'
import { Switch } from '@/components/ui/switch'
import { getZodDefaults } from '@/lib/zod-utils'
import { Textarea } from '@/components/ui/textarea'
import RatingSelect from '@/components/rating-select'

export default function NewBook() {
  return (
    <main className="container-default">
      <header className="page-header">
        <h1 className="page-header-text">New book</h1>
      </header>

      <NewBookForm />
    </main>
  )
}

const tagName = z.string().min(1).max(255)

const formSchema = z.object({
  name: z.string().min(1).max(255).default(''),
  tags: z.array(tagName).max(60).default([]),
  summary: z.string().max(1000).default(''),
  rating: z.enum(['?', 'G', 'PG', 'PG-13', 'R', 'NC-17']).default('?'),
  isPubliclyVisible: z.boolean().default(true),
})

function NewBookForm() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: getZodDefaults(formSchema),
  })

  const createBook = useMutation({
    mutationFn: ({
      name,
      ageRating,
      tags,
      summary,
      isPubliclyVisible,
    }: {
      name: string
      ageRating: AgeRating
      tags: string[]
      summary: string
      isPubliclyVisible: boolean
    }) => {
      return httpCreateBook({ name, ageRating, tags, summary, isPubliclyVisible })
    },
  })

  const disabled = createBook.isPending

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <div className="max-w-[500px] space-y-3">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name of your book</FormLabel>
                <FormControl>
                  <Input disabled={disabled} placeholder="For example, UwUAnimeGirl69" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
          <FormField
            control={form.control}
            name="rating"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Age rating</FormLabel>
                <FormControl>
                  <RatingSelect
                    disabled={field.disabled}
                    value={field.value}
                    onChange={field.onChange}
                  />
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
                <FormLabel>Summary</FormLabel>
                <FormControl>
                  <Textarea placeholder="Short description of what this book is about" {...field} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />

          <FormField
            control={form.control}
            name="isPubliclyVisible"
            render={({ field }) => (
              <FormItem className="flex flex-row items-center justify-between">
                <div className="space-y-0.5">
                  <FormLabel>Make your book publicly accessible</FormLabel>
                  <FormDescription>
                    This will allow other users to find your book. You can always change that in the
                    settings later.
                  </FormDescription>
                </div>

                <FormControl>
                  <Switch checked={field.value} onCheckedChange={field.onChange} />
                </FormControl>
                <FormMessage />
              </FormItem>
            )}
          />
        </div>
        <Button disabled={disabled} variant="outline" type="submit">
          {createBook.isPending && <ButtonSpinner />}
          Create
        </Button>
      </form>
    </Form>
  )

  function onSubmit(values: z.infer<typeof formSchema>) {
    createBook.mutate({
      name: values.name,
      ageRating: values.rating,
      tags: values.tags,
      summary: values.summary,
      isPubliclyVisible: values.isPubliclyVisible,
    })
  }
}
