package keycloak

import pb "golang-grpc-sample/proto"

type KcProvision struct {
	pb.UnimplementedKeycloakProvisionServer
	KcConfig
}

type KcConfig struct {
	MasterRealm   string
	AdminUsername string
	AdminPassword string
	KeycloakURI   string
}

type CreateScopeResponse struct {
	ScopeID   string
	ScopeName string
}

type CreateResourceResponse struct {
	ResourceID   string
	ResourceName string
}

type CreatePolicyResponse struct {
	PolicyID   string
	PolicyName string
}

type CreatePermissionResponse struct {
	PermissionID   string
	PermissionName string
}

type KcProvisionOpts struct {
	// configuration goes here
	KcConfig
}

func NewKcProvision(opt *KcProvisionOpts) *KcProvision {
	return &KcProvision{
		KcConfig: opt.KcConfig,
	}
}
