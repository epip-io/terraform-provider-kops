// Code generated by engine.go; DO NOT EDIT.

package provider

import (
	"k8s.io/kops/pkg/apis/kops"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func expandFileAssetSpecSlice(in interface{}) ([]kops.FileAssetSpec, bool) {
	d := in.([]interface{})
	out := make([]kops.FileAssetSpec , len(d))

  if len(d) < 1 {
    return out, true
  }

	for i := 0; i < len(d); i++ {
		out[i] = kops.FileAssetSpec{}

		if v, ok := d[i].(map[string]interface{})["content"]; ok {
      if value, e := expandString(v); !e {
        out[i].Content = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["is_base_64"]; ok {
      if value, e := expandBool(v); !e {
        out[i].IsBase64 = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["name"]; ok {
      if value, e := expandString(v); !e {
        out[i].Name = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["path"]; ok {
      if value, e := expandString(v); !e {
        out[i].Path = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["roles"]; ok {
      if value, e := expandInstanceGroupRoleSlice(v); !e {
        out[i].Roles = value
      }
		}
	}

	return out, false
}

func expandHookSpecSlice(in interface{}) ([]kops.HookSpec, bool) {
	d := in.([]interface{})
	out := make([]kops.HookSpec , len(d))

  if len(d) < 1 {
    return out, true
  }

	for i := 0; i < len(d); i++ {
		out[i] = kops.HookSpec{}

		if v, ok := d[i].(map[string]interface{})["before"]; ok {
      if value, e := expandStringSlice(v); !e {
        out[i].Before = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["disabled"]; ok {
      if value, e := expandBool(v); !e {
        out[i].Disabled = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["exec_container"]; ok {
      if value, e := expandExecContainerAction(v); !e {
        out[i].ExecContainer = &value
      }
		}

		if v, ok := d[i].(map[string]interface{})["manifest"]; ok {
      if value, e := expandString(v); !e {
        out[i].Manifest = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["name"]; ok {
      if value, e := expandString(v); !e {
        out[i].Name = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["requires"]; ok {
      if value, e := expandStringSlice(v); !e {
        out[i].Requires = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["roles"]; ok {
      if value, e := expandInstanceGroupRoleSlice(v); !e {
        out[i].Roles = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["use_raw_manifest"]; ok {
      if value, e := expandBool(v); !e {
        out[i].UseRawManifest = value
      }
		}
	}

	return out, false
}

func expandInitializers(in interface{}) (v1.Initializers, bool) {
	d := in.([]interface{})
	out := v1.Initializers{}

	if len(d) < 1 {
		return out, true
	}

	m := d[0].(map[string]interface{})

	{
		if value, e := expandInitializerSlice(m["pending"]); !e {
      out.Pending = value
    }
	}

	if v, ok := m["result"]; ok {
		if value, e := expandStatus(v); !e {
      out.Result = &value
    }
	}

  if isEmpty(out) {
    return out, true
  }

	return out, false
}

func expandKubeletConfigSpec(in interface{}) (kops.KubeletConfigSpec, bool) {
	d := in.([]interface{})
	out := kops.KubeletConfigSpec{}

	if len(d) < 1 {
		return out, true
	}

	m := d[0].(map[string]interface{})

	if v, ok := m["api_servers"]; ok {
		if value, e := expandString(v); !e {
      out.APIServers = value
    }
	}

	if v, ok := m["allow_privileged"]; ok {
		if value, e := expandBool(v); !e {
      out.AllowPrivileged = &value
    }
	}

	if v, ok := m["allowed_unsafe_sysctls"]; ok {
		if value, e := expandStringSlice(v); !e {
      out.AllowedUnsafeSysctls = value
    }
	}

	if v, ok := m["anonymous_auth"]; ok {
		if value, e := expandBool(v); !e {
      out.AnonymousAuth = &value
    }
	}

	if v, ok := m["authentication_token_webhook"]; ok {
		if value, e := expandBool(v); !e {
      out.AuthenticationTokenWebhook = &value
    }
	}

	if v, ok := m["authentication_token_webhook_cache_ttl"]; ok {
		if value, e := expandDuration(v); !e {
      out.AuthenticationTokenWebhookCacheTTL = &value
    }
	}

	if v, ok := m["authorization_mode"]; ok {
		if value, e := expandString(v); !e {
      out.AuthorizationMode = value
    }
	}

	if v, ok := m["babysit_daemons"]; ok {
		if value, e := expandBool(v); !e {
      out.BabysitDaemons = &value
    }
	}

	if v, ok := m["bootstrap_kubeconfig"]; ok {
		if value, e := expandString(v); !e {
      out.BootstrapKubeconfig = value
    }
	}

	if v, ok := m["cpu_cfs_quota"]; ok {
		if value, e := expandBool(v); !e {
      out.CPUCFSQuota = &value
    }
	}

	if v, ok := m["cpu_cfs_quota_period"]; ok {
		if value, e := expandDuration(v); !e {
      out.CPUCFSQuotaPeriod = &value
    }
	}

	if v, ok := m["cgroup_root"]; ok {
		if value, e := expandString(v); !e {
      out.CgroupRoot = value
    }
	}

	if v, ok := m["client_ca_file"]; ok {
		if value, e := expandString(v); !e {
      out.ClientCAFile = value
    }
	}

	if v, ok := m["cloud_provider"]; ok {
		if value, e := expandString(v); !e {
      out.CloudProvider = value
    }
	}

	if v, ok := m["cluster_dns"]; ok {
		if value, e := expandString(v); !e {
      out.ClusterDNS = value
    }
	}

	if v, ok := m["cluster_domain"]; ok {
		if value, e := expandString(v); !e {
      out.ClusterDomain = value
    }
	}

	if v, ok := m["configure_cbr_0"]; ok {
		if value, e := expandBool(v); !e {
      out.ConfigureCBR0 = &value
    }
	}

	if v, ok := m["cpu_manager_policy"]; ok {
		if value, e := expandString(v); !e {
      out.CpuManagerPolicy = value
    }
	}

	if v, ok := m["docker_disable_shared_pid"]; ok {
		if value, e := expandBool(v); !e {
      out.DockerDisableSharedPID = &value
    }
	}

	if v, ok := m["enable_custom_metrics"]; ok {
		if value, e := expandBool(v); !e {
      out.EnableCustomMetrics = &value
    }
	}

	if v, ok := m["enable_debugging_handlers"]; ok {
		if value, e := expandBool(v); !e {
      out.EnableDebuggingHandlers = &value
    }
	}

	if v, ok := m["enforce_node_allocatable"]; ok {
		if value, e := expandString(v); !e {
      out.EnforceNodeAllocatable = value
    }
	}

	if v, ok := m["eviction_hard"]; ok {
		if value, e := expandString(v); !e {
      out.EvictionHard = &value
    }
	}

	if v, ok := m["eviction_max_pod_grace_period"]; ok {
		if value, e := expandInt32(v); !e {
      out.EvictionMaxPodGracePeriod = value
    }
	}

	if v, ok := m["eviction_minimum_reclaim"]; ok {
		if value, e := expandString(v); !e {
      out.EvictionMinimumReclaim = value
    }
	}

	if v, ok := m["eviction_pressure_transition_period"]; ok {
		if value, e := expandDuration(v); !e {
      out.EvictionPressureTransitionPeriod = &value
    }
	}

	if v, ok := m["eviction_soft"]; ok {
		if value, e := expandString(v); !e {
      out.EvictionSoft = value
    }
	}

	if v, ok := m["eviction_soft_grace_period"]; ok {
		if value, e := expandString(v); !e {
      out.EvictionSoftGracePeriod = value
    }
	}

	if v, ok := m["experimental_allowed_unsafe_sysctls"]; ok {
		if value, e := expandStringSlice(v); !e {
      out.ExperimentalAllowedUnsafeSysctls = value
    }
	}

	if v, ok := m["fail_swap_on"]; ok {
		if value, e := expandBool(v); !e {
      out.FailSwapOn = &value
    }
	}

	if v, ok := m["feature_gates"]; ok {
		if value, e := expandStringMap(v); !e {
      out.FeatureGates = value
    }
	}

	if v, ok := m["hairpin_mode"]; ok {
		if value, e := expandString(v); !e {
      out.HairpinMode = value
    }
	}

	if v, ok := m["hostname_override"]; ok {
		if value, e := expandString(v); !e {
      out.HostnameOverride = value
    }
	}

	if v, ok := m["image_gc_high_threshold_percent"]; ok {
		if value, e := expandInt32(v); !e {
      out.ImageGCHighThresholdPercent = &value
    }
	}

	if v, ok := m["image_gc_low_threshold_percent"]; ok {
		if value, e := expandInt32(v); !e {
      out.ImageGCLowThresholdPercent = &value
    }
	}

	if v, ok := m["image_pull_progress_deadline"]; ok {
		if value, e := expandDuration(v); !e {
      out.ImagePullProgressDeadline = &value
    }
	}

	if v, ok := m["kube_reserved"]; ok {
		if value, e := expandStringMap(v); !e {
      out.KubeReserved = value
    }
	}

	if v, ok := m["kube_reserved_cgroup"]; ok {
		if value, e := expandString(v); !e {
      out.KubeReservedCgroup = value
    }
	}

	if v, ok := m["kubeconfig_path"]; ok {
		if value, e := expandString(v); !e {
      out.KubeconfigPath = value
    }
	}

	if v, ok := m["kubelet_cgroups"]; ok {
		if value, e := expandString(v); !e {
      out.KubeletCgroups = value
    }
	}

	if v, ok := m["log_level"]; ok {
		if value, e := expandInt32(v); !e {
      out.LogLevel = &value
    }
	}

	if v, ok := m["max_pods"]; ok {
		if value, e := expandInt32(v); !e {
      out.MaxPods = &value
    }
	}

	if v, ok := m["network_plugin_mtu"]; ok {
		if value, e := expandInt32(v); !e {
      out.NetworkPluginMTU = &value
    }
	}

	if v, ok := m["network_plugin_name"]; ok {
		if value, e := expandString(v); !e {
      out.NetworkPluginName = value
    }
	}

	if v, ok := m["node_labels"]; ok {
		if value, e := expandStringMap(v); !e {
      out.NodeLabels = value
    }
	}

	if v, ok := m["node_status_update_frequency"]; ok {
		if value, e := expandDuration(v); !e {
      out.NodeStatusUpdateFrequency = &value
    }
	}

	if v, ok := m["non_masquerade_cidr"]; ok {
		if value, e := expandString(v); !e {
      out.NonMasqueradeCIDR = value
    }
	}

	if v, ok := m["nvidia_gp_uss"]; ok {
		if value, e := expandInt32(v); !e {
      out.NvidiaGPUs = value
    }
	}

	if v, ok := m["pod_cidr"]; ok {
		if value, e := expandString(v); !e {
      out.PodCIDR = value
    }
	}

	if v, ok := m["pod_infra_container_image"]; ok {
		if value, e := expandString(v); !e {
      out.PodInfraContainerImage = value
    }
	}

	if v, ok := m["pod_manifest_path"]; ok {
		if value, e := expandString(v); !e {
      out.PodManifestPath = value
    }
	}

	if v, ok := m["read_only_port"]; ok {
		if value, e := expandInt32(v); !e {
      out.ReadOnlyPort = &value
    }
	}

	if v, ok := m["reconcile_cidr"]; ok {
		if value, e := expandBool(v); !e {
      out.ReconcileCIDR = &value
    }
	}

	if v, ok := m["register_node"]; ok {
		if value, e := expandBool(v); !e {
      out.RegisterNode = &value
    }
	}

	if v, ok := m["register_schedulable"]; ok {
		if value, e := expandBool(v); !e {
      out.RegisterSchedulable = &value
    }
	}

	if v, ok := m["registry_burst"]; ok {
		if value, e := expandInt32(v); !e {
      out.RegistryBurst = &value
    }
	}

	if v, ok := m["registry_pull_qps"]; ok {
		if value, e := expandInt32(v); !e {
      out.RegistryPullQPS = &value
    }
	}

	if v, ok := m["require_kubeconfig"]; ok {
		if value, e := expandBool(v); !e {
      out.RequireKubeconfig = &value
    }
	}

	if v, ok := m["resolv_conf"]; ok {
		if value, e := expandString(v); !e {
      out.ResolverConfig = &value
    }
	}

	if v, ok := m["root_dir"]; ok {
		if value, e := expandString(v); !e {
      out.RootDir = value
    }
	}

	if v, ok := m["runtime_cgroups"]; ok {
		if value, e := expandString(v); !e {
      out.RuntimeCgroups = value
    }
	}

	if v, ok := m["runtime_request_timeout"]; ok {
		if value, e := expandDuration(v); !e {
      out.RuntimeRequestTimeout = &value
    }
	}

	if v, ok := m["seccomp_profile_root"]; ok {
		if value, e := expandString(v); !e {
      out.SeccompProfileRoot = &value
    }
	}

	if v, ok := m["serialize_image_pulls"]; ok {
		if value, e := expandBool(v); !e {
      out.SerializeImagePulls = &value
    }
	}

	if v, ok := m["streaming_connection_idle_timeout"]; ok {
		if value, e := expandDuration(v); !e {
      out.StreamingConnectionIdleTimeout = &value
    }
	}

	if v, ok := m["system_cgroups"]; ok {
		if value, e := expandString(v); !e {
      out.SystemCgroups = value
    }
	}

	if v, ok := m["system_reserved"]; ok {
		if value, e := expandStringMap(v); !e {
      out.SystemReserved = value
    }
	}

	if v, ok := m["system_reserved_cgroup"]; ok {
		if value, e := expandString(v); !e {
      out.SystemReservedCgroup = value
    }
	}

	if v, ok := m["tls_cert_file"]; ok {
		if value, e := expandString(v); !e {
      out.TLSCertFile = value
    }
	}

	if v, ok := m["tls_cipher_suites"]; ok {
		if value, e := expandStringSlice(v); !e {
      out.TLSCipherSuites = value
    }
	}

	if v, ok := m["tls_min_version"]; ok {
		if value, e := expandString(v); !e {
      out.TLSMinVersion = value
    }
	}

	if v, ok := m["tls_private_key_file"]; ok {
		if value, e := expandString(v); !e {
      out.TLSPrivateKeyFile = value
    }
	}

	if v, ok := m["taints"]; ok {
		if value, e := expandStringSlice(v); !e {
      out.Taints = value
    }
	}

	if v, ok := m["volume_plugin_directory"]; ok {
		if value, e := expandString(v); !e {
      out.VolumePluginDirectory = value
    }
	}

	if v, ok := m["volume_stats_agg_period"]; ok {
		if value, e := expandDuration(v); !e {
      out.VolumeStatsAggPeriod = &value
    }
	}

  if isEmpty(out) {
    return out, true
  }

	return out, false
}

func expandManagedFieldsEntrySlice(in interface{}) ([]v1.ManagedFieldsEntry, bool) {
	d := in.([]interface{})
	out := make([]v1.ManagedFieldsEntry , len(d))

  if len(d) < 1 {
    return out, true
  }

	for i := 0; i < len(d); i++ {
		out[i] = v1.ManagedFieldsEntry{}

		if v, ok := d[i].(map[string]interface{})["api_version"]; ok {
      if value, e := expandString(v); !e {
        out[i].APIVersion = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["fields"]; ok {
      if value, e := expandFields(v); !e {
        out[i].Fields = &value
      }
		}

		if v, ok := d[i].(map[string]interface{})["manager"]; ok {
      if value, e := expandString(v); !e {
        out[i].Manager = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["operation"]; ok {
      if value, e := expandManagedFieldsOperationType(v); !e {
        out[i].Operation = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["time"]; ok {
      if value, e := expandTime(v); !e {
        out[i].Time = &value
      }
		}
	}

	return out, false
}

func expandObjectMeta(in interface{}) (v1.ObjectMeta, bool) {
	d := in.([]interface{})
	out := v1.ObjectMeta{}

	if len(d) < 1 {
		return out, true
	}

	m := d[0].(map[string]interface{})

	if v, ok := m["annotations"]; ok {
		if value, e := expandStringMap(v); !e {
      out.Annotations = value
    }
	}

	if v, ok := m["cluster_name"]; ok {
		if value, e := expandString(v); !e {
      out.ClusterName = value
    }
	}

	if v, ok := m["creation_timestamp"]; ok {
		if value, e := expandTime(v); !e {
      out.CreationTimestamp = value
    }
	}

	if v, ok := m["deletion_grace_period_seconds"]; ok {
		if value, e := expandInt64(v); !e {
      out.DeletionGracePeriodSeconds = &value
    }
	}

	if v, ok := m["deletion_timestamp"]; ok {
		if value, e := expandTime(v); !e {
      out.DeletionTimestamp = &value
    }
	}

	if v, ok := m["finalizers"]; ok {
		if value, e := expandStringSlice(v); !e {
      out.Finalizers = value
    }
	}

	if v, ok := m["generate_name"]; ok {
		if value, e := expandString(v); !e {
      out.GenerateName = value
    }
	}

	if v, ok := m["generation"]; ok {
		if value, e := expandInt64(v); !e {
      out.Generation = value
    }
	}

	if v, ok := m["initializers"]; ok {
		if value, e := expandInitializers(v); !e {
      out.Initializers = &value
    }
	}

	if v, ok := m["labels"]; ok {
		if value, e := expandStringMap(v); !e {
      out.Labels = value
    }
	}

	if v, ok := m["managed_fields"]; ok {
		if value, e := expandManagedFieldsEntrySlice(v); !e {
      out.ManagedFields = value
    }
	}

	if v, ok := m["name"]; ok {
		if value, e := expandString(v); !e {
      out.Name = value
    }
	}

	if v, ok := m["namespace"]; ok {
		if value, e := expandString(v); !e {
      out.Namespace = value
    }
	}

	if v, ok := m["owner_references"]; ok {
		if value, e := expandOwnerReferenceSlice(v); !e {
      out.OwnerReferences = value
    }
	}

	if v, ok := m["resource_version"]; ok {
		if value, e := expandString(v); !e {
      out.ResourceVersion = value
    }
	}

	if v, ok := m["self_link"]; ok {
		if value, e := expandString(v); !e {
      out.SelfLink = value
    }
	}

	if v, ok := m["uid"]; ok {
		if value, e := expandUID(v); !e {
      out.UID = value
    }
	}

  if isEmpty(out) {
    return out, true
  }

	return out, false
}

func expandOwnerReferenceSlice(in interface{}) ([]v1.OwnerReference, bool) {
	d := in.([]interface{})
	out := make([]v1.OwnerReference , len(d))

  if len(d) < 1 {
    return out, true
  }

	for i := 0; i < len(d); i++ {
		out[i] = v1.OwnerReference{}

		if v, ok := d[i].(map[string]interface{})["api_version"]; ok {
      if value, e := expandString(v); !e {
        out[i].APIVersion = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["block_owner_deletion"]; ok {
      if value, e := expandBool(v); !e {
        out[i].BlockOwnerDeletion = &value
      }
		}

		if v, ok := d[i].(map[string]interface{})["controller"]; ok {
      if value, e := expandBool(v); !e {
        out[i].Controller = &value
      }
		}

		if v, ok := d[i].(map[string]interface{})["kind"]; ok {
      if value, e := expandString(v); !e {
        out[i].Kind = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["name"]; ok {
      if value, e := expandString(v); !e {
        out[i].Name = value
      }
		}

		if v, ok := d[i].(map[string]interface{})["uid"]; ok {
      if value, e := expandUID(v); !e {
        out[i].UID = value
      }
		}
	}

	return out, false
}

func expandString(in interface{}) (string, bool) {
	d := in.(string)
	r := string(d)
  out := r
  
  if out == "" {
    return out, true
  }

	return out, false
}

func expandStringMap(in interface{}) (map[string]string, bool) {
	out := make(map[string]string)
  d := in.(map[string]interface{})
  
  if len(d) < 1 {
    return out, true
  }

	for k, v := range d {
		out[k] = v.(string)
	}

	return out, false
}

func expandStringSlice(in interface{}) ([]string, bool) {
	d := in.([]interface{})
	out := make([]string , len(d))

  if len(d) < 1 {
    return out, true
  }

	for i := 0; i < len(d); i++ {
		out[i] = d[i].(string)
	}

	return out, false
}

func expandTypeMeta(in interface{}) (v1.TypeMeta, bool) {
	d := in.([]interface{})
	out := v1.TypeMeta{}

	if len(d) < 1 {
		return out, true
	}

	m := d[0].(map[string]interface{})

	if v, ok := m["api_version"]; ok {
		if value, e := expandString(v); !e {
      out.APIVersion = value
    }
	}

	if v, ok := m["kind"]; ok {
		if value, e := expandString(v); !e {
      out.Kind = value
    }
	}

  if isEmpty(out) {
    return out, true
  }

	return out, false
}

func expandUID(in interface{}) (types.UID, bool) {
	d := in.(string)
	r := types.UID(d)
  out := r
  
  if out == "" {
    return out, true
  }

	return out, false
}
