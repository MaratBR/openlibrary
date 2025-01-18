<script lang="ts">
    import {Editor } from '@tiptap/core';
    import { onDestroy, onMount } from "svelte";
    import RatingEditor from "./RatingEditor.svelte";
    import { createEditor, getExistingReviewData } from './ReviewEditor';
    import { ReviewDto } from './api';
    import { _ } from '@/common/i18n';

    let rootElement: HTMLElement | null = null;
    let editorElement: HTMLElement;
    let editor: Editor | undefined = undefined;

    let loading = $state(false);
    let review: ReviewDto | null = $state(null);
    let rating: number = $state(0);
    let active = $state({ bold: false, italic: false });
    
    onMount(() => {
        const data = getExistingReviewData()
        let content = ''

        if (data) {
            content = data.content;
            rating = data.rating;
            review = data;
        }

        editor = createEditor(editorElement, {
            content,
            onTransaction: () => {
                if (editor) {
                    active = {
                        bold: editor.isActive('bold') === true,
                        italic: editor.isActive('italic') === true,
                    }
                } else {
                    active = {
                        bold: false,
                        italic: false,
                    }
                }  
            },
        })
    })

    onDestroy(() => {
        editor?.destroy()
    })

    function handleRatingChange(newRating: number) {
        rating = newRating;
    }

    async function handleSave() {
        // const content = editor?.getHTML() || '';

        // const review = await httpUpdateReview(data.bookId, {
        //     content,
        //     rating: ratingSchema.parse(rating)
        // })
        // rootElement?.dispatchEvent(new CustomEvent('ol-review-updated', { detail: review }));
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
            class:active={ active.bold }
            onclick={() => editor?.chain().focus().toggleBold().run()}>
            <span class="material-symbols-outlined">
                format_bold
            </span>
        </button>

        <button 
            class:active={ active.italic }
            onclick={() => editor?.chain().focus().toggleItalic().run()}>
            <span class="material-symbols-outlined">
                format_italic
            </span>
        </button>
    </div>

    <div class="ol-review-editor__content __user-content __user-content--editor" bind:this={editorElement}></div>

    <button
        disabled={loading} 
        onclick={handleSave} 
        class="ol-btn ol-btn--primary ol-btn--lg rounded-full mt-4">
        {#if loading}
            <span class="loader loader--dark mx-[32px]"></span>            
        {:else}
            { _('common.save') }
        {/if}
    </button>
</div>