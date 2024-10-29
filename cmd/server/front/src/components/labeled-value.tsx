import clsx from "clsx";
import React, { HTMLAttributes } from "react";

export type LabeledValueProps = HTMLAttributes<HTMLDivElement>;

export const LabeledValueLayout = React.forwardRef(
  (
    { className, ...props }: LabeledValueProps,
    ref: React.ForwardedRef<HTMLDivElement>
  ) => {
    return (
      <div
        ref={ref}
        className={clsx(
          "grid grid-rows-[1.5em_1fr] md:grid-rows-1 md:grid-cols-[200px_1fr] gap-3",
          className
        )}
        {...props}
      />
    );
  }
);

export const LabeledValueLabel = React.forwardRef(
  (
    { className, ...props }: LabeledValueProps,
    ref: React.ForwardedRef<HTMLDivElement>
  ) => {
    return <div ref={ref} className={clsx("", className)} {...props} />;
  }
);

export const LabeledValue = React.forwardRef(
  (
    { className, ...props }: LabeledValueProps,
    ref: React.ForwardedRef<HTMLDivElement>
  ) => {
    return <div ref={ref} className={clsx("", className)} {...props} />;
  }
);
