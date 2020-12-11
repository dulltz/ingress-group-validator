package pkg

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Ref https://kubernetes-sigs.github.io/aws-load-balancer-controller/guide/ingress/annotations/#group.name
const groupNameAnnotation = "alb.ingress.kubernetes.io/group.name"

// ValidateGroupName checks the ingress group name is valid
func ValidateGroupName(ing metav1.Object) error {
	groupName, found := ing.GetAnnotations()[groupNameAnnotation]
	if !found {
		return nil
	}

	desiredPrefix := fmt.Sprintf("%s/", ing.GetNamespace())
	if strings.HasPrefix(groupName, desiredPrefix) {
		return nil
	}

	return fmt.Errorf("deny '%s' since the %s annotation does not start with '%s' namespaced", ing.GetName(), groupNameAnnotation, desiredPrefix)
}
