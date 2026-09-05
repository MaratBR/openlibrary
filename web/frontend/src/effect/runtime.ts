import { ApiModuleLive } from '@/features/api/module'
import { Layer, Logger, ManagedRuntime, References } from 'effect'
import { EditorFontsApi } from '@/block-editor/fonts/api'
import { EditorFonts } from '@/block-editor/fonts/service'
import { FontsLoader } from '@/features/fonts-loader/loader'
import { FontsLoaderConfig } from '@/features/fonts-loader/config'

// Application composition root: add more feature modules here.
const LoggerLive = Layer.mergeAll(
  Logger.layer([Logger.consolePretty({ mode: 'browser' })]),
  Layer.succeed(References.MinimumLogLevel, 'Debug'),
)
const FontsLoaderConfigLive = FontsLoaderConfig.layer
const FontsLoaderLive = FontsLoader.layer.pipe(
  Layer.provide(FontsLoaderConfigLive),
  Layer.provide(ApiModuleLive),
)
const EditorFontsLive = EditorFonts.layer.pipe(
  Layer.provide(EditorFontsApi.layer),
  Layer.provide(FontsLoaderLive),
  Layer.provide(ApiModuleLive),
)
const AppLive = Layer.mergeAll(
  LoggerLive,
  ApiModuleLive,
  FontsLoaderConfigLive,
  FontsLoaderLive,
  EditorFontsLive,
)

export const appRuntime = ManagedRuntime.make(AppLive)
