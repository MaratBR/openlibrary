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

export default function ModerationSettings() {
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
          <Label htmlFor="show-adult-content">Show adult content</Label>
          <p className="text-sm text-muted-foreground">
            Show adult content without warnings. <br />
            If this is set, adult content marked as 18+ will be shown to you without any warnings
            and prompts for confirmation. <br /> <br />
            <strong className="font-semibold text-foreground">
              By checking this setting you confirm that you are have reached the age of majority in
              your country (18 years in most countries).
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
          <Label htmlFor="hide-favorites">Blacklisted tags</Label>
          <p className="text-sm text-muted-foreground">List of tags that you do not wish to see.</p>
        </div>

        <Textarea
          placeholder="Each tag name on the new line"
          rows={10}
          value={state.censoredTags.join('\n')}
          onChange={(e) => update({ censoredTags: e.target.value.split('\n') })}
        />
      </div>

      <div className="space-y-4">
        <div className="space-y-0.5">
          <Label htmlFor="hide-email">Books censoring mode</Label>
          <p className="text-sm text-muted-foreground">
            Method of censoring adult books or books that contain blacklisted tags.
          </p>
        </div>

        <Select
          value={state.censoredTagsMode}
          onValueChange={(value) => update({ censoredTagsMode: censorModeSchema.parse(value) })}
        >
          <SelectTrigger className="max-w-64">
            <SelectValue placeholder="Censoring mode">
              {getCensorModeLabel(state.censoredTagsMode)}
            </SelectValue>
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="censor">
              Censor
              <p className="text-sm">Blurs books that contain blacklisted tags or are adult</p>
            </SelectItem>
            <SelectItem value="hide">
              Hide
              <p className="text-sm">
                Hide books from search results if they contain blacklisted tags
              </p>
            </SelectItem>
            <SelectItem value="none">
              None
              <p className="text-sm max-w-[32rem]">
                Do nothing. Keep adult books and books with blacklisted tags visible. But you will
                still get a warning when opening it.
              </p>
            </SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  )
}

function getCensorModeLabel(mode: UserModerationSettings['censoredTagsMode']) {
  switch (mode) {
    case 'censor':
      return 'Censor'
    case 'hide':
      return 'Hide'
    case 'none':
      return 'None'
  }
}
