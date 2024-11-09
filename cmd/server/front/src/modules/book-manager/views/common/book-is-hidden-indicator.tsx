import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { EyeOff } from 'lucide-react'

export default function BookIsHiddenIndicator() {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <div className="badge-alt cursor-help">
          <EyeOff />
          <span className="font-[500]">Hidden</span>
        </div>
      </TooltipTrigger>
      <TooltipContent className="max-w-64">
        This book is not publicly visible to other users. You can change it in the book settings.
      </TooltipContent>
    </Tooltip>
  )
}
