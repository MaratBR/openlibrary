import { OLIsland } from './ol-island';

export const ISLANDS: Record<string, () => Promise<{ default: OLIsland }>> = {
  'review-editor': () => import('../islands/review-editor')
}