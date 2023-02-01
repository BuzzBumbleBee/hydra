// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package consent

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/ory/x/pagination/tokenpagination"

	"github.com/ory/x/httprouterx"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"

	"github.com/ory/fosite"
	"github.com/ory/hydra/v2/driver/config"
	"github.com/ory/hydra/v2/x"
	"github.com/ory/x/errorsx"
	"github.com/ory/x/sqlxx"
	"github.com/ory/x/stringsx"
	"github.com/ory/x/urlx"
)

type Handler struct {
	r InternalRegistry
	c *config.DefaultProvider
}

const (
	LoginPath    = "/oauth2/auth/requests/login"
	ConsentPath  = "/oauth2/auth/requests/consent"
	LogoutPath   = "/oauth2/auth/requests/logout"
	DevicePath   = "/oauth2/auth/requests/device"
	SessionsPath = "/oauth2/auth/sessions"
)

func NewHandler(
	r InternalRegistry,
	c *config.DefaultProvider,
) *Handler {
	return &Handler{
		c: c,
		r: r,
	}
}

func (h *Handler) SetRoutes(admin *httprouterx.RouterAdmin) {
	admin.GET(LoginPath, h.getOAuth2LoginRequest)
	admin.PUT(LoginPath+"/accept", h.acceptOAuth2LoginRequest)
	admin.PUT(LoginPath+"/reject", h.rejectOAuth2LoginRequest)

	admin.GET(ConsentPath, h.getOAuth2ConsentRequest)
	admin.PUT(ConsentPath+"/accept", h.acceptOAuth2ConsentRequest)
	admin.PUT(ConsentPath+"/reject", h.rejectOAuth2ConsentRequest)

	admin.DELETE(SessionsPath+"/login", h.revokeOAuth2LoginSessions)
	admin.GET(SessionsPath+"/consent", h.listOAuth2ConsentSessions)
	admin.DELETE(SessionsPath+"/consent", h.revokeOAuth2ConsentSessions)

	admin.GET(LogoutPath, h.getOAuth2LogoutRequest)
	admin.PUT(LogoutPath+"/accept", h.acceptOAuth2LogoutRequest)
	admin.PUT(LogoutPath+"/reject", h.rejectOAuth2LogoutRequest)

	admin.PUT(DevicePath+"/verify", h.verifyUserCodeRequest)
}

// Revoke OAuth 2.0 Consent Session Parameters
//
// swagger:parameters revokeOAuth2ConsentSessions
type revokeOAuth2ConsentSessions struct {
	// OAuth 2.0 Consent Subject
	//
	// The subject whose consent sessions should be deleted.
	//
	// in: query
	// required: true
	Subject string `json:"subject"`

	// OAuth 2.0 Client ID
	//
	// If set, deletes only those consent sessions that have been granted to the specified OAuth 2.0 Client ID.
	//
	// in: query
	Client string `json:"client"`

	// Revoke All Consent Sessions
	//
	// If set to `true` deletes all consent sessions by the Subject that have been granted.
	//
	// in: query
	All bool `json:"all"`
}

// swagger:route DELETE /admin/oauth2/auth/sessions/consent oAuth2 revokeOAuth2ConsentSessions
//
// # Revoke OAuth 2.0 Consent Sessions of a Subject
//
// This endpoint revokes a subject's granted consent sessions and invalidates all
// associated OAuth 2.0 Access Tokens. You may also only revoke sessions for a specific OAuth 2.0 Client ID.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  204: emptyResponse
//	  default: errorOAuth2
func (h *Handler) revokeOAuth2ConsentSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subject := r.URL.Query().Get("subject")
	client := r.URL.Query().Get("client")
	allClients := r.URL.Query().Get("all") == "true"
	if subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'subject' is not defined but should have been.`)))
		return
	}

	switch {
	case len(client) > 0:
		if err := h.r.ConsentManager().RevokeSubjectClientConsentSession(r.Context(), subject, client); err != nil && !errors.Is(err, x.ErrNotFound) {
			h.r.Writer().WriteError(w, r, err)
			return
		}
	case allClients:
		if err := h.r.ConsentManager().RevokeSubjectConsentSession(r.Context(), subject); err != nil && !errors.Is(err, x.ErrNotFound) {
			h.r.Writer().WriteError(w, r, err)
			return
		}
	default:
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter both 'client' and 'all' is not defined but one of them should have been.`)))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// List OAuth 2.0 Consent Session Parameters
//
// swagger:parameters listOAuth2ConsentSessions
type listOAuth2ConsentSessions struct {
	tokenpagination.RequestParameters

	// The subject to list the consent sessions for.
	//
	// in: query
	// required: true
	Subject string `json:"subject"`
	// The login session id to list the consent sessions for.
	//
	// in: query
	// required: false
	LoginSessionId string `json:"login_session_id"`
}

