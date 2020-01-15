package provider

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/pkg/instancegroups"
	"k8s.io/kops/upup/pkg/fi/cloudup"
)

func resourceInstanceGroupCreate(d *schema.ResourceData, m interface{}) error {
	clientset := m.(*ProviderConfig).clientset

	log.Println("Expanding Metadata")
	metadata := expandObjectMeta(sectionData(d, "metadata"))

	log.Println("Expanding Instance Group Spec")
	spec := expandInstanceGroupSpec(sectionData(d, "spec"))

	ig := &kops.InstanceGroup{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	clusterName := ig.ObjectMeta.Labels[kops.LabelClusterName]
	if clusterName == "" {
		return fmt.Errorf("must specify %q label with cluster name to create instanceGroup", kops.LabelClusterName)
	}

	cluster, err := clientset.GetCluster(clusterName)
	if cluster == nil {
		return fmt.Errorf("cluster %q not found", clusterName)
	}
	if cluster == nil {
		return fmt.Errorf("cluster %q not found", clusterName)
	}

	_, err = clientset.InstanceGroupsFor(cluster).Create(ig)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Errorf("instanceGroup %q already exists", ig.ObjectMeta.Name)
		}
		return fmt.Errorf("error creating instanceGroup: %v", err)
	} else {
		log.Printf("Created instancegroup/%s\n", ig.ObjectMeta.Name)
	}

	d.SetId(ig.ObjectMeta.Name)

	return resourceInstanceGroupRead(d, m)
}

func resourceInstanceGroupRead(d *schema.ResourceData, m interface{}) error {
	ig, err := getInstanceGroup(d, m)
	if err != nil {
		return err
	}
	if err := d.Set("metadata", flattenObjectMeta(ig.ObjectMeta)); err != nil {
		return err
	}
	if err := d.Set("spec", flattenInstanceGroupSpec(ig.Spec)); err != nil {
		return err
	}
	return nil
}

func resourceInstanceGroupUpdate(d *schema.ResourceData, m interface{}) error {
	if ok, _ := resourceInstanceGroupExists(d, m); !ok {
		d.SetId("")
		return nil
	}

	log.Println("Expanding Metadata")
	metadata := expandObjectMeta(sectionData(d, "metadata"))

	log.Println("Expanding Instance Group Spec")
	spec := expandInstanceGroupSpec(sectionData(d, "spec"))

	ig := &kops.InstanceGroup{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	clusterName := ig.ObjectMeta.Labels[kops.LabelClusterName]
	if clusterName == "" {
		return fmt.Errorf("must specify %q label with cluster name to update instanceGroup", kops.LabelClusterName)
	}

	clientset := m.(*ProviderConfig).clientset

	cluster, err := clientset.GetCluster(clusterName)
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("cluster %q not found", clusterName)
		}
		return fmt.Errorf("error fetching cluster %q: %v", clusterName, err)
	}
	// check if the instancegroup exists already
	igName := ig.ObjectMeta.Name
	ig, err = clientset.InstanceGroupsFor(cluster).Get(igName, v1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("instanceGroup: %v does not exist", igName)
		} else {
			return fmt.Errorf("unable to check for instanceGroup: %v", err)
		}
	}

	_, err = clientset.InstanceGroupsFor(cluster).Update(ig)
	if err != nil {
		return fmt.Errorf("error replacing instanceGroup: %v", err)
	}

	return resourceInstanceGroupRead(d, m)
}

func resourceInstanceGroupDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Expanding Metadata")
	metadata := expandObjectMeta(sectionData(d, "metadata"))

	log.Println("Expanding Instance Group Spec")
	spec := expandInstanceGroupSpec(sectionData(d, "spec"))

	ig := &kops.InstanceGroup{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	clusterName := ig.ObjectMeta.Labels[kops.LabelClusterName]
	if clusterName == "" {
		return fmt.Errorf("must specify %q label with cluster name to update instanceGroup", kops.LabelClusterName)
	}

	clientset := m.(*ProviderConfig).clientset

	cluster, err := clientset.GetCluster(clusterName)
	if err != nil {
		if errors.IsNotFound(err) {
			return fmt.Errorf("cluster %q not found", clusterName)
		}
		return fmt.Errorf("error fetching cluster %q: %v", clusterName, err)
	}

	group, err := clientset.InstanceGroupsFor(cluster).Get(ig.ObjectMeta.Name, v1.GetOptions{})
	if err != nil {
		return fmt.Errorf("error reading InstanceGroup %q: %v", ig.ObjectMeta.Name, err)
	}
	if group == nil {
		return fmt.Errorf("InstanceGroup %q not found", ig.ObjectMeta.Name)
	}

	cloud, err := cloudup.BuildCloud(cluster)
	if err != nil {
		return err
	}

	log.Printf("InstanceGroup %q found for deletion\n", ig.ObjectMeta.Name)

	dig := &instancegroups.DeleteInstanceGroup{}
	dig.Cluster = cluster
	dig.Cloud = cloud
	dig.Clientset = clientset

	return dig.DeleteInstanceGroup(group)
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

func getInstanceGroup(d *schema.ResourceData, m interface{}) (*kops.InstanceGroup, error) {
	clientset := m.(*ProviderConfig).clientset

	log.Println("Expanding Metadata")
	metadata := expandObjectMeta(sectionData(d, "metadata"))

	log.Println("Expanding Instance Group Spec")
	spec := expandInstanceGroupSpec(sectionData(d, "spec"))

	ig := &kops.InstanceGroup{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	clusterName := ig.ObjectMeta.Labels[kops.LabelClusterName]
	if clusterName == "" {
		return nil, fmt.Errorf("must specify %q label with cluster name to create instanceGroup", kops.LabelClusterName)
	}

	cluster, err := clientset.GetCluster(clusterName)
	if cluster == nil {
		return nil, fmt.Errorf("cluster %q not found", clusterName)
	}
	if cluster == nil {
		return nil, fmt.Errorf("cluster %q not found", clusterName)
	}

	ig, err = clientset.InstanceGroupsFor(cluster).Get(ig.Name, v1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to read instancegroup %q", ig.Name)
	}

	return ig, nil
}
