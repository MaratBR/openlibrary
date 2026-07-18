import SanitizeHTML from '@/common/SanitizeHTML'

export default function UserContent({ value }: { value: string }) {
  return <SanitizeHTML x-ignore className="user-content contents" value={value} />
}
