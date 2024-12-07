import { useQuery } from '@tanstack/react-query'
import {
  httpGetUserPrivacySettings,
  httpUpdateUserPrivacySettings,
  UserPrivacySettings,
} from '../../api'
import React from 'react'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { debounce } from '@/lib/utils'
import { toast } from 'sonner'
import { CheckCircle2 } from 'lucide-react'

export default function PrivacySettings() {
  const [state, setState] = React.useState<UserPrivacySettings>({
    hideComments: false,
    hideEmail: false,
    hideFavorites: false,
    hideStats: false,
    allowSearching: false,
  })

  const { data, isFetched } = useQuery({
    staleTime: 0,
    gcTime: Infinity,
    queryKey: ['settings', 'privacy'],
    queryFn: () => httpGetUserPrivacySettings(),
  })

  React.useEffect(() => {
    if (data) setState(data)
  }, [data])

  const scheduleUpdate = React.useMemo(
    () =>
      debounce((settings: UserPrivacySettings) => {
        httpUpdateUserPrivacySettings(settings).then(() =>
          toast(
            <div className="flex gap-3 items-center">
              <CheckCircle2 />
              Privacy settings updated
            </div>,
          ),
        )
      }, 1000),
    [],
  )

  const update = (p: Partial<UserPrivacySettings>) => {
    const newState = { ...state, ...p }
    setState(newState)
    scheduleUpdate(newState)
  }

  const disabled = !isFetched

  return (
    <div className="max-w-[600px] space-y-10">
      <div className="flex flex-row items-start justify-between">
        <div className="space-y-0.5">
          <Label htmlFor="hide-comments">Hide your comments from your profile</Label>
          <p className="text-sm text-muted-foreground">
            Hides your comments from your profile. People can still see your comments in
            books/chapter where you posted them.
          </p>
        </div>

        <Switch
          id="hide-comments"
          disabled={disabled}
          checked={state.hideComments}
          onCheckedChange={(value) => update({ hideComments: value })}
        />
      </div>

      <div className="flex flex-row items-start justify-between">
        <div className="space-y-0.5">
          <Label htmlFor="hide-favorites">Hide favorites</Label>
          <p className="text-sm text-muted-foreground">
            Hides list and number of your favorite books from your profile.
          </p>
        </div>

        <Switch
          id="hide-favorites"
          disabled={disabled}
          checked={state.hideFavorites}
          onCheckedChange={(value) => update({ hideFavorites: value })}
        />
      </div>

      <div className="flex flex-row items-start justify-between">
        <div className="space-y-0.5">
          <Label htmlFor="hide-email">Hide your email from your profile</Label>
        </div>

        <Switch
          id="hide-email"
          disabled={disabled}
          checked={state.hideEmail}
          onCheckedChange={(value) => update({ hideEmail: value })}
        />
      </div>

      <div className="flex flex-row items-start justify-between">
        <div className="space-y-0.5">
          <Label htmlFor="hide-stats">Hide stats</Label>
          <p className="text-sm text-muted-foreground">
            Hides your account stats (number of followers and followed users).
          </p>
        </div>

        <Switch
          id="hide-stats"
          disabled={disabled}
          checked={state.hideStats}
          onCheckedChange={(value) => update({ hideStats: value })}
        />
      </div>

      <div className="flex flex-row items-start justify-between">
        <div className="space-y-0.5">
          <Label htmlFor="allow-searching">Allow to search for your account</Label>
          <p className="text-sm text-muted-foreground">
            Allows to search for your account. If this setting is enabled then you can be found
            through search. If disabled, your account can only be found directly.
          </p>
        </div>

        <Switch
          id="allow-searching"
          disabled={disabled}
          checked={state.allowSearching}
          onCheckedChange={(value) => update({ allowSearching: value })}
        />
      </div>
    </div>
  )
}
