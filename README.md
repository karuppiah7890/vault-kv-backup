# vault-kv-backup

# Future Ideas

- Support backing up multiple specific KV v2 secrets engine secrets in a single backup at once by providing a file which contains the mount paths of the secrets engines to be backed up, or by providing the mount paths of the secrets engines as arguments to the CLI, or provide the ability to use either of the two or even both

- Support backing up all the secrets in all the secrets engines in a single backup

- Support backing up the KV v2 secrets engine configuration too apart from the secrets in the secrets engine. https://developer.hashicorp.com/vault/api-docs/secret/kv/kv-v2#configure-the-kv-engine

- Support backing up the metadata of the secrets apart from the secrets the secrets engine. This makes restore a bit tricky - we need to know for which secrets we want to restore the metadata too - or we restore metadata for all or none

- Support backing up all the versions of the secrets in the secrets engine. This makes restore a bit tricky - we need to know which version to restore, or we could restore all versions to start off with and put the latest version as the latest

- Support backing up things with a combination of the above features, that is
  - multiple/all secrets
  - one/multiple/all secrets engines
  - only latest/all versions of secrets
  - metadata of secrets
  - configuration of secrets engines

# Contributing

Please look at https://github.com/karuppiah7890/vault-tooling-contributions for some basic details on how you can contribute
