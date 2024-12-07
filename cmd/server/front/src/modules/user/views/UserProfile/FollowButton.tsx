import { ButtonSpinner } from '@/components/spinner'
import { Button } from '@/components/ui/button'
import { useMutation } from '@tanstack/react-query'
import clsx from 'clsx'
import React from 'react'
import { httpFollowUser, httpUnfollowUser } from '../../api'
import { toast } from 'sonner'
import { useAuthState } from '@/modules/auth/state'
import { useNavigate } from 'react-router'

type FollowButtonProps = {
  userId: string
  isFollowing: boolean
  onFollowingChange: (isFollowing: boolean) => void
}

export default function FollowButton({
  userId,
  isFollowing,
  onFollowingChange,
}: FollowButtonProps) {
  const follow = useMutation({
    mutationFn: async (value: boolean) => {
      if (!isLoggedIn) return
      if (value) {
        await httpFollowUser(userId)
      } else {
        await httpUnfollowUser(userId)
      }
      onFollowingChange(value)
    },
    onError: () => {
      toast('Failed to follow user')
    },
  })

  const isLoggedIn = useAuthState((s) => !!s.user)
  const navigate = useNavigate()

  return (
    <Button
      onClick={() => {
        if (isLoggedIn) {
          follow.mutate(!isFollowing)
        } else {
          navigate('/login')
        }
      }}
      className={clsx('rounded-full transition-none')}
      variant={isFollowing ? 'secondary' : 'default'}
      disabled={follow.isPending}
    >
      {follow.isPending && <ButtonSpinner />}
      {isFollowing ? 'You are following' : 'Follow'}
    </Button>
  )
}
