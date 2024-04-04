package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	conn, err := net.ListenPacket("udp4", ":68")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(5 * time.Second))
	packet := buildDiscoverPacket()

	_, err = conn.WriteTo(packet, &net.UDPAddr{IP: net.IPv4bcast, Port: 67})
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1500)
	_, _, err = conn.ReadFrom(buffer)
	if err != nil {
		log.Fatal(err)
	}

	dhcpOfferOptions := parseOfferPacket(buffer)

	// Extract IP address, subnet mask, and default gateway from DHCP Offer
	ipAddress := dhcpOfferOptions[50]
	subnetMask := dhcpOfferOptions[1]
	defaultGateway := dhcpOfferOptions[3]

	fmt.Println("IP Address:", net.IP(ipAddress))
	fmt.Println("Subnet Mask:", net.IP(subnetMask))
	fmt.Println("Default Gateway:", net.IP(defaultGateway))
}

func buildDiscoverPacket() []byte {
	packet := make([]byte, 300)
	packet[0] = 0x01 // Message Type: DHCP Discover
	packet[1] = 0x01 // Hardware Type: Ethernet
	packet[2] = 0x06 // Hardware Address Length: 6
	packet[3] = 0x00 // Hops
	// Transaction ID
	binary.BigEndian.PutUint32(packet[4:8], uint32(time.Now().Unix()))
	// Flags
	binary.BigEndian.PutUint16(packet[10:12], 0x8000) // Broadcast flag

	// Client Hardware Address
	mac, err := net.ParseMAC("00:00:00:00:00:00")
	if err != nil {
		fmt.Println("Error parsing MAC:", err)
		return nil
	}
	copy(packet[28:34], mac)

	// Magic Cookie
	copy(packet[236:240], []byte{0x63, 0x82, 0x53, 0x63})

	// DHCP Option: End
	packet[299] = 0xFF

	return packet
}

func parseOfferPacket(packet []byte) map[int][]byte {
	options := make(map[int][]byte)

	// Skip DHCP header
	optionsStart := 240

	// Parse options
	for i := optionsStart; i < len(packet); i++ {
		if packet[i] == 0xFF { // End of options
			break
		}
		optionCode := int(packet[i])
		i++
		optionLength := int(packet[i])
		i++
		optionValue := make([]byte, optionLength)
		copy(optionValue, packet[i:i+optionLength])
		options[optionCode] = optionValue
		i += optionLength - 1
	}

	return options
}
