<script lang="ts">
    import {Editor } from '@tiptap/core';
    import { onDestroy, onMount } from "svelte";
    import RatingEditor from "./RatingEditor.svelte";
    import StarterKit from '@tiptap/starter-kit';
    import TextStyle from '@tiptap/extension-text-style';
    import Typography from '@tiptap/extension-typography';
    import HorizontalRule from '@tiptap/extension-horizontal-rule';
    import { httpUpdateReview, ratingSchema, ReviewDto, reviewDtoSchema } from './api';

    let rootElement: HTMLElement | null = null;
    let editorElement: HTMLElement;
    let editor: Editor | undefined = $state(undefined);

    let rating: number = $state(0);

    let {
        data
    }: {
        data: { review: ReviewDto, bookId: string }
    } = $props();

    onMount(() => {
        reviewDtoSchema.parse(data.review);
        rating = data.review.rating;

        editor = new Editor({
            element: editorElement,
            content: data.review.content,
			onTransaction: () => {
				// force re-render so `editor.isActive` works as expected
				editor = editor;
			},
            extensions: [
                StarterKit.configure({
                    horizontalRule: false,
                    codeBlock: false,
                    // paragraph: { HTMLAttributes: { class: 'text-node' } },
                    heading: false,
                    // blockquote: { HTMLAttributes: { class: 'block-node' } },
                    // bulletList: { HTMLAttributes: { class: 'list-node' } },
                    // orderedList: { HTMLAttributes: { class: 'list-node' } },
                    code: { HTMLAttributes: { class: 'inline', spellcheck: 'false' } },
                    dropcursor: { width: 2, class: 'ProseMirror-dropcursor border' },
                }),
                TextStyle,
                // Selection,
                Typography,
                // UnsetAllMarks,
                HorizontalRule,
                // ResetMarksOnEnter,
                // CodeBlockLowlight,
                // Placeholder.configure({ placeholder: () => placeholder }),

            ]
        })
    })

    onDestroy(() => {
        editor?.destroy()
    })

    function handleRatingChange(newRating: number) {
        rating = newRating;
    }

    async function handleSave() {
        const content = editor?.getHTML() || '';

        const review = await httpUpdateReview(data.bookId, {
            content,
            rating: ratingSchema.parse(rating)
        })
        rootElement?.dispatchEvent(new CustomEvent('ol-review-updated', { detail: review }));
    }

    function handleCloseClick() {
        // dispatch dispose event to parent
        rootElement?.dispatchEvent(new CustomEvent('island-request-dispose', { bubbles: true }))
    }

</script>

<div bind:this={rootElement} class="ol-review-editor">
    <RatingEditor scale={0.4} value={rating} onChange={handleRatingChange} />

    <button onclick={handleCloseClick} class="ol-review-editor__close ol-btn ol-btn--icon ol-btn--ghost">
        <span class="material-symbols-outlined">close</span>
    </button>

    <div class="ol-review-editor__toolbar mt-4">
        <button 
            class:active={ editor?.isActive('bold') }
            onclick={() => editor?.chain().focus().toggleBold().run()}>
            <span class="material-symbols-outlined">
                format_bold
            </span>
        </button>

        <button 
            class:active={ editor?.isActive('italic') }
            onclick={() => editor?.chain().focus().toggleItalic().run()}>
            <span class="material-symbols-outlined">
                format_italic
            </span>
        </button>
    </div>

    <div class="ol-review-editor__content" bind:this={editorElement}></div>

    <button onclick={handleSave} class="ol-btn ol-btn--primary ol-btn--lg rounded-full mt-4">
        Save
    </button>
</div>