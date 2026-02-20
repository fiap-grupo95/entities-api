# New Relic - Deploy via Manifestos Kubernetes

Este diretório contém os manifestos Kubernetes para deploy do New Relic Infrastructure e coleta de logs.

## Arquivos

- `newrelic-namespace.yml` - Namespace para recursos do New Relic
- `newrelic-secret.yml` - Secret com a License Key
- `newrelic-configmap.yml` - Configurações gerais do cluster
- `newrelic-rbac.yml` - Service Account, ClusterRole e ClusterRoleBinding
- `newrelic-daemonset-infra.yml` - Infrastructure Agent (métricas dos nodes)
- `newrelic-fluent-bit-configmap.yml` - Configuração do Fluent Bit
- `newrelic-daemonset-fluent-bit.yml` - Fluent Bit para coleta de logs
- `newrelic-kube-state-metrics.yml` - Kube State Metrics
- `newrelic-events-configmap.yml` - Configuração do Events Collector
- `newrelic-deployment-events.yml` - Events Collector

## Pré-requisitos

1. Cluster Kubernetes rodando
2. `kubectl` configurado e conectado ao cluster
3. License Key do New Relic ([obtenha aqui](https://one.newrelic.com/admin-portal/api-keys/home))

## Instalação

### Passo 1: Configurar License Key

Edite o arquivo `newrelic-secret.yml` e substitua a `license` pela sua License Key:

```yaml
stringData:
  license: "sua_license_key_aqui"
```

### Passo 2: Aplicar os Manifestos

```bash
kubectl apply -f newrelic-namespace.yml
kubectl apply -f newrelic-secret.yml
kubectl apply -f newrelic-configmap.yml
kubectl apply -f newrelic-fluent-bit-configmap.yml
kubectl apply -f newrelic-events-configmap.yml
kubectl apply -f newrelic-rbac.yml
kubectl apply -f newrelic-daemonset-infra.yml
kubectl apply -f newrelic-daemonset-fluent-bit.yml
kubectl apply -f newrelic-kube-state-metrics.yml
kubectl apply -f newrelic-deployment-events.yml
```

## Verificação

```bash
kubectl get pods -n newrelic
```

## Logs

```bash
kubectl logs -n newrelic -l app=fluent-bit --tail=50
kubectl logs -n newrelic -l app=newrelic-infra --tail=50
kubectl logs -n newrelic -l app=newrelic-kube-events --tail=50
```

## Acesso no New Relic

- Logs: https://one.newrelic.com/logs
- Kubernetes: https://one.newrelic.com/kubernetes

Use filtros como:
- `cluster.name = "mecanica-xpto-cluster"`
- `k8s.namespaceName = "entities-api"`
