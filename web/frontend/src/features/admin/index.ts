import './styles/common.scss'
import './styles/forms.scss'
import '@/components/dashboard-shell.scss'
import { ReactIsland } from '@/islands/common/react-island'
import App from './components/App'

export const AdminIsland = new ReactIsland(App)
