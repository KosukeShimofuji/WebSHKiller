package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"fmt"
	"github.com/fatih/color"
	"github.com/kr/pretty"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	"github.com/rackspace/gophercloud/openstack/imageservice/v2/images"
	"github.com/rackspace/gophercloud/pagination"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	_ "reflect"
	"strconv"
	"strings"
)

var (
	OPENSTACK_IDENTITY_URI = os.Getenv("OPENSTACK_IDENTITY_URI")
	OPENSTACK_USERNAME     = os.Getenv("OPENSTACK_USERNAME")
	OPENSTACK_PASSWORD     = os.Getenv("OPENSTACK_PASSWORD")
	OPENSTACK_TENANT_NAME  = os.Getenv("OPENSTACK_TENANT_NAME")
	OPENSTACK_REGION       = os.Getenv("OPENSTACK_REGION")
	DEBUG                  = true
)

func usage() {
	c := color.New(color.FgRed).Add(color.Underline)
	c.Printf("*** PLEASE SETTINGS OPENSTACK CREDENTIAL INFOMATION ***\n")
	c = color.New(color.FgYellow)
	c.Printf("[ EXAMPLE ]\n")
	c = color.New(color.FgCyan)
	c.Printf("export OPENSTACK_IDENTITY_URI=\"IDENTITY_URI\"\n")
	c.Printf("export OPENSTACK_USERNAME=\"USERNAME\"\n")
	c.Printf("export OPENSTACK_PASSWORD=\"PASSWORD\"\n")
	c.Printf("export OPENSTACK_TENANT_NAME=\"TENANT_NAME\"\n")
	c.Printf("export OPENSTACK_REGION=\"REGION\"\n")
}

func crit(msg string) {
	c := color.New(color.FgRed)
	c.Printf("[CRIT] %s\n", msg)
	os.Exit(255)
}

func debug(msg string) {
	c := color.New(color.FgGreen)
	c.Printf("[DEBUG] %s\n", msg)
}

func rndstr(len int) string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, len)
}

func view_openstack_flavors(client *gophercloud.ServiceClient) {
	pager := flavors.ListDetail(client, nil)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, err := flavors.ExtractFlavors(page)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		for _, f := range flavorList {
			//fmt.Printf("%# v\n", pretty.Formatter(f))
			fmt.Printf(" * ID:%s Disk:%d RAM:%d NAME:%s RxTxFactor:%f SWAP:%d VCPUs:%d\n",
				f.ID, f.Disk, f.RAM, f.Name, f.RxTxFactor, f.Swap, f.VCPUs)
		}
		return true, nil
	})
}

func view_openstack_keypairs(client *gophercloud.ServiceClient) {
	pager := keypairs.List(client)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		KeypairList, err := keypairs.ExtractKeyPairs(page)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		for _, f := range KeypairList {
			//fmt.Printf("%# v\n", pretty.Formatter(f))
			fmt.Printf(" * NAME:%s FingerPrint:%s UserID:%s\n"+
				"   Public Key: %s\n"+
				"   Private Key: %s\n",
				f.Name, f.Fingerprint, f.UserID,
				strings.TrimRight(f.PublicKey, "\n"),
				strings.TrimRight(f.PrivateKey, "\n"),
			)
		}
		return true, nil
	})
}

func view_openstack_images(client *gophercloud.ServiceClient) {
	pager := images.List(client, nil)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		imageList, err := images.ExtractImages(page)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		for _, f := range imageList {
			//fmt.Printf("%# v\n", pretty.Formatter(f))
			fmt.Printf(" * %s(%s)\n", f.Name, f.ID)
		}
		return true, nil
	})
}

func view_openstack_servers(client *gophercloud.ServiceClient) {
	pager := servers.List(client, nil)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		for _, f := range serverList {
			fmt.Printf("%# v\n", pretty.Formatter(f))
		}
		return true, nil
	})
}

func add_openstack_keypair(client *gophercloud.ServiceClient, name string) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := privateKey.PublicKey
	pub, err := ssh.NewPublicKey(&publicKey)
	pubBytes := ssh.MarshalAuthorizedKey(pub)
	pk := string(pubBytes)

	kp, err := keypairs.Create(client, keypairs.CreateOpts{
		Name:      name,
		PublicKey: pk,
	}).Extract()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%# v\n", pretty.Formatter(kp))
}

func del_openstack_keypair(client *gophercloud.ServiceClient, name string) {
	err := keypairs.Delete(client, name).ExtractErr()

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	debug("Check identity information")
	if OPENSTACK_IDENTITY_URI == "" ||
		OPENSTACK_USERNAME == "" ||
		OPENSTACK_PASSWORD == "" ||
		OPENSTACK_TENANT_NAME == "" ||
		OPENSTACK_REGION == "" {
		usage()
		os.Exit(0)
	}

	debug("Authentication to OpenStack")
	openstack_opts := gophercloud.AuthOptions{
		IdentityEndpoint: OPENSTACK_IDENTITY_URI,
		Username:         OPENSTACK_USERNAME,
		Password:         OPENSTACK_PASSWORD,
		TenantName:       OPENSTACK_TENANT_NAME,
	}
	provider, err := openstack.AuthenticatedClient(openstack_opts)
	if err != nil {
		log.Fatal(err)
	}

	debug("Carete compute node client")
	compute_client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: OPENSTACK_REGION,
	})
	if err != nil {
		log.Fatal(err)
	}

	debug("Get a list of flavors")
	view_openstack_flavors(compute_client)

	debug("Get a list of images")
	view_openstack_images(compute_client)

	debug("Get a list of keypairs")
	view_openstack_keypairs(compute_client)

	debug("Get a list of servers")
	view_openstack_servers(compute_client)

	debug("Add a keypair")
	add_openstack_keypair(compute_client, "some_name")

	debug("Del a keypair")
	del_openstack_keypair(compute_client, "some_name")
}
