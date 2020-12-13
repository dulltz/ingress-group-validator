package api_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const groupNameAnnotation = "alb.ingress.kubernetes.io/group.name"

var _ = Describe("valid cases for Ingress validator", func() {
	It("should allow creating Ingress with namespaced group name", func() {
		ing := &netv1beta1.Ingress{}
		ing.Name = "allow-creating-namespaced"
		ing.Namespace = "default"
		ing.Annotations = map[string]string{groupNameAnnotation: ing.Namespace + "/test"}
		ing.Spec.Backend = &netv1beta1.IngressBackend{ServiceName: "test", ServicePort: intstr.FromInt(8080)}
		err := k8sClient.Create(testCtx, ing)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should allow creating Ingress that does not belong to any group", func() {
		ing := &netv1beta1.Ingress{}
		ing.Name = "allow-creating-empty"
		ing.Namespace = "default"
		ing.Spec.Backend = &netv1beta1.IngressBackend{ServiceName: "test", ServicePort: intstr.FromInt(8080)}
		err := k8sClient.Create(testCtx, ing)
		Expect(err).NotTo(HaveOccurred())
	})
})

var _ = Describe("invalid cases for Ingress", func() {
	It("should deny creating Ingress with namespaced group name", func() {
		ing := &netv1beta1.Ingress{}
		ing.Name = "deny-creating-not-namespaced"
		ing.Namespace = "default"
		ing.Annotations = map[string]string{groupNameAnnotation: "test"}
		ing.Spec.Backend = &netv1beta1.IngressBackend{ServiceName: "test", ServicePort: intstr.FromInt(8080)}
		err := k8sClient.Create(testCtx, ing)
		Expect(err).To(HaveOccurred())
	})

	It("should deny updating Ingress by adding invalid group name", func() {
		ing := &netv1beta1.Ingress{}
		ing.Name = "deny-updating-not-namespaced"
		ing.Namespace = "default"
		ing.Spec.Backend = &netv1beta1.IngressBackend{ServiceName: "test", ServicePort: intstr.FromInt(8080)}
		err := k8sClient.Create(testCtx, ing)
		Expect(err).NotTo(HaveOccurred())

		ing.Annotations = map[string]string{groupNameAnnotation: "test"}
		err = k8sClient.Update(testCtx, ing)
		Expect(err).To(HaveOccurred())
	})
})
