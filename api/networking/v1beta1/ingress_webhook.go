package v1beta1

import (
	"context"
	"net/http"

	"github.com/dulltz/ingress-group-validator/pkg"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-networking-v1beta1-ingress,mutating=false,failurePolicy=fail,groups="networking.k8s.io",resources=ingresses,verbs=create;update,versions=v1beta1,name=vingress.kb.io

// IngressValidator validates Ingresses
type IngressValidator struct {
	decoder *admission.Decoder
}

func (v IngressValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	ing := &netv1beta1.Ingress{}
	if err := v.decoder.Decode(req, ing); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	if err := pkg.ValidateGroupName(ing); err != nil {
		return admission.Denied(err.Error())
	}
	return admission.Allowed("")
}
