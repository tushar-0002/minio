# Configurable maximum expiry for presigned URLs

## Background

MinIO rejects a presigned URL whose `X-Amz-Expires` is longer than 7 days
(604800 seconds). The limit is hard-coded in `cmd/signature-v4-parser.go`.
Some operators want a shorter limit, for example 1 hour, so that a leaked
link stops working quickly. Today they cannot set one.

## Goals

- Let an operator set the longest expiry a presigned URL may have.
- Keep the current behavior for every deployment that does not set it.

## Non-goals

- Per-bucket or per-user limits.
- Changing how a presigned URL is signed or how its signature is checked.

## Requirements

**R1. Setting.** Add the `api` setting `presign_max_expiry` and the
environment variable `MINIO_API_PRESIGN_MAX_EXPIRY`. The value is a Go
duration such as `1h` or `72h`. The default is `168h` (7 days).

**R2. Allowed range.** The value must be between `1m` and `168h`, inclusive.
A value outside this range is rejected when the configuration is loaded:
`mc admin config set` fails with an error that names the allowed range, and
the server keeps the previous value. A value is never silently changed to
fit the range.

**R3. Enforcement.** A presigned request whose `X-Amz-Expires` is greater
than the configured limit fails with `AuthorizationQueryParametersError`
(HTTP 400). A request whose expiry is exactly equal to the limit is accepted.
With the default, a URL signed for 604800 seconds keeps working.

**R4. Error message.** The error message states the configured limit in
seconds, for example "X-Amz-Expires must be at most 3600 seconds". It must
not say "a week" when the limit is something else.

**R5. Signature versions.** The limit applies to both Signature V4
(`X-Amz-Expires`) and Signature V2 (`Expires`, an absolute Unix time)
presigned URLs. For V2, the expiry is measured from the time the request is
received.

**R6. Live updates.** A change made with `mc admin config set` takes effect
on every node without a restart, like the other `api` settings. Reading the
limit while it is being changed must be safe.

**R7. Documentation.** `docs/config/README.md` lists the new setting and the
new environment variable with their default.

## Acceptance

- Unit tests cover the default, a custom value, both range bounds, and the
  equal-to-limit case.
- With no setting, presigned URLs behave exactly as before this change.
