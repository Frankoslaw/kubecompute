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
	Namespace       string `yaml:"namespace"`
	Name            string `yaml:"name"`
	ResourceVersion int    `yaml:"resourceVersion,omitempty"`
}

func (om *ObjectMeta) NamespacedName() NamespacedName {
	return NamespacedName{
		Namespace: om.Namespace,
		Name:      om.Name,
	}
}
