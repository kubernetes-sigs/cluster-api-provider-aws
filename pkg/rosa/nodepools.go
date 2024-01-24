package rosa

import cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"

// CreateNodePool adds a new node pool to the cluster.
func (c *RosaClient) CreateNodePool(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		Add().Body(nodePool).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

// GetNodePools retrieves the list of node pools in the cluster.
func (c *RosaClient) GetNodePools(clusterID string) ([]*cmv1.NodePool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		List().Page(1).Size(-1).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Items().Slice(), nil
}

// GetNodePool retrieves the details of the specified node pool.
func (c *RosaClient) GetNodePool(clusterID string, nodePoolID string) (*cmv1.NodePool, bool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().
		NodePool(nodePoolID).
		Get().
		Send()
	if response.Status() == 404 {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, handleErr(response.Error(), err)
	}
	return response.Body(), true, nil
}

// UpdateNodePool updates the specified node pool.
func (c *RosaClient) UpdateNodePool(clusterID string, nodePool *cmv1.NodePool) (*cmv1.NodePool, error) {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().NodePool(nodePool.ID()).
		Update().Body(nodePool).
		Send()
	if err != nil {
		return nil, handleErr(response.Error(), err)
	}
	return response.Body(), nil
}

// DeleteNodePool deletes the specified node pool.
func (c *RosaClient) DeleteNodePool(clusterID string, nodePoolID string) error {
	response, err := c.ocm.ClustersMgmt().V1().
		Clusters().Cluster(clusterID).
		NodePools().NodePool(nodePoolID).
		Delete().
		Send()
	if err != nil {
		return handleErr(response.Error(), err)
	}
	return nil
}
