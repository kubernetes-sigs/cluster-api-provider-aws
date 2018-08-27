# AWS History and Lessons Learned

In an attempt to gather our knowledge of deploying and managing Kubernetes on AWS, we ask that if you have a lesson to share - to please share it

Please add anything that has been notably challenging, difficult, or surprising that you have experienced while deploying or managing Kubernetes on AWS

* [Nova] The rate limits on the AWS API need to be taken into consideration at the software level
* [Nova] Not every resource in AWS can be “tagged”
* [Dolezal] Working with IAM roles to secure access has been difficult to implement on a deployment or pod level basis.
* [Ashish] Using ASGs for worker nodes has worked well with the cluster-autoscaler simplifies machine-controller implementation in scenarios where machines (EC2 instances) fail status checks. However these are some of the things to consider:
  * ASG per zone (e.g. us-west-2a) to tolerate zone outages
  * Storage class per zone with using EBS volumes for persistent disks
  * Cluster autoscaler by default is maks ASG api calls too aggressively and can get throttled error.
* [Ashish] Security in clusters:
  * Credentials to the cluster:
    * Certs issued by cluster cert authority- not scalable but may get us out of the door sooner
    * JWT credentials using OIDC and group based RBAC (group being a field in the credentials JWT) for namespace admins
  * IAM roles as pod identity, using Kube2Iam, to control access to AWS resources. Works in multi-account environments. Kube2Iam doesn’t work very well under load and we’ve seen aws clients fail with “unable to load credentials” and had to mitigate by increasing timeouts and retry count- works most of the time.
    * [Matt Reynolds] Try [Kiam](https://github.com/uswitch/kiam) it's resolved load issues we had with kube2iam
* [Ashish] Single account environments don’t work very well as users may run into resource limits. (not totally specific to Kubernetes though, but something to keep in mind). Similarly, account peering doesn’t scale.
* [Ashish] SSL termination at the loadbalancer in multi-account environments hasn’t worked well.
* [Ashish] Ability to deploy support applications that make the cluster usable- this would include things like: monitoring and alerting stack, kube2iam (if we choose it), credential getting service. May be useful in solving this upstream than specific to AWS
* [Naadir] Subnet allocation:
  * Use at least /19s for VPCs. at least /22s for subnets.
  * For multi-region preparedness, have a larger supernetwork from which the VPCs are carved out of.
  * Use RFC3531 to carve smaller CIDR ranges from bigger ones (e.g. subnets from VPCs) so that they can be resized/shared with DCs later.
    * We used netaddr-rb 1.x (the newer one & golang port doesn’t have it)
* [Naadir] Bucket / DNS names:
  * For multi-region preparedness, copy AWS endpoint naming schemes, e.g.:
    * `<service>.<region>.<account identifier>.<domain_name>`
    * `k8s-apiserver.eu-west-1.123456789012.lucernepublishing.com`
    * `k8s-apiserver.eu-west-1.prod.adatum.com`
