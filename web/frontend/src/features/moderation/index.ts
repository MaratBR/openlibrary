import './styles'
import '@/components/dashboard-shell.scss'
import { ReactIsland } from '@/islands/common/react-island'
import Portal from './components/Portal'

export const ModerationPortalIsland = new ReactIsland(Portal)
