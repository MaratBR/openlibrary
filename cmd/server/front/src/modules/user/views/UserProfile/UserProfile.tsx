import { useParams } from 'react-router'

export default function UserProfile() {
  const { userId } = useParams<{ userId: string }>()

  const id = `iframe${userId?.replace(/-/g, '')}`

  return <iframe id={id} className="w-full" src={`/user/__profile?userId=${userId}`} />
}
