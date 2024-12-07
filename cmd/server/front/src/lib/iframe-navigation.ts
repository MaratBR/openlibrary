import {
  AgnosticDataRouteObject,
  AgnosticRouteObject,
  Blocker,
  BlockerFunction,
  Fetcher,
  FutureConfig,
  GetScrollPositionFunction,
  GetScrollRestorationKeyFunction,
  Location,
  Path,
  Router,
  RouterFetchOptions,
  RouterState,
  RouterSubscriber,
  To,
  UNSAFE_DeferredData,
} from '@remix-run/router'
import { z } from 'zod'

const routerNavigationMessage = z.object({
  type: z.literal('router-navigate-call'),
  args: z.array(z.unknown()),
})

export function initIframeRouter(router: Router) {
  if (window.top !== window.self) {
    return
  }

  window.addEventListener('message', (event) => {
    const result = routerNavigationMessage.safeParse(event.data)
    if (result.success) {
      // @ts-expect-error ts(2349)
      router.navigate(...result.data.args)
    }
  })
}

export function wrapRouter(router: Router): Router {
  if (window.self === window.top) {
    return router
  }
  return new IframeRouter(router)
}

class IframeRouter implements Router {
  private readonly _inner: Router

  constructor(inner: Router) {
    if (window.self === window.top) {
      throw new Error('Cannot use IframeRouter outside of an iframe')
    }
    this._inner = inner

    this.initialize = this.initialize.bind(this)
    this.subscribe = this.subscribe.bind(this)
    this.enableScrollRestoration = this.enableScrollRestoration.bind(this)
    this.navigate = this.navigate.bind(this)
    this.fetch = this.fetch.bind(this)
    this.revalidate = this.revalidate.bind(this)
    this.createHref = this.createHref.bind(this)
    this.encodeLocation = this.encodeLocation.bind(this)
    this.getFetcher = this.getFetcher.bind(this)
    this.deleteFetcher = this.deleteFetcher.bind(this)
    this.dispose = this.dispose.bind(this)
    this.getBlocker = this.getBlocker.bind(this)
    this.deleteBlocker = this.deleteBlocker.bind(this)
    this.patchRoutes = this.patchRoutes.bind(this)
    this._internalSetRoutes = this._internalSetRoutes.bind(this)
  }

  get basename(): string | undefined {
    return this._inner.basename
  }
  get future(): FutureConfig {
    return this._inner.future
  }
  get state(): RouterState {
    return this._inner.state
  }
  get routes(): AgnosticDataRouteObject[] {
    return this._inner.routes
  }
  get window(): Window | undefined {
    return this._inner.window
  }
  initialize(): Router {
    return this._inner.initialize()
  }
  subscribe(fn: RouterSubscriber): () => void {
    return this._inner.subscribe(fn)
  }
  enableScrollRestoration(
    savedScrollPositions: Record<string, number>,
    getScrollPosition: GetScrollPositionFunction,
    getKey?: GetScrollRestorationKeyFunction,
  ): () => void {
    return this._inner.enableScrollRestoration(savedScrollPositions, getScrollPosition, getKey)
  }
  navigate(...args: unknown[]): Promise<void> {
    ;(this.window ?? window).parent.postMessage({
      type: 'router-navigate-call',
      args,
    } satisfies z.infer<typeof routerNavigationMessage>)
    return Promise.resolve()
  }
  fetch(key: string, routeId: string, href: string | null, opts?: RouterFetchOptions): void {
    return this._inner.fetch(key, routeId, href, opts)
  }
  revalidate(): void {
    return this._inner.revalidate()
  }
  createHref(location: Location | URL): string {
    return this._inner.createHref(location)
  }
  encodeLocation(to: To): Path {
    return this._inner.encodeLocation(to)
  }
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  getFetcher<TData = any>(key: string): Fetcher<TData> {
    return this._inner.getFetcher(key)
  }
  deleteFetcher(key: string): void {
    return this._inner.deleteFetcher(key)
  }
  dispose(): void {
    return this._inner.dispose()
  }
  getBlocker(key: string, fn: BlockerFunction): Blocker {
    return this._inner.getBlocker(key, fn)
  }
  deleteBlocker(key: string): void {
    this._inner.deleteBlocker(key)
  }
  patchRoutes(routeId: string | null, children: AgnosticRouteObject[]): void {
    this._inner.patchRoutes(routeId, children)
  }
  _internalSetRoutes(routes: AgnosticRouteObject[]): void {
    this._inner._internalSetRoutes(routes)
  }
  get _internalFetchControllers(): Map<string, AbortController> {
    return this._inner._internalFetchControllers
  }
  set _internalFetchControllers(value: Map<string, AbortController>) {
    this._inner._internalFetchControllers = value
  }
  get _internalActiveDeferreds(): Map<string, UNSAFE_DeferredData> {
    return this._inner._internalActiveDeferreds
  }
  set _internalActiveDeferreds(value: Map<string, UNSAFE_DeferredData>) {
    this._inner._internalActiveDeferreds = value
  }
}
