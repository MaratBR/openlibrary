import { SelfUserDto } from '@/modules/users/api'
import { create } from 'zustand'

type AuthState = {
  user: SelfUserDto | null

  reset: () => void
  init: (user: SelfUserDto | null) => void
  logout(): void
}

export const useAuthState = create<AuthState>()((set) => ({
  user: __server__.user ?? null,
  reset() {
    set({ user: null })
  },
  init(user: SelfUserDto | null) {
    set({ user })
  },
  logout() {
    set({ user: null })
  },
}))

export function useCurrentUser() {
  return useAuthState((x) => x.user)
}
