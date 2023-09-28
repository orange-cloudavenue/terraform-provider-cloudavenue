package client

import (
	"fmt"

	netbackupclient "github.com/orange-cloudavenue/netbackup-sdk-go"
)

// Netbackup is the main struct for the NetBackup client.
type NetBackup struct {
	Client   *netbackupclient.Client
	URL      string
	User     string
	Password string
}

// NewNetBackup creates a new NetBackup client.
func (c *CloudAvenue) NewNetBackupClient() (err error) {
	c.NetBackup.Client, err = netbackupclient.New(netbackupclient.Opts{
		APIEndpoint: c.URL,
		Username:    c.User,
		Password:    c.Password,
		Debug:       false,
	})
	if err != nil {
		return fmt.Errorf("%w : %w", ErrConfigureNetBackup, err)
	}
	return
}

// IsDefined checks if the NetBackup configuration is defined.
func (nB *NetBackup) IsDefined() bool {
	if nB.URL != "" || nB.User != "" || nB.Password != "" {
		return true
	}
	return false
}
