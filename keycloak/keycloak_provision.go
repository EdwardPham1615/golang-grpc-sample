package keycloak

import (
	"context"
	"crypto/tls"
	pb "golang-grpc-sample/proto"
	"github.com/Nerzal/gocloak/v8"
)

var basicRoles = []string{
	"admin",
	"normal",
}

func (c *KcConfig) newKeycloakClient() gocloak.GoCloak {
	client := gocloak.NewClient(c.KeycloakURI)
	restyClient := client.RestyClient()
	restyClient.SetDebug(true)
	restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	return client
}

func (c *KcConfig) newKeycloakToken(ctx context.Context, client gocloak.GoCloak) (*gocloak.JWT, error) {
	return client.LoginAdmin(ctx, c.AdminUsername, c.AdminPassword, c.MasterRealm)
}

func CreateRealm(ctx context.Context, token gocloak.JWT, kcClient gocloak.GoCloak, tenantID string) (string, error) {
	realm := gocloak.RealmRepresentation{
		Realm:   gocloak.StringP(tenantID),
		Enabled: gocloak.BoolP(true),
	}

	var realmInfo string
	realmInfo, err := kcClient.CreateRealm(ctx, token.AccessToken, realm)
	if err != nil {
		return "", err
	}
	return realmInfo, nil
}

