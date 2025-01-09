import './button.css';
import React from 'react';
import { Slot } from '@radix-ui/react-slot';
import { cn } from '../utils';

export type ButtonVariant =
  | 'default'
  | 'destructive'
  | 'outline'
  | 'outline2'
  | 'outline3'
  | 'secondary'
  | 'ghost'
  | 'link';

export type ButtonSize = 'default' | 'sm' | 'lg' | 'icon';

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  asChild?: boolean;
  variant?: ButtonVariant;
  size?: ButtonSize;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  (
    {
      className,
      variant = 'default',
      size = 'default',
      asChild = false,
      ...props
    },
    ref,
  ) => {
    const Comp = asChild ? Slot : 'button';
    return (
      <Comp
        className={cn(
          `ol-button ol-button--${variant} ol-button--${size}`,
          className,
        )}
        ref={ref}
        {...props}
      />
    );
  },
);
Button.displayName = 'Button';
