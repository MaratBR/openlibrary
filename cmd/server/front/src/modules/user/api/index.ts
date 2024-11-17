import { httpClient } from '@/modules/common/api'

export type UserDetailDto = {
  id: string
  name: string
  avatar: {
    lg: string
    md: string
  }
  joinedAt: string
  isBlocked: boolean
  isAdmin: boolean
  hasCustomTheme: boolean
  about: {
    status: string
    bio: string
    gender: string
  }
}

export function httpGetUser(id: string): Promise<UserDetailDto> {
  return httpClient.get(`/api/users/${id}`).then((r) => r.json())
}
