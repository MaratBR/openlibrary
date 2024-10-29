import clsx from "clsx";

export type AdultIndicatorProps = {
  className?: string;
};

export default function AdultIndicator({ className }: AdultIndicatorProps) {
  return (
    <div
      className={clsx(
        "font-bold px-[0.25em] py-[0.125em] bg-red-800 text-white rounded-sm text-[1.2em] h-7 inline-flex items-center justify-center",
        className
      )}
    >
      Adult
    </div>
  );
}
