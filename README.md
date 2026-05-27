# workflow-plugin-okta

> ⚠️ **Experimental** — This plugin compiles and passes its unit tests but has not been validated in any active GoCodeAlone-internal production deployment. Use with caution. Please [open an issue](https://github.com/GoCodeAlone/workflow-plugin-okta/issues/new) if you adopt it so we can promote it to **verified** status.

Okta identity and access management plugin

## Auth Provider Descriptor

The plugin exposes `step.okta_auth_provider_describe` so an admin surface can
discover Okta management capabilities without hard-coded UI assumptions. The
descriptor advertises only capabilities backed by existing Okta steps:

- identity management for users, groups, group rules, applications, and
  lifecycle actions
- OAuth/OIDC authorization server administration for scopes, claims, policies,
  policy rules, and keys
- enterprise SSO administration for applications, identity providers, domains,
  brands, and org settings
- MFA/authenticator administration
- audit log and hook administration

The descriptor also publishes selectable Okta management scopes and required
configuration fields for API token or OAuth private-key operation. These fields
are emitted through strict protobuf contracts in `workflow.plugins.okta.v1`.

## Install

```sh
wfctl plugin install workflow-plugin-okta
```

## License

MIT
