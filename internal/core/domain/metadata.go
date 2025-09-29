package domain

type NamespacedName struct {
	Namespace string `json:"namespace" yaml:"namespace"`
	Name      string `json:"name" yaml:"name"`
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
	Namespace       string `json:"namespace" yaml:"namespace"`
	Name            string `json:"name" yaml:"name"`
	ResourceVersion int    `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
}

func (om *ObjectMeta) NamespacedName() NamespacedName {
	return NamespacedName{
		Namespace: om.Namespace,
		Name:      om.Name,
	}
}
