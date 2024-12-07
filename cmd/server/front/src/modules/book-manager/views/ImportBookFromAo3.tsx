import { ButtonSpinner } from '@/components/spinner'
import { Button } from '@/components/ui/button'
import {
  Form,
  FormControl,
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
import { httpImportBookFromAo3 } from '../api'
import { useNavigate } from 'react-router'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import { useState } from 'react'
import { Textarea } from '@/components/ui/textarea'

export default function ImportBookFromAo3() {
  return (
    <main className="container-default">
      <header className="page-header">
        <h1 className="page-header-text">New book</h1>
      </header>

      <ImportBookFromAo3Form />
    </main>
  )
}

const formSchema = z.object({
  id: z.string().min(1),
})

function ImportBookFromAo3Form() {
  const [isMultiple, setMultiple] = useState(false)

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      id: '',
    },
  })

  const navigate = useNavigate()

  const createBook = useMutation({
    mutationFn: ({ id }: { id: string }) => {
      return httpImportBookFromAo3({ id })
    },
  })

  const disabled = createBook.isPending

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <div className="max-w-[400px]">
          <div className="flex gap-3 flex-col mb-8">
            <Label htmlFor="auto-apply-filters-switch">Import multiple</Label>
            <Switch checked={isMultiple} onCheckedChange={setMultiple} />
          </div>

          {isMultiple ? (
            <FormField
              control={form.control}
              name="id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>ID of the books on Ao3</FormLabel>
                  <FormControl>
                    <Textarea
                      placeholder="Put each book ID in a new line"
                      rows={10}
                      disabled={disabled}
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          ) : (
            <FormField
              control={form.control}
              name="id"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>ID of the book on Ao3</FormLabel>
                  <FormControl>
                    <Input disabled={disabled} {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
          )}
        </div>
        <Button disabled={disabled} variant="outline" type="submit">
          {createBook.isPending && <ButtonSpinner />}
          Import
        </Button>
      </form>
    </Form>
  )

  async function onSubmit(values: z.infer<typeof formSchema>) {
    if (isMultiple) {
      const failedIds: string[] = []

      const ids = values.id
        .split('\n')
        .map((x) => x.trim())
        .map(Number)
      for (const id of ids) {
        if (Number.isNaN(id)) continue

        try {
          await createBook.mutateAsync({
            id: id + '',
          })
        } catch {
          failedIds.push(id + '')
        }
        await new Promise((r) => setTimeout(r, 500))
      }
      if (failedIds.length > 0) {
        form.setValue('id', failedIds.join('\n'))
      } else {
        navigate('/manager/books')
      }
    } else {
      await createBook
        .mutateAsync({
          id: values.id,
        })
        .then(({ id }) => {
          navigate(`/book/${id}`)
        })
    }
  }
}
