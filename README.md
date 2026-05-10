# ArgonVault

A local, offline secret manager for the terminal. Secrets are encrypted with
**AES-256-GCM** under a key derived from your master password using
**Argon2id**. Nothing ever leaves your machine — no server, no cloud, no
telemetry. Without your master password, the vault is mathematically
unreadable.

## Features

- AES-256-GCM authenticated encryption for every secret
- Argon2id key derivation (memory-hard, GPU-resistant) with a unique salt per vault
- SQLite storage under your platform's user config directory (WAL mode, foreign keys on)
- Multiple named vaults, each with independent salts
- Atomic master-password rotation that re-encrypts every secret across every vault
- JSON export of a vault's decrypted contents
- Tamper-evident audit log of every action

## Install

Requires Go 1.25+.

```sh
git clone https://github.com/HIJOdelIDANII/ArgonVault.git
cd ArgonVault
go build -o argonv .
```

Move the resulting `argonv` binary somewhere on your `PATH`.

## Usage

```sh
argonv init   --name personal                 # create a vault
argonv store  --vault personal --name aws_key --data AKIA...
argonv get    --vault personal --name aws_key
argonv list                                   # list vaults
argonv list   --vault personal                # list secrets in a vault
argonv delete --vault personal --name aws_key
argonv delete --vault personal                # delete the vault itself
argonv export --vault personal --output personal.json
argonv logs                                   # audit trail
argonv logs   --vault personal
argonv rotate-key                             # change the master password
```

You will be prompted for the master password the first time it is needed in a
session.

## Storage location

The SQLite database lives under your OS's user config directory:

| OS     | Path                                              |
| ------ | ------------------------------------------------- |
| Linux  | `~/.config/argonvault/storage.db`                 |
| macOS  | `~/Library/Application Support/argonvault/storage.db` |

The directory is created with `0700` permissions.

## Security notes

- Keys are derived per-vault: `Argon2id(password, vault_salt) -> 32-byte key`
- Each secret encrypts with a fresh 12-byte random nonce; GCM provides integrity
- Master-password rotation generates a new salt for every vault and re-encrypts every secret in a single SQLite transaction — either everything rotates or nothing does
- Wrong passwords surface as authenticated-decryption failures, not silent corruption

## License

MIT
