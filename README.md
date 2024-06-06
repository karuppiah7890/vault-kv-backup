# vault-kv-backup

Using this CLI tool, you can backup Vault KV v2 Secrets Engine Secrets from a Vault instance to your local machine! :D

Note: The tool is written in Golang and uses the Vault Official Golang API. The Official Vault Golang API documentation is here - https://pkg.go.dev/github.com/hashicorp/vault/api

Note: The tool needs Vault credentials of a user/account that has access to Vault, to read and list the Vault KV v2 Secrets Engine Secrets that you want to backup. Look at [Authorization Details for the Vault Credentials](#authorization-details-for-the-vault-credentials) for more details

Note: We have tested this only with some versions of Vault (like v1.15.x). So beware to test this in a testing environment with whatever version of Vault you are using, before using this in critical environments like production! Also, ensure that the testing environment is as close to your production environment as possible so that your testing makes sense

## Building

```bash
CGO_ENABLED=0 go build -v
```

or

```bash
make
```

## Authorization Details for the Vault Credentials

As mentioned before in a note, the tool needs Vault credentials of a user/account that has access to Vault, to read and list the Vault KV v2 Secrets Engine Secrets that you want to backup.

An example Vault Policy that's required to backup all the secrets in a Vault KV v2 Secrets Engine is -

```hcl
# Vault KV v2 Secrets Engine mount path is "secret"
path "secret/*" {
  capabilities = ["read", "list"]
}
```

You can use a similar Vault Policy based on the mount path of the Vault KV v2 Secrets Engine that you are using and want to backup. You can create a Vault Token that has this Vault Policy attached to it and use that token to backup the Vault KV v2 Secrets Engine Secrets using the `vault-kv-backup` tool :)

## Usage

```bash
$ ./vault-kv-backup --help

usage: vault-kv-backup [-quiet|--quiet] [-file|-file <vault-kv-backup-json-file-path>] <kv-mount-path>

Note that the flags MUST come before the arguments

arguments of ./vault-kv-backup:

  <kv-mount-path> string
    vault kv v2 secrets engine mount path for backing up the
    vault kv v2 secrets engine secrets present in that mount
    path

flags of ./vault-kv-backup:

  -file / --file string
      vault kv backup json file path (default "vault_kv_backup.json")

  -quiet / --quiet
      quiet progress (default false).
      By default vault-kv-backup CLI will show a lot of details
      about the backup process and detailed progress during the
      backup process

  -h / -help / --help
      show help

examples:

# show help
vault-kv-backup -h
vault-kv-backup --help

# backs up all vault KV v2 Secrets Engine Secrets to the JSON file
vault-kv-backup -file <path-to-vault-kv-backup-json-file> <kv-mount-path>

# OR you can use --file too instead of -file

vault-kv-backup --file <path-to-vault-kv-backup-json-file> <kv-mount-path>

# quietly backs up all vault KV v2 Secrets Engine Secrets to the JSON file
# this will just show dots (.) for progress
vault-kv-backup -quiet -file <path-to-vault-kv-backup-json-file> <kv-mount-path>

# OR you can use --quiet too instead of -quiet

vault-kv-backup --quiet --file <path-to-vault-kv-backup-json-file> <kv-mount-path>
```

# Demo

I created a new dummy local Vault instance in developer mode for this demo. I ran the Vault server like this -

```bash
vault server -dev -dev-root-token-id root -dev-listen-address 127.0.0.1:8200
```

I'm going to create two dummy demo secrets in the KV v2 Secrets Engine, which is mounted at `secret/`. I'll be using the Vault CLI to do this but you can do it in any way you want. I'll be using the `root` Vault API token to create the two secrets, but it's not necessary to use the root token, you can use any token with less privileges too, following the Principle of Least Privilege. Ensure that your token is safe and secure, regardless of it being root token or not

Initially the Vault looks like this -

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault secrets list
Path          Type         Accessor              Description
----          ----         --------              -----------
cubbyhole/    cubbyhole    cubbyhole_b6e93602    per-token private secret storage
identity/     identity     identity_53f9b46c     identity store
secret/       kv           kv_af2e5d33           key/value secret storage
sys/          system       system_3d3d77e6       system endpoints used for control, policy and debugging

$ vault kv list secret
No value found at secret/metadata

