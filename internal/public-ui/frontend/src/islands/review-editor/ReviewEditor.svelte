<script lang="ts">
    import {Editor } from '@tiptap/core';
    import { onDestroy, onMount } from "svelte";
    import RatingEditor from "./RatingEditor.svelte";
    import StarterKit from '@tiptap/starter-kit';
    import TextStyle from '@tiptap/extension-text-style';
    import Typography from '@tiptap/extension-typography';
    import HorizontalRule from '@tiptap/extension-horizontal-rule';

    let editorElement: HTMLElement;

    let editor: Editor | undefined
    export let data: unknown;

    onMount(() => {
        editor = new Editor({
            element: editorElement,
            content: (data as any).content,
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

</script>

<div class="ol-review-editor">
    <RatingEditor scale={0.4} value={6} />

    <button on:click={(e) => e.target?.dispatchEvent(new CustomEvent('island-request-dispose', { bubbles: true }))} class="ol-review-editor__close ol-btn ol-btn--icon ol-btn--ghost">
        <span class="material-symbols-outlined">close</span>
    </button>

    <div class="ol-review-editor__toolbar mt-4">
        <button 
            class:active={ editor?.isActive('bold') }
            on:click={() => editor?.chain().focus().toggleBold().run()}>
            <span class="material-symbols-outlined">
                format_bold
            </span>
        </button>

        <button 
            class:active={ editor?.isActive('italic') }
            on:click={() => editor?.chain().focus().toggleItalic().run()}>
            <span class="material-symbols-outlined">
                format_italic
            </span>
        </button>
    </div>

    <div class="ol-review-editor__content" bind:this={editorElement}></div>

    <button class="ol-btn ol-btn--primary ol-btn--lg rounded-full mt-4">
        Save
    </button>
</div>