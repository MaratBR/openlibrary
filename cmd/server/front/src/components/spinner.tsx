import clsx from "clsx";
import React from "react";

export type SpinnerProps = {
  className?: string;
  style?: React.CSSProperties;
  size?: number;
};

export default function Spinner({ style, className, size = 24 }: SpinnerProps) {
  return (
    <div
      style={{ borderWidth: size / 8, height: size, width: size, ...style }}
      className={clsx(
        "animate-spin inline-block  border-[0px] border-current border-t-transparent text-primary rounded-full",
        className
      )}
      role="status"
      aria-label="loading"
    >
      <span className="sr-only">Loading...</span>
    </div>
  );
}

export function ButtonSpinner() {
  return <Spinner size={16} className="text-foreground opacity-80" />;
}
