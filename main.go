package main

import (
    "fmt"
    "os"
    "github.com/rackspace/gophercloud"
    "github.com/rackspace/gophercloud/openstack"
    "github.com/fatih/color"
)

var (
	OPENSTACK_IDENTITY_URI   = os.Getenv("OPENSTACK_IDENTITY_URI")
	OPENSTACK_USERNAME       = os.Getenv("OPENSTACK_USERNAME")
	OPENSTACK_PASSWORD       = os.Getenv("OPENSTACK_PASSWORD")
        OPENSTACK_TENANT_NAME    = os.Getenv("OPENSTACK_TENANT_NAME") 
)

func usage() {
    c := color.New(color.FgRed).Add(color.Underline)
    c.Printf("*** PLEASE SETTINGS OPENSTACK CREDENTIAL INFOMATION ***\n")
    c = color.New(color.FgYellow)
    c.Printf("[ EXAMPLE ]\n")
    c = color.New(color.FgCyan)
    c.Printf("export OPENSTACK_IDENTITY_URI=\"https://identity.tyo1.conoha.io/v2.0\"\n")
    c.Printf("export OPENSTACK_USERNAME=\"USERNAME\"\n")
    c.Printf("export OPENSTACK_PASSWORD=\"PASSWORD\"\n")
    c.Printf("export OPENSTACK_TENANT_NAME=\"TENANT_NAME\"\n")
}

func crit(msg string){
    c := color.New(color.FgRed)
    c.Printf("[CRIT] %s\n", msg)
    os.Exit(255)
}

func main() {
    // 環境変数が設定されていないならusageを出力して終了する
    if OPENSTACK_IDENTITY_URI == "" || OPENSTACK_USERNAME == "" || 
       OPENSTACK_PASSWORD == "" || OPENSTACK_TENANT_NAME == "" {
        usage()
        os.Exit(0)
    }

    // 環境変数から受け取ったOpenStackの認証情報を用いて認証を実施する
    opts := gophercloud.AuthOptions{
        IdentityEndpoint: OPENSTACK_IDENTITY_URI,
        Username:         OPENSTACK_USERNAME,
        Password:         OPENSTACK_PASSWORD,
        TenantName:       OPENSTACK_TENANT_NAME,
    }
    provider, err := openstack.AuthenticatedClient(opts)
    if err != nil{
        fmt.Printf("OpenStack authentication error : %s", err.Error()) 
    }
    fmt.Println(provider)
}


