# cgbot


Bot for CyGame Telegram chats.

## Skip deploy

Add `!skip` to commit message.

## Migrations

```bash
curl -sSf https://atlasgo.sh | sh
```

### Add migration

To add migration named `some-migration-name`:

```console
atlas migrate --env dev diff some-migration-name
```

## Golden files

In package directory:

```console
go test -update
```