// swagger:route GET /admin/oauth2/auth/sessions/consent oAuth2 listOAuth2ConsentSessions
//
// # List OAuth 2.0 Consent Sessions of a Subject
//
// This endpoint lists all subject's granted consent sessions, including client and granted scope.
// If the subject is unknown or has not granted any consent sessions yet, the endpoint returns an
// empty JSON array with status code 200 OK.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2ConsentSessions
//	  default: errorOAuth2
func (h *Handler) listOAuth2ConsentSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subject := r.URL.Query().Get("subject")
	if subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'subject' is not defined but should have been.`)))
		return
	}
	loginSessionId := r.URL.Query().Get("login_session_id")

	page, itemsPerPage := x.ParsePagination(r)

	var s []AcceptOAuth2ConsentRequest
	var err error
	if len(loginSessionId) == 0 {
		s, err = h.r.ConsentManager().FindSubjectsGrantedConsentRequests(r.Context(), subject, itemsPerPage, itemsPerPage*page)
	} else {
		s, err = h.r.ConsentManager().FindSubjectsSessionGrantedConsentRequests(r.Context(), subject, loginSessionId, itemsPerPage, itemsPerPage*page)
	}
	if errors.Is(err, ErrNoPreviousConsentFound) {
		h.r.Writer().Write(w, r, []OAuth2ConsentSession{})
		return
	} else if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	var a []OAuth2ConsentSession
	for _, session := range s {
		session.ConsentRequest.Client = sanitizeClient(session.ConsentRequest.Client)
		a = append(a, OAuth2ConsentSession(session))
	}

	if len(a) == 0 {
		a = []OAuth2ConsentSession{}
	}

	n, err := h.r.ConsentManager().CountSubjectsGrantedConsentRequests(r.Context(), subject)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	x.PaginationHeader(w, r.URL, int64(n), itemsPerPage, itemsPerPage*page)
	h.r.Writer().Write(w, r, a)
}

// Revoke OAuth 2.0 Consent Login Sessions Parameters
//
// swagger:parameters revokeOAuth2LoginSessions
type revokeOAuth2LoginSessions struct {
	// OAuth 2.0 Subject
	//
	// The subject to revoke authentication sessions for.
	//
	// in: query
	// required: true
	Subject string `json:"subject"`
}

// swagger:route DELETE /admin/oauth2/auth/sessions/login oAuth2 revokeOAuth2LoginSessions
//
// # Revokes All OAuth 2.0 Login Sessions of a Subject
//
// This endpoint invalidates a subject's authentication session. After revoking the authentication session, the subject
// has to re-authenticate at the Ory OAuth2 Provider. This endpoint does not invalidate any tokens and
// does not work with OpenID Connect Front- or Back-channel logout.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  204: emptyResponse
//	  default: errorOAuth2
func (h *Handler) revokeOAuth2LoginSessions(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	subject := r.URL.Query().Get("subject")
	if subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'subject' is not defined but should have been.`)))
		return
	}

	if err := h.r.ConsentManager().RevokeSubjectLoginSession(r.Context(), subject); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get OAuth 2.0 Login Request
//
// swagger:parameters getOAuth2LoginRequest
type getOAuth2LoginRequest struct {
	// OAuth 2.0 Login Request Challenge
	//
	// in: query
	// required: true
	Challenge string `json:"login_challenge"`
}

