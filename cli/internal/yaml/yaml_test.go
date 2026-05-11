package yaml_test

import (
	"testing"

	"github.com/Pumahawk/simpl-agents-controller/internal/yaml"
)

func TestUpdateAttributeSimple(t *testing.T) {
	expected := `
obj1: test
obj2: test3
`
	content := `
obj1: test
obj2: test2
`
	obj := yaml.NewObj([]byte(content))
	obj.UpdateAttribute("test3", "obj2")

	result := string(obj.Bytes())
	if result != expected {
		t.Errorf("result != expected, result=%q, expected=%q", result, expected)
	}

}

func TestUpdateAttributeComplex(t *testing.T) {
	expected := `
keycloak:
  repo_URL: 'https://charts.bitnami.com/bitnami'
  chart_name: keycloak
  targetRevision: 25.2.0
  keycloak_authenticator:
    targetRevision: 2.7.1
  eidas_demo_keycloak_extension:
    targetRevision: 0.5.0-rc.190
  nameOverride: keycloak
  resourcesPreset: "none"
  resources:
    requests:
      cpu: "500m"
      memory: "1G"
    limits:
      cpu: "1"
      memory: "2G"
  valueFiles:
    - values.yaml
ejbca:
  repo_URL: 'https://keyfactor.github.io/ejbca-community-helm'
  targetRevision: 1.0.7
  chart_name: ejbca-community-helm
  fullnameOverride: ejbca-community-helm
  resources:
    requests:
      cpu: "500m"
      memory: "1024Mi"
    limits:
      cpu: "500m"
      memory: "1024Mi"
  valueFiles:
    - values.yaml
tier1_gateway:
  projectID: "772"
  targetRevision: 2.14.0-rc.1136
  chart_name: tier1-gateway
  resources:
    requests:
      cpu: "400m"
      memory: "1Gi"
    limits:
      cpu: "400m"
      memory: "1Gi"
  valueFiles:
    - values.yaml
sap:
  projectID: "861"
  targetRevision: 2.13.0-rc.1336
  chart_name: security-attributes-provider
  resources:
    requests:
      cpu: "400m"
      memory: "1Gi"
    limits:
      cpu: "400m"
      memory: "1Gi"
  valueFiles:
    - values.yaml
users_roles:
  enabled: true
  projectID: "771"
  targetRevision: 2.12.5-new
  chart_name: users-roles
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
  valueFiles:
    - values.yaml
fe_auth_provider:
  enabled: true
  projectID: "1308"
  targetRevision: 2.9.2
  chart_name: fe-authentication-provider
  resources:
    requests:
      cpu: "250m"
      memory: "512Mi"
    limits:
      cpu: "250m"
      memory: "512Mi"
fe_users_roles:
  enabled: true
  projectID: "1310"
  targetRevision: 2.12.4
  chart_name: fe-users-and-roles
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
fe_identity_provider:
  enabled: true
  projectID: "1311"
  targetRevision: 2.10.4
  chart_name: fe-identity-provider
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
fe_onboarding:
  enabled: true
  projectID: "1307"
  targetRevision: 2.12.1
  chart_name: fe-onboarding
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
fe_security_attribute_provider:
  enabled: true
  projectID: "1309"
  targetRevision: 2.9.1
  chart_name: fe-security-attribute-provider
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
onboarding:
  enabled: true
  projectID: "770"
  targetRevision: 2.12.6
  chart_name: onboarding
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "1"
      memory: "1Gi"
  valueFiles:
    - values.yaml
tier2_gateway:
  enabled: true
  projectID: "860"
  targetRevision: 2.11.3
  chart_name: tier2-gateway
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
  valueFiles:
    - values.yaml
    - values-routes-authority.yaml
identity_provider:
  enabled: true
  projectID: "913"
  targetRevision: 2.11.3
  chart_name: identity-provider
  fullnameOverride: identity-provider
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "1"
      memory: "1Gi"
  valueFiles:
    - values.yaml
auth_provider:
  enabled: true
  projectID: "939"
  targetRevision: 2.12.3
  chart_name: authentication-provider
  resources:
    requests:
      cpu: "400m"
      memory: "400Mi"
    limits:
      cpu: "500m"
      memory: "1Gi"
  valueFiles:
    - values.yaml
redis:
  enabled: true
  repo_URL: 'https://charts.bitnami.com/bitnami'
  targetRevision: 19.6.0
  chart_name: redis
  nameOverride: redis
  resourcesPreset: "none"
  resources:
    requests:
      cpu: "0"
      memory: "0"
  valueFiles:
    - values.yaml
redis_commander:
  enabled: false
  repo_URL: 'https://github.com/joeferner/redis-commander.git'
  targetRevision: master
  nameOverride: redis-commander
  path: 'k8s/helm-chart/redis-commander'
  resources:
    requests:
      cpu: "500m"
      memory: "512Mi"
    limits:
      cpu: "500m"
      memory: "512Mi"
  valueFiles:
    - values.yaml
tier2_proxy:
  enabled: true
  projectID: "1112"
  targetRevision: 1.5.2
  chart_name: tier2-proxy
  resources:
    requests:
      cpu: "500m"
      memory: "500Mi"
    limits:
      cpu: "500m"
      memory: "500Mi"
  valueFiles:
    - values.yaml
`

	content := `
keycloak:
  repo_URL: 'https://charts.bitnami.com/bitnami'
  chart_name: keycloak
  targetRevision: 25.2.0
  keycloak_authenticator:
    targetRevision: 2.7.0
  eidas_demo_keycloak_extension:
    targetRevision: 0.5.0-rc.190
  nameOverride: keycloak
  resourcesPreset: "none"
  resources:
    requests:
      cpu: "500m"
      memory: "1G"
    limits:
      cpu: "1"
      memory: "2G"
  valueFiles:
    - values.yaml
ejbca:
  repo_URL: 'https://keyfactor.github.io/ejbca-community-helm'
  targetRevision: 1.0.7
  chart_name: ejbca-community-helm
  fullnameOverride: ejbca-community-helm
  resources:
    requests:
      cpu: "500m"
      memory: "1024Mi"
    limits:
      cpu: "500m"
      memory: "1024Mi"
  valueFiles:
    - values.yaml
tier1_gateway:
  projectID: "772"
  targetRevision: 2.14.0-rc.1135
  chart_name: tier1-gateway
  resources:
    requests:
      cpu: "400m"
      memory: "1Gi"
    limits:
      cpu: "400m"
      memory: "1Gi"
  valueFiles:
    - values.yaml
sap:
  projectID: "861"
  targetRevision: 2.13.0-rc.1336
  chart_name: security-attributes-provider
  resources:
    requests:
      cpu: "400m"
      memory: "1Gi"
    limits:
      cpu: "400m"
      memory: "1Gi"
  valueFiles:
    - values.yaml
users_roles:
  enabled: true
  projectID: "771"
  targetRevision: 2.12.5
  chart_name: users-roles
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
  valueFiles:
    - values.yaml
fe_auth_provider:
  enabled: true
  projectID: "1308"
  targetRevision: 2.9.2
  chart_name: fe-authentication-provider
  resources:
    requests:
      cpu: "250m"
      memory: "512Mi"
    limits:
      cpu: "250m"
      memory: "512Mi"
fe_users_roles:
  enabled: true
  projectID: "1310"
  targetRevision: 2.12.4
  chart_name: fe-users-and-roles
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
fe_identity_provider:
  enabled: true
  projectID: "1311"
  targetRevision: 2.10.4
  chart_name: fe-identity-provider
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
fe_onboarding:
  enabled: true
  projectID: "1307"
  targetRevision: 2.12.1
  chart_name: fe-onboarding
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
fe_security_attribute_provider:
  enabled: true
  projectID: "1309"
  targetRevision: 2.9.1
  chart_name: fe-security-attribute-provider
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
onboarding:
  enabled: true
  projectID: "770"
  targetRevision: 2.12.6
  chart_name: onboarding
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "1"
      memory: "1Gi"
  valueFiles:
    - values.yaml
tier2_gateway:
  enabled: true
  projectID: "860"
  targetRevision: 2.11.3
  chart_name: tier2-gateway
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "400m"
      memory: "512Mi"
  valueFiles:
    - values.yaml
    - values-routes-authority.yaml
identity_provider:
  enabled: true
  projectID: "913"
  targetRevision: 2.11.3
  chart_name: identity-provider
  fullnameOverride: identity-provider
  resources:
    requests:
      cpu: "400m"
      memory: "512Mi"
    limits:
      cpu: "1"
      memory: "1Gi"
  valueFiles:
    - values.yaml
auth_provider:
  enabled: true
  projectID: "939"
  targetRevision: 2.12.3
  chart_name: authentication-provider
  resources:
    requests:
      cpu: "400m"
      memory: "400Mi"
    limits:
      cpu: "500m"
      memory: "1Gi"
  valueFiles:
    - values.yaml
redis:
  enabled: true
  repo_URL: 'https://charts.bitnami.com/bitnami'
  targetRevision: 19.6.0
  chart_name: redis
  nameOverride: redis
  resourcesPreset: "none"
  resources:
    requests:
      cpu: "0"
      memory: "0"
  valueFiles:
    - values.yaml
redis_commander:
  enabled: false
  repo_URL: 'https://github.com/joeferner/redis-commander.git'
  targetRevision: master
  nameOverride: redis-commander
  path: 'k8s/helm-chart/redis-commander'
  resources:
    requests:
      cpu: "500m"
      memory: "512Mi"
    limits:
      cpu: "500m"
      memory: "512Mi"
  valueFiles:
    - values.yaml
tier2_proxy:
  enabled: true
  projectID: "1112"
  targetRevision: 1.5.2
  chart_name: tier2-proxy
  resources:
    requests:
      cpu: "500m"
      memory: "500Mi"
    limits:
      cpu: "500m"
      memory: "500Mi"
  valueFiles:
    - values.yaml
`
	obj := yaml.NewObj([]byte(content))
	if ok, err := obj.UpdateAttribute("2.7.1", "keycloak", "keycloak_authenticator", "targetRevision"); err != nil {
		t.Errorf("error: %s", err)
		return
	} else if !ok {
		t.Errorf("not updated")
		return
	}
	if ok, err := obj.UpdateAttribute("2.14.0-rc.1136", "tier1_gateway", "targetRevision"); err != nil {
		t.Errorf("error: %s", err)
		return
	} else if !ok {
		t.Errorf("not updated")
		return
	}
	if ok, err := obj.UpdateAttribute("2.12.5-new", "users_roles", "targetRevision"); err != nil {
		t.Errorf("error: %s", err)
		return
	} else if !ok {
		t.Errorf("not updated")
		return
	}

	result := string(obj.Bytes())
	if result != expected {
		t.Errorf("result != expected,\n"+
			"------ result ------\n\n"+
			"%s\n"+
			"------ expected ------\n\n"+
			"%s\n"+
			"", result, expected)
	}

}
