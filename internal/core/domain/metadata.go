package domain

type NamespacedName struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
}

func (nn NamespacedName) String() string {
	return nn.Namespace + "-" + nn.Name
}

func (nn NamespacedName) ObjectMeta() ObjectMeta {
	return ObjectMeta{
		Namespace: nn.Namespace,
		Name:      nn.Name,
	}
}

type ObjectMeta struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	// to prevent warning from linters about NamespacedName and ObjectMeta being identical
	_ struct{}
}

func (om *ObjectMeta) NamespacedName() NamespacedName {
	return NamespacedName{
		Namespace: om.Namespace,
		Name:      om.Name,
	}
}
