import ReactDOM from 'react-dom'
import { useParams } from 'react-router'

export default function UserProfile() {
  const { userId } = useParams<{ userId: string }>()

  const id = `iframe${userId?.replace(/-/g, '')}`

  return (
    <>
      {ReactDOM.createPortal(
        <iframe id={id} className="w-full" src={`/user/${userId}/__profile`} />,
        document.body,
      )}
    </>
  )
}
