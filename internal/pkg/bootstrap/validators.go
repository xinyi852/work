package bootstrap

import "plesk/internal/api/validators"

func SetupValidators() {
	validators.Init()
}