$ vault kv list -mount=secret
No value found at secret/metadata
```

Let's put in some data

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault kv put -mount=secret foo bar=baz blah=bloo blee=bley

# OR you can also use

$ vault kv put secret/foo bar=baz blah=bloo blee=bley

$ cat dummy-data.json
{
  "something": {
    "another-thing": {
      "yet-another-thing": {
        "and-then-something": "okay",
        "and-the-one-more-thing": "haha, okay, right!"
      }
    }
  }
}

$ vault kv put -mount secret something/over/here @dummy-data.json

# OR you can also use

$ vault kv put secret/something/over/here @dummy-data.json
```

When you run the above, you will find output like this -

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault kv put -mount=secret foo bar=baz blah=bloo blee=bley
= Secret Path =
secret/data/foo

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-05T17:56:31.362351Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

$ vault kv put -mount secret something/over/here @dummy-data.json
========= Secret Path =========
secret/data/something/over/here

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-05T17:56:36.019262Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1
```

By the way, all this is just for demonstrating (Demo) and teaching purposes only. I'm not an expert in all the nitty gritty details but I'll do my best :) :D

Let's look at the secret data we put / stored inside the Vault KV v2 Secrets Engine. And then we can do a backup :)

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault kv get -mount secret foo
= Secret Path =
secret/data/foo

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-05T17:56:31.362351Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

==== Data ====
Key     Value
---     -----
bar     baz
blah    bloo
blee    bley

$ vault kv get -format json -mount secret foo
{
  "request_id": "f0df730c-6866-a610-4060-dccc86dd80f1",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "data": {
      "bar": "baz",
      "blah": "bloo",
      "blee": "bley"
    },
    "metadata": {
      "created_time": "2024-06-05T17:56:31.362351Z",
      "custom_metadata": null,
      "deletion_time": "",
      "destroyed": false,
      "version": 1
    }
  },
  "warnings": null
}

$ vault kv get -mount secret something/over/here
========= Secret Path =========
secret/data/something/over/here

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-05T17:56:36.019262Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

====== Data ======
Key          Value
---          -----
something    map[another-thing:map[yet-another-thing:map[and-the-one-more-thing:haha, okay, right! and-then-something:okay]]]

$ vault kv get -format json -mount secret something/over/here
{
  "request_id": "4d3f9dcd-71e7-0f7b-5a42-f1852fb75007",
  "lease_id": "",
  "lease_duration": 0,
  "renewable": false,
  "data": {
    "data": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    },
    "metadata": {
      "created_time": "2024-06-05T17:56:36.019262Z",
      "custom_metadata": null,
      "deletion_time": "",
      "destroyed": false,
      "version": 1
    }
  },
  "warnings": null
}
```

So, we have two examples of secrets here. One is a secret with nested JSON values and another a simple set of key-value pairs which is flat without any nesting and has only one level of keys and immediate scalar values for each key and not some object/map value

Now let's create a token which has the least privilege to read secrets from the Vault KV v2 Secrets Engine mounted at `secret/` mount path

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ cat /Users/karuppiah.n/every-day-log/allow_read_secrets.hcl
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["read", "list"]
}

$ vault policy write read_kv_secrets /Users/karuppiah.n/every-day-log/allow_read_secrets.hcl
Success! Uploaded policy: read_kv_secrets

$ vault policy read read_kv_secrets
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["read", "list"]
}

