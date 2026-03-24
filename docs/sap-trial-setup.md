# SAP ABAP Platform Trial 2023 — Self-Hosted Setup Guide

This document describes how to run the SAP ABAP Platform Trial 2023 Docker
container on a Linux server (e.g. Hetzner Cloud), configure it for ADT access,
and connect the integration test suite and GitHub Actions CI to it.

> **Security note:** This guide intentionally omits the server IP/hostname.
> Never commit connection URLs, credentials, or license keys to the repository.

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Server Setup](#server-setup)
3. [SAP ABAP Trial Container](#sap-abap-trial-container)
   - [Pulling the Image](#pulling-the-image)
   - [Starting the Container](#starting-the-container)
   - [Disk Space Warning](#disk-space-warning)
4. [SAP System Configuration](#sap-system-configuration)
   - [License Installation](#license-installation)
   - [Work Process Tuning](#work-process-tuning)
   - [User Access](#user-access)
   - [Unlocking the DEVELOPER User](#unlocking-the-developer-user)
5. [HTTPS / Reverse Proxy Setup](#https--reverse-proxy-setup)
6. [Integration Tests](#integration-tests)
   - [Running Locally](#running-locally)
   - [Known Test Failures](#known-test-failures)
7. [GitHub Actions CI](#github-actions-ci)
   - [Workflow Overview](#workflow-overview)
   - [GitHub Secrets Setup](#github-secrets-setup)
8. [Troubleshooting](#troubleshooting)

---

## Prerequisites

- A Linux server with at least **16 GB RAM**, **4 CPU cores**, and **150 GB
  disk** (the SAP container image is ~80 GB compressed).
- Docker (or Podman) installed.
- Root or `sudo` access on the server.
- A DNS A record pointing a subdomain at the server IP (for HTTPS).
- A SAP license file for your hardware key (obtain from the SAP trial portal).

---

## Server Setup

### Install Docker

```bash
# Debian / Ubuntu
apt-get update
apt-get install -y docker.io
systemctl enable --now docker
```

### Verify disk space

The SAP image is large. Confirm there is sufficient space before pulling:

```bash
df -h /var/lib/docker
```

---

## SAP ABAP Trial Container

### Pulling the Image

The official SAP ABAP Cloud Developer Trial image is available from Docker Hub
under the `sapse` organisation:

```bash
docker pull sapse/abap-cloud-developer-trial:2023
```

> **Note:** `podman` can be used as a drop-in replacement. If you get a disk-
> full error from podman's `/var/tmp` overlay, ensure the underlying partition
> has enough space or reconfigure the podman storage driver.

### Starting the Container

```bash
docker run -d \
  --name a4h \
  --hostname vhcala4hci \
  -p 50000:50000 \
  -p 50001:50001 \
  -p 8443:8443 \
  -p 30213:30213 \
  --sysctl net.ipv4.ip_local_port_range="40000 60999" \
  --sysctl kernel.shmmax=21474836480 \
  --sysctl kernel.shmmni=32768 \
  --sysctl kernel.shmall=5242880 \
  -v /data/sap/sysvol:/sysvol \
  sapse/abap-cloud-developer-trial:2023
```

Key parameters:

| Parameter | Purpose |
|-----------|---------|
| `--hostname vhcala4hci` | SAP requires a specific hostname |
| `-p 50000:50000` | SAP ICM HTTP port (ADT, browser access) |
| `-p 50001:50001` | SAP ICM HTTPS port |
| `-p 8443:8443` | Alternative HTTPS |
| `-p 30213:30213` | HANA SQL port (multitenant tenant DB) |
| `--sysctl ...` | Required kernel parameters for SAP/HANA |
| `-v /data/sap/sysvol:/sysvol` | Persistent volume for SAP data |

### Disk Space Warning

If you see:
```
Error: copying file write /var/tmp/podman934593548: no space left on device
```
This means the partition hosting `/var/tmp` or the podman overlay is full.
Either free space or move Docker/Podman storage to a larger partition.

### Verifying the Container is Up

The SAP system takes 5-10 minutes to fully start. Check readiness:

```bash
# Watch SAP startup progress
docker exec a4h /usr/sap/hostctrl/exe/sapcontrol -nr 00 -function GetProcessList

# Quick HTTP ping (expects 403 when SAP is up)
curl -s -o /dev/null -w "%{http_code}" http://localhost:50000/sap/bc/ping
```

SAP is ready when `sapcontrol GetProcessList` shows all processes as **Running**.

---

## SAP System Configuration

### License Installation

The trial image ships without a permanent license. Obtain a permanent license
from the SAP trial portal for your hardware key.

**Find your hardware key:**

```bash
docker exec a4h /usr/sap/A4H/SYS/exe/run/saplikey \
  pf=/usr/sap/A4H/SYS/profile/A4H_D00_vhcala4hci \
  -get
```

Note the `Hardware Key` from the output and request a license file from the
SAP trial portal.

**Install the license:**

```bash
# Copy license file into container
docker cp /path/to/A4H_license.txt a4h:/tmp/A4H_license.txt

# Install all keys from the file
docker exec a4h /usr/sap/A4H/SYS/exe/run/saplikey \
  pf=/usr/sap/A4H/SYS/profile/A4H_D00_vhcala4hci \
  -install /tmp/A4H_license.txt

# Verify installation
docker exec a4h /usr/sap/A4H/SYS/exe/run/saplikey \
  pf=/usr/sap/A4H/SYS/profile/A4H_D00_vhcala4hci \
  -get
```

The correct profile path inside the container is:
```
/usr/sap/A4H/SYS/profile/A4H_D00_vhcala4hci
```

> **Common mistake:** The profile is `A4H_D00_vhcala4hci`, not
> `A4H_DVEBMGS00_vhcala4hci`. List `ls /usr/sap/A4H/SYS/profile/` to confirm
> the correct filename if `saplikey` reports a missing profile error.

### Work Process Tuning

The default SAP profile only allocates **7 dialog work processes**. Running the
full integration test suite (34 tests) exhausts these quickly and causes 503
errors. Increase them:

**Edit the instance profile inside the container:**

```bash
docker exec -it a4h bash
vi /usr/sap/A4H/SYS/profile/A4H_D00_vhcala4hci
```

Change:
```
rdisp/wp_no_dia = 7
```
To:
```
rdisp/wp_no_dia = 15
rdisp/wp_no_btc = 5
rdisp/wp_no_vb  = 1
```

**Restart the ABAP application server (not the whole container):**

```bash
# Stop ABAP only
docker exec a4h /usr/sap/hostctrl/exe/sapcontrol -nr 00 -function Stop
# Wait ~60s for full stop
docker exec a4h /usr/sap/hostctrl/exe/sapcontrol -nr 00 -function Start
```

> **Note:** `RestartInstance` did not work reliably; use explicit `Stop` then
> `Start`.

### User Access

The trial system ships with these pre-configured users:

| User | Default Password | Role |
|------|-----------------|------|
| `DEVELOPER` | `ABAPtr2023#00` | ABAP developer (S_DEVELOP auth) |
| `DDIC` | `ABAPtr2023#00` | Data dictionary admin |
| `BWDEVELOPER` | `ABAPtr2023#00` | BW developer |

**Use `DEVELOPER` for ADT and integration tests.** `DDIC` does not have the
`S_DEVELOP` authorization object required to create/edit ABAP objects via ADT.

### Unlocking the DEVELOPER User

After many failed login attempts, the `DEVELOPER` user gets locked
(`UFLAG = 128` in `USR02`). This manifests as HTTP 401 from ADT endpoints even
though `/sap/bc/ping` returns 403 (ping uses a lighter auth check).

**Unlock via HANA SQL (no HANA SYSTEM password required):**

The `a4hadm` OS user has a pre-configured HANA userstore key that connects as
the ABAP schema owner (`SAPA4H`):

```bash
docker exec -it a4h bash
su - a4hadm

# Connect to the HANA tenant DB as SAPA4H
hdbsql -U DEFAULT -d HDB

# Unlock DEVELOPER
UPDATE SAPA4H.USR02 SET UFLAG = 0 WHERE BNAME = 'DEVELOPER';

# Verify
SELECT BNAME, UFLAG, PWDSTATE FROM SAPA4H.USR02
  WHERE BNAME IN ('DEVELOPER', 'DDIC');
\q
```

`UFLAG = 0` means unlocked. `UFLAG = 128` means locked by too many failed
logon attempts. `PWDSTATE = 1` means the user must change password on next
login (leave as-is; ADT handles this transparently).

> **Alternative:** If you have access to SAP GUI or ABAP Developer Tools,
> use transaction `SU01` to unlock users without direct HANA access.

---

## HTTPS / Reverse Proxy Setup

Expose the SAP system over HTTPS via Nginx and Let's Encrypt.

### Install Nginx and Certbot

```bash
apt-get install -y nginx certbot python3-certbot-nginx
```

### Configure Nginx Reverse Proxy

Create `/etc/nginx/sites-available/<your-subdomain>`:

```nginx
server {
    listen 80;
    server_name <your-subdomain>;

    location / {
        proxy_pass         http://localhost:50000;
        proxy_set_header   Host              $host;
        proxy_set_header   X-Real-IP         $remote_addr;
        proxy_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
        proxy_set_header   X-Forwarded-Proto $scheme;
        proxy_read_timeout 300s;
        proxy_send_timeout 300s;
        client_max_body_size 50m;
    }
}
```

Enable the site:

```bash
ln -s /etc/nginx/sites-available/<your-subdomain> /etc/nginx/sites-enabled/
nginx -t && systemctl reload nginx
```

### Obtain Let's Encrypt Certificate

```bash
certbot --nginx -d <your-subdomain>
```

Certbot will automatically update the Nginx config with SSL settings and set
up auto-renewal via a systemd timer.

> **DNS propagation:** Run certbot only after the DNS A record has propagated
> (verify with `dig <your-subdomain>`). Let's Encrypt will fail with a challenge
> error if the record hasn't propagated yet.

---

## Integration Tests

### Running Locally

The integration tests are gated by the `integration` build tag and require four
environment variables:

```bash
export SAP_URL=https://<your-subdomain>   # or http://<ip>:50000
export SAP_USER=DEVELOPER
export SAP_PASSWORD='ABAPtr2023#00'
export SAP_CLIENT=001

go test -tags=integration -v -count=1 ./pkg/adt/
```

The tests:
- Create temporary ABAP objects in the `$TMP` package
- Exercise the full ADT API surface (read, write, activate, unit tests, etc.)
- Clean up all created objects after the test run

### Known Test Failures

| Test | Status | Reason |
|------|--------|--------|
| `TestIntegration_RAP_E2E_OData` | Expected FAIL | Requires a pre-existing RAP service binding (`ZTEST_MCP_SB_FLIGHT`) that is not created on a fresh trial system. The test exercises the end-to-end OData pipeline; skip it unless you have the RAP infrastructure set up. |

All other 33 tests should pass on a correctly configured trial system.

### Running a Specific Test

```bash
go test -tags=integration -v -run TestIntegration_CRUD_FullWorkflow ./pkg/adt/
```

---

## GitHub Actions CI

### Workflow Overview

The workflow is defined in `.github/workflows/test.yml`:

```
push / pull_request / workflow_dispatch
      │
      ├── unit (ubuntu-latest)
      │     ├── go test ./... -count=1 -race    ← all unit tests
      │     └── go build ./cmd/vsp              ← verify binary builds
      │
      └── integration (ubuntu-latest) [needs: unit]
            ├── condition: PR, push to main, or manual dispatch
            ├── environment: sap-trial          ← uses GitHub environment secrets
            └── go test -tags=integration -v ./pkg/adt/
```

The integration job only runs when:
- A pull request is opened/updated
- A push lands on `main`
- The workflow is manually dispatched with `run_integration: true`

### GitHub Secrets Setup

Integration tests read credentials from GitHub environment secrets in the
`sap-trial` environment.

**Create the environment:**

```bash
gh api --method PUT repos/<owner>/<repo>/environments/sap-trial
```

**Set the four secrets:**

```bash
gh secret set SAP_URL      --env sap-trial --repo <owner>/<repo> --body "https://<your-subdomain>"
gh secret set SAP_USER     --env sap-trial --repo <owner>/<repo> --body "DEVELOPER"
gh secret set SAP_PASSWORD --env sap-trial --repo <owner>/<repo> --body "ABAPtr2023#00"
gh secret set SAP_CLIENT   --env sap-trial --repo <owner>/<repo> --body "001"
```

**Verify:**

```bash
gh secret list --env sap-trial --repo <owner>/<repo>
```

**Trigger a manual run with integration tests:**

```bash
gh workflow run test.yml --repo <owner>/<repo> --field run_integration=true
```

---

## Troubleshooting

### SAP returns 401 on ADT but 403 on `/sap/bc/ping`

The user is locked (`UFLAG=128`). ADT enforces strict auth and rejects locked
users immediately; the lightweight `/sap/bc/ping` service returns 403 (auth
succeeded but no authorisation) for the same locked user.

Fix: [Unlock the DEVELOPER user via HANA SQL](#unlocking-the-developer-user).

### Integration tests fail with 503 during parallel execution

Too few dialog work processes. [Increase `rdisp/wp_no_dia`](#work-process-tuning)
to at least 15 and restart ABAP.

### `saplikey: profile not found`

List the actual profile files:

```bash
ls /usr/sap/A4H/SYS/profile/
```

Use the `A4H_D00_vhcala4hci` file, not `A4H_DVEBMGS00_vhcala4hci`.

### HANA SYSTEM password unknown

Use the `a4hadm` userstore key instead. It connects as `SAPA4H` (the ABAP
schema owner) without needing the SYSTEM password:

```bash
su - a4hadm
hdbsql -U DEFAULT -d HDB
```

### Build fails: undefined debugger types in integration tests

If you see errors like:
```
pkg/adt/integration_test.go:1642:28: client.GetExternalBreakpoints undefined
```

The `pkg/adt/debugger.go` file is missing from your working tree. Ensure it is
committed and present — it defines the external breakpoint API
(`GetExternalBreakpoints`, `SetExternalBreakpoint`, `BreakpointRequest`, etc.).

### Container starts but SAP is not ready after 10 minutes

Check the container logs:

```bash
docker logs a4h --tail 100
```

Look for HANA startup errors. Common causes:
- Insufficient shared memory (`kernel.shmmax` sysctl not set)
- Disk full during HANA startup

### HTTPS certificate errors in integration tests

If using a self-signed cert or testing against HTTP, set:

```bash
export SAP_INSECURE=true
```

or use the plain HTTP URL (`http://server-ip:50000`). Let's Encrypt certificates
do not require `SAP_INSECURE`.
