import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useBookManager } from "./book-manager-context";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useForm } from "react-hook-form";
import RatingSelect from "@/components/rating-select";
import React from "react";
import { Textarea } from "@/components/ui/textarea";

const formSchema = z.object({
  name: z.string().min(1).max(50),
  rating: z.enum(["?", "G", "PG", "PG-13", "R", "NC-17"]).default("?"),
  tags: z.array(z.string()).max(50),
  summary: z.string().max(1000).optional(),
});

export default function EditBookForm() {
  const { book } = useBookManager();
  const [isEditing, setIsEditing] = React.useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    disabled: !isEditing,
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: book.name,
      rating: book.ageRating,
      summary: "", // TODO
      tags: [], // TODO
    },
  });

  async function onSubmit(values: z.infer<typeof formSchema>) {}

  return (
    <Form {...form}>
      <div className="mb-4">
        {isEditing ? (
          <div className="space-x-2">
            <Button
              type="reset"
              variant="outline"
              onClick={(e) => {
                e.preventDefault();
                setIsEditing(false);
                form.reset();
              }}
            >
              Cancel
            </Button>
            <Button type="submit">Save</Button>
          </div>
        ) : (
          <Button variant="outline" onClick={() => setIsEditing(true)}>
            Edit
          </Button>
        )}
      </div>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4">
        <div className="grid grid-cols-2">
          <div className="space-y-2">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Name of the book</FormLabel>
                  <FormControl>
                    <Input
                      placeholder="For example, UwUAnimeGirl69"
                      {...field}
                    />
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
        </div>
      </form>
    </Form>
  );
}
