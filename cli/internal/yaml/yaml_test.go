package yaml_test

import (
	"testing"

	"github.com/Pumahawk/simpl-agents-controller/internal/yaml"
)

func TestUpdateAttributeSimple(t *testing.T) {
	expected := `
obj1: test
obj2: "test3"
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
    targetRevision: "2.7.1"
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
  targetRevision: "2.14.0-rc.1136"
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
