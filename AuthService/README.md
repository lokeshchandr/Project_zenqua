
1. Set env vars (example):

```powershell
$env:JWT_SECRET = "super-secret";
$env:JWT_EXPIRY = "24h";
$env:SUPERADMIN_EMAIL = "root@root.com";
$env:SUPERADMIN_PASSWORD = "Root123";
$env:SUPERADMIN_NAME = "Root";
$env:DB_PATH = "auth.db";
```

2. Build and run:
```powershell
# from AuthService directory
go mod tidy;
go run .
```