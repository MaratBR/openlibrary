import './UserProfileInner.css'
import { useQuery } from '@tanstack/react-query'
import { httpGetUser, UserDetailsDto } from '../../api'
import exampleCoverUrl from './example-cover.jpg'
import Timestamp from '@/components/timestamp'
import { useQueryParam } from '@/lib/router-utils'
import FollowButton from './FollowButton'
import React from 'react'
import { useAuthState } from '@/modules/auth/state'

export default function UserProfileInner() {
  const [userId] = useQueryParam('userId')

  const { data } = useQuery({
    queryKey: ['user', userId],
    enabled: !!userId,
    queryFn: () => httpGetUser(userId!),
  })

  if (!userId || !data) return null

  return (
    <div className="user-page">
      <UserHeader user={data} />
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

function UserHeader({ user }: { user: UserDetailsDto }) {
  const currentUserId = useAuthState((s) => s.user?.id)
  const [isFollowing, setIsFollowing] = React.useState(user.isFollowing)

  return (
    <header className="profile-header">
      <div className="profile-cover" style={{ backgroundImage: `url(${exampleCoverUrl})` }}></div>
      <div className="profile-top">
        <UserAvatar url={user.avatar.lg} />

        <div className="profile-top__follow">
          {user.id !== currentUserId && (
            <FollowButton
              isFollowing={isFollowing}
              onFollowingChange={setIsFollowing}
              userId={user.id}
            />
          )}
        </div>
      </div>
      <div className="profile-info">
        <h1 className="profile-info__username">{user.name}</h1>
        <p className="profile-info__joined">
          Joined <Timestamp value={user.joinedAt} />
        </p>
        <p className="profile-info__bio">{user.about.bio}</p>

        <div className="user-stats">
          <div className="user-stat">
            <div className="user-stat__value">{user.booksTotal.toLocaleString()}</div>
            <div className="user-stat__label">books</div>
          </div>

          <div className="user-stat">
            <div className="user-stat__value">
              {user.followers < 0 ? '--' : user.followers.toLocaleString()}
            </div>
            <div className="user-stat__label">followers</div>
          </div>

          <div className="user-stat">
            <div className="user-stat__value">
              {user.favorites < 0 ? '--' : user.favorites.toLocaleString()}
            </div>
            <div className="user-stat__label">favorites</div>
          </div>
        </div>
      </div>
    </header>
  )
}