// swagger:route GET /admin/oauth2/auth/requests/login oAuth2 getOAuth2LoginRequest
//
// # Get OAuth 2.0 Login Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, Ory asks the login provider
// to authenticate the subject and then tell the Ory OAuth2 Service about it.
//
// Per default, the login provider is Ory itself. You may use a different login provider which needs to be a web-app
// you write and host, and it must be able to authenticate ("show the subject a login screen")
// a subject (in OAuth2 the proper name for subject is "resource owner").
//
// The authentication challenge is appended to the login provider URL to which the subject's user-agent (browser) is redirected to. The login
// provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2LoginRequest
//	  410: oAuth2RedirectTo
//	  default: errorOAuth2
func (h *Handler) getOAuth2LoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("login_challenge"),
		r.URL.Query().Get("challenge"),
	)

	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	request, err := h.r.ConsentManager().GetLoginRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	if request.WasHandled {
		h.r.Writer().WriteCode(w, r, http.StatusGone, &OAuth2RedirectTo{
			RedirectTo: request.RequestURL,
		})
		return
	}

	request.Client = sanitizeClient(request.Client)
	h.r.Writer().Write(w, r, request)
}

// Accept OAuth 2.0 Login Request
//
// swagger:parameters acceptOAuth2LoginRequest
type acceptOAuth2LoginRequest struct {
	// OAuth 2.0 Login Request Challenge
	//
	// in: query
	// required: true
	Challenge string `json:"login_challenge"`

	// in: body
	Body HandledLoginRequest
}

// swagger:route PUT /admin/oauth2/auth/requests/login/accept oAuth2 acceptOAuth2LoginRequest
//
// # Accept OAuth 2.0 Login Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, Ory asks the login provider
// to authenticate the subject and then tell the Ory OAuth2 Service about it.
//
// The authentication challenge is appended to the login provider URL to which the subject's user-agent (browser) is redirected to. The login
// provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.
//
// This endpoint tells Ory that the subject has successfully authenticated and includes additional information such as
// the subject's ID and if Ory should remember the subject's subject agent for future authentication attempts by setting
// a cookie.
//
// The response contains a redirect URL which the login provider should redirect the user-agent to.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2RedirectTo
//	  default: errorOAuth2
func (h *Handler) acceptOAuth2LoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("login_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p HandledLoginRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithWrap(err).WithHintf("Unable to decode body because: %s", err)))
		return
	}

	if p.Subject == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Field 'subject' must not be empty.")))
		return
	}

	p.ID = challenge
	ar, err := h.r.ConsentManager().GetLoginRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	} else if ar.Subject != "" && p.Subject != ar.Subject {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Field 'subject' does not match subject from previous authentication.")))
		return
	}

	if ar.Skip {
		p.Remember = true // If skip is true remember is also true to allow consecutive calls as the same user!
		p.AuthenticatedAt = ar.AuthenticatedAt
	} else {
		p.AuthenticatedAt = sqlxx.NullTime(time.Now().UTC().
			// Rounding is important to avoid SQL time synchronization issues in e.g. MySQL!
			Truncate(time.Second))
		ar.AuthenticatedAt = p.AuthenticatedAt
	}
	p.RequestedAt = ar.RequestedAt

	request, err := h.r.ConsentManager().HandleLoginRequest(r.Context(), challenge, &p)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	ru, err := url.Parse(request.RequestURL)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &OAuth2RedirectTo{
		RedirectTo: urlx.SetQuery(ru, url.Values{"login_verifier": {request.Verifier}}).String(),
	})
}

// Reject OAuth 2.0 Login Request
//
// swagger:parameters rejectOAuth2LoginRequest
type rejectOAuth2LoginRequest struct {
	// OAuth 2.0 Login Request Challenge
	//
	// in: query
	// required: true
	Challenge string `json:"login_challenge"`

	// in: body
	Body RequestDeniedError
}

// swagger:route PUT /admin/oauth2/auth/requests/login/reject oAuth2 rejectOAuth2LoginRequest
//
// # Reject OAuth 2.0 Login Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, Ory asks the login provider
// to authenticate the subject and then tell the Ory OAuth2 Service about it.
//
// The authentication challenge is appended to the login provider URL to which the subject's user-agent (browser) is redirected to. The login
// provider uses that challenge to fetch information on the OAuth2 request and then accept or reject the requested authentication process.
//
// This endpoint tells Ory that the subject has not authenticated and includes a reason why the authentication
// was denied.
//
// The response contains a redirect URL which the login provider should redirect the user-agent to.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2RedirectTo
//	  default: errorOAuth2
func (h *Handler) rejectOAuth2LoginRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("login_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p RequestDeniedError
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithWrap(err).WithHintf("Unable to decode body because: %s", err)))
		return
	}

	p.valid = true
	p.SetDefaults(loginRequestDeniedErrorName)
	ar, err := h.r.ConsentManager().GetLoginRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	request, err := h.r.ConsentManager().HandleLoginRequest(r.Context(), challenge, &HandledLoginRequest{
		Error:       &p,
		ID:          challenge,
		RequestedAt: ar.RequestedAt,
	})
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	ru, err := url.Parse(request.RequestURL)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &OAuth2RedirectTo{
		RedirectTo: urlx.SetQuery(ru, url.Values{"login_verifier": {request.Verifier}}).String(),
	})
}

