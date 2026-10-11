# Initial setup

The public setup middleware needs a registered application SetupService.

## Requirements

- Setup is required when SiteSetup.IsSetup is false or no administrator exists.
- An administrator is a user with the admin role; a system user does not substitute for one.
- SiteSetup.IsSetup defaults to false and is persisted as a versioned siteconfig row.
- IsSetup caches successful checks, including false; errors are retried.
- GetSetupState reads current state without using or changing the IsSetup cache and reports errors.
- GetSetupState reports AdminUserExist so callers can skip administrator creation.
- Setup requires a nonempty site name and, when no admin exists, a valid username and nonempty password satisfying site password requirements.
- Setup creates an admin only when none exists, then saves the site name and setup flag and invalidates the cache.
- Completed sites reject repeated Setup calls. Failed configuration saves preserve the prior flag and can be retried without creating another admin.

## Setup page

- GET /setup renders a library name form and, only when an admin is missing, administrator username, password, and password confirmation fields.
- The page displays configured password requirements and works without JavaScript.
- POST /setup uses existing CSRF middleware and submits to SetupService.
- Invalid input returns a rendered form with an error, retaining the name and username but never passwords.
- Successful setup redirects to login, and normal public pages become accessible.
- Completed sites redirect GET and POST /setup to login without invoking Setup again.
- Setup pages and redirects are not cached. State lookup failures return 503 and save failures show a generic retry message without exposing internal errors.
- Unsupported methods return 405 with an Allow header.
- Setup checks precede public authorization and content preferences middleware.
- New English page copy lives in the setup translation namespace.