$ vault token create -policy read_kv_secrets
Key                  Value
---                  -----
token                hvs.CAESIL6jlmcSpvGz5Y_Muae4mkjWyFyIrWbMUGCDRayNskCPGh4KHGh2cy5GOVhQZk1uUVpraFcxa2dnMjA2MGh4WGI
token_accessor       VLxax0wwDSfwGI5tpasCcbAR
token_duration       768h
token_renewable      true
token_policies       ["default" "read_kv_secrets"]
identity_policies    []
policies             ["default" "read_kv_secrets"]
```

Note that the above token has two Vault policies attached to it - one is `default` policy and another is our custom policy `read_kv_secrets`. You can choose to modify the `default` policy to ensure how much access you want to give by default to a token. In this case, I'm fine with whatever `default` policy Vault is giving by default

If you don't want the token to have the `default` policy attached to it, you can use `-no-default-policy` flag while creating the token. It will look something like this -

```bash
$ vault token create -no-default-policy -policy read_kv_secrets
Key                  Value
---                  -----
token                hvs.CAESIBK45gan5zjmAEm1Lg-mfu-RICWB5zOcga3CG8pnPZRbGh4KHGh2cy5NdEc4eVdKOE0zQldmcUI2SVdUN2lFTk8
token_accessor       UvhVu6SvrZIQ5p9hcS72g6hs
token_duration       768h
token_renewable      true
token_policies       ["read_kv_secrets"]
identity_policies    []
policies             ["read_kv_secrets"]
```

Now, let's use the first token we created to create a backup of all the Secrets in the Vault KV v2 Secrets Engine mounted at `secret/`

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="hvs.CAESIL6jlmcSpvGz5Y_Muae4mkjWyFyIrWbMUGCDRayNskCPGh4KHGh2cy5GOVhQZk1uUVpraFcxa2dnMjA2MGh4WGI"

$ vault token lookup
Key                 Value
---                 -----
accessor            VLxax0wwDSfwGI5tpasCcbAR
creation_time       1717610614
creation_ttl        768h
display_name        token
entity_id           n/a
expire_time         2024-07-07T23:33:34.284724+05:30
explicit_max_ttl    0s
id                  hvs.CAESIL6jlmcSpvGz5Y_Muae4mkjWyFyIrWbMUGCDRayNskCPGh4KHGh2cy5GOVhQZk1uUVpraFcxa2dnMjA2MGh4WGI
issue_time          2024-06-05T23:33:34.284728+05:30
meta                <nil>
num_uses            0
orphan              false
path                auth/token/create
policies            [default read_kv_secrets]
renewable           true
ttl                 767h52m58s
type                service

$ vault token capabilities secret
deny

$ vault token capabilities secret/
list, read

$ vault token capabilities secret/*
list, read

$ ./vault-kv-backup -quiet -file my_secret_backup.json secret

$ cat my_secret_backup.json
{"secrets":{"foo":{"bar":"baz","blah":"bloo","blee":"bley"},"something/over/here":{"something":{"another-thing":{"yet-another-thing":{"and-the-one-more-thing":"haha, okay, right!","and-then-something":"okay"}}}}}}

$ cat my_secret_backup.json | jq
{
  "secrets": {
    "foo": {
      "bar": "baz",
      "blah": "bloo",
      "blee": "bley"
    },
    "something/over/here": {
      "something": {
        "another-thing": {
          "yet-another-thing": {
            "and-the-one-more-thing": "haha, okay, right!",
            "and-then-something": "okay"
          }
        }
      }
    }
  }
}
```

As you can see, all the secrets in the Vault KV v2 Secrets Engine mounted at `secret/` mount path have been backed up :) The secrets have the secret path and the secret content / data itself :)

# Possible Errors

There are quite some possible errors you can face. Mostly relating to one of the following

- DNS Resolution issues. If you are accessing Vault using it's domain name (DNS record), and not an IP address, then ensure that the DNS resolution works well
- Connectivity issues with Vault. Ensure you have good network connectivity to the Vault system. Ensure the IP you are connecting to is right and belong to the Vault API server, and also check the API server port too.
- Access / Authorization issues. Ensure you have enough access to list and read the Vault KV v2 Secrets Engine that you want to take backup of

Example access errors / authorization errors / permission errors -

I'll create a token with half baked access and not full access, let's see how that goes -

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="root"

$ vault policy write half_baked_read_access1 -
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["list"]
}

$ vault policy read half_baked_read_access1
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["list"]
}

$ vault token create -no-default-policy -policy half_baked_read_access1
Key                  Value
---                  -----
token                hvs.CAESIN-lzzO7k9Ohmixj7U-z3yLgvQ0hJwR0Il2kIlc1jvTAGh4KHGh2cy4zcEE4MzVUa2dCZ3N3VmFuUU1WV0wycUU
token_accessor       2gLBsSST5pAJ1F7hw1d5W6v2
token_duration       768h
token_renewable      true
token_policies       ["half_baked_read_access1"]
identity_policies    []
policies             ["half_baked_read_access1"]
```

Now, let's use this token to do a backup of secrets in the Vault KV v2 Secrets Engine mounted at `secret/`

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="hvs.CAESIN-lzzO7k9Ohmixj7U-z3yLgvQ0hJwR0Il2kIlc1jvTAGh4KHGh2cy4zcEE4MzVUa2dCZ3N3VmFuUU1WV0wycUU"

$ ./vault-kv-backup -quiet -file my_another_secret_backup.json secret
.2024/06/05 23:53:32 error occurred while getting latest version of the secret at path `foo`: error encountered while reading secret at secret/data/foo: Error making API request.

URL: GET http://127.0.0.1:8200/v1/secret/data/foo
Code: 403. Errors:

* 1 error occurred:
	* permission denied
```

