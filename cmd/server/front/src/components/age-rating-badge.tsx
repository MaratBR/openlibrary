import { AgeRating } from "@/modules/book/api";
import clsx from "clsx";
import { Tooltip, TooltipContent, TooltipTrigger } from "./ui/tooltip";
import { useAgeRatingsInfo } from "./age-rating-util";

export type AgeRatingProps = {
  value: AgeRating;
  disableTooltip?: boolean;
  className?: string;
};

export default function AgeRatingBadge({
  value,
  disableTooltip = false,
  className,
}: AgeRatingProps) {
  const ratings = useAgeRatingsInfo();

  if (disableTooltip) {
    return <AgeRatingBadgeNoTooltip className={className} value={value} />;
  }

  return (
    <Tooltip>
      <TooltipTrigger>
        <AgeRatingBadgeNoTooltip className={className} value={value} />
      </TooltipTrigger>
      <TooltipContent className="!animate-none rounded-[4px]">
        <div className="grid grid-cols-[auto_1fr] gap-4 w-[400px] pt-3">
          <AgeRatingBadgeNoTooltip value={value} />
          <div className="pb-3">
            <h3 className="font-semibold">{ratings[value].title}</h3>
            <p>{ratings[value].summary}</p>
          </div>
        </div>
      </TooltipContent>
    </Tooltip>
  );
}

function AgeRatingBadgeNoTooltip({ value, className }: AgeRatingProps) {
  return (
    <div
      className={clsx(
        className,
        {
          "text-white bg-[#006835] w-7": value === "G",
          "text-white bg-[#f15a24] w-9": value === "PG",
          "text-white bg-[#955ea9] w-16": value === "PG-13",
          "text-white bg-[#d8121a] w-7": value === "R",
          "text-white bg-[#1b3e9b] w-16": value === "NC-17",
          "text-white bg-gray-500 w-7": value === "?",
        },
        "font-semibold text-[1.2em] rounded-[5px] h-7 flex flex-col align-middle items-center justify-center"
      )}
    >
      {value}
    </div>
  );
}

function AgeRatingTooltipContent({ value }: { value: AgeRating }) {
  return (
    <div className="w-[300px]">
      <AgeRatingBadgeNoTooltip value={value} />
    </div>
  );
}
