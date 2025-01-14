<script lang="ts">
    export let value: number;
    export let scale: number = 1;

</script>

<div class="relative" style={`width:${540*scale}px;height:${100*scale}px`}>
    <div class="ol-star-background h-full w-full opacity-15"></div>
    <div class="absolute left-0 top-0 ol-star-background ol-star-background--filled h-full" style={`width:${value * 10}%`}></div>
</div>