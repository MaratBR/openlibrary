import { useQuery } from '@tanstack/react-query'
import {
  censorModeSchema,
  httpGetUserModerationSettings,
  httpUpdateUserModerationSettings,
  UserModerationSettings,
} from '../../api'
import { debounce } from '@/lib/utils'
import { toast } from 'sonner'
import { CheckCircle2 } from 'lucide-react'
import React from 'react'
import { Label } from '@radix-ui/react-label'
import { Switch } from '@/components/ui/switch'
import { Textarea } from '@/components/ui/textarea'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useTranslation } from 'react-i18next'

export default function ModerationSettings() {
  const { t } = useTranslation()

  const [state, setState] = React.useState<UserModerationSettings>({
    showAdultContent: false,
    censoredTags: [],
    censoredTagsMode: 'censor',
  })

  const { data, isFetched } = useQuery({
    staleTime: 0,
    gcTime: Infinity,
    queryKey: ['settings', 'moderation'],
    queryFn: () => httpGetUserModerationSettings(),
  })

  React.useEffect(() => {
    if (data) setState(data)
  }, [data])

  const scheduleUpdate = React.useMemo(
    () =>
      debounce((settings: UserModerationSettings) => {
        httpUpdateUserModerationSettings(settings).then(() =>
          toast(
            <div className="flex gap-3 items-center">
              <CheckCircle2 />
              Moderation settings
            </div>,
          ),
        )
      }, 500),
    [],
  )

  const update = (p: Partial<UserModerationSettings>) => {
    const newState = { ...state, ...p }
    setState(newState)
    scheduleUpdate(newState)
  }

  const disabled = !isFetched

  return (
    <div className="max-w-[600px] space-y-10">
      <div className="flex flex-row items-start justify-between">
        <div className="space-y-0.5">
          <Label htmlFor="show-adult-content">{t('settings.moderation.showAdultContent')}</Label>
          <p className="text-sm text-muted-foreground">
            <span className="whitespace-pre-wrap">
              {t('settings.moderation.showAdultContentDescription')}
            </span>{' '}
            <br /> <br />
            <strong className="font-semibold text-foreground">
              {t('settings.moderation.showAdultContentWarning')}
            </strong>
          </p>
        </div>

        <Switch
          id="show-adult-content"
          disabled={disabled}
          checked={state.showAdultContent}
          onCheckedChange={(value) => update({ showAdultContent: value })}
        />
      </div>

      <div className="space-y-2">
        <div className="space-y-0.5">
          <Label htmlFor="hide-favorites">{t('settings.moderation.blacklistedTags')}</Label>
          <p className="text-sm text-muted-foreground">
            {t('settings.moderation.blacklistedTagsDescription')}
          </p>
        </div>

        <Textarea
          placeholder={t('settings.moderation.blacklistedTagsPlaceholder')}
          rows={10}
          value={state.censoredTags.join('\n')}
          onChange={(e) => update({ censoredTags: e.target.value.split('\n') })}
        />
      </div>

      <div className="space-y-4">
        <div className="space-y-0.5">
          <Label htmlFor="hide-email">{t('settings.moderation.booksCensoring.title')}</Label>
          <p className="text-sm text-muted-foreground">
            {t('settings.moderation.booksCensoring.description')}
          </p>
        </div>

        <Select
          value={state.censoredTagsMode}
          onValueChange={(value) => update({ censoredTagsMode: censorModeSchema.parse(value) })}
        >
          <SelectTrigger className="max-w-64">
            <SelectValue placeholder="Censoring mode">
              {t(`settings.moderation.booksCensoring.${state.censoredTagsMode}`)}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="censor">
              {t('settings.moderation.booksCensoring.censor')}
              <p className="text-sm">{t('settings.moderation.booksCensoring.censorDescription')}</p>
            </SelectItem>
            <SelectItem value="hide">
              {t('settings.moderation.booksCensoring.hide')}
              <p className="text-sm">{t('settings.moderation.booksCensoring.hideDescription')}</p>
            </SelectItem>
            <SelectItem value="none">
              {t('settings.moderation.booksCensoring.none')}
              <p className="text-sm max-w-[32rem]">
                {t('settings.moderation.booksCensoring.noneDescription')}
              </p>
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  )
}
