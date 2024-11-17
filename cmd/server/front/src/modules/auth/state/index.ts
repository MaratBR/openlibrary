import { create } from 'zustand'

type AuthState = {
  user: {
    id: string
    name: string
    sessionExpiresAt: string
  } | null

  reset: () => void
}

const useAuthState = create<AuthState>()((set) => ({
  user: SERVER_DATA.session
    ? {
        id: SERVER_DATA.session.id,
        name: SERVER_DATA.session.username,
        sessionExpiresAt: SERVER_DATA.session.expiresAt,
      }
    : null,
  reset() {
    set({ user: null })
  },
}))

export default useAuthState
