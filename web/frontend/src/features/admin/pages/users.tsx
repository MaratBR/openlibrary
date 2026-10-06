import { useState } from 'react'
import { Link, LoaderFunctionArgs, useLoaderData, useLocation, useSearchParams } from 'react-router'
import { AdminUser, loadUser, loadUsers, roles, save } from '@/features/admin/api'
import {
  Card,
  Empty,
  Field,
  PageHeader,
  Pager,
  SaveBar,
  t,
  useDirtyForm,
  useSave,
} from '@/features/admin/components/ui'
import { evaluatePasswordStrength } from '@/lib/password'

export function usersLoader({ request }: LoaderFunctionArgs) {
  return loadUsers(new URL(request.url).search, request.signal)
}
export function userLoader({ params, request }: LoaderFunctionArgs) {
  return loadUser(params.id || '', request.signal)
}
export function Users() {
  const data = useLoaderData<typeof usersLoader>()
  const [params, setParams] = useSearchParams()
  const location = useLocation()
  const editState = { back: `/users${location.search}` }
  return (
    <>
      <PageHeader title={t('admin.users.title')} description={t('admin.ui.usersDescription')} />
      <Card className="admin-card--table">
        <form
          key={params.toString()}
          className="Admin-toolbar"
          onSubmit={(event) => {
            event.preventDefault()
            const fields = new FormData(event.currentTarget)
            const q = new URLSearchParams()
            const search = String(fields.get('q') || '').trim()
            const role = String(fields.get('role') || '')
            if (search) q.set('q', search)
            if (role) q.set('usersFilter.role', role)
            setParams(q)
          }}
        >
          <Field id="user-search" label={t('admin.ui.searchUsers')}>
            <input
              className="input"
              id="user-search"
              name="q"
              type="search"
              defaultValue={params.get('q') || ''}
              placeholder={t('common.search')}
            />
          </Field>
          <Field id="role-filter" label={t('admin.users.role')}>
            <select
              className="Select"
              name="role"
              id="role-filter"
              defaultValue={params.get('usersFilter.role') || ''}
            >
              <option value="">{t('admin.ui.allRoles')}</option>
              {roles.map((role) => (
                <option value={role} key={role}>
                  {t(`role.${role}`)}
                </option>
              ))}
            </select>
          </Field>
          <button type="submit" className="Btn Btn--primary">
            {t('common.search')}
          </button>
        </form>
        <div className="Admin-resultCount">
          {t('admin.ui.results', { count: data.total.toLocaleString() })}
        </div>
        {data.users.length === 0 ? (
          <Empty kind="users" />
        ) : (
          <div
            className="Admin-tableScroll"
            tabIndex={0}
            role="region"
            aria-label={t('admin.users.title')}
          >
            <table className="table Admin-userTable">
              <thead>
                <tr>
                  <th>{t('admin.users.name')}</th>
                  <th>{t('admin.users.role')}</th>
                  <th>{t('admin.ui.status')}</th>
                  <th>{t('admin.users.actions')}</th>
                </tr>
              </thead>
              <tbody>
                {data.users.map((user) => (
                  <tr key={user.id}>
                    <td>
                      <Link
                        to={`/users/${user.id}`}
                        state={editState}
                        className="Admin-userIdentity"
                      >
                        <span className="Admin-avatar">
                          {user.avatar ? (
                            <img src={user.avatar} alt="" loading="lazy" />
                          ) : (
                            user.name.slice(0, 1)
                          )}
                        </span>
                        <span>
                          <strong>{user.name}</strong>
                          <small>
                            {t('admin.ui.joined')} {user.joinedAt}
                          </small>
                        </span>
                      </Link>
                    </td>
                    <td>
                      <span className="Admin-badge">{t(`role.${user.role}`)}</span>
                    </td>
                    <td>
                      <span
                        className={`Admin-badge ${user.isBanned ? 'Admin-badge--danger' : 'Admin-badge--accent'}`}
                      >
                        {t(user.isBanned ? 'admin.ui.banned' : 'admin.ui.active')}
                      </span>
                    </td>
                    <td>
                      <Link
                        to={`/users/${user.id}`}
                        state={editState}
                        className="Btn Btn--outline Btn--sm"
                      >
                        {t('common.edit')}
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>
      <Pager page={data.page} total={data.totalPages} />
    </>
  )
}
export function UserEdit() {
  const user = useLoaderData<typeof userLoader>()
  const location = useLocation()
  const back =
    typeof location.state?.back === 'string' && location.state.back.startsWith('/users')
      ? location.state.back
      : '/users'
  return <UserForm key={user.id} user={user} back={back} />
}
function newPassword() {
  // Cryptographic randomness for administrative password resets.
  const alphabet = 'ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz23456789!@#$%&*'
  const values = crypto.getRandomValues(new Uint8Array(24))
  return Array.from(values, (value) => alphabet[value % alphabet.length]).join('')
}
function UserForm({ user, back }: { user: AdminUser; back: string }) {
  const initial = { about: user.bio, gender: user.gender, role: user.role, password: '' }
  const [values, setValues] = useState(initial)
  const [baseline, setBaseline] = useState(initial)
  const [genderType, setGenderType] = useState(
    ['', 'male', 'female'].includes(user.gender) ? user.gender : 'other',
  )
  const [reset, setReset] = useState(false)
  const [showPassword, setShowPassword] = useState(false)
  const state = useSave()
  const prompt = useDirtyForm(JSON.stringify(values) !== JSON.stringify(baseline), state.pending)
  const change = (next: Partial<typeof values>) => {
    setValues((value) => ({ ...value, ...next }))
    state.changed()
  }
  return (
    <>
      <Link className="Admin-back" to={back}>
        ← {t('admin.ui.backUsers')}
      </Link>
      <PageHeader
        title={user.name}
        description={t('admin.users.generalInformation')}
        actions={
          <>
            <a className="Btn Btn--outline" data-admin-external href={`/users/${user.id}`}>
              {t('admin.users.profile')}
            </a>
            <a className="Btn Btn--outline" data-admin-external href={`/mod/user/${user.id}`}>
              {t('admin.users.moderation')}
            </a>
          </>
        }
      />
      <form
        onSubmit={(event) =>
          state.submit(event, async () => {
            await save(`/users/${user.id}`, values)
            setBaseline(values)
            if (reset) {
              setReset(false)
              setValues({ ...values, password: '' })
              setBaseline({ ...values, password: '' })
            }
          })
        }
      >
        <fieldset disabled={state.pending} className="Admin-workspace">
          <div className="Admin-stack">
            <Card title={t('admin.users.generalInformation')}>
              <Field id="user-name" label={t('admin.users.name')} help={t('admin.ui.nameReadOnly')}>
                <input
                  className="input"
                  id="user-name"
                  disabled
                  value={user.name}
                  aria-describedby="user-name-help"
                />
              </Field>
              <Field id="user-gender" label={t('admin.users.gender')}>
                <select
                  id="user-gender"
                  className="Select"
                  value={genderType}
                  onChange={(event) => {
                    const next = event.target.value
                    setGenderType(next)
                    change({ gender: next === 'other' ? '' : next })
                  }}
                >
                  <option value="">{t('gender.ratherNotSay')}</option>
                  <option value="male">{t('gender.m')}</option>
                  <option value="female">{t('gender.f')}</option>
                  <option value="other">{t('gender.o')}</option>
                </select>
              </Field>
              {genderType === 'other' && (
                <Field id="custom-gender" label={t('admin.ui.customGender')}>
                  <input
                    id="custom-gender"
                    className="input"
                    value={values.gender}
                    onChange={(event) => change({ gender: event.target.value })}
                  />
                </Field>
              )}
              <Field id="user-about" label={t('admin.users.about')}>
                <textarea
                  id="user-about"
                  className="textarea"
                  rows={7}
                  value={values.about}
                  onChange={(event) => change({ about: event.target.value })}
                />
              </Field>
            </Card>
            <Card
              title={t('admin.users.accountManagement')}
              description={t('admin.ui.accountDescription')}
            >
              {!reset ? (
                <button
                  className="Btn Btn--outline"
                  type="button"
                  onClick={() => {
                    setReset(true)
                    change({ password: newPassword() })
                  }}
                >
                  {t('admin.users.resetPassword')}
                </button>
              ) : (
                <>
                  <Field
                    id="new-password"
                    label={t('admin.ui.password')}
                    help={t('admin.ui.passwordHelp')}
                  >
                    <input
                      id="new-password"
                      className="input"
                      type={showPassword ? 'text' : 'password'}
                      autoComplete="new-password"
                      required
                      minLength={8}
                      value={values.password}
                      aria-describedby="new-password-help"
                      onChange={(event) => change({ password: event.target.value })}
                    />
                    <p className="Admin-help">
                      {t(`passwordStrength.${evaluatePasswordStrength(values.password).strength}`)}
                    </p>
                  </Field>
                  <div className="Admin-actions">
                    <button
                      type="button"
                      className="Btn Btn--outline"
                      onClick={() => setShowPassword(!showPassword)}
                    >
                      {t(showPassword ? 'admin.ui.hidePassword' : 'admin.ui.showPassword')}
                    </button>
                    <button
                      type="button"
                      className="Btn Btn--outline"
                      onClick={() => change({ password: newPassword() })}
                    >
                      {t('admin.ui.generatePassword')}
                    </button>
                    <button
                      type="button"
                      className="Btn Btn--ghost"
                      onClick={() => {
                        setReset(false)
                        change({ password: '' })
                      }}
                    >
                      {t('admin.ui.cancelReset')}
                    </button>
                  </div>
                </>
              )}
            </Card>
          </div>
          <div className="Admin-stack">
            <Card title={t('admin.ui.profile')}>
              <div className="Admin-profile">
                <span className="Admin-profileAvatar">
                  {user.avatar && <img src={user.avatar} alt="" />}
                </span>
                <strong>{user.name}</strong>
                <span
                  className={`Admin-badge ${user.isBanned ? 'Admin-badge--danger' : 'Admin-badge--accent'}`}
                >
                  {t(user.isBanned ? 'admin.ui.banned' : 'admin.ui.active')}
                </span>
                <p className="Admin-help">
                  {t('admin.ui.joined')} {user.joinedAt}
                </p>
              </div>
            </Card>
            <Card
              title={t('admin.users.rolesAndPermissions')}
              description={t('admin.ui.permissionsDescription')}
            >
              <Field id="user-role" label={t('admin.users.role')}>
                <select
                  id="user-role"
                  className="Select"
                  value={values.role}
                  onChange={(event) => change({ role: event.target.value as AdminUser['role'] })}
                >
                  {roles.map((role) => (
                    <option key={role} value={role}>
                      {t(`role.${role}`)}
                    </option>
                  ))}
                </select>
              </Field>
            </Card>
          </div>
        </fieldset>
        <SaveBar state={state} cancelTo={back} />
      </form>
      {prompt}
    </>
  )
}
