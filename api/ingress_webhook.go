package api

import (
	"context"
	"net/http"

	"github.com/dulltz/ingress-group-validator/pkg"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-networking-v1beta1-ingress,mutating=false,failurePolicy=fail,groups="networking.k8s.io";"extensions",resources=ingresses,verbs=create;update,versions=v1beta1,name=vingress.kb.io

// IngressValidator validates Ingresses
type IngressValidator struct {
	Decoder *admission.Decoder
}

var handleLog = ctrl.Log.WithName("handle")

func (v IngressValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	ing := &netv1beta1.Ingress{}
	handleLog.Info("decoding to networking.k8s.io/Ingress")
	if err := v.Decoder.Decode(req, ing); err != nil {
		handleLog.Info("failed to decode input to networking.k8s.io/Ingress, so trying to decode to extensions/Ingress")
		ing := &extv1beta1.Ingress{}
		if err := v.Decoder.Decode(req, ing); err != nil {
			return admission.Errored(http.StatusBadRequest, err)
		}
	}
	if err := pkg.ValidateGroupName(ing); err != nil {
		return admission.Denied(err.Error())
	}
	return admission.Allowed("")
}
