import { SvelteIsland } from '../common'
import ReviewEditor from './ReviewEditor.svelte'
import './ReviewEditor.scss';

window.OLIslandsRegistry.instance.register('review-editor', new SvelteIsland(ReviewEditor))