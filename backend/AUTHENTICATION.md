# Authentication with OIDC/Authentik

## Quick Setup

### Current Configuration (No Auth)
The app is set to `AUTH_MODE=none` for testing. All endpoints work without authentication.

### Enable OIDC with Authentik

1. **In Authentik**: Create OAuth2/OpenID Provider
   - Redirect URI: `http://localhost:8080/api/auth/oidc/callback`
   - Scopes: `openid`, `email`, `profile`

2. **Update `.env`**:
   ```bash
   AUTH_MODE=oauth
   OIDC_ISSUER_URL=https://authentik.company/application/o/nom-database/
   OIDC_CLIENT_ID=your_client_id
   OIDC_CLIENT_SECRET=your_secret
   OIDC_REDIRECT_URL=http://localhost:8080/api/auth/oidc/callback
   ```

3. **Restart** and visit: `http://localhost:8080/api/auth/oidc/login`

## API Endpoints

- `GET /api/auth/oidc/login` - Start OIDC login
- `GET /api/auth/oidc/callback` - OIDC callback
- `POST /api/auth/register` - Local registration
- `POST /api/auth/login` - Local login

## Supported Providers

✅ Authentik ✅ Keycloak ✅ Auth0 ✅ Okta ✅ Azure AD

Any OIDC-compliant provider works!

## Documentation

- [Authentik Docs](https://goauthentik.io/docs/)
- [OpenID Connect](https://openid.net/specs/openid-connect-core-1_0.html)
