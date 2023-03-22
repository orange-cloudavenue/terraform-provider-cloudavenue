# If `vDC` is not specified, the default `vDC` will be used
# The `affinityRuleIdentifier` can be either a name or an ID. If it is a name, it will succeed only if the name is unique.
terraform import cloudavenue_vm_affinity_rule.example affinityRuleIdentifier

# or you can specify the vDC
terraform import cloudavenue_vm_affinity_rule.example myVDC.affinityRuleIdentifier
