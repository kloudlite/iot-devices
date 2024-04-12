package networkmanager

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/vishvananda/netlink"
)

const dnsServer = "1.1.1.1"
const marker = "# Added by kloudlite"

// AddExternalDns adds the DNS server 1.1.1.1 to /etc/resolv.conf if it's not already present
func AddExternalDns() error {
	filePath := "/etc/resolv.conf"
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Check if DNS server is already set
	if strings.Contains(string(contents), dnsServer) {
		return nil // DNS server already set
	}

	// Open file in append mode
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Add DNS server and marker
	_, err = f.WriteString("\nnameserver " + dnsServer + " " + marker + "\n")
	if err != nil {
		return err
	}

	return nil
}

// CleanExternalDns removes the DNS entry added by SetDNS from /etc/resolv.conf
func CleanExternalDns() error {
	filePath := "/etc/resolv.conf"
	tempFilePath := "/etc/resolv.conf.tmp"

	// Open the original file for reading
	originalFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer originalFile.Close()

	// Create a temporary file for writing
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(originalFile)
	writer := bufio.NewWriter(tempFile)

	// Copy lines to temp file, omitting the added DNS line
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, dnsServer) || !strings.Contains(line, marker) {
			writer.WriteString(line + "\n")
		}
	}
	writer.Flush()

	if err := scanner.Err(); err != nil {
		return err
	}

	// Replace the original file with the modified temporary file
	if err := os.Rename(tempFilePath, filePath); err != nil {
		return err
	}

	return nil
}

func getDefaultLink() (netlink.Link, error) {
	// Fetch all routes
	routes, err := netlink.RouteList(nil, netlink.FAMILY_ALL)
	if err != nil {
		return nil, fmt.Errorf("failed to list routes: %v", err)
	}
	// Find the default route (destination is nil)
	for _, route := range routes {
		if route.Dst == nil {
			// Return the link associated with the default route
			return netlink.LinkByIndex(route.LinkIndex)
		}
	}
	return nil, fmt.Errorf("default route not found")
}

func AddRoute(ip, gateway string, priority int) error {
	link, err := getDefaultLink()
	if err != nil {
		return fmt.Errorf("error fetching default interface: %v", err)
	}
	dst, err := netlink.ParseIPNet(ip)
	if err != nil {
		return fmt.Errorf("error parsing IP address: %v", err)
	}

	fmt.Println("[#] adding route", dst, "via", gateway, "on", link.Attrs().Name)
	gw := net.ParseIP(gateway)
	if gw == nil {
		return fmt.Errorf("invalid gateway IP")
	}
	route := &netlink.Route{
		LinkIndex: link.Attrs().Index,
		Dst:       dst,
		Gw:        gw,
		Priority:  priority,
	}
	if err := netlink.RouteAdd(route); err != nil {
		return fmt.Errorf("error adding route: %v", err)
	}
	return nil
}

func DeleteRoute(ip, gateway string) error {
	dst, err := netlink.ParseIPNet(ip)
	if err != nil {
		return fmt.Errorf("error parsing IP address: %v", err)
	}
	gw := net.ParseIP(gateway)
	if gw == nil {
		return fmt.Errorf("invalid gateway IP")
	}
	route := &netlink.Route{
		Dst: dst,
		Gw:  gw,
	}

	fmt.Println("[#] deleting route", dst, "via", gw)
	if err := netlink.RouteDel(route); err != nil {
		return fmt.Errorf("error deleting route: %v", err)
	}
	return nil
}

func ListRoutes(ip string) ([]netlink.Route, error) {
	dst, err := netlink.ParseIPNet(ip)
	if err != nil {
		return nil, fmt.Errorf("error parsing IP address: %v", err)
	}
	routes, err := netlink.RouteListFiltered(netlink.FAMILY_ALL, &netlink.Route{Dst: dst}, netlink.RT_FILTER_DST)
	if err != nil {
		return nil, fmt.Errorf("error listing routes: %v", err)
	}
	return routes, nil
}
