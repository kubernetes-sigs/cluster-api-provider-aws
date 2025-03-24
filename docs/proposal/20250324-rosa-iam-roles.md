# Propose ROSA HCP IAM Account Roles, IAM Operator Roles, and OIDC Configuration as CRD in CAPA

## Table of Contents
- [Summary](#summary)
- [Motivation](#motivation)
  - [Problem Statement](#problem-statement)
  - [Goals](#goals)
  - [Non-Goals](#non-goals)
- [Proposal](#proposal)
  - [User Stories](#user-stories)
  - [Implementation Details](#implementation-details)
  - [Risks and Mitigations](#risks-and-mitigations)
- [Alternatives](#alternatives)
- [Implementation Plan](#implementation-plan)
- [Resources](#resources)

## Summary
The ROSA HCP in CAPA (Cluster API Provider AWS) project currently lacks built-in support for managing IAM account roles, IAM operator roles, and OIDC configuration as Custom Resource Definitions (CRDs). This proposal suggests integrating the [ROSA library](https://github.com/openshift/rosa) to extend ROSA-HCP's functionality, enabling declarative IAM and OIDC management.

## Motivation
### Problem Statement
Managing ROSA HCP IAM roles and OIDC configurations manually can be complex and error-prone. ROSA-HCP in CAPA lacks an automated mechanism to provision and reconcile these IAM resources, making it challenging for OpenShift ROSA HCP clusters deployed via CAPA to comply with security best practices and AWS integration requirements.

### Goals
- Introduce ROSARoleConfig CRD for ROSA HCP IAM account roles, IAM operator roles, and OIDC configuration in CAPA.
- Enable automated reconciliation of these resources using the ROSA library.
- Simplify cluster lifecycle management by aligning IAM/OIDC management with Kubernetes declarative paradigms.
- Reduce operational overhead for users deploying OpenShift ROSA HCP clusters on AWS via CAPA.

### Non-Goals
- Intorducing new IAM management workflows out of CAPA. 
- Introducing non-AWS IAM solutions.

## Proposal
### User Stories
#### 1. Declarative ROSA-HCP IAM Role Management
As a ROSA-HCP CAPA user, I want to define IAM account roles, operator roles and OIDC configurations using Kubernetes CRDs so that they are automatically created, seamlessly integrate with identity providers and reconciled by CAPA.

### Implementation Details
1. **Define CRDs for ROSA-HCP IAM roles and OIDC resources**
   ROSARoleConfig is a cluster scope CRD defines the required IAM account roles, operator roles and OIDC configurations to create ROSA HCP cluster. Creating the Roles and OIDC config should be in order as follows;
   1. Create ROSA HCP account roles.
   1. Create ROSA HCP OIDC config & OIDC provider.
   1. Create ROSA HCP operator roles.

   It’s important to follow the specified order, as the account roles are required to create the OIDC configuration, and the OIDC config ID is required for setting up the operator roles. The oidcConfig contains the externalAuthProviders configurations, which can be utilized by the RosaControlPlane. The ROSA lib will serve as the reference for creating the IAM roles and OIDC provider. Below is an example of a ROSARoleConfig CR with detailed descriptions.

    ```yaml
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta2
    kind: ROSARoleConfig
    metadata:
      name: rosa-role-config
    spec:
      accountRoleConfig:
        # User-defined prefix for the generated AWS IAM roles.
        # required
        prefix: capa

        # The ARN of the policy that is used to set the permissions boundary for the account roles.
        # optional
        permissionsBoundaryARN: ""

        # The arn path for the account/operator roles as well as their policies.
        # optional
        path: ""
        
        # Version of OpenShift that will be used to setup policy tag, for example "4.18"
        # required
        version: "4.18"

        sharedVPCConfig:
          # Role ARN associated with the private hosted zone used for Hosted Control Plane cluster shared VPC, this role contains policies to be used with Route 53
          # optional
          routeRoleARN: ""
          
          # Role ARN associated with the shared VPC used for Hosted Control Plane clusters, this role contains policies to be used with the VPC endpoint.
          # optional
          vpcEndpointRoleArn: ""
      operatorRoleConfig:
        # User-defined prefix for generated AWS operator policies.
        # required
        prefix: capa
        
        # The ARN of the policy that is used to set the permissions boundary for the operator roles.
        # optional
        permissionsBoundaryARN: ""
        
        # oidcConfigId is registered OIDC configuration ID to add its issuer URL as the trusted relationship to the operator roles.
        # Cannot be used while oidcConfig->createManagedOIDC is enabled
        # optional
        oidcConfigId: ""

        sharedVPCConfig:
          # AWS IAM Role Arn with policy attached, associated with shared VPC. Grants permission necessary to handle route53 operations associated with a cross-account VPC.
          # optional
          routeRoleARN: ""

          # AWS IAM Role ARN with policy attached, associated with the shared VPC. Grants permissions necessary to communicate with and handle a Hosted Control Plane cross-account VPC.
          # optional
          vpcEndpointRoleArn: ""
      oidcConfig:
        # CreateManagedOIDC flag, default is enabled to create a managed oidc-config and oidc-provider. The created oidc-config-id wil be used to create the operator roles.
        createManagedOIDC: enabled
        # example of externalAuthProviders similar to RosaControlPlane->externalAuthProviders.
        externalAuthProviders:
        - name: my-oidc-provider
          issuer:
            issuerURL: https://login.myoidc.com/<tenant-id>/v2.0
            audiences:  # audiences that will be trusted by the kube-apiserver
            - "audience1"
          claimMappings:
            username:
              claim: email
              prefixPolicy: ""
            groups:
              claim: groups
          oidcClients:
            - componentName: console
              componentNamespace: openshift-console
              clientID: <console-client-id>
              clientSecret:
                name: console-client-secret
    status:
      conditions:
        - type: Ready
          status: "True"
          message: "ROSA HCP Roles & OIDC created"
          reason: "Created"
      # Created oidc config and provider.      
      oidcID: "123456789123456789abcdefgh"
      oidcProviderARN: arn:aws:iam::12345678910:oidc-provider/oidc.ocp.dev.org/123456789123456789abcdefgh

      # Created ROSA-HCP managed account roles
      accountRolesRef:
        installerRoleARN: "arn:aws:iam::12345678910:role/capa-HCP-ROSA-Installer-Role"
        supportRoleARN: "arn:aws:iam::12345678910:role/capa-HCP-ROSA-Support-Role"
        workerRoleARN: "arn:aws:iam::12345678910:role/capa-HCP-ROSA-Worker-Role"      

      # Created ROSA-HCP managed operator roles
      operatorRolesRef:
        ingressARN: "arn:aws:iam::12345678910:role/capa-openshift-ingress-operator-cloud-credentials"
        imageRegistryARN: "arn:aws:iam::12345678910:role/capa-openshift-image-registry-installer-cloud-credentials"
        storageARN: "arn:aws:iam::12345678910:role/capa-openshift-cluster-csi-drivers-ebs-cloud-credentials"
        networkARN: "arn:aws:iam::12345678910:role/capa-openshift-cloud-network-config-controller-cloud-credent"
        kubeCloudControllerARN: "arn:aws:iam::12345678910:role/capa-kube-system-kube-controller-manager"
        nodePoolManagementARN: "arn:aws:iam::12345678910:role/capa-kube-system-capa-controller-manager"
        controlPlaneOperatorARN: "arn:aws:iam::12345678910:role/capa-kube-system-control-plane-operator"
        kmsProviderARN: "arn:aws:iam::12345678910:role/capa-kube-system-kms-provider"      
    ```

    The ROSARoleConfig status conditions should present the roles creation order and failed status. Below are examples of condition state.

    ```yaml
    status:
      conditions:
        - type: Ready
          status: "False"
          message: "Creating OIDC provider"
          reason: "Creating"
    ```

    ```yaml
    status:
      conditions:
        - type: Ready
          status: "True"
          message: "ROSA HCP Roles & OIDC created"
          reason: "Created"
    ```

    ```yaml
    status:
      conditions:
        - type: Ready
          status: "False"
          message: "Failed to create IAM operator roles oidc-config not exist."
          reason: "Failed"
    ```

2. **Implement a CAPA controller leveraging the ROSA library**
   - Manage CRUD operations for ROSA HCP IAM roles and OIDC configurations.
   - Ensure proper IAM permissions and security best practices.

3. **Validate and test integration**
   - Develop unit tests for ROSA HCP IAM and OIDC CRD reconciliation.
   - Conduct integration tests in AWS environments.
   - Introduce CI job that cover the ROSARoleConfig creation and validation.

4. **RosaControlPlane reference**
   - The `ROSAControlPlane` custom resource (CR) will include a reference to the `ROSARoleConfig`, allowing the end user to reference the ROSA IAM roles and OIDC configuration across multiple ROSA-HCP clusters. Ex;
       ```yaml
       kind: ROSAControlPlane
         metadata:
           name: rosa-hcp-stage
           namespace: default
       spec:
         rosaRoleConfigRef:
           name: stage-rosaRoleConfig
       ```

    - The `ROSAControlPlane` should be able to use the IAM role ARN and OIDC config ID from the `ROSARoleConfig` CR status, and it will validate the status of the `ROSARoleConfig` through the conditions outlined above.

    - A validation should be implemented in the ROSAControlPlane to prevent the end user from setting both `rosaRoleConfigRef` and the `rolesRef` and `oidcID` in the ROSAControlPlane at the same time. 

5. **Deleting ROSARoleConfig**
   - During the deletion of a ROSARoleConfig, a validation should be performed to check if any `ROSAControlPlane` CR references the deleted `ROSARoleConfig` CR.
   
### Risks and Mitigations
- **Complexity of IAM management**  Address through clear documentation and best practices.
- **Dependency on ROSA library**  Maintain compatibility by tracking upstream changes.

## Alternatives
- Using Terraform or manual IAM configurations (not scalable or Kubernetes-native).
- Developing a custom IAM management solution for CAPA (reinventing existing functionality in ROSA HCP).

## Implementation Plan
1. Research ROSA library capabilities.
2. Develop and integrate the new controller for ROSARoleConfig CRD.
3. Validate with unit and integration tests.
4. Document usage and best practices.
5. Engage the CAPA and OpenShift ROSA HCP community for feedback.

## Resources
- [ROSA Library](https://github.com/openshift/rosa)
- [Cluster API Provider AWS (CAPA)](https://github.com/kubernetes-sigs/cluster-api-provider-aws)
- [ROSA HCP External Auth Providers](https://cluster-api-aws.sigs.k8s.io/topics/rosa/external-auth)
- [AWS IAM Best Practices](https://docs.aws.amazon.com/IAM/latest/UserGuide/best-practices.html)