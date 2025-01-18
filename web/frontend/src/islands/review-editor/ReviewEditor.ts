import StarterKit from '@tiptap/starter-kit';
import TextStyle from '@tiptap/extension-text-style';
import Typography from '@tiptap/extension-typography';
import HorizontalRule from '@tiptap/extension-horizontal-rule';
import { Editor, EditorOptions } from '@tiptap/core';
import { ReviewDto, reviewDtoSchema } from './api';

export function calcPerc(value: number): number {
  return (500 * (value / 10) + Math.floor(value / 2) * 10) / 540 * 100;
}

export function createEditor(editorElement: HTMLElement, options?: Partial<EditorOptions>) {
  return new Editor({
      element: editorElement,
      content: '',
      extensions: [
          StarterKit.configure({
              horizontalRule: false,
              codeBlock: false,
              heading: false,
              code: { HTMLAttributes: { class: 'inline', spellcheck: 'false' } },
              dropcursor: { width: 2, class: 'ProseMirror-dropcursor border' },
          }),
          TextStyle,
          Typography,
          HorizontalRule,
      ],
      ...options
  })
}

/**
 * Finds a hidden element containing review data in JSON format and parses it
 * into a ReviewDto object using `reviewDtoSchema`. If the element is not found,
 * null is returned.
 *
 * This function is used to pre-fill the review text area with existing review
 * data when the review editor is rehydrated on the server.
 *
 * @returns {ReviewDto | null}
 */

export function getExistingReviewData(): ReviewDto | null {
  const el = document.getElementById('island-review-editor-data');

  if (el instanceof HTMLTemplateElement) {
    const data = JSON.parse(el.content.textContent || '')
    return reviewDtoSchema.parse(data);
  }

  return null;
}