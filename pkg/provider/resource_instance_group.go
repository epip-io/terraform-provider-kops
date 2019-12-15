package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/upup/pkg/fi/cloudup"

	"github.com/epip-io/terraform-provider-kops/pkg/api"
)

type instanceGroupID struct {
	clusterName       string
	instanceGroupName string
}

func resourceInstanceGroupCreate(d *schema.ResourceData, m interface{}) error {
	clusterName := d.Get("cluster_name").(string)
	clientset := m.(*ProviderConfig).clientset
	cluster, err := clientset.GetCluster(clusterName)
	if err != nil {
		return err
	}

	instanceGroup, err := clientset.InstanceGroupsFor(cluster).Create(&kops.InstanceGroup{
		ObjectMeta: api.MarshalObjectMeta(sectionData(d, "metadata")),
		Spec:       api.MarshalInstanceGroupSpec(sectionData(d, "spec")),
	})
	if err != nil {
		return err
	}

	channel, err := cloudup.ChannelForCluster(cluster)
	if err != nil {
		return err
	}

	fullInstanceGroup, err := cloudup.PopulateInstanceGroupSpec(cluster, instanceGroup, channel)
	if err != nil {
		return err
	}

	_, err = clientset.InstanceGroupsFor(cluster).Update(fullInstanceGroup)
	if err != nil {
		return err
	}

	d.SetId(instanceGroupID{
		clusterName:       clusterName,
		instanceGroupName: instanceGroup.ObjectMeta.Name,
	}.String())

	return resourceInstanceGroupRead(d, m)
}

func resourceInstanceGroupRead(d *schema.ResourceData, m interface{}) error {
	instanceGroup, err := getInstanceGroup(d, m)
	if err != nil {
		return err
	}
	if err := d.Set("metadata", api.UnmarshalObjectMeta(instanceGroup.ObjectMeta)); err != nil {
		return err
	}
	if err := d.Set("spec", api.UnmarshalInstanceGroupSpec(instanceGroup.Spec)); err != nil {
		return err
	}
	return nil
}

func resourceInstanceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	if ok, _ := resourceInstanceGroupExists(d, m); !ok {
		d.SetId("")
		return nil
	}

	clusterName := d.Get("cluster_name").(string)
	clientset := m.(*ProviderConfig).clientset
	cluster, err := clientset.GetCluster(clusterName)
	if err != nil {
		return err
	}

	_, err = clientset.InstanceGroupsFor(cluster).Update(&kops.InstanceGroup{
		ObjectMeta: api.MarshalObjectMeta(sectionData(d, "metadata")),
		Spec:       api.MarshalInstanceGroupSpec(sectionData(d, "spec")),
	})
	if err != nil {
		return err
	}

	return resourceInstanceGroupRead(d, m)
}

func resourceInstanceGroupDelete(d *schema.ResourceData, m interface{}) error {
	groupID := parseInstanceGroupID(d.Id())
	clientset := m.(*ProviderConfig).clientset
	cluster, err := clientset.GetCluster(groupID.clusterName)
	if err != nil {
		return err
	}
	return clientset.InstanceGroupsFor(cluster).Delete(groupID.instanceGroupName, &v1.DeleteOptions{})
}

func resourceInstanceGroupExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, err := getInstanceGroup(d, m)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (i instanceGroupID) String() string {
	return fmt.Sprintf("%s/%s", i.clusterName, i.instanceGroupName)
}

func parseInstanceGroupID(id string) instanceGroupID {
	split := strings.Split(id, "/")
	if len(split) == 2 {
		return instanceGroupID{
			clusterName:       split[0],
			instanceGroupName: split[1],
		}
	}
	return instanceGroupID{}
}

func getInstanceGroup(d *schema.ResourceData, m interface{}) (*kops.InstanceGroup, error) {
	groupID := parseInstanceGroupID(d.Id())
	clientset := m.(*ProviderConfig).clientset
	cluster, err := clientset.GetCluster(groupID.clusterName)
	if err != nil {
		return nil, err
	}
	return clientset.InstanceGroupsFor(cluster).Get(groupID.instanceGroupName, v1.GetOptions{})
}
