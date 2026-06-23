import * as React from 'react'
import { cn } from './cn'

export interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'default' | 'secondary' | 'ghost' | 'outline' | 'destructive'
  size?: 'sm' | 'md' | 'lg' | 'icon'
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = 'default', size = 'md', ...props }, ref) => {
    return (
      <button
        className={cn(
          'inline-flex items-center justify-center rounded-md text-sm font-medium transition-colors',
          'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
          'disabled:pointer-events-none disabled:opacity-50',
          {
            'bg-neutral-900 text-white hover:bg-neutral-800 dark:bg-neutral-100 dark:text-neutral-900 dark:hover:bg-neutral-200':
              variant === 'default',
            'bg-neutral-100 text-neutral-900 hover:bg-neutral-200 dark:bg-neutral-800 dark:text-neutral-100 dark:hover:bg-neutral-700':
              variant === 'secondary',
            'hover:bg-neutral-100 dark:hover:bg-neutral-800': variant === 'ghost',
            'border border-neutral-200 bg-transparent hover:bg-neutral-100 dark:border-neutral-700 dark:hover:bg-neutral-800':
              variant === 'outline',
            'bg-red-600 text-white hover:bg-red-700 dark:bg-red-700 dark:hover:bg-red-800':
              variant === 'destructive',
          },
          {
            'h-7 px-2 text-xs': size === 'sm',
            'h-9 px-4': size === 'md',
            'h-11 px-8': size === 'lg',
            'h-9 w-9': size === 'icon',
          },
          className,
        )}
        ref={ref}
        {...props}
      />
    )
  },
)
Button.displayName = 'Button'

export { Button }
