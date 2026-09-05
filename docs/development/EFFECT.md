# Effect in the frontend

This document records how the TypeScript frontend currently uses
[Effect](https://effect.website/). It is a project convention guide, not a general Effect
tutorial. Check the installed Effect version and its local types when an API is uncertain; the
project may use a release candidate whose API differs from older examples online.

Effect currently has two related uses in this repository:

1. `Effect`, `Context`, `Layer`, and `ManagedRuntime` provide dependency injection, controlled
   asynchronous work, errors, cancellation, caching, and logging.
2. `Schema` validates data at network and other trust boundaries. Schema can be used independently
   in code that does not need the application runtime.

The composition root is `web/frontend/src/effect/runtime.ts`. Shared service implementations live
with their owning feature, while reusable API services are collected in
`web/frontend/src/features/api/module.ts`.

## Defining services

Use `Context.Service` for dependencies and stateful capabilities. Give every service a stable,
project-qualified identifier and describe its public interface explicitly:

```ts
export class ExampleApi extends Context.Service<
  ExampleApi,
  {
    readonly load: (id: string) => Effect.Effect<Item, unknown>
  }
>()('openlibrary/ExampleApi') {
  static readonly layer = Layer.effect(
    this,
    Effect.gen(function* () {
      const httpClient = yield* HttpClient

      const load = Effect.fn('ExampleApi.load')(function* (id: string) {
        // implementation
      })

      return ExampleApi.of({ load })
    }),
  )
}
```

Conventions:

- Use `Layer.succeed` for an already-created synchronous value, as `HttpClient` does.
- Use `Layer.effect` when construction reads dependencies, creates state, or may eventually become
  asynchronous.
- Resolve dependencies during layer construction with `yield* SomeService`. Do not import a
  concrete singleton to bypass the service boundary.
- Keep mutable state inside the service layer closure. A layer supplied once to `ManagedRuntime`
  gives the application one managed service instance.
- Return the implementation with `ServiceName.of(...)`.
- Use `Effect.fn('Service.method')` for meaningful asynchronous or fallible service operations.
  The name improves tracing and diagnostics.
- Public service methods should expose their actual error and dependency types. Several existing
  services currently use `unknown` for errors; prefer a more specific domain error when practical.
- Configuration is also a service when consumers should not own its source or lifetime. See
  `FontsLoaderConfig` for the current pattern.

API services should handle transport and boundary decoding. Stateful feature services may combine
multiple APIs and expose application state or operations, but UI concerns should remain in React,
Jotai, or the relevant view layer.

## Effects and errors

Use the constructor that matches the side effect:

- `Effect.sync` for synchronous mutations or callbacks that must be sequenced.
- `Effect.try` for synchronous code that can throw.
- `Effect.tryPromise` for promises that can reject.
- `Effect.fail` for an explicit failure.
- `Effect.gen` to sequence services and effects.
- `Effect.all` for independent work; state the concurrency policy when ordering or resource use
  matters.

Do not hide promise work inside an untracked callback when it can be represented as an Effect.
Convert it at the boundary and compose it from there.

Use `Effect.tapError` for error reporting or state updates that should preserve the original
failure. Use `Effect.ensuring` for cleanup or state restoration that must happen on success,
failure, or interruption. A typical stateful operation looks like:

```ts
const initialize = Effect.fn('Feature.initialize')(function* () {
  yield* setState({ loading: true, error: undefined })
  yield* work.pipe(
    Effect.tapError((error) => setState({ error })),
    Effect.ensuring(setState({ loading: false })),
  )
})
```

Use `Schema.decodeUnknownEffect` after parsing untrusted JSON. Declare every field callers need:
schema object decoding strips fields that are not declared.

## Layers and the runtime

`ApiModuleLive` groups API services and supplies their shared `HttpClient` dependency. The
application runtime then composes API, configuration, feature, and infrastructure layers.

When adding a service:

1. Define its service interface and live layer next to the owning feature.
2. Provide the layer's dependencies with `Layer.provide`.
3. Merge the fully provided layer into `AppLive` if code at the runtime boundary needs to resolve
   it directly.
4. Add reusable API services to `ApiModuleLive` when they belong to that shared module.

Use `Layer.provideMerge` when building a reusable module that should retain both the supplied
services and the services produced by the layer. Use ordinary `Layer.provide` when satisfying an
implementation detail that should not be exposed from that composed layer.

Do not call `Effect.runPromise` or `Effect.runSync` directly for application work. Use the shared
`appRuntime`; otherwise the program will not receive application services, logger configuration,
or the runtime-managed service instances.

## React and state-library boundaries

React event handlers, effects, and state-library actions are execution boundaries. Build and
compose Effects below that boundary, then execute them through `appRuntime`:

```ts
const abortController = new AbortController()

void appRuntime.runPromise(program, { signal: abortController.signal }).then(onSuccess, onFailure)

return () => abortController.abort()
```

Cancellation should follow component lifetime for requests and other interruptible work. Do not
update React state after cancellation.

Resolve a long-lived service with `appRuntime.runSync(Service)` only when an imperative integration
needs to retain its methods, such as an external-store adapter or DOM lifecycle bridge. Otherwise,
prefer resolving it inside the program:

```ts
const program = Effect.gen(function* () {
  const api = yield* ExampleApi
  return yield* api.load(id)
})
```

Jotai write atoms may return Effects. The React caller remains responsible for running the returned
Effect through `appRuntime`. Wrap external state mutation in `Effect.sync` so its ordering and
cleanup are part of the program.

For subscriptions exposed by an Effect service, pair `subscribe` and `getState` with React's
`useSyncExternalStore`. Ensure service initialization is idempotent because React effects may mount
more than once in development.

## Logging

The application runtime installs a browser pretty-console logger and currently sets
`References.MinimumLogLevel` to `Debug`.

Use Effect logging rather than direct `console` calls inside Effects:

```ts
return Effect.logDebug('Loaded items').pipe(
  Effect.annotateLogs({ service: 'ExampleApi', count: items.length }),
)
```

Guidelines:

- Write a short event as the message and put structured context in `Effect.annotateLogs`.
- Include a `service` annotation for service-level logs.
- Use debug logs for lifecycle, cache, request, and count diagnostics.
- Use `Effect.logError` for failures that are logged at the point they are handled or propagated.
- Do not log secrets, CSRF tokens, full user content, or unnecessary response bodies.
- Avoid logging the same failure at every layer. Log where useful context is available or where the
  failure is finally handled.

## Resource lifetime and browser APIs

Represent acquisition and cleanup together where possible. Components that attach browser
resources to a stateful service must detach them on unmount. If setup must finish before dependent
UI is safe to render, await the Effect at that lifecycle boundary.

Browser APIs that expose promises, events, or DOM mutation should be wrapped in `Effect.tryPromise`
or `Effect.async` rather than leaking raw promises from service methods. Preserve interruption and
cleanup semantics when the underlying API allows it.

## Testing and verification

Service layers make substitution possible: provide a test layer for HTTP, configuration, or other
dependencies and run the program with that layer. Keep business behavior in service Effects so it
can be exercised without mounting React.

After changing Effect-based frontend code, run:

```sh
pnpm exec tsc --noEmit
pnpm run build
git diff --check
```

Report unrelated baseline failures rather than changing unrelated files. For asynchronous services,
test success, failure, cleanup/interruption, idempotency, and concurrency or cache behavior when
those properties matter.

## Existing examples

- `web/frontend/src/effect/runtime.ts`: application layer composition and logger setup.
- `web/frontend/src/features/api/module.ts`: reusable API module composition.
- `web/frontend/src/features/http-client/client.ts`: synchronous service layer.
- `web/frontend/src/features/book-manager/api.ts`: named API operations and HTTP failures.
- `web/frontend/src/features/search/api.ts`: service-owned Effect cache.
- `web/frontend/src/features/search/use-tags-search.ts`: React execution and cancellation boundary.
- `web/frontend/src/block-editor/fonts/service.ts`: stateful service and external-store interface.
- `web/frontend/src/features/fonts-loader`: configuration, mutable service state, logging, queueing,
  and DOM resource integration.
