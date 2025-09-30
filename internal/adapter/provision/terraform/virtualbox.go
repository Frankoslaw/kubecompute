package provision_terraform

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"kubecompute/internal/core/domain"
	"kubecompute/internal/core/port"
	"kubecompute/internal/core/util"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

type TerraformNodeProvider struct {
	tfDir string
}

func NewVirtualboxNodeProvider() port.NodeProvider {
	baseDir := util.GetOsWorkDir()
	tfDir := filepath.Join(baseDir, "terraform")
	return &TerraformNodeProvider{tfDir: tfDir}
}

func (d *TerraformNodeProvider) EnsureNode(ctx context.Context, desired *domain.Node) (bool, error) {
	dir := filepath.Join(d.tfDir, desired.ObjectMeta.NamespacedName().String())
	slog.Info("ensuring terraform node", "dir", dir)

	// Make sure directory exists
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return false, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	prog, err := createTerraformProgram(desired)
	if err != nil {
		return false, err
	}

	progFile := filepath.Join(dir, "main.tf")
	if err := os.WriteFile(progFile, []byte(prog), 0o644); err != nil {
		return false, fmt.Errorf("failed to write terraform file: %w", err)
	}

	return false, errors.New("not implemented")
}

func (d *TerraformNodeProvider) DeleteNode(ctx context.Context, name domain.NamespacedName) (bool, error) {
	return false, errors.New("not implemented")
}

func createTerraformProgram(node *domain.Node) (string, error) {
	// TODO: improve template handling
	_, filename, _, _ := runtime.Caller(0) // 0 = this function
	dir := filepath.Dir(filename)
	tmplPath := filepath.Join(dir, "node.tf.tmpl")

	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}

	tfNodeSpec := struct {
		Name   string
		Image  string
		CPUs   int
		Memory string
	}{
		Name:   node.ObjectMeta.NamespacedName().String(),
		Image:  "https://app.vagrantup.com/ubuntu/boxes/bionic64/versions/20180903.0.0/providers/virtualbox.box",
		CPUs:   1,
		Memory: "512 mib",
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tfNodeSpec); err != nil {
		return "", err
	}
	return buf.String(), nil
}
