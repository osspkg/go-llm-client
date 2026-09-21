# Codex authentication with storage and UI

Use this reference when an application needs to combine the repository's
Codex device-code authentication with a database, secret manager, desktop
keyring, or web UI. The library owns the OAuth/device-code protocol; the
application owns identity, persistence, cookies, and UI state.

## Storage boundary

Keep persistence outside `pkg/auth`:

```go
type CodexTokenStore interface {
	Load(context.Context, string) (auth.CodexTokens, error)
	Save(context.Context, string, auth.CodexTokens) error
	Delete(context.Context, string) error
}
```

The subject key should be an authenticated application user or account ID.
Use a protected database column, KMS-backed secret store, or OS keyring. Treat
all fields of `auth.CodexTokens` as secrets. Do not serialize tokens into a
browser response, URL, log, metric label, or diagnostic error.

The store should distinguish these outcomes:

- no credentials have been connected for the subject;
- credentials exist but cannot be decrypted or loaded;
- credentials were deleted by logout or revocation.

Do not silently replace a storage error with a new login flow: that can hide
an availability or integrity failure.

## Web UI flow

Use two backend endpoints or equivalent RPC methods.

### Start

The backend authenticates the application user, calls
`RequestDeviceCode`, stores the complete value with a short expiration, and
returns only display data:

```go
deviceCode, err := codexAuth.RequestDeviceCode(ctx)
if err != nil {
	return err
}

if err := pendingStore.Save(ctx, subjectID, deviceCode); err != nil {
	return err
}

return StartResponse{
	VerificationURL: deviceCode.VerificationURL,
	UserCode:        deviceCode.UserCode,
	ExpiresAt:        deviceCode.ExpiresAt,
}
```

`pendingStore` is an application-owned store. It may use a short-lived cache or
database row, but it must bind the value to the subject and delete it after
completion, cancellation, or expiry.

Do not send the complete Go value to the browser and later reconstruct it.
`CodexDeviceCode` contains private protocol state required by `Complete`, in
addition to the URL and user code shown to the user.

### Complete

The UI opens `VerificationURL`, asks the user to complete authorization, and
then calls the backend. The backend loads the pending value and completes the
exchange:

```go
deviceCode, err := pendingStore.Load(ctx, subjectID)
if err != nil {
	return err
}

session, err := codexAuth.Complete(ctx, deviceCode)
if err != nil {
	return err
}

if err := tokenStore.Save(ctx, subjectID, session.Tokens()); err != nil {
	return err
}
if err := pendingStore.Delete(ctx, subjectID); err != nil {
	return err
}
```

Return an authorization status to the UI. The browser should receive an
opaque application session cookie with `HttpOnly`, `Secure`, and an appropriate
`SameSite` policy, not a Codex token.

Protect both endpoints with the application's authentication and
authorization. Add CSRF protection where cookie authentication is used, rate
limit repeated start/complete attempts, and verify that the pending flow
belongs to the current subject.

## Restoring and using a session

After a process restart, load the protected token set and construct a new
refreshable session:

```go
tokens, err := tokenStore.Load(ctx, subjectID)
if err != nil {
	return err
}

session, err := codexAuth.NewSession(tokens)
if err != nil {
	return err
}

client, err := openai.New(
	openai.WithAuthProvider(session.HeaderProvider()),
)
```

The session is safe for concurrent use. `HeaderProvider` can refresh an
expiring access token; it keeps the current token set in memory and exposes a
snapshot through `session.Tokens()`.

## Persisting refreshed credentials

The auth package cannot know which application store to update. Persist the
snapshot after login and after client operations that may have refreshed the
session:

```go
response, err := client.Chat().Create(ctx, request)
if err != nil {
	return err
}

if err := tokenStore.Save(ctx, subjectID, session.Tokens()); err != nil {
	return err
}
return response
```

If losing a rotated refresh token is unacceptable, wrap the session provider
at the application boundary. The wrapper should compare token snapshots and
write only when they change; it must preserve the outbound request context and
must decide whether storage failure should fail the request. Never put token
values in the wrapper's error text.

## Desktop UI

For a desktop application, the UI can call a local backend or Go service that
owns the device-code flow. Store `auth.CodexTokens` in the operating system
keyring, not in renderer storage, a plaintext file, or `localStorage`. Expose
only login status, verification URL, user code, and errors safe for display.

On logout:

1. delete the stored token set;
2. discard the in-memory `CodexSession` and client;
3. invalidate the application session;
4. clear pending device-code state.

The library does not revoke the provider credential automatically; implement
provider-specific revocation separately if the product requires it.
