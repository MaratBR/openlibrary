import { useMutation, useQuery } from '@tanstack/react-query'
import {
  httpGetUserAboutSettings,
  httpUpdateUserAboutSettings,
  UserAboutSettings,
  userAboutSettingsSchema,
} from '../../api'
import React from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from '@/components/ui/form'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { ButtonSpinner } from '@/components/spinner'
import GenderInput from './GenderInput'

export default function AboutMeSettings() {
  const [ready, setReady] = React.useState(false)
  const form = useForm<z.infer<typeof userAboutSettingsSchema>>({
    resolver: zodResolver(userAboutSettingsSchema),
    defaultValues: {
      about: '',
      gender: '',
    },
  })

  const { isPending } = useQuery({
    staleTime: 0,
    gcTime: Infinity,
    queryKey: ['settings', 'about'],
    queryFn: () =>
      httpGetUserAboutSettings().then((r) => {
        form.setValue('about', r.about)
        form.setValue('gender', r.gender)
        setReady(true)
      }),
  })

  const update = useMutation({
    mutationFn: (settings: UserAboutSettings) => httpUpdateUserAboutSettings(settings),
    mutationKey: ['settings', 'about'],
  })

  function onSubmit(values: z.infer<typeof userAboutSettingsSchema>) {
    update.mutate(values)
  }

  if (!ready) {
    return null
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 max-w-[600px]">
        <FormField
          control={form.control}
          name="about"
          render={({ field }) => (
            <FormItem>
              <FormLabel>About</FormLabel>
              <FormControl>
                <Textarea {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <FormField
          control={form.control}
          name="gender"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Gender</FormLabel>
              <FormControl>
                <GenderInput {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button disabled={update.isPending} type="submit">
          {update.isPending && <ButtonSpinner />}
          Save
        </Button>
      </form>
    </Form>
  )
}
