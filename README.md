## dieichsipi - an attempt to manually get a DHCP lease from a DHCP server

I still need to figure out why I get the following error:
```bash
read udp4 0.0.0.0:68: i/o timeout
```

I built this because I need to manually assign IPs on static devices that are not connected to the internet. I have a DHCP server running on my local network and I want to get a lease from it by sending a DHCPDISCOVER packet and receiving a DHCPOFFER packet.

-

### Challenge
I don't want to use any external packages to do this. I want to build this from scratch using only the standard library.
