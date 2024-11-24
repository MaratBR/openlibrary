import { useQuery } from '@tanstack/react-query'
import { httpGetUserCustomizationSettings } from '../../api'

export default function CustomizationSettings() {
  const { data } = useQuery({
    staleTime: 0,
    gcTime: Infinity,
    queryKey: ['settings', 'customization'],
    queryFn: () => httpGetUserCustomizationSettings(),
  })

  return <pre>{JSON.stringify(data, null, 2)}</pre>
}
