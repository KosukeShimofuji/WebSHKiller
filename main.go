package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/kr/pretty"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	_ "github.com/rackspace/gophercloud/openstack/compute/v2/servers"
	_ "github.com/rackspace/gophercloud/openstack/imageservice/v2/images"
	"github.com/rackspace/gophercloud/pagination"
	"os"
	_ "reflect"
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

func view_index_flavor(client *gophercloud.ServiceClient) {
	pager := flavors.ListDetail(client, nil)
	pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, err := flavors.ExtractFlavors(page)
		if err != nil {
			fmt.Println(err)
			return false, err
		}
		for _, f := range flavorList {
			fmt.Printf("%# v\n", pretty.Formatter(f))
		}
		return true, nil
	})
}

func main() {
	// 環境変数が設定されていないならusageを出力して終了する
	debug("Check identity information")
	if OPENSTACK_IDENTITY_URI == "" ||
		OPENSTACK_USERNAME == "" ||
		OPENSTACK_PASSWORD == "" ||
		OPENSTACK_TENANT_NAME == "" ||
		OPENSTACK_REGION == "" {
		usage()
		os.Exit(0)
	}

	// 環境変数から受け取ったOpenStackの認証情報を用いて認証を実施する
	debug("Authentication to OpenStack")
	openstack_opts := gophercloud.AuthOptions{
		IdentityEndpoint: OPENSTACK_IDENTITY_URI,
		Username:         OPENSTACK_USERNAME,
		Password:         OPENSTACK_PASSWORD,
		TenantName:       OPENSTACK_TENANT_NAME,
	}
	provider, err := openstack.AuthenticatedClient(openstack_opts)
	if err != nil {
		fmt.Printf("OpenStack authentication error : %s", err.Error())
	}

	// compute node clientを作成する
	debug("Carete compute node client")
	compute_client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: OPENSTACK_REGION,
	})
	if err != nil {
		fmt.Printf("Create OpenStack instance error : %s", err.Error())
	}

	// flavorの一覧表示を行う
	debug("Get a list of flavor")
	view_index_flavor(compute_client)
}
