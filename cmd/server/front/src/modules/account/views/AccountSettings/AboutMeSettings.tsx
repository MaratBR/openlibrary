import { useQuery } from '@tanstack/react-query'
import { httpGetUserAboutSettings } from '../../api'

export default function AboutMeSettings() {
  const { data } = useQuery({
    staleTime: 0,
    gcTime: Infinity,
    queryKey: ['settings', 'about'],
    queryFn: () => httpGetUserAboutSettings(),
  })

  return <pre>{JSON.stringify(data, null, 2)}</pre>
}
