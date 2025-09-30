package provision_terraform

import (
	"bytes"
	"context"
	"kubecompute/internal/core/domain"
	"kubecompute/internal/core/port"
	"kubecompute/internal/core/util"
	"path/filepath"
	"runtime"
	"text/template"

	"github.com/hashicorp/terraform-exec/tfexec"
)

var (
	forceDryRun = false
)

type TerraformNodeProvider struct {
	tfDir string
}

func NewVirtualboxNodeProvider() port.NodeProvider {
	kubecomputeDir := util.GetOsWorkDir()
	tfDir := filepath.Join(kubecomputeDir, "terraform")

	return &TerraformNodeProvider{tfDir: tfDir}
}

func (d *TerraformNodeProvider) EnsureNode(ctx context.Context, desired *domain.Node) (bool, error) {
	dir, err := ensureTfexecWorkdir(desired.ObjectMeta.NamespacedName())
	if err != nil {
		return false, err
	}

	tf, err := getTfexec(dir)
	if err != nil {
		return false, err
	}

	prog, err := d.createTerraformProgram(desired)
	if err != nil {
		return false, err
	}
	initTerraformProject(tf, prog)

	drifted, err := d.terraformApply(tf, forceDryRun)
	if err != nil {
		return false, err
	}

	return drifted, nil
}

func (d *TerraformNodeProvider) DeleteNode(ctx context.Context, name domain.NamespacedName) (bool, error) {
	dir, err := ensureTfexecWorkdir(name)
	if err != nil {
		return false, err
	}

	tf, err := getTfexec(dir)
	if err != nil {
		return false, err
	}

	drifted, err := d.terraformDestroy(tf, forceDryRun)
	if err != nil {
		return false, err
	}

	return drifted, nil
}

func (d *TerraformNodeProvider) createTerraformProgram(node *domain.Node) (string, error) {
	_, filename, _, _ := runtime.Caller(0) // 0 = this function
	dir := filepath.Dir(filename)
	tmplPath := filepath.Join(dir, "node.tf.tmpl")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	tfNodeSpec := struct {
		Name string
	}{
		Name: node.ObjectMeta.NamespacedName().String(),
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tfNodeSpec); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (d *TerraformNodeProvider) terraformApply(tf *tfexec.Terraform, dryRun bool) (bool, error) {
	drifted, err := tf.Plan(context.Background())
	if err != nil {
		return false, err
	}
	if dryRun || !drifted {
		return false, nil
	}

	err = tf.Apply(context.Background())
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *TerraformNodeProvider) terraformDestroy(tf *tfexec.Terraform, dryRun bool) (bool, error) {
	drifted, err := tf.Plan(context.Background(), tfexec.Destroy(true))
	if err != nil {
		return false, err
	}
	if dryRun || !drifted {
		return false, nil
	}

	err = tf.Destroy(context.Background())
	if err != nil {
		return false, err
	}
	return true, nil
}
