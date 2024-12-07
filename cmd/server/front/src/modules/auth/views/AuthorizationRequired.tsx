import React from 'react'
import { useAuthState } from '../state'
import { Navigate } from 'react-router'

export function AuthorizationRequired({
  children,
  fallback,
}: React.PropsWithChildren<{
  fallback: React.ReactNode | React.ReactNode
}>) {
  const hasUser = useAuthState((s) => !!s.user)

  if (!hasUser) return fallback

  return <>{children}</>
}

export function AuthorizationRequiredRedirect({ children }: React.PropsWithChildren) {
  return (
    <AuthorizationRequired fallback={<Navigate to="/login" />}>{children}</AuthorizationRequired>
  )
}
