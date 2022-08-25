// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the EARLY ACCESS SOFTWARE LICENSE AGREEMENT
// available at http://github.com/namespacelabs/foundation

package eks

import (
	"context"
	"encoding/base64"
	"os"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"namespacelabs.dev/foundation/internal/fnerrors"
	"namespacelabs.dev/foundation/schema"
)

func Kubeconfig(awsCluster *AwsCluster, env *schema.Environment) (*clientcmdapi.Config, error) {
	cluster := awsCluster.Cluster
	if cluster.Name == nil {
		return nil, fnerrors.BadInputError("cluster name is missing")
	}

	if cluster.Endpoint == nil {
		return nil, fnerrors.BadInputError("cluster endpoint is missing")
	}

	if cluster.CertificateAuthority == nil || cluster.CertificateAuthority.Data == nil {
		return nil, fnerrors.BadInputError("cluster certificateauthority is missing")
	}

	decoded, err := base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, fnerrors.BadInputError("failed to parse certificateauthority: %w", err)
	}

	clusterName := *cluster.Name
	contextName := clusterName
	return &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server:                   *cluster.Endpoint,
				CertificateAuthorityData: decoded,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: contextName,
			},
		},
		CurrentContext: contextName,
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			contextName: {
				Exec: &clientcmdapi.ExecConfig{
					APIVersion: "client.authentication.k8s.io/v1",
					Command:    os.Args[0],
					Args: []string{"eks", "generate-token", "--exec_credential",
						"--env", env.Name, clusterName},
					Env:             []clientcmdapi.ExecEnvVar{},
					InteractiveMode: clientcmdapi.NeverExecInteractiveMode,
				},
			},
		},
	}, nil
}

func KubeconfigFromCluster(ctx context.Context, s *Session, clusterName string) (*clientcmdapi.Config, error) {
	cluster, err := DescribeCluster(ctx, s, clusterName)
	if err != nil {
		return nil, err
	}

	return Kubeconfig(cluster, s.env)
}
