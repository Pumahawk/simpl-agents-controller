package main

var prIdsDemux = projectIdsDemux{
	"microbe": {
		"authentication-provider",
		"identity-provider",
		"onboarding",
		"security-attributes-provider",
		"tier1-gateway",
		"tier2-gateway",
		"users-roles",
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
	"common": 796,
	"cm":     796,
	"com":    796,

	// Backend Authentication provider
	"authentication-provider": 939,
	"auth":                    939,

	// Plugin Eidas keycloak
	"eidas-keycloak": 1313,
	"eidas-k":        1313,

	// Eidas Node
	"eidas-node": 1312,
	"eidas-n":    1312,

	// Frontend Authentication provider
	"fe-authentication-provider": 1308,
	"fe-auth":                    1308,

	// Frontend Identity provider
	"fe-identity-provider": 1311,
	"fe-ide":               1311,

	// Frontend Onboarding
	"fe-onboarding": 1307,
	"fe-onb":        1307,

	// Frontend Security attribute provider
	"fe-security-attribute-provider": 1309,
	"fe-sap":                         1309,

	// Frontend Users and roles
	"fe-users-and-roles": 1310,
	"fe-usr":             1310,

	// Backend Identity provider
	"identity-provider": 913,
	"ide":               913,

	// Plugin keycloak authenticator
	"keycloak-authenticator": 915,
	"k-auth":                 915,

	// Backend Onboarding
	"onboarding": 770,
	"onb":        770,

	// Backend Security attribute provider
	"security-attributes-provider": 861,
	"sap":                          861,

	// Backend Lib Http client
	"simpl-http-client": 859,
	"http":              859,

	// Plugin Tier1 authenticator
	"tier1-authentication": 1457,
	"t1-auth":              1457,

	// Backend Tier1 gateway
	"tier1-gateway": 772,
	"t1g":           772,

	// Backend Tier2 gateway
	"tier2-gateway": 860,
	"t2g":           860,

	// Backend Tier2 proxy
	"tier2-proxy": 1112,
	"t2x":         1112,

	// Backend Users roles
	"users-roles": 771,
	"usr":         771,

	// Chart authority
	"ch-authority": 1402,
	"ch-auth":      1402,

	// Chart consumer
	"ch-consumer": 1404,
	"ch-con":      1404,

	// Chart consumer
	"ch-provider": 1403,
	"ch-pro":      1403,
}

type projectNameSvT map[string]int

func (p projectNameSvT) Get(key string) int {
	return p[key]
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
