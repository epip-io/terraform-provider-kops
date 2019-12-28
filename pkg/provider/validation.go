package provider

var (
	V1CauseType = []string{
		"FieldValueNotFound",
		"FieldValueRequired",
		"FieldValueDuplicate",
		"FieldValueInvalid",
		"FieldValueNotSupported",
		"UnexpectedServerResponse",
		"FieldManagerConflict",
	}

	V1ManagedFieldsOperationType = []string{
		"Apply",
		"Update",
	}

	V1StatusReason = []string{
		"",
		"Unauthorized",
		"Forbidden",
		"NotFound",
		"AlreadyExists",
		"Conflict",
		"Gone",
		"Invalid",
		"ServiceTimeout",
		"Timeout",
		"TooManyRequests",
		"BadRequest",
		"MethodNotAllowed",
		"NotAcceptable",
		"RequestEntityTooLarge",
		"UnsupportedMediaType",
		"InternalError",
		"Expired",
		"ServiceUnavailable",
	}

	KopsDNSType = []string{
		"Public",
		"Private",
	}

	KopsEtcdProviderType = []string{
		"Manager",
		"Legacy",
	}

	KopsLoadBalancerType = []string{
		"Public",
		"Internal",
	}

	KopsSubnetType = []string{
		"Public",
		"Private",
		"Utility",
	}

	KopsInstanceGroupRole = []string{
		"Master",
		"Node",
		"Bastion",
	}
)
