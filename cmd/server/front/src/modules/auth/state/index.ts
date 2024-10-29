import { create } from "zustand";

type AuthState = {
  user: {
    id: string;
    name: string;
    sessionExpiresAt: string;
  } | null;

  reset: () => void;
};

const useAuthState = create<AuthState>()((set) => ({
  user: window.SERVER_DATA.session
    ? {
        id: window.SERVER_DATA.session.id,
        name: window.SERVER_DATA.session.username,
        sessionExpiresAt: window.SERVER_DATA.session.expiresAt,
      }
    : null,
  reset() {
    set({ user: null });
  },
}));

export default useAuthState;
