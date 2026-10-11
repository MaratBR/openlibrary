import { DashboardContent } from '@/components/dashboard-layout-components'

export default function PlaceholderPage({ title }: { title: string }) {
  return (
    <DashboardContent.Root>
      <DashboardContent.Header title={title} />
    </DashboardContent.Root>
  )
}
