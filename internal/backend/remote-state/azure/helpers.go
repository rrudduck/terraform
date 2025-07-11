// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// This file is copied from terraform-provider-azurerm: internal/provider/helpers.go

package azure

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform/internal/backend/backendbase"
)

// logEntry avoids log entries showing up in test output
func logEntry(f string, v ...interface{}) {
	if os.Getenv("TF_LOG") == "" {
		return
	}

	if os.Getenv("TF_ACC") != "" {
		return
	}

	log.Printf(f, v...)
}

func decodeCertificate(clientCertificate string) ([]byte, error) {
	var pfx []byte
	if clientCertificate != "" {
		out := make([]byte, base64.StdEncoding.DecodedLen(len(clientCertificate)))
		n, err := base64.StdEncoding.Decode(out, []byte(clientCertificate))
		if err != nil {
			return pfx, fmt.Errorf("could not decode client certificate data: %v", err)
		}
		pfx = out[:n]
	}
	return pfx, nil
}

func getOidcToken(d *backendbase.SDKLikeData) (*string, error) {
	idToken := strings.TrimSpace(d.String("oidc_token"))

	if path := d.String("oidc_token_file_path"); path != "" {
		fileTokenRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading OIDC Token from file %q: %v", path, err)
		}

		fileToken := strings.TrimSpace(string(fileTokenRaw))

		if idToken != "" && idToken != fileToken {
			return nil, fmt.Errorf("mismatch between supplied OIDC token and supplied OIDC token file contents - please either remove one or ensure they match")
		}

		idToken = fileToken
	}

	if d.Bool("use_aks_workload_identity") && os.Getenv("AZURE_FEDERATED_TOKEN_FILE") != "" {
		path := os.Getenv("AZURE_FEDERATED_TOKEN_FILE")
		fileTokenRaw, err := os.ReadFile(os.Getenv("AZURE_FEDERATED_TOKEN_FILE"))

		if err != nil {
			return nil, fmt.Errorf("reading OIDC Token from file %q provided by AKS Workload Identity: %v", path, err)
		}

		fileToken := strings.TrimSpace(string(fileTokenRaw))

		if idToken != "" && idToken != fileToken {
			return nil, fmt.Errorf("mismatch between supplied OIDC token and OIDC token file contents provided by AKS Workload Identity - please either remove one, ensure they match, or disable use_aks_workload_identity")
		}

		idToken = fileToken
	}

	return &idToken, nil
}

func getClientId(d *backendbase.SDKLikeData) (*string, error) {
	clientId := strings.TrimSpace(d.String("client_id"))

	if path := d.String("client_id_file_path"); path != "" {
		fileClientIdRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading Client ID from file %q: %v", path, err)
		}

		fileClientId := strings.TrimSpace(string(fileClientIdRaw))

		if clientId != "" && clientId != fileClientId {
			return nil, fmt.Errorf("mismatch between supplied Client ID and supplied Client ID file contents - please either remove one or ensure they match")
		}

		clientId = fileClientId
	}

	if d.Bool("use_aks_workload_identity") && os.Getenv("AZURE_CLIENT_ID") != "" {
		aksClientId := os.Getenv("AZURE_CLIENT_ID")
		if clientId != "" && clientId != aksClientId {
			return nil, fmt.Errorf("mismatch between supplied Client ID and that provided by AKS Workload Identity - please remove, ensure they match, or disable use_aks_workload_identity")
		}
		clientId = aksClientId
	}

	return &clientId, nil
}

func getClientSecret(d *backendbase.SDKLikeData) (*string, error) {
	clientSecret := strings.TrimSpace(d.String("client_secret"))

	if path := d.String("client_secret_file_path"); path != "" {
		fileSecretRaw, err := os.ReadFile(path)

		if err != nil {
			return nil, fmt.Errorf("reading Client Secret from file %q: %v", path, err)
		}

		fileSecret := strings.TrimSpace(string(fileSecretRaw))

		if clientSecret != "" && clientSecret != fileSecret {
			return nil, fmt.Errorf("mismatch between supplied Client Secret and supplied Client Secret file contents - please either remove one or ensure they match")
		}

		clientSecret = fileSecret
	}

	return &clientSecret, nil
}

func getTenantId(d *backendbase.SDKLikeData) (*string, error) {
	tenantId := strings.TrimSpace(d.String("tenant_id"))

	if d.Bool("use_aks_workload_identity") && os.Getenv("AZURE_TENANT_ID") != "" {
		aksTenantId := os.Getenv("AZURE_TENANT_ID")
		if tenantId != "" && tenantId != aksTenantId {
			return nil, fmt.Errorf("mismatch between supplied Tenant ID and that provided by AKS Workload Identity - please remove, ensure they match, or disable use_aks_workload_identity")
		}
		tenantId = aksTenantId
	}

	return &tenantId, nil
}
