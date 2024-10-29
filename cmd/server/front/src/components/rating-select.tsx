import { AGE_RATINGS_LIST, AgeRating } from "@/modules/book/api";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "./ui/select";
import { useAgeRatingsInfo } from "./age-rating-util";
import AgeRatingBadge from "./age-rating-badge";
import { isAgeRatingAdult } from "@/modules/book/utils";
import AdultIndicator from "./adult-indicator";

export type RatingSelectProps = {
  value?: AgeRating;
  disabled?: boolean;
  onChange?: (value: AgeRating) => void;
};

export default function RatingSelect({
  value,
  onChange,
  disabled = false,
}: RatingSelectProps) {
  const ratings = useAgeRatingsInfo();

  return (
    <Select value={value} onValueChange={onChange} disabled={disabled}>
      <SelectTrigger>
        <SelectValue placeholder="Rating">
          <div className="flex items-center gap-1">
            {isAgeRatingAdult(value ?? "?") && <AdultIndicator />}
            <AgeRatingBadge disableTooltip value={value ?? "?"} />
          </div>
        </SelectValue>
      </SelectTrigger>

      <SelectContent>
        {AGE_RATINGS_LIST.map((rating) => (
          <SelectItem key={rating} value={rating}>
            <div>{ratings[rating].title}</div>
            <p className="max-w-[600px] text-muted-foreground">
              {ratings[rating].summary}
            </p>
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
