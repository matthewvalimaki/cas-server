package spec

import (
	"fmt"
	"net/http"

	"github.com/matthewvalimaki/cas-server/security"
	"github.com/matthewvalimaki/cas-server/spec/xml"
	"github.com/matthewvalimaki/cas-server/types"
	"github.com/matthewvalimaki/cas-server/validators"
)

// HandleValidateV2 handles `/serviceValidate` request
func HandleValidateV2(w http.ResponseWriter, r *http.Request) {
	// do format check here as it will affect response format
	format := r.URL.Query().Get("format")
	if len(format) > 0 {
		err := validators.ValidateFormat(format)

		if err != nil {
			validateResponseV2("XML", err, nil, w, r)
			return
		}
	}
	if len(format) == 0 {
		format = "XML"
	}

	_, err := runValidators(w, r)
	if err != nil {
		validateResponseV2(format, err, nil, w, r)
		return
	}

	_, pgtURL, proxyGrantingTicket, proxyGrantingTicketIOU, err := runValidatorsV2(w, r)
	if err != nil {
		validateResponseV2(format, err, nil, w, r)
		return
	}

	// see: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html#servicevalidate-cas-20
	if pgtURL != "" {
		strg.SaveTicket(proxyGrantingTicket)

		validateResponseV2(format, nil, proxyGrantingTicketIOU, w, r)
		return
	}

	validateResponseV2(format, nil, nil, w, r)
}

func runValidatorsV2(w http.ResponseWriter, r *http.Request) (service *types.Service, pgtURL string, proxyGrantingTicket *types.Ticket, proxyGrantingTicketIOU *types.Ticket, err *types.CasError) {
	pgtURL = r.URL.Query().Get("pgtUrl")
	if len(pgtURL) > 0 {
		serviceParameter := r.URL.Query().Get("service")

		// make sure that pgtURL can be used with service
		service, err := validators.ValidateProxyGrantingURL(config, serviceParameter, pgtURL)
		if err != nil {
			return nil, "", nil, nil, err
		}

		// Make sure endpoint can be reached and uses SSL as dictated by CAS spec
		// see: https://jasig.github.io/cas/4.2.x/protocol/CAS-Protocol-Specification.html#head2.5.4
		err = validators.ValidateProxyURLEndpoint(pgtURL)
		if err != nil {
			return nil, "", nil, nil, err
		}

		// Generate PGT (ProxyGrantingTicket) and PGTIOU (ProxyGgrantingTicketIOU)
		proxyGrantingTicket, err := security.CreateNewProxyGrantingTicket()
		if err != nil {
			return nil, "", nil, nil, err
		}

		proxyGrantingTicketIOU, err := security.CreateNewProxyGrantingTicketIOU()
		if err != nil {
			return nil, "", nil, nil, err
		}

		strg.SaveTicket(proxyGrantingTicket)
		strg.SaveTicket(proxyGrantingTicketIOU)

		// reach out to proxy and then validate behavior
		err = validators.SendAndValidateProxyIDAndIOU(pgtURL, proxyGrantingTicket, proxyGrantingTicketIOU)
		if err != nil {
			return nil, "", nil, nil, err
		}

		return service, pgtURL, proxyGrantingTicket, proxyGrantingTicketIOU, nil
	}

	return nil, "", nil, nil, nil
}

func validateResponseV2(format string, casError *types.CasError, proxyGrantingTicketIOU *types.Ticket, w http.ResponseWriter, r *http.Request) {
	if format == "XML" {
		w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
	} else {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	}

	if casError != nil {
		fmt.Fprintf(w, xml.V2ValidationFailure(casError, format))
		return
	}

	fmt.Fprintf(w, xml.V2ValidationSuccess("test", proxyGrantingTicketIOU, format))
}

// HandleProxyV2 handles `proxy` request
func HandleProxyV2(w http.ResponseWriter, r *http.Request) {
	// do format check here as it will affect response format
	format := r.URL.Query().Get("format")
	if len(format) > 0 {
		err := validators.ValidateFormat(format)

		if err != nil {
			proxyResponseV2("XML", nil, err, w, r)
			return
		}
	}
	if len(format) == 0 {
		format = "XML"
	}

	err := runProxyValidatorsV2(w, r)
	if err != nil {
		proxyResponseV2(format, nil, err, w, r)
		return
	}

	proxyTicket, err := security.CreateNewProxyTicket()
	if err != nil {
		proxyResponseV2(format, nil, err, w, r)
		return
	}

	strg.SaveTicket(proxyTicket)

	proxyResponseV2(format, proxyTicket, nil, w, r)
}

func runProxyValidatorsV2(w http.ResponseWriter, r *http.Request) *types.CasError {
	pgt := r.URL.Query().Get("pgt")
	err := validators.ValidateTicket(pgt)
	if err != nil {
		return err
	}

	targetService := r.URL.Query().Get("targetService")
	ticket := &types.Ticket{Service: targetService, Ticket: pgt}
	err = security.ValidateProxyGrantingTicket(strg, ticket)
	if err != nil {
		return err
	}

	return nil
}

func proxyResponseV2(format string, proxyTicket *types.Ticket, casError *types.CasError, w http.ResponseWriter, r *http.Request) {
	if format == "XML" {
		w.Header().Set("Content-Type", "application/xml;charset=UTF-8")
	} else {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	}

	if casError != nil {
		fmt.Fprintf(w, xml.V2ProxyFailure(casError, format))
		return
	}

	fmt.Fprintf(w, xml.V2ProxySuccess(proxyTicket, format))
}
