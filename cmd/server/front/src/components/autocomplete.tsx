import { CommandGroup, CommandItem, CommandList } from './ui/command'
import { Command as CommandPrimitive } from 'cmdk'

import { Skeleton } from './ui/skeleton'

import { cn } from '../lib/utils'
import React, { useRef } from 'react'
import { Popover } from './ui/popover'
import { PopoverContent, PopoverTrigger } from '@radix-ui/react-popover'

type AutoCompleteProps<TOption> = {
  options: TOption[]
  getKey: (option: TOption) => React.Key
  getLabel: (option: TOption) => string
  itemComponent: React.ComponentType<{ value: TOption }>
  onInputValueChange?: (value: string) => void
  onValueChange?: (value: TOption) => void
  emptyMessage: string
  isLoading?: boolean
  disabled?: boolean
  placeholder?: string

  onOpen?: () => void
  onClosed?: () => void
}

export const AutoComplete = function <TOption>({
  options,
  placeholder,
  emptyMessage,
  onValueChange,
  onInputValueChange,
  getKey,
  getLabel,
  disabled,
  isLoading = false,
  itemComponent: ItemComponent,
  onClosed,
  onOpen,
}: AutoCompleteProps<TOption>) {
  const inputRef = React.useRef<HTMLInputElement>(null)

  const [isOpen, setOpen] = React.useState(false)
  const [popupWidth, setPopupWidth] = React.useState(0)

  React.useEffect(() => {
    if (isOpen) {
      propsRef.current.onOpen?.()
    } else {
      propsRef.current.onClosed?.()
    }
  }, [isOpen])

  const [inputValue, setInputValue] = React.useState<string>('')
  const propsRef = useRef({ getKey, getLabel, onClosed, onOpen })
  propsRef.current = { getKey, getLabel, onClosed, onOpen }

  const handleKeyDown = React.useCallback(
    (event: React.KeyboardEvent<HTMLDivElement>) => {
      const input = inputRef.current
      if (!input) {
        return
      }

      // Keep the options displayed when the user is typing
      if (!isOpen) {
        setOpen(true)
      }

      // This is not a default behavior of the <input /> field
      if (event.key === 'Enter' && input.value !== '') {
        const optionToSelect = options.find(
          (option) => propsRef.current.getLabel(option) === input.value,
        )
        if (optionToSelect) {
          onValueChange?.(optionToSelect)
        }
      }

      if (event.key === 'Escape') {
        input.blur()
      }
    },
    [isOpen, options, onValueChange],
  )

  const handleBlur = React.useCallback(() => {
    setOpen(false)
  }, [])

  const handleFocus = React.useCallback(() => {
    setOpen(true)
    if (inputRef.current) {
      window.requestAnimationFrame(() => {
        if (inputRef.current) {
          const width = inputRef.current.getBoundingClientRect().width
          setPopupWidth(width)
        }
      })
    }
  }, [])

  const handleSelectOption = React.useCallback(
    (selectedOption: TOption) => {
      setInputValue('')
      onValueChange?.(selectedOption)

      // This is a hack to prevent the input from being focused after the user selects an option
      // We can call this hack: "The next tick"
      // setTimeout(() => {
      //   inputRef?.current?.blur();
      // }, 0);
    },
    [onValueChange],
  )

  const handleInputValueChange = React.useCallback(
    (value: string) => {
      setInputValue(value)
      onInputValueChange?.(value)
    },
    [onInputValueChange],
  )

  return (
    <CommandPrimitive onKeyDown={handleKeyDown}>
      <Popover open={isOpen}>
        <PopoverTrigger asChild>
          <div>
            <CommandPrimitive.Input
              ref={inputRef}
              value={inputValue}
              onValueChange={handleInputValueChange}
              onBlur={handleBlur}
              onFocus={handleFocus}
              placeholder={placeholder}
              disabled={disabled}
              className="flex h-9 w-full rounded-md bg-background px-3 py-2 text-base ring-offset-background transition-colors 
          file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-foreground 
          placeholder:text-muted-foreground 
          hover:bg-muted
          focus-visible:outline-none focus-visible:bg-muted
          disabled:cursor-not-allowed disabled:opacity-50"
            />
          </div>
        </PopoverTrigger>
        <PopoverContent
          align="start"
          className="z-10"
          onOpenAutoFocus={(e) => e.preventDefault()}
          style={{ width: popupWidth ? `${popupWidth}px` : undefined }}
        >
          <div
            className={cn(
              'animate-in fade-in-0 zoom-in-95 z-10 w-full rounded-xl bg-card outline-none',
              // isOpen ? 'block' : 'hidden',
            )}
          >
            <CommandList className="rounded-lg ring-1 ring-border bg-background">
              {isLoading ? (
                <CommandPrimitive.Loading>
                  <div className="p-1">
                    <Skeleton className="h-8 w-full" />
                  </div>
                </CommandPrimitive.Loading>
              ) : null}
              {options.length > 0 && !isLoading ? (
                <CommandGroup>
                  {options.map((option) => {
                    return (
                      <CommandItem
                        key={propsRef.current.getKey(option)}
                        value={propsRef.current.getLabel(option)}
                        onMouseDown={(event) => {
                          event.preventDefault()
                          event.stopPropagation()
                        }}
                        onSelect={() => handleSelectOption(option)}
                        className={cn('flex w-full items-center gap-2 pl-8')}
                      >
                        <ItemComponent value={option} />
                      </CommandItem>
                    )
                  })}
                </CommandGroup>
              ) : null}
              {!isLoading ? (
                <CommandPrimitive.Empty className="select-none rounded-sm px-2 py-3 text-center text-sm">
                  {emptyMessage}
                </CommandPrimitive.Empty>
              ) : null}
            </CommandList>
          </div>
        </PopoverContent>
      </Popover>
    </CommandPrimitive>
  )
}
