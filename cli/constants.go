package main

var prAuthenticationProvider = PrInfo{939, "helm", "authentication-provider"}
var prChAuthority = PrInfo{1402, "helm", "ch-authority"}
var prChConsumer = PrInfo{1404, "helm", "ch-consumer"}
var prChProvider = PrInfo{1403, "helm", "ch-provider"}
var prCommon = PrInfo{796, "maven", "common"}
var prConsentManagementBe = PrInfo{1723, "helm", "consent-management-be"}
var prEidasKeycloak = PrInfo{1313, "maven", "eidas-keycloak"}
var prEidasNode = PrInfo{1312, "helm", "eidas-node"}
var prFeAuthenticationProvider = PrInfo{1308, "helm", "fe-authentication-provider"}
var prFeIdentityProvider = PrInfo{1311, "helm", "fe-identity-provider"}
var prFeOnboarding = PrInfo{1307, "helm", "fe-onboarding"}
var prFeSecurityAttributeProvider = PrInfo{1309, "helm", "fe-security-attribute-provider"}
var prFeUsersAndRoles = PrInfo{1310, "helm", "fe-users-and-roles"}
var prIdentityProvider = PrInfo{913, "helm", "identity-provider"}
var prKeycloakAuthenticator = PrInfo{915, "maven", "keycloak-authenticator"}
var prOnboarding = PrInfo{770, "helm", "onboarding"}
var prSecurityAttributesProvider = PrInfo{861, "helm", "security-attributes-provider"}
var prSimplHttpClient = PrInfo{859, "helm", "simpl-http-client"}
var prTier1Authentication = PrInfo{1457, "helm", "tier1-authentication"}
var prTier1Gateway = PrInfo{772, "helm", "tier1-gateway"}
var prTier2Gateway = PrInfo{860, "helm", "tier2-gateway"}
var prTier2Proxy = PrInfo{1112, "helm", "tier2-proxy"}
var prUsersRoles = PrInfo{771, "helm", "users-roles"}

var prIdsDemux = projectIdsDemux{
	"microbe": {
		"authentication-provider",
		"identity-provider",
		"onboarding",
		"security-attributes-provider",
		"tier1-gateway",
		"tier2-gateway",
		"users-roles",
		"consent-management-be",
	},
	"microfe": {
		"fe-authentication-provider",
		"fe-identity-provider",
		"fe-onboarding",
		"fe-security-attribute-provider",
		"fe-users-and-roles",
	},
	"lib": {
		"common",
		"eidas-keycloak",
		"eidas-node",
		"keycloak-authenticator",
		"simpl-http-client",
	},
	"misc": {
		"tier2-proxy",
	},
	"charts": {
		"ch-authority",
		"ch-consumer",
		"ch-provider",
	},
}

var prIds = projectNameSvT{
	// Backend Common
	"common": prCommon,
	"cm":     prCommon,
	"com":    prCommon,

	// Backend Authentication provider
	"authentication-provider": prAuthenticationProvider,
	"auth":                    prAuthenticationProvider,

	// Plugin Eidas keycloak
	"eidas-keycloak": prEidasKeycloak,
	"eidas-k":        prEidasKeycloak,

	// Consent Management Be
	"consent-management-be": prConsentManagementBe,
	"cmb":                   prConsentManagementBe,

	// Eidas Node
	"eidas-node": prEidasNode,
	"eidas-n":    prEidasNode,

	// Frontend Authentication provider
	"fe-authentication-provider": prFeAuthenticationProvider,
	"fe-auth":                    prFeAuthenticationProvider,

	// Frontend Identity provider
	"fe-identity-provider": prFeIdentityProvider,
	"fe-ide":               prFeIdentityProvider,

	// Frontend Onboarding
	"fe-onboarding": prFeOnboarding,
	"fe-onb":        prFeOnboarding,

	// Frontend Security attribute provider
	"fe-security-attribute-provider": prFeSecurityAttributeProvider,
	"fe-sap":                         prFeSecurityAttributeProvider,

	// Frontend Users and roles
	"fe-users-and-roles": prFeUsersAndRoles,
	"fe-usr":             prFeUsersAndRoles,

	// Backend Identity provider
	"identity-provider": prIdentityProvider,
	"ide":               prIdentityProvider,

	// Plugin keycloak authenticator
	"keycloak-authenticator": prKeycloakAuthenticator,
	"k-auth":                 prKeycloakAuthenticator,

	// Backend Onboarding
	"onboarding": prOnboarding,
	"onb":        prOnboarding,

	// Backend Security attribute provider
	"security-attributes-provider": prSecurityAttributesProvider,
	"sap":                          prSecurityAttributesProvider,

	// Backend Lib Http client
	"simpl-http-client": prSimplHttpClient,
	"http":              prSimplHttpClient,

	// Plugin Tier1 authenticator
	"tier1-authentication": prTier1Authentication,
	"t1-auth":              prTier1Authentication,

	// Backend Tier1 gateway
	"tier1-gateway": prTier1Gateway,
	"t1g":           prTier1Gateway,

	// Backend Tier2 gateway
	"tier2-gateway": prTier2Gateway,
	"t2g":           prTier2Gateway,

	// Backend Tier2 proxy
	"tier2-proxy": prTier2Proxy,
	"t2x":         prTier2Proxy,

	// Backend Users roles
	"users-roles": prUsersRoles,
	"usr":         prUsersRoles,

	// Chart authority
	"ch-authority": prChAuthority,
	"ch-auth":      prChAuthority,

	// Chart consumer
	"ch-consumer": prChConsumer,
	"ch-con":      prChConsumer,

	// Chart consumer
	"ch-provider": prChProvider,
	"ch-pro":      prChProvider,
}

type projectNameSvT map[string]PrInfo

type PrInfo struct {
	Id   int
	Type string
	Name string
}

func (p projectNameSvT) Get(key string) (PrInfo, bool) {
	pr, ok := p[key]
	return pr, ok
}

type projectIdsDemux map[string][]string

func (p projectIdsDemux) Demux(v []string) []string {
	strs := make([]string, 0, len(v))
	for _, s := range v {
		if prs, ok := p[s]; ok {
			strs = append(strs, prs...)
		} else {
			strs = append(strs, s)
		}
	}
	return strs
}
