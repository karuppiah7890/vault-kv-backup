package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/hashicorp/vault/api"
)

var usage = `
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
`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout, "%s", usage)
		os.Exit(0)
	}
	quietProgress := flag.Bool("quiet", false, "quiet progress")
	vaultKvBackupJsonFileName := flag.String("file", "vault_kv_backup.json", "vault kv backup json file path")
	flag.Parse()

	if !(flag.NArg() == 1) {
		fmt.Fprintf(os.Stderr, "invalid number of arguments: %d. expected 1 argument.\n\n", flag.NArg())
		flag.Usage()
	}

	config := api.DefaultConfig()

	client, err := api.NewClient(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating vault client: %s\n", err)
		os.Exit(1)
	}

	kvMountPath := flag.Args()[0]

	allSecretPathsAndSecrets := walkVaultKvMountPathAndGetSecrets(kvMountPath, "", client, *quietProgress)

	vaultKvBackup := VaultKvBackup{
		Secrets: allSecretPathsAndSecrets,
	}

	vaultKvBackupJSON, err := convertVaultKvBackupToJSON(vaultKvBackup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error converting vault kv backup to json: %s\n", err)
		os.Exit(1)
	}
	err = writeToFile(vaultKvBackupJSON, *vaultKvBackupJsonFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing vault kv backup to json file: %s\n", err)
		os.Exit(1)
	}
}

func walkVaultKvMountPathAndGetSecrets(kvMounthPath, kvSecretsPath string, client *api.Client, quietProgress bool) map[string]map[string]interface{} {
	logicalClient := client.Logical()

	listPath := path.Join(kvMounthPath, "metadata", kvSecretsPath)

	kvSecrets, err := logicalClient.List(listPath)
	if err != nil {
		log.Fatalf("error occurred while listing metadata at path `%s`: %v", listPath, err)
	}

	if kvSecrets == nil {
		if !quietProgress {
			fmt.Fprintf(os.Stdout, "getting secrets at `%s`\n\n", kvSecretsPath)
		} else {
			fmt.Fprintf(os.Stdout, ".")
		}
		return map[string]map[string]interface{}{
			kvSecretsPath: getSecrets(kvMounthPath, kvSecretsPath, client),
		}
	}

	data := kvSecrets.Data
	if data == nil {
		log.Fatalf("no data found at path `%s`", listPath)
	}

	keys, ok := data["keys"]
	if !ok {
		log.Fatalf("no data found at path `%s`", listPath)
	}

	combinedSecretPathsAndSecrets := map[string]map[string]interface{}{}

	// TODO: `keys` can be of non-array type too. So, type assertion is required.
	// No problems if `keys` is of array type.
	// If `keys` is not an array, then it will panic. So, handle this issue.
	for _, key := range keys.([]interface{}) {
		newKvSecretsPath := path.Join(kvSecretsPath, key.(string))
		secretPathsAndSecrets := walkVaultKvMountPathAndGetSecrets(kvMounthPath, newKvSecretsPath, client, quietProgress)

		for secretsPath, secrets := range secretPathsAndSecrets {
			combinedSecretPathsAndSecrets[secretsPath] = secrets
		}
	}

	return combinedSecretPathsAndSecrets
}

func getSecrets(kvMounthPath, kvSecretsPath string, client *api.Client) map[string]interface{} {
	kvClient := client.KVv2(kvMounthPath)
	kvSecrets, err := kvClient.Get(context.TODO(), kvSecretsPath)

	if err != nil {
		log.Fatalf("error occurred while getting latest version of the secret at path `%s`: %v", kvSecretsPath, err)
	}

	if kvSecrets == nil {
		log.Fatalf("no secret found at path `%s`", kvSecretsPath)
	}

	return kvSecrets.Data
}
