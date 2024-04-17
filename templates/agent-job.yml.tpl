{{- $accountName := get . "accountName" }}
{{- $clusterToken := get . "clusterToken" }}
{{- $clusterName := get . "clusterName" }}
{{- $version := get . "version" }}

apiVersion: v1
kind: Namespace
metadata:
  name: kloudlite

---
apiVersion: v1
kind: ServiceAccount
metadata:
  annotations:
    kloudlite.io/description: Service account used by helm charts to run helm release jobs
  name: helm-job-svc-account
  namespace: kloudlite

---

apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  annotations:
    kloudlite.io/description: Cluster role binding used by helm charts to run helm
      release jobs
  name: helm-job-svc-account-rb
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
- kind: ServiceAccount
  name: helm-job-svc-account
  namespace: kloudlite

---
apiVersion: batch/v1
kind: Job
metadata:
  labels:
    kloudlite.io/chart-install-or-upgrade-job: "true"
    kloudlite.io/helm-chart.name: kloudlite-agent
  name: helm-job-kloudlite-agent
  namespace: kloudlite
spec:
  backoffLimit: 1
  completionMode: NonIndexed
  completions: 1
  parallelism: 1
  suspend: false
  template:
    metadata:
      annotations:
        kloudlite.io/job_name: helm-job-kloudlite-agent
        kloudlite.io/job_type: helm-install
      labels:
        job-name: helm-job-kloudlite-agent
    spec:
      containers:
      - command:
        - bash
        - -c
        - |+
          set -o nounset
          set -o pipefail
          set -o errexit

          helm repo add helm-repo https://kloudlite.github.io/helm-charts
          helm repo update helm-repo
          echo "running pre-install job script"

          kubectl apply -f https://github.com/kloudlite/helm-charts/releases/download/{{ $version }}/crds-all.yml --server-side


          cat > values.yml <<EOF

          accountName: {{ $accountName }}
          agent:
            enabled: true
            name: kl-agent
            nodeSelector: {}
            tolerations: []
          cloudProvider: aws
          clusterIdentitySecretName: kl-cluster-identity
          clusterInternalDNS: cluster.local
          clusterName: {{ $clusterName }}
          clusterToken: {{ $clusterToken }}
          helmCharts:
            certManager:
              enabled: false
            clusterAutoscaler:
              enabled: false
            ingressNginx:
              enabled: false
            vector:
              debugOnStdout: false
              enabled: false
          imagePullPolicy: Always
          messageOfficeGRPCAddr: message-office.kloudlite.io:443
          nodeSelector:
            node-role.kubernetes.io/master: "true"
          operators:
            agentOperator:
              configuration:
                nodepools:
                  enabled: false
                routers:
                  enabled: false
                wireguard:
                  enabled: false
              enabled: true
              name: kl-agent-operator
              nodeSelector: {}
              tolerations: []
          preferOperatorsOnMasterNodes: true
          svcAccountName: sa
          tolerations:
          - effect: NoSchedule
            key: node-role.kubernetes.io/master
            operator: Exists

          EOF

          helm upgrade --install kloudlite-agent helm-repo/kloudlite-agent --namespace kloudlite --version {{ $version }} --values values.yml 2>&1 | tee /dev/termination-log
          echo "running post-install job script"

          if kubectl get ns kloudlite-tmp;
          then
            kubectl delete ns kloudlite-tmp
          fi



        image: ghcr.io/kloudlite/kloudlite/operator/workers/helm-job-runner:{{ $version }}
        imagePullPolicy: IfNotPresent
        name: helm
        resources: {}
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Never
      schedulerName: default-scheduler
      securityContext: {}
      serviceAccount: helm-job-svc-account
      serviceAccountName: helm-job-svc-account
      terminationGracePeriodSeconds: 30