func CreateRealmRoles(ctx context.Context, token gocloak.JWT, kcClient gocloak.GoCloak, tenantID string) error {
	for _, roleName := range basicRoles {
		role := gocloak.Role{
			Name: gocloak.StringP(roleName),
		}

		_, err := kcClient.CreateRealmRole(ctx, token.AccessToken, tenantID, role)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateClient(ctx context.Context, token gocloak.JWT, kcClient gocloak.GoCloak, tenantID string) (string, error) {
	apiClient := gocloak.Client{
		ClientID:                     gocloak.StringP("api"),
		DirectAccessGrantsEnabled:    gocloak.BoolP(true),
		Enabled:                      gocloak.BoolP(true),
		RedirectURIs:                 &[]string{"*"},
		StandardFlowEnabled:          gocloak.BoolP(true),
		AuthorizationServicesEnabled: gocloak.BoolP(true),
		ServiceAccountsEnabled:       gocloak.BoolP(true),
	}

	frontendClient := gocloak.Client{
		PublicClient:              gocloak.BoolP(true),
		ClientID:                  gocloak.StringP("frontend"),
		DirectAccessGrantsEnabled: gocloak.BoolP(true),
		Enabled:                   gocloak.BoolP(true),
		RedirectURIs:              &[]string{"*"},
		StandardFlowEnabled:       gocloak.BoolP(true),
	}

	var apiClientID string
	apiClientID, err := kcClient.CreateClient(ctx, token.AccessToken, tenantID, apiClient)
	if err != nil {
		return "", err
	}

	_, err = kcClient.CreateClient(ctx, token.AccessToken, tenantID, frontendClient)
	if err != nil {
		return "", err
	}

	return apiClientID, nil
}

func CreateScope(ctx context.Context, token gocloak.JWT, kcClient gocloak.GoCloak, tenantID string, clientID string) (*CreateScopeResponse, error) {
	kcScope := gocloak.ScopeRepresentation{
		DisplayName: gocloak.StringP("all"),
		Name:        gocloak.StringP("all"),
	}

	var scopeInfo *gocloak.ScopeRepresentation
	scopeInfo, err := kcClient.CreateScope(ctx, token.AccessToken, tenantID, clientID, kcScope)
	if err != nil {
		return nil, err
	}

	return &CreateScopeResponse{
		ScopeID:   *scopeInfo.ID,
		ScopeName: *scopeInfo.Name,
	}, nil
}

func CreateResource(ctx context.Context, token gocloak.JWT, kcClient gocloak.GoCloak, tenantID string, clientID string, scopeID string, scopeName string) (*CreateResourceResponse, error) {
	kcResource := gocloak.ResourceRepresentation{
		DisplayName: gocloak.StringP("api"),
		Name:        gocloak.StringP("api"),
		Owner: &gocloak.ResourceOwnerRepresentation{
			ID:   gocloak.StringP(clientID),
			Name: gocloak.StringP("api"),
		},
		Scopes: &[]gocloak.ScopeRepresentation{
			{
				ID:   gocloak.StringP(scopeID),
				Name: gocloak.StringP(scopeName),
			},
		},
	}

	var resourceInfo *gocloak.ResourceRepresentation
	resourceInfo, err := kcClient.CreateResource(ctx, token.AccessToken, tenantID, clientID, kcResource)
	if err != nil {
		return nil, err
	}

	return &CreateResourceResponse{
		ResourceID:   *resourceInfo.ID,
		ResourceName: *resourceInfo.Name,
	}, nil
}

func CreatePolicy(ctx context.Context, token gocloak.JWT, kcClient gocloak.GoCloak, tenantID string, clientID string) (*CreatePolicyResponse, error) {
	kcPolicy := gocloak.PolicyRepresentation{
		DecisionStrategy: gocloak.UNANIMOUS,
		Logic:            gocloak.POSITIVE,
		Name:             gocloak.StringP("Is normal"),
		Type:             gocloak.StringP("role"),
		RolePolicyRepresentation: gocloak.RolePolicyRepresentation{
			Roles: &[]gocloak.RoleDefinition{
				{
					ID:       gocloak.StringP("normal"),
					Required: gocloak.BoolP(true),
				},
			},
		},
	}

	var policyInfo *gocloak.PolicyRepresentation
	policyInfo, err := kcClient.CreatePolicy(ctx, token.AccessToken, tenantID, clientID, kcPolicy)
	if err != nil {
		return nil, err
	}

	return &CreatePolicyResponse{
		PolicyID:   *policyInfo.ID,
		PolicyName: *policyInfo.Name,
	}, nil
}

func CreatePermission(ctx context.Context, token gocloak.JWT, kcClient gocloak.GoCloak, tenantID string, clientID string) (*CreatePermissionResponse, error) {
	kcPermisison := gocloak.PermissionRepresentation{
		DecisionStrategy: gocloak.AFFIRMATIVE,
		Logic:            gocloak.POSITIVE,
		Name:             gocloak.StringP("Access all API"),
		Policies:         &[]string{"Is normal"},
		Resources:        &[]string{"api"},
		Scopes:           &[]string{"all"},
		Type:             gocloak.StringP("scope"),
	}

	var permissionInfo *gocloak.PermissionRepresentation
	permissionInfo, err := kcClient.CreatePermission(ctx, token.AccessToken, tenantID, clientID, kcPermisison)
	if err != nil {
		return nil, err
	}

	return &CreatePermissionResponse{
		PermissionID:   *permissionInfo.ID,
		PermissionName: *permissionInfo.Name,
	}, nil
}

func CreateUser(ctx context.Context, token gocloak.JWT, kcClient gocloak.GoCloak, tenantID string, userName string, passWord string) (string, error) {
	user := gocloak.User{
		Username: gocloak.StringP(userName),
		Enabled:  gocloak.BoolP(true),
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Temporary: gocloak.BoolP(false),
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP(passWord),
			},
		},
		RealmRoles: &basicRoles,
	}

	var userID string
	userID, err := kcClient.CreateUser(ctx, token.AccessToken, tenantID, user)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func AssignRoles(ctx context.Context, token gocloak.JWT, kcClient gocloak.GoCloak, tenantID string, userID string) error {

	var realmRoles []*gocloak.Role
	realmRoles, err := kcClient.GetRealmRoles(ctx, token.AccessToken, tenantID)
	if err != nil {
		return err
	}
	roles := make([]gocloak.Role, 0)
	for _, role := range realmRoles {
		roles = append(roles, *role)
	}

	err = kcClient.AddRealmRoleToUser(ctx, token.AccessToken, tenantID, userID, roles)
	if err != nil {
		return err
	}

	return nil
}

func (s *KcProvision) InitializeKeyCloak(ctx context.Context, request *pb.InitializeKeycloakRealmRequest) (*pb.InitializeKeycloakRealmResponse, error) {
	keyCloakConfig := s.KcConfig
	keyCloakClient := keyCloakConfig.newKeycloakClient()
	keyCloakToken, err := keyCloakConfig.newKeycloakToken(ctx, keyCloakClient)
	if err != nil {
		return nil, err
	}

	_, err = CreateRealm(ctx, *keyCloakToken, keyCloakClient, request.TenantId)
	if err != nil {
		return nil, err
	}

	err = CreateRealmRoles(ctx, *keyCloakToken, keyCloakClient, request.TenantId)
	if err != nil {
		return nil, err
	}

	apiClientID, err := CreateClient(ctx, *keyCloakToken, keyCloakClient, request.TenantId)
	if err != nil {
		return nil, err
	}

	var scopeInfo *CreateScopeResponse
	scopeInfo, err = CreateScope(ctx, *keyCloakToken, keyCloakClient, request.TenantId, apiClientID)
	if err != nil {
		return nil, err
	}

	//var resourceInfo *CreateResourceResponse
	_, err = CreateResource(ctx, *keyCloakToken, keyCloakClient, request.TenantId, apiClientID, scopeInfo.ScopeID, scopeInfo.ScopeName)
	if err != nil {
		return nil, err
	}

	//var policyInfo *CreatePolicyResponse
	_, err = CreatePolicy(ctx, *keyCloakToken, keyCloakClient, request.TenantId, apiClientID)
	if err != nil {
		return nil, err
	}

	_, err = CreatePermission(ctx, *keyCloakToken, keyCloakClient, request.TenantId, apiClientID)
	if err != nil {
		return nil, err
	}

	return &pb.InitializeKeycloakRealmResponse{
		Message: "Finish initialize new tenant.",
	}, nil
}

func (s *KcProvision) InitializeKeyCloakUser(ctx context.Context, request *pb.InitializeKeycloakUserRequest) (*pb.InitializeKeycloakUserResponse, error) {
	keyCloakConfig := s.KcConfig
	keyCloakClient := keyCloakConfig.newKeycloakClient()
	keyCloakToken, err := keyCloakConfig.newKeycloakToken(ctx, keyCloakClient)
	if err != nil {
		return nil, err
	}

	userID, err := CreateUser(ctx, *keyCloakToken, keyCloakClient, request.TenantId, request.UserName, request.PassWord)
	if err != nil {
		return nil, err
	}

	err = AssignRoles(ctx, *keyCloakToken, keyCloakClient, request.TenantId, userID)
	if err != nil {
		return nil, err
	}

	return &pb.InitializeKeycloakUserResponse{
		Message: "Finish setup new user for tenant.",
	}, nil

}
