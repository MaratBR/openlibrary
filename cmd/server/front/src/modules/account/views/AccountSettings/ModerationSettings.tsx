import { useQuery } from '@tanstack/react-query'
import { httpGetUserModerationSettings } from '../../api'

export default function ModerationSettings() {
  const { data } = useQuery({
    staleTime: 0,
    gcTime: Infinity,
    queryKey: ['settings', 'moderation'],
    queryFn: () => httpGetUserModerationSettings(),
  })

  return <pre>{JSON.stringify(data, null, 2)}</pre>
}
