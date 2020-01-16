package provider

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/klog"
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/kops/pkg/sshcredentials"
	"k8s.io/kops/upup/pkg/fi"
)

const SecretTypeSSHPublicKey = kops.KeysetType("SSHPublicKey")

func resourceSSHCredentialCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("Expanding Metadata")
	metadata := expandObjectMeta(sectionData(d, "metadata"))

	log.Println("Expanding SSH Credential Spec")
	spec := expandSSHCredentialSpec(sectionData(d, "spec"))

	sshc := &kops.SSHCredential{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	clusterName := sshc.ObjectMeta.Labels[kops.LabelClusterName]
	if clusterName == "" {
		return fmt.Errorf("must specify %q label with cluster name to create SSHCredential", kops.LabelClusterName)
	}
	if sshc.Spec.PublicKey == "" {
		return fmt.Errorf("spec.PublicKey is required")
	}

	clientset := m.(*ProviderConfig).clientset

	cluster, err := clientset.GetCluster(clusterName)
	if err != nil {
		return err
	}

	sshCredentialStore, err := clientset.SSHCredentialStore(cluster)
	if err != nil {
		return err
	}

	sshKeyArr := []byte(sshc.Spec.PublicKey)
	err = sshCredentialStore.AddSSHPublicKey("admin", sshKeyArr)
	if err != nil {
		return err
	}

	fingerprint, err := sshcredentials.Fingerprint(sshc.Spec.PublicKey)
	if err != nil {
		klog.Error("unable to compute fingerprint for public key")
	}

	d.SetId(fingerprint)

	return resourceSSHCredentialRead(d, m)
}

func resourceSSHCredentialRead(d *schema.ResourceData, m interface{}) error {
	sshc, err := getSSHCredential(d, m)
	if err != nil {
		return err
	}
	if err := d.Set("metadata", flattenObjectMeta(sshc.ObjectMeta)); err != nil {
		return err
	}
	if err := d.Set("spec", flattenSSHCredentialSpec(sshc.Spec)); err != nil {
		return err
	}
	return nil
}

func resourceSSHCredentialUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSSHCredentialRead(d, m)
}

func resourceSSHCredentialDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("Expanding Metadata")
	metadata := expandObjectMeta(sectionData(d, "metadata"))

	log.Println("Expanding SSH Credential Spec")
	spec := expandSSHCredentialSpec(sectionData(d, "spec"))

	sshc := &kops.SSHCredential{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	clusterName := sshc.ObjectMeta.Labels[kops.LabelClusterName]
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

	sshCredentialStore, err := clientset.SSHCredentialStore(cluster)
	if err != nil {
		return err
	}

	sshc, err = getSSHCredential(d, m)
	if err != nil {
		return err
	}

	return sshCredentialStore.DeleteSSHCredential(sshc)
}

func resourceSSHCredentialExists(d *schema.ResourceData, m interface{}) (bool, error) {
	_, err := getSSHCredential(d, m)
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func getSSHCredential(d *schema.ResourceData, m interface{}) (*kops.SSHCredential, error) {
	log.Println("Expanding Metadata")
	metadata := expandObjectMeta(sectionData(d, "metadata"))

	log.Println("Expanding SSH Credential Spec")
	spec := expandSSHCredentialSpec(sectionData(d, "spec"))

	sshc := &kops.SSHCredential{
		ObjectMeta: metadata,
		Spec:       spec,
	}

	clusterName := sshc.ObjectMeta.Labels[kops.LabelClusterName]
	if clusterName == "" {
		return nil, fmt.Errorf("must specify %q label with cluster name to update instanceGroup", kops.LabelClusterName)
	}

	clientset := m.(*ProviderConfig).clientset

	cluster, err := clientset.GetCluster(clusterName)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, fmt.Errorf("cluster %q not found", clusterName)
		}
		return nil, fmt.Errorf("error fetching cluster %q: %v", clusterName, err)
	}

	fingerprint, err := sshcredentials.Fingerprint(sshc.Spec.PublicKey)
	if err != nil {
		klog.Error("unable to compute fingerprint for public key")
	}

	sshCredentialStore, err := clientset.SSHCredentialStore(cluster)
	if err != nil {
		return nil, err
	}

	l, err := sshCredentialStore.ListSSHCredentials()
	if err != nil {
		return nil, fmt.Errorf("error listing SSH credentials %v", err)
	}

	secrets := make([]*fi.KeystoreItem, 0)
	for i := range l {
		id, err := sshcredentials.Fingerprint(l[i].Spec.PublicKey)
		if err != nil {
			klog.Warningf("unable to compute fingerprint for public key %q", l[i].Name)
		}
		item := &fi.KeystoreItem{
			Name: l[i].Name,
			Id:   id,
			Type: SecretTypeSSHPublicKey,
		}
		if l[i].Spec.PublicKey != "" {
			item.Data = []byte(l[i].Spec.PublicKey)
		}

		secrets = append(secrets, item)
	}

	if fingerprint != "" {
		var matches []*fi.KeystoreItem
		for _, s := range secrets {
			if s.Id == fingerprint {
				matches = append(matches, s)
			}
		}
		secrets = matches
	}

	if len(secrets) == 0 {
		return nil, fmt.Errorf("secret not found")
	}

	sshCredential := &kops.SSHCredential{}
	sshCredential.Name = secrets[0].Name
	if secrets[0].Data != nil {
		sshCredential.Spec.PublicKey = string(secrets[0].Data)
	}

	return sshCredential, nil
}
