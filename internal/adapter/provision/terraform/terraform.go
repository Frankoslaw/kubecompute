package provision_terraform

import (
	"context"
	"kubecompute/internal/core/domain"
	"kubecompute/internal/core/util"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
)

func ensureTfexecWorkdir(name domain.NamespacedName) (string, error) {
	kubecomputeDir := util.GetOsWorkDir()
	tfDir := filepath.Join(kubecomputeDir, "terraform", name.String())
	if err := os.MkdirAll(tfDir, 0o755); err != nil {
		return "", err
	}
	return tfDir, nil
}

func getTfexec(dir string) (*tfexec.Terraform, error) {
	tfPath, err := exec.LookPath("terraform")
	if err != nil {
		return nil, err
	}

	tf, err := tfexec.NewTerraform(dir, tfPath)
	if err != nil {
		return nil, err
	}
	return tf, nil
}

func initTerraformProject(tf *tfexec.Terraform, prog string) error {
	progFile := filepath.Join(tf.WorkingDir(), "main.tf")
	if err := os.WriteFile(progFile, []byte(prog), 0o644); err != nil {
		return err
	}
	return tf.Init(context.Background(), tfexec.Upgrade(true))
}
