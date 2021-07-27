package fw

import (
	"fmt"

	"github.com/aquasecurity/tfsec/pkg/result"
	"github.com/aquasecurity/tfsec/pkg/severity"

	"github.com/aquasecurity/tfsec/pkg/provider"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/cidr"
	"github.com/aquasecurity/tfsec/internal/app/tfsec/hclcontext"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/block"

	"github.com/aquasecurity/tfsec/pkg/rule"

	"github.com/aquasecurity/tfsec/internal/app/tfsec/scanner"
)

func init() {
	scanner.RegisterCheckRule(rule.Rule{
		Service:   "fw",
		ShortCode: "no-public-access",
		Documentation: rule.RuleDocumentation{
			Summary:     "A firewall rule allows traffic from/to the public internet",
			Explanation: `Opening up ports to the public internet is generally to be avoided. You should restrict access to IP addresses or ranges that explicitly require it where possible.`,
			Impact:      "Exposure of infrastructure to the public internet",
			Resolution:  "Employ more restrictive firewall rules",
			BadExample: `
resource "openstack_fw_rule_v1" "rule_1" {
	name             = "my_rule"
	description      = "let anyone in"
	action           = "allow"
	protocol         = "tcp"
	destination_port = "22"
	enabled          = "true"
}
			`,
			GoodExample: `
resource "openstack_fw_rule_v1" "rule_1" {
	name                   = "my_rule"
	description            = "don't let just anyone in"
	action                 = "allow"
	protocol               = "tcp"
	destination_ip_address = "10.10.10.1"
	source_ip_address      = "10.10.10.2"
	destination_port       = "22"
	enabled                = "true"
}
			`,
			Links: []string{
				"https://registry.terraform.io/providers/terraform-provider-openstack/openstack/latest/docs/resources/fw_rule_v1",
			},
		},
		Provider:        provider.OpenStackProvider,
		RequiredTypes:   []string{"resource"},
		RequiredLabels:  []string{"openstack_fw_rule_v1"},
		DefaultSeverity: severity.Medium,
		CheckFunc: func(set result.Set, resourceBlock block.Block, _ *hclcontext.Context) {

			if enabledAttr := resourceBlock.GetAttribute("enabled"); enabledAttr != nil && enabledAttr.IsFalse() {
				return
			}

			if actionAttr := resourceBlock.GetAttribute("action"); actionAttr != nil && actionAttr.Equals("deny") {
				return
			}

			if destinationIP := resourceBlock.GetAttribute("destination_ip_address"); destinationIP == nil || destinationIP.Equals("") {
				set.Add(
					result.New(resourceBlock).
						WithDescription(fmt.Sprintf("Resource '%s' defines a firewall rule with no restriction on destination IP", resourceBlock)),
				)
			} else if cidr.IsOpen(destinationIP) {
				set.Add(
					result.New(resourceBlock).
						WithDescription(fmt.Sprintf("Resource '%s' defines a firewall rule with a public destination CIDR", resourceBlock)),
				)
			}

			if sourceIP := resourceBlock.GetAttribute("source_ip_address"); sourceIP == nil || sourceIP.Equals("") {
				set.Add(
					result.New(resourceBlock).
						WithDescription(fmt.Sprintf("Resource '%s' defines a firewall rule with no restriction on source IP", resourceBlock)),
				)
			} else if cidr.IsOpen(sourceIP) {
				set.Add(
					result.New(resourceBlock).
						WithDescription(fmt.Sprintf("Resource '%s' defines a firewall rule with a public source CIDR", resourceBlock)),
				)
			}
		},
	})
}