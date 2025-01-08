import SearchFilters from './SearchFilters'
import SearchResults from './SearchResult'
import './SearchPage.css'
import { Search } from 'lucide-react'
import { useSearchState } from './search-params'

export default function SearchPage() {
  const tags = useSearchState((s) => s.activeFilters.include.tags)

  return (
    <main className="container-default">
      <header className="page-header">
        <h1 className="page-header-text flex items-center gap-2">
          <Search />
          <span>
            Search
            {tags.length === 1 && <span className="text-muted-foreground">: {tags[0].name}</span>}
          </span>
        </h1>
      </header>
      <div className="search-page-grid">
        <div className="search-filters">
          <SearchFilters />
        </div>
        <div className="search-results">
          <SearchResults />
        </div>
      </div>
    </main>
  )
}
