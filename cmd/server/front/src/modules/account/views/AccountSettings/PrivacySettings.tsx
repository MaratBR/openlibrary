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
import { useTranslation } from 'react-i18next'

export default function PrivacySettings() {
  const { t } = useTranslation()

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
          <Label htmlFor="hide-comments">{t('settings.privacy.hideComments')}</Label>
          <p className="text-sm text-muted-foreground">
            {t('settings.privacy.hideCommentsDescription')}
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
          <Label htmlFor="hide-favorites">{t('settings.privacy.hideFavorites')}</Label>
          <p className="text-sm text-muted-foreground">
            {t('settings.privacy.hideFavoritesDescription')}
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
          <Label htmlFor="hide-email">{t('settings.privacy.hideEmail')}</Label>
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
          <Label htmlFor="hide-stats">{t('settings.privacy.hideStats')}</Label>
          <p className="text-sm text-muted-foreground">
            {t('settings.privacy.hideStatsDescription')}
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
          <Label htmlFor="allow-searching">{t('settings.privacy.allowSearch')}</Label>
          <p className="text-sm text-muted-foreground">
            {t('settings.privacy.allowSearchDescription')}
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
