import { SvelteIsland } from '../common'
import ReviewEditor from './ReviewEditor.svelte'
import './ReviewEditor.scss';

const reviewEditorIsland = new SvelteIsland(ReviewEditor)

export default reviewEditorIsland;