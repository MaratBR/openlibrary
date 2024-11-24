import { useForm } from 'react-hook-form'
import { z } from 'zod'
import { zodResolver } from '@hookform/resolvers/zod'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { useMutation } from '@tanstack/react-query'
import { httpSignUp, SignUpRequest } from '../api'
import { ButtonSpinner } from '@/components/spinner'

const formSchema = z.object({
  username: z.string().min(1).max(50),
  password: z.string().min(1).max(50),
})

export default function SignUpForm() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: '',
      password: '',
    },
  })

  const signUpMutation = useMutation({
    mutationFn: (req: SignUpRequest) => {
      return httpSignUp(req)
    },
  })

  function onSubmit(values: z.infer<typeof formSchema>) {
    signUpMutation.mutate({
      username: values.username,
      password: values.password,
    })
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <FormField
          control={form.control}
          name="username"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Username</FormLabel>
              <FormControl>
                <Input placeholder="For example, UwUAnimeGirl69" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Password</FormLabel>
              <FormControl>
                <Input type="password" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" disabled={signUpMutation.isPending}>
          {signUpMutation.isPending && <ButtonSpinner />}
          Sign up
        </Button>
      </form>
    </Form>
  )
}
