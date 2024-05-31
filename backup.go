package main

type VaultKvBackup struct {
	Secrets map[string]interface{} `json:"secrets"`
}

func convertVaultKvBackupToJSON(vaultKvBackup VaultKvBackup) ([]byte, error) {
	vaultKvBackupJSON, err := toJSON(vaultKvBackup)
	if err != nil {
		return nil, err
	}
	return vaultKvBackupJSON, nil
}
