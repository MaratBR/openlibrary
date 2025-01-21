import { PreactIsland } from '../common'
import SearchFilters from './SearchFilters'

const island = new PreactIsland(SearchFilters)

window.OLIslandsRegistry.instance.register('search-filters', island)
