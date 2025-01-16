<script lang="ts">
    export let value: number;
    export let scale: number = 1;
    export let onChange: (value: number) => void;

    let rootElement: HTMLDivElement | null = null;

    function handleClick(event: MouseEvent) {
        if (!rootElement) return;

        const rect = rootElement.getBoundingClientRect();
        const x = event.clientX - rect.left;
        const width = rect.width;
        const newValue = Math.max(Math.min(Math.ceil((x / width) * 10), 10), 1);
        if (newValue !== value) {
            onChange(newValue);
        }
    }

</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div bind:this={rootElement} class="relative cursor-pointer" on:click={handleClick} style={`width:${540*scale}px;height:${100*scale}px`}>
    <div class="ol-star-background h-full w-full opacity-15"></div>
    <div class="absolute left-0 top-0 ol-star-background ol-star-background--filled h-full" style={`width:${value * 10}%;background-size:auto ${scale * 100}px`}></div>
</div>