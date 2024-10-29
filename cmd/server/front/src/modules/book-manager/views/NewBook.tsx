import Spinner, { ButtonSpinner } from "@/components/spinner";
import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { PlusIcon } from "lucide-react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { AgeRating } from "../../book/api";
import { httpCreateBook } from "../api";

export default function NewBook() {
  return (
    <main className="container-default">
      <header className="page-header">
        <h1 className="page-header-text">New book</h1>
      </header>

      <NewBookForm />
    </main>
  );
}

const tagName = z.string().min(1).max(255);

const formSchema = z.object({
  name: z.string().min(1).max(255),
  tags: z.array(tagName).max(60),
  ageRating: z.enum(["?", "G", "PG", "PG-13", "R", "NC-17"]),
});

function NewBookForm() {
  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      tags: [],
      ageRating: "?",
    },
  });

  const createBook = useMutation({
    mutationFn: ({
      name,
      ageRating,
      tags,
    }: {
      name: string;
      ageRating: AgeRating;
      tags: string[];
    }) => {
      return httpCreateBook({ name, ageRating, tags });
    },
  });

  const disabled = createBook.isPending;

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <div className="max-w-[400px]">
          <FormField
            control={form.control}
            name="name"
            render={({ field }) => (
              <FormItem>
                <FormLabel>Name of your book</FormLabel>
                <FormControl>
                  <Input
                    disabled={disabled}
                    placeholder="For example, UwUAnimeGirl69"
                    {...field}
                  />
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
  );

  function onSubmit(values: z.infer<typeof formSchema>) {
    createBook.mutate({
      name: values.name,
      ageRating: values.ageRating,
      tags: values.tags,
    });
  }
}
