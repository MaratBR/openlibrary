import { Schema } from 'effect'

export const managerBookChapterDetailsDtoSchema = Schema.Struct({
  id: Schema.String,
  name: Schema.String,
  createdAt: Schema.String,
  words: Schema.Number,
  summary: Schema.String,
  order: Schema.Int,
  content: Schema.String,
  isPubliclyVisible: Schema.Boolean,
  fonts: Schema.Array(Schema.String),
})
