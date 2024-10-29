import { Card } from "@/components/ui/card";
import { BookDetailsDto } from "../../api";
import {
  LabeledValue,
  LabeledValueLabel,
  LabeledValueLayout,
} from "@/components/labeled-value";
import AgeRatingBadge from "@/components/age-rating-badge";
import { Separator } from "@/components/ui/separator";
import WordsCount from "@/components/words-count";
import Tag from "../Tag";

export default function BookInfoCard({ book }: { book: BookDetailsDto }) {
  book = {
    ...book,
    tags: [
      {
        id: "rw",
        name: "Hermione Granger/Ron Weasley",
        isAdult: false,
        isDefined: true,
        isSpoiler: false,
        category: "relationship",
      },
      {
        id: "rw",
        name: "Hermione Granger/Harry Potter",
        isAdult: false,
        isDefined: true,
        isSpoiler: false,
        category: "relationship",
      },
    ],
  };

  return (
    <Card>
      <div className="space-y-2 py-3">
        <LabeledValueLayout className="px-4 py-1">
          <LabeledValueLabel>Age rating</LabeledValueLabel>
          <LabeledValue>
            <AgeRatingBadge value={book.ageRating} />
          </LabeledValue>
        </LabeledValueLayout>
        <Separator />
        <LabeledValueLayout className="px-4 py-1">
          <LabeledValueLabel>Tags</LabeledValueLabel>
          <LabeledValue>
            <div className="flex flex-wrap gap-2">
              {book.tags.map((tag) => {
                return <Tag key={tag.id} tag={tag} />;
              })}
            </div>
          </LabeledValue>
        </LabeledValueLayout>
        <Separator />
        <LabeledValueLayout className="px-4 py-1">
          <LabeledValueLabel>Words in total</LabeledValueLabel>
          <LabeledValue>
            <WordsCount value={book.words} />
          </LabeledValue>
        </LabeledValueLayout>
        <LabeledValueLayout className="px-4 py-1">
          <LabeledValueLabel>Words per chapter (average)</LabeledValueLabel>
          <LabeledValue>
            <WordsCount value={book.wordsPerChapter} />
          </LabeledValue>
        </LabeledValueLayout>
        <Separator />
        <LabeledValueLayout className="px-4 py-1">
          <LabeledValueLabel>Collections</LabeledValueLabel>
          <LabeledValue>
            <pre>{JSON.stringify(book.collections, null, 2)}</pre>
          </LabeledValue>
        </LabeledValueLayout>
      </div>
    </Card>
  );
}
