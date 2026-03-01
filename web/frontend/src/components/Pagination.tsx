import clsx from 'clsx'
import { ComponentChildren, ComponentType, HTMLAttributes, JSX } from 'preact'
import React, { ForwardedRef, forwardRef, useMemo } from 'preact/compat'
import { NavLink } from 'react-router'

const Pagination_Root = forwardRef(
  (
    { children, class: class_, className, ...props }: HTMLAttributes<HTMLElement>,
    ref: React.ForwardedRef<HTMLElement>,
  ) => {
    return (
      <nav ref={ref} role="listbox" class={clsx('pagination', class_, className)} {...props}>
        {children}
      </nav>
    )
  },
)

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type ElementType = keyof JSX.IntrinsicElements | ComponentType<any>

type PropsOf<T extends ElementType> = T extends keyof JSX.IntrinsicElements
  ? JSX.IntrinsicElements[T]
  : T extends ComponentType<infer P>
    ? P
    : never

type PaginationItemProps<T extends ElementType> = {
  as?: T
  children?: ComponentChildren
  active?: boolean
} & Omit<PropsOf<T>, 'class' | 'role' | 'ref'>

const Pagination_Item = forwardRef(
  <T extends ElementType = 'button'>(
    { as, children, active = false, ...rest }: PaginationItemProps<T>,
    ref: ForwardedRef<T>,
  ) => {
    const Component = (as || 'button') as ElementType

    return (
      <Component
        ref={ref}
        role="listbox"
        className={clsx('pagination__item', {
          'pagination__item--active': active,
        })}
        {...rest}
      >
        {children}
      </Component>
    )
  },
) as <T extends ElementType = 'button'>(props: PaginationItemProps<T>) => JSX.Element

type PaginatioFacadeProps = {
  page: number
  totalPages: number
  size: number
  disabled?: boolean
}

function Pagination_Facade({ page, totalPages, size, disabled = false }: PaginatioFacadeProps) {
  const order = useMemo(() => getPaginationOrder(page, totalPages, size), [page, totalPages, size])

  return (
    <Pagination.Root>
      {order.map((p) =>
        p === page ? (
          <Pagination.Item key={`${page}_current`} active as="button">
            {p}
          </Pagination.Item>
        ) : (
          <Pagination.Item to={{ search: `?page=${p}` }} key={`${page}_current`} as={NavLink}>
            {p}
          </Pagination.Item>
        ),
      )}
    </Pagination.Root>
  )
}

export const Pagination = {
  Root: Pagination_Root,
  Item: Pagination_Item,
  Facade: Pagination_Facade,
}

function getPaginationOrder(page: number, totalPages: number, size: number): number[] {
  if (totalPages === 0) return []

  let remaining = size - 1
  let left = Math.floor(remaining / 2)

  if (left > page - 1) {
    left = page - 1
  }

  remaining -= left

  let right = remaining
  if (page + right > totalPages) {
    right = totalPages - page
  }

  remaining -= right

  if (remaining > 0 && left < page - 1) {
    const need = page - 1 - left
    left += Math.min(remaining, need)
  }

  const result: number[] = []

  // first page
  if (page !== 1) {
    result.push(1)
  }

  // left side
  for (let i = page - left + 1; i < page; i++) {
    result.push(i)
  }

  // current page
  result.push(page)

  // right side
  for (let i = page + 1; i <= page + right - 1; i++) {
    result.push(i)
  }

  // last page
  if (page < totalPages) {
    result.push(totalPages)
  }

  return result
}
