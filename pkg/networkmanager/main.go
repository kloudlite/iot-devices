package networkmanager

// import (
// 	"errors"
// 	"fmt"
// 	"net"
// 	"os"
//
// 	"github.com/godbus/dbus/v5"
// 	"github.com/vishvananda/netlink"
// )
//
// // NetworkInterface manages network interfaces
// type NetworkInterface interface {
// 	CreateInterface(name string) error
// 	SetDNS(interfaceName string, dnsServers []string) error
// 	GetDNS(interfaceName string) ([]string, error)
// 	AddOrUpdateIP(interfaceName string, ip netlink.Addr) error
// 	InterfaceExists(interfaceName string) (bool, error)
// 	DeleteInterface(interfaceName string) error
//
// 	SetGateway(interfaceName string, gatewayIP string) error
// 	AddPolicyBasedRoutingNoIP(interfaceName string, gatewayIP string, tableID int) error
//
// 	BringUpInterface(interfaceName string) error
// }
//
// // defaultManager is a default implementation of the NetworkInterface
// type defaultManager struct{}
//
// // CreateInterface creates a new network interface
// func (dm *defaultManager) CreateInterface(name string) error {
//
// 	if b, _ := dm.InterfaceExists(name); b {
// 		if err := dm.DeleteInterface(name); err != nil {
// 			return err
// 		}
// 	}
//
// 	link := &netlink.Dummy{
// 		LinkAttrs: netlink.LinkAttrs{Name: name},
// 	}
// 	return netlink.LinkAdd(link)
// }
//
// // SetDNS sets the DNS for a network interface using systemd-resolved
// func (dm *defaultManager) SetDNS(interfaceName string, dnsServers []string) error {
// 	conn, err := dbus.SystemBus()
// 	if err != nil {
// 		return err
// 	}
// 	resolver := conn.Object("org.freedesktop.resolve1", "/org/freedesktop/resolve1")
// 	call := resolver.Call("org.freedesktop.resolve1.Manager.SetLinkDNS", 0, getLinkIndex(interfaceName), createDNSArray(dnsServers))
// 	return call.Err
// }
//
// // GetDNS retrieves the DNS settings for a network interface
// func (dm *defaultManager) GetDNS(interfaceName string) ([]string, error) {
// 	conn, err := dbus.SystemBus()
// 	if err != nil {
// 		return nil, err
// 	}
// 	resolver := conn.Object("org.freedesktop.resolve1", dbus.ObjectPath("/org/freedesktop/resolve1"))
// 	var linkData []map[string]dbus.Variant
// 	err = resolver.Call("org.freedesktop.DBus.Properties.GetAll", 0, "org.freedesktop.resolve1.Link").Store(&linkData)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return parseDNSData(linkData, interfaceName), nil
// }
//
// // AddOrUpdateIP adds or updates an IP address to a network interface
// func (dm *defaultManager) AddOrUpdateIP(interfaceName string, ip netlink.Addr) error {
// 	link, err := netlink.LinkByName(interfaceName)
// 	if err != nil {
// 		return err
// 	}
// 	return netlink.AddrReplace(link, &ip)
// }
//
// // Helper functions to interact with systemd-resolved via DBus
// func getLinkIndex(interfaceName string) int32 {
// 	link, err := netlink.LinkByName(interfaceName)
// 	if err != nil {
// 		return -1
// 	}
// 	return int32(link.Attrs().Index)
//
// }
//
// func createDNSArray(dnsServers []string) [][]byte {
// 	var dnsArray [][]byte
// 	for _, dns := range dnsServers {
// 		ip := parseIP(dns)
// 		if ip != nil {
// 			dnsArray = append(dnsArray, ip)
// 		}
// 	}
// 	return dnsArray
// }
//
// func parseIP(ipStr string) []byte {
// 	ip := net.ParseIP(ipStr)
// 	if ip == nil {
// 		return nil
// 	}
// 	if ip4 := ip.To4(); ip4 != nil {
// 		return ip4
// 	}
// 	return ip.To16()
// }
//
// func parseDNSData(linkData []map[string]dbus.Variant, interfaceName string) []string {
// 	for _, data := range linkData {
// 		if data["Name"].Value().(string) == interfaceName {
// 			var dnsServers []string
// 			dnsRecords := data["DNS"].Value().([][]byte)
// 			for _, record := range dnsRecords {
// 				dnsServers = append(dnsServers, net.IP(record).String())
// 			}
// 			return dnsServers
// 		}
// 	}
// 	return nil
// }
//
// // InterfaceExists checks if a network interface already exists
// func (dm *defaultManager) InterfaceExists(interfaceName string) (bool, error) {
// 	_, err := netlink.LinkByName(interfaceName)
// 	if err != nil {
// 		if errors.Is(err, netlink.LinkNotFoundError{}) {
// 			return false, nil
// 		}
// 		return false, err
// 	}
// 	return true, nil
// }
//
// // DeleteInterface deletes a network interface
// func (dm *defaultManager) DeleteInterface(interfaceName string) error {
// 	link, err := netlink.LinkByName(interfaceName)
// 	if err != nil {
// 		return err
// 	}
// 	return netlink.LinkDel(link)
// }
//
// // BringUpInterface brings up a network interface
// func (dm *defaultManager) BringUpInterface(interfaceName string) error {
// 	link, err := netlink.LinkByName(interfaceName)
// 	if err != nil {
// 		return fmt.Errorf("failed to find interface '%s': %v", interfaceName, err)
// 	}
//
// 	// Bring the interface up
// 	if err := netlink.LinkSetUp(link); err != nil {
// 		return fmt.Errorf("failed to bring up interface '%s': %v", interfaceName, err)
// 	}
//
// 	return nil
// }
//
// // SetGateway sets the default gateway for a network interface
// func (dm *defaultManager) SetGateway(interfaceName string, gatewayIP string) error {
// 	// if err := dm.AddPolicyBasedRoutingNoIP(interfaceName, gatewayIP, 200); err != nil {
// 	// 	return err
// 	// }
//
// 	link, err := netlink.LinkByName(interfaceName)
// 	if err != nil {
// 		return err
// 	}
//
// 	// Parse the gateway IP
// 	gateway := net.ParseIP(gatewayIP)
// 	if gateway == nil {
// 		return fmt.Errorf("invalid gateway IP address format")
// 	}
//
// 	// Define the route
// 	route := &netlink.Route{
// 		LinkIndex: link.Attrs().Index,
// 		Scope:     netlink.SCOPE_UNIVERSE,
// 		// Gw:        net.ParseIP("192.168.0.105/24"),
// 		Dst: &net.IPNet{IP: gateway, Mask: net.CIDRMask(32, 32)},
// 	}
//
// 	// Add the route
// 	if err := netlink.RouteAdd(route); err != nil {
// 		fmt.Println("Error adding route:", err)
// 		if !os.IsExist(err) {
// 			return err
// 		}
// 	}
//
// 	return nil
// }
//
// func NewManager() NetworkInterface {
// 	return &defaultManager{}
// }
//
// // func example() {
// // 	manager := NewManager()
// //
// // 	// Create interface
// // 	if err := manager.CreateInterface("myinterface"); err != nil {
// // 		fmt.Println("Error creating interface:", err)
// // 		return
// // 	}
// //
// // 	// Set DNS
// // 	dnsServers := []string{"8.8.8.8", "8.8.4.4"}
// // 	if err := manager.SetDNS("myinterface", dnsServers); err != nil {
// // 		fmt.Println("Error setting DNS:", err)
// // 		return
// // 	}
// //
// // 	// Get DNS
// // 	currentDNS, err := manager.GetDNS("myinterface")
// // 	if err != nil {
// // 		fmt.Println("Error getting DNS:", err)
// // 		return
// // 	}
// // 	fmt.Println("Current DNS Servers:", currentDNS)
// //
// // 	// Add or Update IP
// // 	ip, _ := netlink.ParseAddr("192.168.1.1/24")
// // 	if err := manager.AddOrUpdateIP("myinterface", *ip); err != nil {
// // 		fmt.Println("Error adding/updating IP:", err)
// // 		return
// // 	}
// // }
//
// // AddPolicyBasedRoutingNoIP setups policy-based routing for a specific interface using its name
// func (dm *defaultManager) AddPolicyBasedRoutingNoIP(interfaceName string, gatewayIP string, tableID int) error {
// 	link, err := netlink.LinkByName(interfaceName)
// 	if err != nil {
// 		return fmt.Errorf("failed to get link by name: %v", err)
// 	}
//
// 	// Ensure the routing table is clear before setting up new routes
// 	if err := netlink.RouteDel(&netlink.Route{Table: tableID}); err != nil {
// 		return fmt.Errorf("failed to clear old routes from table %d: %v", tableID, err)
// 	}
//
// 	// Add a rule to look up the custom table for all packets from this interface
// 	rule := netlink.NewRule()
// 	rule.IifName = interfaceName
// 	rule.Table = tableID
// 	if err := netlink.RuleAdd(rule); err != nil {
// 		return fmt.Errorf("failed to add routing rule: %v", err)
// 	}
//
// 	// Add default route to the new table to route all traffic via the specified gateway
// 	gwIP := net.ParseIP(gatewayIP)
// 	if gwIP == nil {
// 		return fmt.Errorf("invalid gateway IP address format")
// 	}
// 	route := &netlink.Route{
// 		LinkIndex: link.Attrs().Index,
// 		Scope:     netlink.SCOPE_UNIVERSE,
// 		Gw:        gwIP,
// 		Table:     tableID,
// 	}
// 	if err := netlink.RouteAdd(route); err != nil {
// 		return fmt.Errorf("failed to add route to routing table: %v", err)
// 	}
//
// 	return nil
// }
//
// func Test() error {
// 	ifName := "myinterface"
// 	manager := NewManager()
//
// 	if err := manager.CreateInterface(ifName); err != nil {
// 		return err
// 	}
//
// 	if err := manager.BringUpInterface(ifName); err != nil {
// 		return err
// 	}
//
// 	if err := manager.SetGateway(ifName, "192.168.0.105"); err != nil {
// 		return err
// 	}
//
// 	// ip, _ := netlink.ParseAddr("34.93.62.238/32")
// 	// if err := manager.AddOrUpdateIP(ifName, *ip); err != nil {
// 	// 	return err
// 	// }
//
// 	return nil
// }
