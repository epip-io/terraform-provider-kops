package provider

import (
	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/pkg/assets"
	"k8s.io/kops/upup/pkg/fi/cloudup"
	"log"

	"github.com/epip-io/terraform-provider-kops/pkg/convert"
)

func resourceClusterCreate(d *schema.ResourceData, m interface{}) error {
	clientset := m.(*ProviderConfig).clientset

	log.Println("Expanding Metadata")
	metadata := convert.ExpandObjectMeta(sectionData(d, "metadata"))

	log.Println("Expanding Cluster Spec")
	spec := convert.ExpandClusterSpec(sectionData(d, "spec"))

	log.Println("Creating Cluster")
	cluster, err := clientset.CreateCluster(&kops.Cluster{
		ObjectMeta: metadata,
		Spec:       spec,
	})
	if err != nil {
		return err
	}

	cluster, err = clientset.GetCluster(cluster.Name)
	if err != nil {
		return err
	}

	assetBuilder := assets.NewAssetBuilder(cluster, "")
	fullCluster, err := cloudup.PopulateClusterSpec(clientset, cluster, assetBuilder)
	if err != nil {
		return err
	}

	log.Println("Updating Cluster")
	_, err = clientset.UpdateCluster(fullCluster, nil)
	if err != nil {
		return err
	}

	d.SetId(cluster.Name)

	return resourceClusterRead(d, m)
}

func resourceClusterRead(d *schema.ResourceData, m interface{}) error {
	cluster, err := getCluster(d, m)
	if err != nil {
		return err
	}
	if err := d.Set("metadata", convert.FlattenObjectMeta(cluster.ObjectMeta)); err != nil {
		return err
	}
	if err := d.Set("spec", convert.FlattenClusterSpec(cluster.Spec)); err != nil {
		return err
	}
	return nil
}

func resourceClusterUpdate(d *schema.ResourceData, m interface{}) error {
	if ok, _ := resourceClusterExists(d, m); !ok {
		d.SetId("")
		return nil
	}

	clientset := m.(*ProviderConfig).clientset

	_, err := clientset.UpdateCluster(&kops.Cluster{
		ObjectMeta: convert.ExpandObjectMeta(sectionData(d, "metadata")),
		Spec:       convert.ExpandClusterSpec(sectionData(d, "spec")),
	}, nil)

	if err != nil {
		return err
	}

	return resourceClusterRead(d, m)
}

func resourceClusterDelete(d *schema.ResourceData, m interface{}) error {
	clientset := m.(*ProviderConfig).clientset
	cluster, err := getCluster(d, m)
	if err != nil {
		return err
	}

	return clientset.DeleteCluster(cluster)
}

func resourceClusterExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, err := getCluster(d, m)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func getCluster(d *schema.ResourceData, m interface{}) (*kops.Cluster, error) {
	clientset := m.(*ProviderConfig).clientset
	cluster, err := clientset.GetCluster(d.Id())
	return cluster, err
}

func sectionData(d *schema.ResourceData, section string) map[string]interface{} {
	return d.Get(section).([]interface{})[0].(map[string]interface{})
}
