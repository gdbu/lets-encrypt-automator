# Let's Encrypt Server
Let's Encrypt Server is a stand-alone server for auto renewing SSL certificates.

## Usage
```bash
go build
./lets-encrypt-server
```

### Configuration variables
```toml
[env]
email = "[Contact email]"
domain = "[Domain]"
directory = "./my/tls/dir" # Optional, default is "./tls"
port = "80" # Optional, default is "80"
tls-port = "443" # Optional, default is "443"
```

### Example output
```
● Let's Encrypt :: Certificate is expired (or expiring soon), executing renewal process
● Let's Encrypt :: Client created
● Let's Encrypt :: [INFO] acme: Registering account for [User email]
● Let's Encrypt :: User registered
● Let's Encrypt :: [INFO] [your-domain.com] acme: Obtaining bundled SAN certificate
● Let's Encrypt :: [INFO] [your-domain.com] AuthURL: https://acme-v02.api.letsencrypt.org/acme/authz-v3/[Cert Authorization URL]
● Let's Encrypt :: [INFO] [your-domain.com] acme: use tls-alpn-01 solver
● Let's Encrypt :: [INFO] [your-domain.com] acme: Trying to solve TLS-ALPN-01
● Let's Encrypt :: [INFO] [your-domain.com] The server validated our request
● Let's Encrypt :: [INFO] [your-domain.com] acme: Validations succeeded; requesting certificates
● Let's Encrypt :: [INFO] [your-domain.com] Server responded with a certificate.
● Let's Encrypt :: Certificates obtained
● Let's Encrypt :: Certificate renewal process complete
```