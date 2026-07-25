import { DashboardContent } from '@/components/dashboard-layout-components'

export default function PlaceholderPage({ title }: { title: string }) {
  return (
    <DashboardContent.Root>
      <DashboardContent.StickyHeader title={title} />
    </DashboardContent.Root>
  )
}
