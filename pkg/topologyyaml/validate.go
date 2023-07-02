package topologyyaml

import (
	"errors"
	"fmt"
)

func Validate(topo *Topology) error {
	errs := make([]error, 0)

	err := validateGlobal(topo.Global)
	if err != nil {
		errs = append(errs, fmt.Errorf("global: %s", err))
	}

	for i, fe := range topo.FEs {
		err := validateCommon(fe.ComponentSpec)
		if err != nil {
			errs = append(errs, fmt.Errorf("fes[%d]: %s", i, err))
		}
		err = validateFe(fe)
		if err != nil {
			errs = append(errs, fmt.Errorf("fes[%d]: %s", i, err))
		}
	}

	for i, be := range topo.BEs {
		err := validateCommon(be.ComponentSpec)
		if err != nil {
			errs = append(errs, fmt.Errorf("bes[%d]: %s", i, err))
		}
		err = validateBe(be)
		if err != nil {
			errs = append(errs, fmt.Errorf("bes[%d]: %s", i, err))
		}
	}

	// TODO: more checks
	// - all fe.http_port should be identical
	// - ensure no port conflict
	// - ensure each host has at most 1 be

	return errors.Join(errs...)
}

func validateGlobal(global GlobalSpec) error {
	errs := make([]error, 0)

	if len(global.DeployUser) == 0 {
		errs = append(errs, errors.New("deploy_user is required"))
	}
	if len(global.SSHPrivateKeyPath) == 0 {
		errs = append(errs, errors.New("ssh_private_key_path is required"))
	}

	return errors.Join(errs...)
}

func validateCommon(comp ComponentSpec) error {
	errs := make([]error, 0)

	if len(comp.Host) == 0 {
		errs = append(errs, errors.New("host is required"))
	}

	if comp.SSHPort <= 0 {
		errs = append(errs, errors.New("ssh_port must be positive"))
	}

	if len(comp.DeployDir) == 0 {
		errs = append(errs, errors.New("deploy_dir is required"))
	}

	return errors.Join(errs...)
}

func validateFe(fe FESpec) error {
	errs := make([]error, 0)

	if fe.FeConfig.EditLogPort <= 0 {
		errs = append(errs, errors.New("edit_log_port must be positive"))
	}

	if len(fe.PackageURL) == 0 {
		errs = append(errs, errors.New("package_url is required"))
	}
	if len(fe.Digest) == 0 {
		errs = append(errs, errors.New("digest is required"))
	}

	if fe.InstallJava {
		if len(fe.JavaPackageURL) == 0 {
			errs = append(errs, errors.New("java_package_url is required"))
		}
		if len(fe.JavaDigest) == 0 {
			errs = append(errs, errors.New("java_digest is required"))
		}
	}

	return errors.Join(errs...)
}

func validateBe(be BESpec) error {
	errs := make([]error, 0)

	if len(be.BeConfig.StorageRootPath) == 0 {
		errs = append(errs, errors.New("storage_root_path is required"))
	}

	if len(be.PackageURL) == 0 {
		errs = append(errs, errors.New("package_url is required"))
	}
	if len(be.Digest) == 0 {
		errs = append(errs, errors.New("digest is required"))
	}

	if be.InstallJava {
		if len(be.JavaPackageURL) == 0 {
			errs = append(errs, errors.New("java_package_url is required"))
		}
		if len(be.JavaDigest) == 0 {
			errs = append(errs, errors.New("java_digest is required"))
		}
	}

	return errors.Join(errs...)
}
