//go:build cloud
// +build cloud

package gcore

import (
	"fmt"
	"testing"

	"github.com/G-Core/gcorelabscloud-go/gcore/secret/v1/secrets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	secretName       = "test-secret"
	privateKey       = "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDQ4E6U0vql4EST\n8o41TlHRz6MKmMhddVUjM2juTKjxv4WuB4T3z/wokznEjQg4H7gfYEKeCJqelrfq\ntdOtbPsznSceMOXB5uA2Sc9WVKwk7owoRJxPd4LQeOcarVOFdIzudzkgSK/oV7Za\nL8Y2hylsB4SX2cfbULtmW/WDePp3YZAL6zYV1fXJSnK+hL2iUSqikiViEGRta+47\nnaTKZnnmSgojdshzsw0wlF/PgRJ/Anf9j9J8ratdJP81yAG5daU3L2NdJ3qx9UbV\ntKnSq2z2u4yx6xdb4t4WFQBKNjC6+YZN/gI5lp96p3FNTNS4PKYxAAUrnCwf0EE3\n7dOR4eWlAgMBAAECggEBALPm3ge0h4li1e4PVYh4AmSRT74KxVgpfMCqwM+uWzyM\nVpkDhPTjwC06UOEHD3M3bqAninkOtA2vhoyzOrP+T4Wu70hDmUAemDJp9BhJKVNN\n2o28Olz/dD4WRAZoDq29Kr0hFqTFtiyJj1eyGihQ1c5j00HuowI0UJPi1Fz+T8uN\nPwukUtTPYwEds6SApii3v9VKjmvbRDmsbHU3KkUoaeqpRnRagyp1vtoLXigezUcK\nrQcoh6wlKtvj0YLR2lxq9Wmj1nn6m3F5Bom54X8o18tcOmFSRudRb+Fxjb0jnqSK\nAsyVlZg4alTBQUmx9gIKv0oSJAIh2nXdclECkGjs8WkCgYEA9xvdDWephsbv+X3k\nndnDG9JTxfrR6HMHPrUrTaZ8/VD+Qw4zuReoNGkcQbV3Cb26egprWQWfYc9+l6mU\nAWgOjFgeGie1uwOwkhv6CfhE/iVvotJ3hOOsC5pLEhz4vRpO75C9wSehjfTYkP1m\nXEAhRTRbgMnvzChWyh5CEjosX5sCgYEA2GRHrG0JVxsYSCugLPKf9fSK4CQDm0bK\nywBwZtAWX0xhiHO/BW6PeK1Mqx2nbiWl1hXNpZKJNS9bnrZWym/yUqOvg2XJKjb6\nhHBvwAD1MOQ8Ysby4JHGCrMBEwlcDpI2wpMpXkKhU3X0XWjkqrhqCH/TETFKkqLt\nfJX/c9PTQ78CgYAEPek0grQJST7zVHLpNsS/pIOloWGbEOZt8CQ3KAV7P7mtov/G\nTJ6pj6hZhGjvtN8Pm0Aufgc3YZ11swaEY6nkRNr3bfkTpcORLoPDSgy9JB1feSdu\nE45vgI2LWQ34CQyT1jM7rpd6XVqeWos4SC2KB5UOh+ji40piG9TchT0fwwKBgA/M\nmpMTTvhGKSqzzLkbaeR6W11sI7tFmu7hdFN9Y/THTeO5l7vcy6ri9FMWEjBvnUEZ\nTG+HWG9CquzWoVWcgNPZ0anFV7+2Teo3j2E0cLKGJ4aKwhb1bcFAOpbaOxdxQ4BH\nYGDaeo7ucM4VJ4TzfAJs2stJjwlPzgknpoQddjJfAoGBAIFfnU8x/SrNhAqZrG9d\n3kpJ5LmbVswOYtj01KHM+KpEwOQVF+s2NOeHqyC7QUIWrue00+1MT88F9cNHDeWk\n0dEOJNWCfzcV85l8A+0p6/4qAW7h7RNiFqeA8GyVKCT8f7fu/7WpYw8D0aq8w5X/\nKZl+AjB+MzYFs71+SC4ohTlI\n-----END PRIVATE KEY-----"
	certificate      = "-----BEGIN CERTIFICATE-----\nMIIDpDCCAoygAwIBAgIJAIUvym0uaBHbMA0GCSqGSIb3DQEBCwUAMD0xCzAJBgNV\nBAYTAlJVMQ8wDQYDVQQIDAZNT1NDT1cxCzAJBgNVBAoMAkNBMRAwDgYDVQQDDAdS\nT09UIENBMB4XDTIxMDczMDE1MTU0NVoXDTMxMDcyODE1MTU0NVowTDELMAkGA1UE\nBhMCQ0ExDTALBgNVBAgMBE5vbmUxCzAJBgNVBAcMAk5CMQ0wCwYDVQQKDAROb25l\nMRIwEAYDVQQDDAlsb2NhbGhvc3QwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK\nAoIBAQDQ4E6U0vql4EST8o41TlHRz6MKmMhddVUjM2juTKjxv4WuB4T3z/wokznE\njQg4H7gfYEKeCJqelrfqtdOtbPsznSceMOXB5uA2Sc9WVKwk7owoRJxPd4LQeOca\nrVOFdIzudzkgSK/oV7ZaL8Y2hylsB4SX2cfbULtmW/WDePp3YZAL6zYV1fXJSnK+\nhL2iUSqikiViEGRta+47naTKZnnmSgojdshzsw0wlF/PgRJ/Anf9j9J8ratdJP81\nyAG5daU3L2NdJ3qx9UbVtKnSq2z2u4yx6xdb4t4WFQBKNjC6+YZN/gI5lp96p3FN\nTNS4PKYxAAUrnCwf0EE37dOR4eWlAgMBAAGjgZcwgZQwVwYDVR0jBFAwTqFBpD8w\nPTELMAkGA1UEBhMCUlUxDzANBgNVBAgMBk1PU0NPVzELMAkGA1UECgwCQ0ExEDAO\nBgNVBAMMB1JPT1QgQ0GCCQCectJTETy4lTAJBgNVHRMEAjAAMAsGA1UdDwQEAwIE\n8DAhBgNVHREEGjAYgglsb2NhbGhvc3SCCyoubG9jYWxob3N0MA0GCSqGSIb3DQEB\nCwUAA4IBAQBqzJcwygLsVCTPlReUpcKVn84aFqzfZA0m7hYvH+7PDH/FM8SbX3zg\nteBL/PgQAZw1amO8xjeMc2Pe2kvi9VrpfTeGqNia/9axhGu3q/NEP0tyDFXAE2bR\njBdGhd5gCmg+X4WdHigCgn51cz5r2k3fSOIWP+TQWHqc8Yt+vZXnkwnQkRA1Ki7N\nWOiJjj/ae5RWwma/kJNmShTZn754gbQn06bAjNbPjclsHRLkawmLqikd1rYUhIdk\nOr1Nrl+CWMx3CXg0TVVdJ6rH3dO31uyvb+3qEY7WnL+HhZyr08ay8gJsEKPuPFA2\nxvveXqt9ceU5qh+8T7mHwGALEUw96QcP\n-----END CERTIFICATE-----"
	certificateChain = "-----BEGIN CERTIFICATE-----\nMIIC9jCCAd4CCQCectJTETy4lTANBgkqhkiG9w0BAQsFADA9MQswCQYDVQQGEwJS\nVTEPMA0GA1UECAwGTU9TQ09XMQswCQYDVQQKDAJDQTEQMA4GA1UEAwwHUk9PVCBD\nQTAeFw0yMTA3MzAxNTExMzVaFw0yNDA1MTkxNTExMzVaMD0xCzAJBgNVBAYTAlJV\nMQ8wDQYDVQQIDAZNT1NDT1cxCzAJBgNVBAoMAkNBMRAwDgYDVQQDDAdST09UIENB\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAo6tZ0NV6QIR/mvsqtAII\nzTTuBMrZR5OTwKvcGnhe4GVDwzJ/OgEWkghLAzOojcJvkfzJOtWwOXqwgphksc+7\n+vwIPTPt3iWjbQUzXK8pFLkjxrO8px/QxPuUrp+U6DTVvvgQesjMZ9jQRUFKOiCc\nu0st1N5Q/CJR4VOJxtYoLy1ZUlsABhwJ+6trkoOFTLRPlMUX1EIG57jYAotHvQFo\nc8UNx3KzvJsJJ56SniXCIkeu61IOt8aOXHU+3TLYhZnPiP311cMbXA0J3vGPRZwz\n25BZjF3IF/ShXlfzz76FjWUTAThc0+HA8lzx53xD4/n8HN+sGubGx9TvLyZimG/U\nGwIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQAnK8Wzw33fR6R6pqV05XI9Yu8J+BwC\nCn2bKxxYwwQWZyX1as+UIlGuvyBRJba9W2UGMj95FQfWVdDyFC98spUur+O/5yL+\nNHH+dxGnkxIRc6RMIy+GXJwPrLiB/t70hSvwgVa249zNJVcwYN/5SGX5wLaJKnim\neY99xm75nr03O/RJK/DR8HvWysH7zxvrMWs0ppfwxkxrwOcg0Cb9xODVkg/wyClw\nLiHWlmH/eyC8nkiLYJKmV7566VWCV+gy+hC/DRstVVjIMG6LsqaPq6ycm7N8EV8s\nBb5uXIVHW6w5a20c40+W9G4EDYiQjdgEaf0FoMAWGDnOEaPsvjQk2/z5\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIDPDCCAiQCCQDxA75ydLHVoTANBgkqhkiG9w0BAQsFADBgMQswCQYDVQQGEwJS\nVTEPMA0GA1UECAwGTU9TQ09XMQ8wDQYDVQQHDAZNT1NDT1cxFTATBgNVBAoMDElO\nVEVSTUVESUFURTEYMBYGA1UEAwwPSU5URVJNRURJQVRFIENBMB4XDTIxMDczMDE1\nMTIyMloXDTI0MDUxOTE1MTIyMlowYDELMAkGA1UEBhMCUlUxDzANBgNVBAgMBk1P\nU0NPVzEPMA0GA1UEBwwGTU9TQ09XMRUwEwYDVQQKDAxJTlRFUk1FRElBVEUxGDAW\nBgNVBAMMD0lOVEVSTUVESUFURSBDQTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC\nAQoCggEBAKOrWdDVekCEf5r7KrQCCM007gTK2UeTk8Cr3Bp4XuBlQ8MyfzoBFpII\nSwMzqI3Cb5H8yTrVsDl6sIKYZLHPu/r8CD0z7d4lo20FM1yvKRS5I8azvKcf0MT7\nlK6flOg01b74EHrIzGfY0EVBSjognLtLLdTeUPwiUeFTicbWKC8tWVJbAAYcCfur\na5KDhUy0T5TFF9RCBue42AKLR70BaHPFDcdys7ybCSeekp4lwiJHrutSDrfGjlx1\nPt0y2IWZz4j99dXDG1wNCd7xj0WcM9uQWYxdyBf0oV5X88++hY1lEwE4XNPhwPJc\n8ed8Q+P5/BzfrBrmxsfU7y8mYphv1BsCAwEAATANBgkqhkiG9w0BAQsFAAOCAQEA\ngOHvrh66+bQoG3Lo8bfp7D1Xvm/Md3gJq2nMotl2BH1TvNzMV93fCXygRX8J8rTL\n7xjUC2SbOrFDWFq2hNJQagdecAeuG+U55BY6Wi8SsHw+fhgxQyl9wtXWwotQPmsD\nuRhR1rL3vEphgPLbxNBzA7Lvj+P89Ar988Qy+o5AiUzHMUuqZbGOqs8UcKCQP7e/\nIX+zqqFwqyI8f90SVySGgs574jo8jQFy3l5fnp6yK0MPWg2cBCjpa5H1A+5DADF+\nnryV6Ie/m/wfxmitZZN+YCJu+8Bmmdl/FCwbmiH+HCLhrO8gonH3K21cQujMyFF5\nc7OFj86hvhqbr4kzz1J8lg==\n-----END CERTIFICATE-----"
)

func TestAccSecret(t *testing.T) {
	fullName := "gcore_secret.acctest"
	kpTemplate := fmt.Sprintf(`
	resource "gcore_secret" "acctest" {
	  %s
      %s
      name = "%s"
      private_key = %q
      certificate = %q
      certificate_chain = %q
      expiration = "2025-12-28T19:14:44.213"
	}
	`, projectInfo(), regionInfo(), secretName, privateKey, certificate, certificateChain)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccSecretDestroy,
		Steps: []resource.TestStep{
			{
				Config: kpTemplate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExists(fullName),
					resource.TestCheckResourceAttr(fullName, "name", secretName),
				),
			},
		},
	})
}

func testAccSecretDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := CreateTestClient(config.Provider, secretPoint, versionPointV1)
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gcore_secret" {
			continue
		}

		_, err := secrets.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("secret still exists")
		}
	}

	return nil
}
