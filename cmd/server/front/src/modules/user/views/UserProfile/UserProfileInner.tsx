import { useParams } from 'react-router'
import './UserProfileInner.css'
import { useQuery } from '@tanstack/react-query'
import { httpGetUser } from '../../api'
import Timestamp from '@/components/timestamp'

export default function UserProfileInner() {
  const { userId } = useParams<{ userId: string }>()

  const { data } = useQuery({
    queryKey: ['user', userId],
    enabled: !!userId,
    queryFn: () => httpGetUser(userId!),
  })

  if (!userId || !data) return null

  return (
    <div className="user-page">
      <div className="user-layout">
        <aside className="user-card">
          <UserAvatar url={data.avatar.lg} />
          <p className="user-card__joined-at">
            Joined at <Timestamp value={data.joinedAt} />
          </p>
        </aside>
        <p>
          Lorem ipsum dolor sit amet consectetur adipisicing elit. Dignissimos neque qui molestiae
          veritatis ratione perferendis quos aliquam ab ex distinctio excepturi, corporis sequi,
          fuga at fugit a earum ad delectus!
        </p>
      </div>
    </div>
  )
}

function UserAvatar({ url }: { url: string }) {
  return (
    <div className="user-avatar">
      <img className="user-avatar__img" src={url} alt="user's avatar" />
    </div>
  )
}
