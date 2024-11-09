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
  id: z.string().min(1).max(255),
})

function ImportBookFromAo3Form() {
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
        </div>
        <Button disabled={disabled} variant="outline" type="submit">
          {createBook.isPending && <ButtonSpinner />}
          Import
        </Button>
      </form>
    </Form>
  )

  function onSubmit(values: z.infer<typeof formSchema>) {
    createBook
      .mutateAsync({
        id: values.id,
      })
      .then(({ id }) => {
        navigate(`/book/${id}`)
      })
  }
}
