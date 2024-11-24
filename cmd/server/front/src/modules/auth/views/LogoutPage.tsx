import { Button } from '@/components/ui/button'
import { useMutation } from '@tanstack/react-query'
import { httpSignOut } from '../api'
import { ButtonSpinner } from '@/components/spinner'
import { useAuthState } from '../state'

export default function LogoutPage() {
  const state = useAuthState()

  const logout = useMutation({
    mutationFn: () => httpSignOut(),
    onSuccess() {
      state.logout()
      window.location.pathname = '/home?from=logout'
    },
  })

  return (
    <main className="container-default pt-8">
      {state.user ? (
        <>
          <p className="mb-4">
            You are logged in as <b>{state.user.name}</b>
          </p>
          <Button variant="outline" onClick={() => logout.mutate()} disabled={logout.isPending}>
            {logout.isPending && <ButtonSpinner />}
            Log out
          </Button>
        </>
      ) : (
        <>
          <p>You are not logged in</p>
        </>
      )}
    </main>
  )
}
