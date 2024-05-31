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
usage: vault-kv-backup <kv-mount-path>

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

func walkVaultKvMountPathAndGetSecrets(kvMounthPath, kvSecretsPath string, client *api.Client, quietProgress bool) map[string]interface{} {
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
		return map[string]interface{}{
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

	combinedSecretPathsAndSecrets := map[string]interface{}{}

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
