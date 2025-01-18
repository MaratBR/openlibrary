import { PreactIsland } from '../common'
import ReviewEditor from './ReviewEditor.tsx'
import './ReviewEditor.scss'

window.OLIslandsRegistry.instance.register('review-editor', new PreactIsland(ReviewEditor))