// Get OAuth 2.0 Consent Request
//
// swagger:parameters getOAuth2ConsentRequest
type getOAuth2ConsentRequest struct {
	// OAuth 2.0 Consent Request Challenge
	//
	// in: query
	// required: true
	Challenge string `json:"consent_challenge"`
}

// swagger:route GET /admin/oauth2/auth/requests/consent oAuth2 getOAuth2ConsentRequest
//
// # Get OAuth 2.0 Consent Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, Ory asks the login provider
// to authenticate the subject and then tell Ory now about it. If the subject authenticated, he/she must now be asked if
// the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject's behalf.
//
// The consent challenge is appended to the consent provider's URL to which the subject's user-agent (browser) is redirected to. The consent
// provider uses that challenge to fetch information on the OAuth2 request and then tells Ory if the subject accepted
// or rejected the request.
//
// The default consent provider is available via the Ory Managed Account Experience. To customize the consent provider, please
// head over to the OAuth 2.0 documentation.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2ConsentRequest
//	  410: oAuth2RedirectTo
//	  default: errorOAuth2
func (h *Handler) getOAuth2ConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("consent_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	request, err := h.r.ConsentManager().GetConsentRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}
	if request.WasHandled {
		h.r.Writer().WriteCode(w, r, http.StatusGone, &OAuth2RedirectTo{
			RedirectTo: request.RequestURL,
		})
		return
	}

	if request.RequestedScope == nil {
		request.RequestedScope = []string{}
	}

	if request.RequestedAudience == nil {
		request.RequestedAudience = []string{}
	}

	request.Client = sanitizeClient(request.Client)
	h.r.Writer().Write(w, r, request)
}

// Accept OAuth 2.0 Consent Request
//
// swagger:parameters acceptOAuth2ConsentRequest
type acceptOAuth2ConsentRequest struct {
	// OAuth 2.0 Consent Request Challenge
	//
	// in: query
	// required: true
	Challenge string `json:"consent_challenge"`

	// in: body
	Body AcceptOAuth2ConsentRequest
}

// swagger:route PUT /admin/oauth2/auth/requests/consent/accept oAuth2 acceptOAuth2ConsentRequest
//
// # Accept OAuth 2.0 Consent Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, Ory asks the login provider
// to authenticate the subject and then tell Ory now about it. If the subject authenticated, he/she must now be asked if
// the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject's behalf.
//
// The consent challenge is appended to the consent provider's URL to which the subject's user-agent (browser) is redirected to. The consent
// provider uses that challenge to fetch information on the OAuth2 request and then tells Ory if the subject accepted
// or rejected the request.
//
// This endpoint tells Ory that the subject has authorized the OAuth 2.0 client to access resources on his/her behalf.
// The consent provider includes additional information, such as session data for access and ID tokens, and if the
// consent request should be used as basis for future requests.
//
// The response contains a redirect URL which the consent provider should redirect the user-agent to.
//
// The default consent provider is available via the Ory Managed Account Experience. To customize the consent provider, please
// head over to the OAuth 2.0 documentation.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2RedirectTo
//	  default: errorOAuth2
func (h *Handler) acceptOAuth2ConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("consent_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p AcceptOAuth2ConsentRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errorsx.WithStack(err))
		return
	}

	cr, err := h.r.ConsentManager().GetConsentRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	p.ID = challenge
	p.RequestedAt = cr.RequestedAt
	p.HandledAt = sqlxx.NullTime(time.Now().UTC())

	hr, err := h.r.ConsentManager().HandleConsentRequest(r.Context(), &p)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	} else if hr.Skip {
		p.Remember = false
	}

	ru, err := url.Parse(hr.RequestURL)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &OAuth2RedirectTo{
		RedirectTo: urlx.SetQuery(ru, url.Values{"consent_verifier": {hr.Verifier}}).String(),
	})
}

