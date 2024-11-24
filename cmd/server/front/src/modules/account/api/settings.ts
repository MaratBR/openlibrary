import { httpClient, responseSchema } from '@/modules/common/api'
import { z } from 'zod'

const censorModeSchema = z.enum(['hide', 'censor', 'none']) // Replace with actual CensorMode values

export const userPrivacySettingsSchema = z.object({
  hideStats: z.boolean(),
  hideFavorites: z.boolean(),
  hideComments: z.boolean(),
  hideEmail: z.boolean(),
  allowSearching: z.boolean(),
})

export type UserPrivacySettings = z.infer<typeof userPrivacySettingsSchema>

export const userModerationSettingsSchema = z.object({
  showAdultContent: z.boolean(),
  censoredTags: z.array(z.string()),
  censoredTagsMode: censorModeSchema,
})

export type UserModerationSettings = z.infer<typeof userModerationSettingsSchema>

export const userAboutSettingsSchema = z.object({
  about: z.string(),
  status: z.string(),
  gender: z.string(),
})

export type UserAboutSettings = z.infer<typeof userAboutSettingsSchema>

export const userCustomizationSettingSchema = z.object({
  profileCss: z.string(),
  defaultTheme: z.string(),
  enableProfileCss: z.boolean(),
})

export type UserCustomizationSettings = z.infer<typeof userCustomizationSettingSchema>

export function httpGetUserPrivacySettings() {
  return httpClient.get('/api/settings/privacy').then(responseSchema(userPrivacySettingsSchema))
}

export function httpGetUserAboutSettings() {
  return httpClient.get('/api/settings/about').then(responseSchema(userAboutSettingsSchema))
}

export function httpGetUserCustomizationSettings() {
  return httpClient
    .get('/api/settings/customization')
    .then(responseSchema(userCustomizationSettingSchema))
}

export function httpGetUserModerationSettings() {
  return httpClient
    .get('/api/settings/moderation')
    .then(responseSchema(userModerationSettingsSchema))
}

export function httpUpdateUserPrivacySettings(settings: UserPrivacySettings) {
  return httpClient.put('/api/settings/privacy', { json: settings })
}

export function httpUpdateUserAboutSettings(settings: UserAboutSettings) {
  return httpClient.put('/api/settings/about', { json: settings })
}

export function httpUpdateUserCustomizationSettings(settings: UserCustomizationSettings) {
  return httpClient.put('/api/settings/customization', { json: settings })
}

export function httpUpdateUserModerationSettings(settings: UserModerationSettings) {
  return httpClient.put('/api/settings/moderation', { json: settings })
}
