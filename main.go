package main

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		publicKeyPath := cfg.Require("publicKeyPath")
		privateKeyPath := cfg.Require("privateKeyPath")

		publicKeyBytes, err := os.ReadFile(publicKeyPath)
		if err != nil {
			return err
		}
		publicKey := pulumi.String(string(publicKeyBytes))

		privateKeyBytes, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return err
		}
		privateKey := pulumi.ToSecret(string(privateKeyBytes)).(pulumi.StringOutput)

		dropKeypair, err := digitalocean.NewSshKey(ctx, "default", &digitalocean.SshKeyArgs{
			Name:      pulumi.String("Example"),
			PublicKey: pulumi.String(publicKey),
		})
		if err != nil {
			return err
		}

		dropVPC, err := digitalocean.NewVpc(ctx, "vpc-new", &digitalocean.VpcArgs{
			Name:    pulumi.String("pulumi-vpc"),
			Region:  pulumi.String("sgp1"),
			IpRange: pulumi.String("10.13.2.0/24"),
		})
		if err != nil {
			return err
		}

		droplet, err := digitalocean.NewDroplet(ctx, "pulumi-start", &digitalocean.DropletArgs{
			Image:   pulumi.String("centos-stream-9-x64"),
			Name:    pulumi.String("pulumi-start"),
			Region:  dropVPC.Region,
			VpcUuid: dropVPC.ID(),
			Size:    pulumi.String(digitalocean.DropletSlugDropletS2VCPU2GB),
			SshKeys: pulumi.StringArray{
				dropKeypair.Fingerprint,
			},
		})
		if err != nil {
			return err
		}

		droplet.Name.ApplyT(func(name interface{}) error {
			dropletName := name.(string)
			dropMeta, err := digitalocean.LookupDroplet(ctx, &digitalocean.LookupDropletArgs{
				Name: pulumi.StringRef(dropletName),
			}, nil)
			if err != nil {
				return err
			}

			renderConfig, err := local.NewCommand(ctx, "renderConfig", &local.CommandArgs{
				Create: pulumi.String("cat ./template/server.tftpl | envsubst > server.yml"),
				Environment: pulumi.StringMap{
					"IP_ADDRESS": pulumi.String(dropMeta.Ipv4Address),
				},
			})
			if err != nil {
				return err
			}

			updatePythonCmd, err := remote.NewCommand(ctx, "updatePythonCmd", &remote.CommandArgs{
				Connection: &remote.ConnectionArgs{
					Host:       pulumi.String(dropMeta.Ipv4Address),
					Port:       pulumi.Float64(22),
					User:       pulumi.String("root"),
					PrivateKey: privateKey,
				},
				Create: pulumi.String("(sudo yum update -y || true);" +
					"(sudo yum install python3 -y)\n"),
			})
			if err != nil {
				return err
			}

			_, err = local.NewCommand(ctx, "playAnsiblePlaybookCmd", &local.CommandArgs{
				Create: pulumi.String(fmt.Sprintf(
					"ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook "+
						"-u root "+
						"-i '%v,' "+
						"--private-key %v "+
						"ansible_playbook.yaml",
					dropMeta.Ipv4Address, privateKeyPath,
				)),
			}, pulumi.DependsOn([]pulumi.Resource{
				renderConfig,
				updatePythonCmd,
			}))

			ctx.Export("dropletIP", pulumi.String(dropMeta.Ipv4Address))
			return nil
		})

		return nil
	})
}
