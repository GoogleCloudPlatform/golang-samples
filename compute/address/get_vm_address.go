// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package snippets

// [START compute_ip_address_get_vm_address]
import (
	"fmt"
	"io"
	"strings"

	"cloud.google.com/go/compute/apiv1/computepb"
)

// getInstanceIPAddresses retrieves the specified type of IP address (IPv4 or IPv6, internal or external) of a specified Compute Engine instance.
func getInstanceIPAddresses(w io.Writer, instance *computepb.Instance, addressType computepb.Address_AddressType, isIPV6 bool) []string {
	//instance := &computepb.Instance{  // or query compute.NewImagesRESTClient actual instance
	//	NetworkInterfaces: []*computepb.NetworkInterface{
	//		{
	//			NetworkIP: proto.String("10.128.0.1"),
	//			AccessConfigs: []*computepb.AccessConfig{
	//				{
	//					Type:  proto.String(computepb.AccessConfig_ONE_TO_ONE_NAT.String()),
	//					NatIP: proto.String("34.68.123.45"),
	//				},
	//				{
	//					Type:  proto.String(computepb.AccessConfig_ONE_TO_ONE_NAT.String()),
	//					NatIP: proto.String("34.68.123.46"),
	//				},
	//			},
	//			Ipv6AccessConfigs: []*computepb.AccessConfig{
	//				{
	//					Type:         proto.String(computepb.AccessConfig_DIRECT_IPV6.String()),
	//					ExternalIpv6: proto.String("2600:1901:0:1234::"),
	//				},
	//			},
	//			Ipv6Address: proto.String("2600:1901:0:5678::"),
	//		},
	//	},
	//}
	// addressType := computepb.Address_EXTERNAL
	// isIPV6 := false

	var ips []string

	if instance.GetNetworkInterfaces() == nil {
		return ips
	}

	for _, iface := range instance.GetNetworkInterfaces() {
		if isIPV6 {
			// Handle IPv6 addresses
			if addressType == computepb.Address_EXTERNAL {
				if ipv6Configs := iface.GetIpv6AccessConfigs(); ipv6Configs != nil {
					for _, ipv6Config := range ipv6Configs {
						if ipv6Config.GetType() == "DIRECT_IPV6" {
							ips = append(ips, ipv6Config.GetExternalIpv6())
						}
					}
				}
			} else if addressType == computepb.Address_INTERNAL {
				if internalIPv6 := iface.GetIpv6Address(); internalIPv6 != "" {
					ips = append(ips, internalIPv6)
				}
			}
		} else {
			// Handle IPv4 addresses
			if addressType == computepb.Address_EXTERNAL {
				for _, config := range iface.GetAccessConfigs() {
					if config.GetType() == "ONE_TO_ONE_NAT" {
						ips = append(ips, config.GetNatIP())
					}
				}
			} else if addressType == computepb.Address_INTERNAL {
				ips = append(ips, iface.GetNetworkIP())
			}
		}
	}

	fmt.Fprintf(w, "Received list of IPS: [%s]", strings.Join(ips, ", "))

	return ips
}

// [END compute_ip_address_get_vm_address]
