import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useBookManager, useBookManagerUpdateMutation } from './book-manager-context'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import RatingSelect from '@/components/rating-select'
import React from 'react'
import { Textarea } from '@/components/ui/textarea'
import TagsField from '@/modules/book/components/tags-field'
import { useForm } from 'react-hook-form'
import { ButtonSpinner } from '@/components/spinner'
import { definedTagDtoSchema } from '@/modules/book/api'

const formSchema = z.object({
  name: z.string().min(1).max(50),
  rating: z.enum(['?', 'G', 'PG', 'PG-13', 'R', 'NC-17']).default('?'),
  tags: z.array(definedTagDtoSchema).min(0).max(50),
  summary: z.string().max(1000).default(''),
})

export default function EditBookForm() {
  const { book } = useBookManager()
  const [isEditing, setIsEditing] = React.useState(false)

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: book.name,
      rating: book.ageRating,
      summary: '', // TODO
      tags: book.tags, // TODO
    },
  })

  const updateBook = useBookManagerUpdateMutation()
  const disableFields = !isEditing || updateBook.isPending

  async function onSubmit(values: z.infer<typeof formSchema>) {
    updateBook.mutate({
      name: values.name,
      ageRating: values.rating,
      tags: values.tags.map((x) => x.name),
      summary: values.summary,
    })
  }

  function startEditing() {
    setIsEditing(true)
  }

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit, () => console.log(form.getValues()))}
        className="space-y-4"
      >
        <div className="mb-4">
          {isEditing ? (
            <div className="space-x-2">
              <Button
                type="reset"
                variant="outline"
                onClick={(e) => {
                  e.preventDefault()
                  setIsEditing(false)
                  form.reset()
                }}
              >
                Cancel
              </Button>
              <Button type="submit">
                {updateBook.isPending && <ButtonSpinner />}
                Save
              </Button>
            </div>
          ) : (
            <Button
              variant="outline"
              onClick={(e) => {
                e.preventDefault()
                startEditing()
              }}
            >
              Edit
            </Button>
          )}
        </div>

        <div className="grid grid-cols-2 gap-x-5">
          <div className="space-y-2">
            <FormField
              disabled={disableFields}
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name of the book</FormLabel>
                  <FormControl>
                    <Input placeholder="For example, UwUAnimeGirl69" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <FormField
              control={form.control}
              disabled={disableFields}
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
              disabled={disableFields}
              name="summary"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Summary</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Short description of what this book is about"
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
          <div className="space-y-2">
            <FormField
              control={form.control}
              disabled={disableFields}
              name="tags"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Tags</FormLabel>
                  <FormControl>
                    <div className="p-2 border border-input rounded-md">
                      <TagsField
                        disabled={field.disabled}
                        value={field.value}
                        onChange={field.onChange}
                      />
                    </div>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          </div>
        </div>
      </form>
    </Form>
  )
}