// Reject OAuth 2.0 Consent Request
//
// swagger:parameters rejectOAuth2ConsentRequest
type adminRejectOAuth2ConsentRequest struct {
	// OAuth 2.0 Consent Request Challenge
	//
	// in: query
	// required: true
	Challenge string `json:"consent_challenge"`

	// in: body
	Body RequestDeniedError
}

// swagger:route PUT /admin/oauth2/auth/requests/consent/reject oAuth2 rejectOAuth2ConsentRequest
//
// # Reject OAuth 2.0 Consent Request
//
// When an authorization code, hybrid, or implicit OAuth 2.0 Flow is initiated, Ory asks the login provider
// to authenticate the subject and then tell Ory now about it. If the subject authenticated, he/she must now be asked if
// the OAuth 2.0 Client which initiated the flow should be allowed to access the resources on the subject's behalf.
//
// The consent challenge is appended to the consent provider's URL to which the subject's user-agent (browser) is redirected to. The consent
// provider uses that challenge to fetch information on the OAuth2 request and then tells Ory if the subject accepted
// or rejected the request.
//
// This endpoint tells Ory that the subject has not authorized the OAuth 2.0 client to access resources on his/her behalf.
// The consent provider must include a reason why the consent was not granted.
//
// The response contains a redirect URL which the consent provider should redirect the user-agent to.
//
// The default consent provider is available via the Ory Managed Account Experience. To customize the consent provider, please
// head over to the OAuth 2.0 documentation.
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2RedirectTo
//	  default: errorOAuth2
func (h *Handler) rejectOAuth2ConsentRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("consent_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p RequestDeniedError
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteErrorCode(w, r, http.StatusBadRequest, errorsx.WithStack(err))
		return
	}

	p.valid = true
	p.SetDefaults(consentRequestDeniedErrorName)
	hr, err := h.r.ConsentManager().GetConsentRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	request, err := h.r.ConsentManager().HandleConsentRequest(r.Context(), &AcceptOAuth2ConsentRequest{
		Error:       &p,
		ID:          challenge,
		RequestedAt: hr.RequestedAt,
		HandledAt:   sqlxx.NullTime(time.Now().UTC()),
	})
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(err))
		return
	}

	ru, err := url.Parse(request.RequestURL)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &OAuth2RedirectTo{
		RedirectTo: urlx.SetQuery(ru, url.Values{"consent_verifier": {request.Verifier}}).String(),
	})
}

// Accept OAuth 2.0 Logout Request
//
// swagger:parameters acceptOAuth2LogoutRequest
type acceptOAuth2LogoutRequest struct {
	// OAuth 2.0 Logout Request Challenge
	//
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`
}

// swagger:route PUT /admin/oauth2/auth/requests/logout/accept oAuth2 acceptOAuth2LogoutRequest
//
// # Accept OAuth 2.0 Session Logout Request
//
// When a user or an application requests Ory OAuth 2.0 to remove the session state of a subject, this endpoint is used to confirm that logout request.
//
// The response contains a redirect URL which the consent provider should redirect the user-agent to.
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2RedirectTo
//	  default: errorOAuth2
func (h *Handler) acceptOAuth2LogoutRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("logout_challenge"),
		r.URL.Query().Get("challenge"),
	)

	c, err := h.r.ConsentManager().AcceptLogoutRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	h.r.Writer().Write(w, r, &OAuth2RedirectTo{
		RedirectTo: urlx.SetQuery(urlx.AppendPaths(h.c.PublicURL(r.Context()), "/oauth2/sessions/logout"), url.Values{"logout_verifier": {c.Verifier}}).String(),
	})
}

// Reject OAuth 2.0 Logout Request
//
// swagger:parameters rejectOAuth2LogoutRequest
type rejectOAuth2LogoutRequest struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`
}

// swagger:route PUT /admin/oauth2/auth/requests/logout/reject oAuth2 rejectOAuth2LogoutRequest
//
// # Reject OAuth 2.0 Session Logout Request
//
// When a user or an application requests Ory OAuth 2.0 to remove the session state of a subject, this endpoint is used to deny that logout request.
// No HTTP request body is required.
//
// The response is empty as the logout provider has to chose what action to perform next.
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  204: emptyResponse
//	  default: errorOAuth2
func (h *Handler) rejectOAuth2LogoutRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("logout_challenge"),
		r.URL.Query().Get("challenge"),
	)

	if err := h.r.ConsentManager().RejectLogoutRequest(r.Context(), challenge); err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get OAuth 2.0 Logout Request
