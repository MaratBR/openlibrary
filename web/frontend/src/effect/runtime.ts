import { ApiModuleLive } from '@/features/api/module'
import { Layer, ManagedRuntime } from 'effect'

// Application composition root: add more feature modules here.
const AppLive = Layer.mergeAll(ApiModuleLive)

export const appRuntime = ManagedRuntime.make(AppLive)