Notice how it was able to list the secrets and understand that there's `foo` in it, but not able to read it? This is because of the lack of access. Similar example below -

```bash
$ vault policy write half_baked_read_access2 -
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["read"]
}
Success! Uploaded policy: half_baked_read_access2

$ vault policy read half_baked_read_access2
# KV v2 secrets engine mount path is "secret"
path "secret/*" {
  capabilities = ["read"]
}

$ vault token create -no-default-policy -policy half_baked_read_access2
Key                  Value
---                  -----
token                hvs.CAESIIxGBE-N0fR6OlHa_yz-r2cNqIJptc97R07HDgZFKwWqGh4KHGh2cy5Sazd0NE5FODc5bXhvZmhRS005RmNBbnk
token_accessor       MJxiU0eHd0eHs9FgVUAfWrRr
token_duration       768h
token_renewable      true
token_policies       ["half_baked_read_access2"]
identity_policies    []
policies             ["half_baked_read_access2"]
```

Now, let's use this token to do a backup of secrets in the Vault KV v2 Secrets Engine mounted at `secret/`

```bash
$ export VAULT_ADDR='http://127.0.0.1:8200'
$ export VAULT_TOKEN="hvs.CAESIIxGBE-N0fR6OlHa_yz-r2cNqIJptc97R07HDgZFKwWqGh4KHGh2cy5Sazd0NE5FODc5bXhvZmhRS005RmNBbnk"

$ ./vault-kv-backup -quiet -file my_another_secret_backup.json secret
2024/06/06 00:02:35 error occurred while listing metadata at path `secret/metadata`: Error making API request.

URL: GET http://127.0.0.1:8200/v1/secret/metadata?list=true
Code: 403. Errors:

* 1 error occurred:
	* permission denied

$ vault kv get -mount secret foo
= Secret Path =
secret/data/foo

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-05T18:20:42.190213Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

==== Data ====
Key     Value
---     -----
bar     baz
blah    bloo
blee    bley

$ vault kv get -mount secret something/over/here
========= Secret Path =========
secret/data/something/over/here

======= Metadata =======
Key                Value
---                -----
created_time       2024-06-05T18:20:52.733659Z
custom_metadata    <nil>
deletion_time      n/a
destroyed          false
version            1

====== Data ======
Key          Value
---          -----
something    map[another-thing:map[yet-another-thing:map[and-the-one-more-thing:haha, okay, right! and-then-something:okay]]]
```

As you can see, it can read secrets in Vault KV v2 Secrets Engine at `secret/` mount path. But it can't list stuff, like list metadata, which is important for us to know all the secrets present in the mount path and inside it in a recursive manner, so the `vault-kv-backup` tool fails to work without that access / permission and that's clear from the above error too

Also, if there's not enough permissions in other ways also it won't work - like, access to list and read some secrets but not others, all within the same Vault KV v2 Secrets Engine mounted at a particular path. Also, if there's some explicit `deny` on some secret / secrets or secret path / secrets path, then that will also cause problems for the tool. As of now, the tool abruptly stops when there's such an error

# Future Ideas

- Any and all issues in the [GitHub Issues](https://github.com/karuppiah7890/vault-kv-backup/issues) section

- Allow user to say "It's okay if the tool cannot backup some secrets and/ some secret paths, due to permission issues. Just backup the secrets the tool can" and be able to skip intermittent errors here and there and ignore the errors than abruptly stop at errors

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

- Allow for flags to come even after the arguments. This would require using better CLI libraries / frameworks like [`spf13/pflag`](https://github.com/spf13/pflag), [`spf13/cobra`](https://github.com/spf13/cobra), [`urfave/cli`](https://github.com/urfave/cli) etc than using the basic Golang's built-in `flag` standard library / package

# Contributing

Please look at https://github.com/karuppiah7890/vault-tooling-contributions for some basic details on how you can contribute
