import assert from 'node:assert/strict'
import test from 'node:test'

import {
  componentName,
  convertBemToken,
  discoverComponents,
  parseBemToken,
  transformSource,
} from './bem-to-suit.mjs'

test('converts canonical BEM tokens to canonical SUIT names', () => {
  assert.equal(componentName('book-reader'), 'BookReader')
  assert.equal(convertBemToken('book-reader__settings-button'), 'BookReader-settingsButton')
  assert.equal(
    convertBemToken('book-reader__chapter-link--previous'),
    'BookReader-chapterLink--previous',
  )
  assert.equal(convertBemToken('btn--rounded-md'), 'Btn--roundedMd')
})

test('preserves an existing PascalCase component and supports overrides', () => {
  assert.equal(convertBemToken('UserMenu__item-icon'), 'UserMenu-itemIcon')
  assert.equal(
    convertBemToken('ol-comments__heading-row', {
      components: new Map([['ol-comments', 'OLComments']]),
    }),
    'OLComments-headingRow',
  )
})

test('parses element modifiers without treating the modifier as part of the element', () => {
  assert.deepEqual(parseBemToken('btn__icon--left'), {
    block: 'btn',
    element: 'icon',
    modifier: 'left',
  })
})

test('transforms nested and compound SCSS selectors', () => {
  const source = `.carousel {
  &--on-background &__curtain {}
  &__button {
    &--scroll-left {}
  }
  &[data-ready] &__button--scroll-left {}
}`
  const components = discoverComponents([source, 'class="carousel Carousel-button"'])
  assert.equal(
    transformSource(source, '.scss', components),
    `.Carousel {
  &--onBackground &-curtain {}
  &-button {
    &--scrollLeft {}
  }
  &[data-ready] &-button--scrollLeft {}
}`,
  )
})

test('transforms static markup classes but leaves ordinary prose alone', () => {
  const source = `<div class="error error__message"></div>\n<p>error</p>`
  const components = discoverComponents([source])
  assert.equal(
    transformSource(source, '.templ', components),
    `<div class="Error Error-message"></div>\n<p>error</p>`,
  )
})

test('can preserve kebab case for descendants and modifiers', () => {
  assert.equal(
    convertBemToken('book-reader__settings-button--extra-large', { descendantCase: 'kebab' }),
    'BookReader-settings-button--extra-large',
  )
})
