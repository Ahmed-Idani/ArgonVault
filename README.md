```
       █████╗ ██████╗  ██████╗  ██████╗ ███╗   ██╗██╗   ██╗ █████╗ ██╗   ██╗██╗  ████████╗
      ██╔══██╗██╔══██╗██╔════╝ ██╔═══██╗████╗  ██║██║   ██║██╔══██╗██║   ██║██║  ╚══██╔══╝
      ███████║██████╔╝██║  ███╗██║   ██║██╔██╗ ██║██║   ██║███████║██║   ██║██║     ██║
      ██╔══██║██╔══██╗██║   ██║██║   ██║██║╚██╗██║╚██╗ ██╔╝██╔══██║██║   ██║██║     ██║
      ██║  ██║██║  ██║╚██████╔╝╚██████╔╝██║ ╚████║ ╚████╔╝ ██║  ██║╚██████╔╝███████╗██║
      ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝  ╚═════╝ ╚═╝  ╚═══╝  ╚═══╝  ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝
```

### **A zero-trust, offline-first secrets manager for your terminal.**
### *AES-256-GCM encryption · Argon2id key derivation · pure-Go binary.*

---

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![SQLite](https://img.shields.io/badge/SQLite-WAL-003B57?style=for-the-badge&logo=sqlite&logoColor=white)
![AES-256-GCM](https://img.shields.io/badge/AES--256--GCM-Authenticated-1f6feb?style=for-the-badge&logo=letsencrypt&logoColor=white)
![Argon2id](https://img.shields.io/badge/Argon2id-Memory--Hard-6f42c1?style=for-the-badge&logo=keybase&logoColor=white)
![CLI](https://img.shields.io/badge/CLI-Cobra-0a7d2c?style=for-the-badge&logo=gnubash&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-black?style=for-the-badge)
![Platform](https://img.shields.io/badge/Linux%20·%20macOS-supported-111?style=for-the-badge&logo=linux&logoColor=white)

---

## Overview

**ArgonVault** is a security-first command-line secrets manager written in Go.
It stores credentials, tokens, API keys, and configuration secrets in a local
SQLite database, encrypted with **AES-256-GCM** under keys derived from your
master password using **Argon2id**.

There is no server, no cloud sync, no telemetry, no daemon. The vault is a
single file on your disk. Without your master password, that file is
mathematically unreadable — even to you.

ArgonVault is built for engineers who want the convenience of a CLI password
manager without trusting a third party with the cryptographic root of trust.

---

## Features

- **Strong cryptography by default** — AES-256-GCM authenticated encryption with a fresh 96-bit nonce per write.
- **Argon2id key derivation** — memory-hard, GPU-resistant, parameters tuned for interactive use.
- **Per-vault salts** — every vault carries an independent random salt; key reuse across vaults is impossible.
- **Atomic master-password rotation** — re-encrypts every secret in every vault inside a single SQLite transaction.
- **Audit log** — every action (`STORE`, `GET`, `DELETE`, `ROTATE`, `EXPORT`) is recorded with status and timestamp.
- **Encrypted JSON export** — exfiltrate a vault as readable JSON when you need to migrate or back up.
- **Modern terminal UI** — colored output, aligned tables, sensible error messages.
- **Single static binary** — no runtime, no daemon, no background service.
- **Offline by design** — zero network code paths.

---

## Tech Stack

| Layer            | Technology                                    |
| ---------------- | --------------------------------------------- |
| Language         | Go 1.25+                                      |
| Symmetric Cipher | AES-256-GCM (`crypto/aes` + `crypto/cipher`)  |
| KDF              | Argon2id (`golang.org/x/crypto/argon2`)       |
| Storage          | SQLite (WAL, foreign keys, busy timeout)      |
| CLI Framework    | Cobra                                         |
| Terminal UI      | Lipgloss / charmbracelet                      |
| Driver           | `mattn/go-sqlite3`                            |

---

## Installation

### Build from source

```sh
git clone https://github.com/HIJOdelIDANII/ArgonVault.git
cd ArgonVault
go build -o argonv .
```

Move the resulting `argonv` binary somewhere on your `PATH`:

```sh
sudo mv argonv /usr/local/bin/
```

### Verify

```sh
argonv --help
```

---

## Usage

### Initialize a vault

```sh
argonv init --name personal
```

You will be prompted to set your master password the first time it is needed.

### Store a secret

```sh
argonv store --vault personal --name aws_access_key --data AKIAIOSFODNN7EXAMPLE
```

### Retrieve a secret

```sh
argonv get --vault personal --name aws_access_key
```

### List vaults and secrets

```sh
argonv list                       # list all vaults
argonv list --vault personal      # list secrets in a vault
```

### Delete

```sh
argonv delete --vault personal --name aws_access_key   # delete a secret
argonv delete --vault personal                         # delete the entire vault
```

### Export a vault as JSON

```sh
argonv export --vault personal --output personal.json
```

### Rotate the master password

```sh
argonv rotate-key
```

Every vault gets a new salt, every secret is re-encrypted, all atomically.

### Audit trail

```sh
argonv logs                       # full history
argonv logs --vault personal      # filter by vault
```

---

## Security Model

ArgonVault is engineered around a **zero-trust local threat model**: assume an
attacker may eventually obtain a copy of your vault file. The only secret that
must never leak is your master password.

### Encryption

- **Cipher:** AES-256-GCM — authenticated encryption with associated data.
- **Nonce:** 96-bit random IV generated from `crypto/rand` for **every write**.
- **Integrity:** GCM authentication tag detects any tampering with ciphertext or IV; decryption fails closed.

### Key Derivation

- **Function:** Argon2id — winner of the Password Hashing Competition, resistant to GPU and side-channel attacks.
- **Parameters:** `time=3`, `memory=64 MiB`, `threads=4`, `keyLen=32` — interactive-cost defaults.
- **Salt:** 16 random bytes per **vault**, stored alongside the vault row. Salts are not secret, but uniqueness guarantees no rainbow-table or cross-vault key reuse.
- **Key:** `key = Argon2id(masterPassword, vaultSalt)` — the key never touches disk.

### Key Handling

- Keys exist only in process memory for the duration of a single command.
- The master password is read from the controlling TTY with echo disabled.
- No swap-resident caching, no key-agent process, no daemon — each invocation derives keys from scratch.

### Storage

- SQLite database under the OS user-config directory (`~/.config/argonvault/storage.db` on Linux), created with `0700` permissions.
- WAL journal mode for crash safety.
- Foreign keys enforced — orphan secrets are impossible.

### What ArgonVault does **not** protect against

- A compromised host with a keylogger.
- An attacker who already has root and can read process memory.
- A weak master password — Argon2id raises the cost, but cannot rescue `password123`.

---

## Architecture

```
                          ┌─────────────────────────────┐
                          │           cmd/              │
                          │  Cobra commands (CLI verbs) │
                          │  init · store · get · …     │
                          └──────────────┬──────────────┘
                                         │
                          ┌──────────────▼──────────────┐
                          │        internal/            │
                          │                             │
                          │  auth.go     ── master pwd  │
                          │  password.go ── TTY prompt  │
                          │  crypto.go   ── AES-GCM +   │
                          │                 Argon2id    │
                          │  storage.go  ── CRUD layer  │
                          │  rotate.go   ── atomic re-  │
                          │                 encryption  │
                          │  models.go   ── schema      │
                          │  dbconn.go   ── SQLite open │
                          │  ui/         ── lipgloss    │
                          └──────────────┬──────────────┘
                                         │
                          ┌──────────────▼──────────────┐
                          │       SQLite (WAL)          │
                          │  vaults · secrets · configs │
                          │  · audit_logs               │
                          └─────────────────────────────┘
```

**Data flow for a `store` command**

1. Cobra parses flags and dispatches to the `store` handler.
2. `EnsureMasterPassword` prompts the user once per process.
3. `GetVault` loads the target vault and its salt.
4. `Encrypt` derives `Argon2id(password, salt)`, generates a 96-bit IV, seals the plaintext with AES-256-GCM.
5. `SaveSecret` writes the ciphertext + IV inside a transaction.
6. `LogAction` appends a tamper-evident audit entry.

Every write path is symmetric to a corresponding read path; failures are logged with `FAILURE` status and surface as typed errors.

---

## License

Released under the **MIT License**. See [`LICENSE`](LICENSE) for details.

---

> Built with care for engineers who would rather hold their own keys.