//
// swagger:parameters getOAuth2LogoutRequest
type getOAuth2LogoutRequest struct {
	// in: query
	// required: true
	Challenge string `json:"logout_challenge"`
}

// swagger:route GET /admin/oauth2/auth/requests/logout oAuth2 getOAuth2LogoutRequest
//
// # Get OAuth 2.0 Session Logout Request
//
// Use this endpoint to fetch an Ory OAuth 2.0 logout request.
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2LogoutRequest
//	  410: oAuth2RedirectTo
//	  default: errorOAuth2
func (h *Handler) getOAuth2LogoutRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("logout_challenge"),
		r.URL.Query().Get("challenge"),
	)

	request, err := h.r.ConsentManager().GetLogoutRequest(r.Context(), challenge)
	if err != nil {
		h.r.Writer().WriteError(w, r, err)
		return
	}

	// We do not want to share the secret so remove it.
	if request.Client != nil {
		request.Client.Secret = ""
	}

	if request.WasHandled {
		h.r.Writer().WriteCode(w, r, http.StatusGone, &OAuth2RedirectTo{
			RedirectTo: request.RequestURL,
		})
		return
	}

	h.r.Writer().Write(w, r, request)
}

// Verify OAuth 2.0 User Code Request
//
// swagger:parameters verifyUserCodeRequest
type verifyUserCodeRequest struct {
	// in: query
	// required: true
	Challenge string `json:"device_challenge"`

	// in: body
	Body DeviceGrantVerifyUserCodeRequest
}

// swagger:route PUT /admin/oauth2/auth/requests/device/verify oAuth2 verifyUserCodeRequest
//
// # Verifies a device grant request
//
// Verifies a device grant request
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http, https
//
//	Responses:
//	  200: oAuth2RedirectTo
//	  default: errorOAuth2
func (h *Handler) verifyUserCodeRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	challenge := stringsx.Coalesce(
		r.URL.Query().Get("device_challenge"),
		r.URL.Query().Get("challenge"),
	)
	if challenge == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint(`Query parameter 'challenge' is not defined but should have been.`)))
		return
	}

	var p DeviceGrantVerifyUserCodeRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&p); err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithWrap(err).WithHintf("Unable to decode body because: %s", err)))
		return
	}

	if p.UserCode == "" {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrInvalidRequest.WithHint("Field 'user_code' must not be empty.")))
		return
	}

	userCodeSignature, err := h.r.RFC8628HMACStrategy().UserCodeSignature(r.Context(), p.UserCode)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithHint(`'user_code' signature could not be computed`)))
		return
	}
	userCodeRequest, err := h.r.OAuth2Storage().GetUserCodeSession(r.Context(), userCodeSignature, &fosite.DefaultSession{})
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrNotFound.WithWrap(err).WithHint(`'user_code' session not found`)))
		return
	}

	clientId := userCodeRequest.GetClient().GetID()
	// UserCode & DeviceCode Request shares the same RequestId as it's the same request;
	deviceRequestId := userCodeRequest.GetID()
	requestedScopes := userCodeRequest.GetRequestedScopes()
	requestedAudience := userCodeRequest.GetRequestedAudience()

	err = h.r.OAuth2Storage().InvalidateUserCodeSession(r.Context(), userCodeSignature)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithHint(`Could not invalidate 'user_code'`)))
		return
	}

	// req.GetID() is actually the DeviceCodeSignature
	grantRequest, err := h.r.ConsentManager().AcceptDeviceGrantRequest(r.Context(), challenge, deviceRequestId, clientId, requestedScopes, requestedAudience)
	if err != nil {
		h.r.Writer().WriteError(w, r, errorsx.WithStack(fosite.ErrServerError.WithWrap(err).WithHint(`Could not accept device grant request`)))
		return
	}

	h.r.Writer().Write(w, r, &OAuth2RedirectTo{
		RedirectTo: urlx.SetQuery(h.c.OAuth2DeviceAuthorisationURL(r.Context()), url.Values{"device_verifier": {grantRequest.Verifier}, "client_id": {clientId}}).String(),
	})
}
