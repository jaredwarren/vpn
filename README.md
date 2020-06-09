# VPN
OSX VPN manager / connector
basically just a wrapper for `scutil`


## Usage
1. add VPN in System Preferences
2. Connect
### connect
```go
	vpn, err := vpn.Dial("My VPN")
``` 
