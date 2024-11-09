import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { Ban } from 'lucide-react'
import { NavLink } from 'react-router-dom'

export default function BookIsBannedIndicator({ bookId }: { bookId: string }) {
  return (
    <Tooltip>
      <TooltipTrigger asChild>
        <NavLink to={`/manager/book/banned/${bookId}`}>
          <div className="badge-alt">
            <Ban className="text-red-600" />
            <span className="font-[500]">Banned</span>
          </div>
        </NavLink>
      </TooltipTrigger>
      <TooltipContent className="max-w-64">
        This book has been banned by moderation team. Click here to see more information.
      </TooltipContent>
    </Tooltip>
  )
}
