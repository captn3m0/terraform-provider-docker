package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	dc "github.com/docker/docker/client"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestAccDockerVolume_basic(t *testing.T) {
	var v types.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDockerVolumeConfig,
				Check: resource.ComposeTestCheckFunc(
					checkDockerVolume("docker_volume.foo", &v),
					resource.TestCheckResourceAttr("docker_volume.foo", "id", "testAccDockerVolume_basic"),
					resource.TestCheckResourceAttr("docker_volume.foo", "name", "testAccDockerVolume_basic"),
				),
			},
		},
	})
}

func checkDockerVolume(n string, volume *types.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		ctx := context.Background()
		client := testAccProvider.Meta().(*ProviderConfig).DockerClient
		v, err := client.VolumeInspect(ctx, n)
		// volumes, err := client.VolumeList(ctx, filters.Args{})
		if err != nil {
			if dc.IsErrNotFound(err) {
				return fmt.Errorf("Volume not found: %s", rs.Primary.ID)
			}
			return err
		}

		*volume = v

		return nil
	}
}

const testAccDockerVolumeConfig = `
resource "docker_volume" "foo" {
	name = "testAccDockerVolume_basic"
}
`
