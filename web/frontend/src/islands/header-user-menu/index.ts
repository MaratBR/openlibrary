import { createElement } from 'react'
import { ReactIsland } from '@/islands/common/react-island'
import { HeaderActions } from '@/components/header-actions'

export default new ReactIsland(() => createElement(HeaderActions))
