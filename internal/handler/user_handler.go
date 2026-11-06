package handler

import (
	"net/http"
	"strings"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// UserHandler handles HTTP requests for users.
type UserHandler struct {
	service *service.UserService
}

// NewUserHandler creates a UserHandler that uses userService.
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{service: userService}
}

// GetAll handles GET /api/v1/users[?filterRole=<role name>][&companyId=<customer id>][&membersOnly=true]:
// the users the current user may see, e.g. ?filterRole=customer for the customer list and
// ?membersOnly=true for the users list (company members only: no customers, no superadmin).
// filterRole narrows the list by the listed users' role; the caller is always the session user.
//
// @Summary      Get users
// @Description  Returns the users the user may see, narrowed by filter: every user for the superadmin, the company's members for a customer or admin. filterRole=customer is the superadmin's customer list (Account Center › Customer); membersOnly=true leaves out customers and the superadmin, so only the users created for a company are listed (Masters › Users). A customer or admin asking for another company gets Forbidden.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Customers / Users list, Add / edit customer or user, Add / edit branch, Add / edit machine, Clientele, Add / edit part, Add / edit production schedule, Shifts, Breaks, Holidays / week off, Downtime / idle reasons, Licence keys
// @Tags         Users
// @ID           getUsers
// @Produce      json
// @Security     cookieAuth
// @Param        companyId query int false "Company id"
// @Param        filterRole query string false "Filter role"
// @Param        membersOnly query bool false "Only company members (users created by a customer), no customers"
// @Success      200 {object} response.Envelope{data=[]model.User}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/users [get]
func (userHandler *UserHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	companyID, err := queryID(request, "companyId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	membersOnly, err := queryBool(request, "membersOnly")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	filter := dto.UserListFilter{
		Role:        strings.ToLower(strings.TrimSpace(request.URL.Query().Get("filterRole"))),
		CompanyID:   companyID,
		MembersOnly: membersOnly,
	}

	users, err := userHandler.service.GetAll(request.Context(), actor, filter)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, users)
}

// Me handles GET /api/v1/me: the current user with its customer details.
//
// @Summary      Get current user
// @Description  Returns the logged-in user (userID) with its details.
// @Description
// @Description  Used by pages: Customers / Users list, Add / edit customer or user, Customer / user detail, App shell (session, logout, theme)
// @Tags         Users
// @ID           getCurrentUser
// @Produce      json
// @Security     cookieAuth
// @Success      200 {object} response.Envelope{data=dto.UserResponse}
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/me [get]
func (userHandler *UserHandler) Me(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	user, err := userHandler.service.GetCurrent(request.Context(), userID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, user)
}

// UpdateTheme handles PUT /api/v1/me/theme: saves the current user's UI theme.
//
// @Summary      Save my theme
// @Description  Sets the UI theme of userID (the logged-in user).
// @Description  Validation: Please enter valid theme.
// @Description
// @Description  Used by pages: Customers / Users list, Add / edit customer or user, Customer / user detail, App shell (session, logout, theme)
// @Tags         Users
// @ID           saveMyTheme
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.UpdateThemeRequest true "Request body"
// @Success      200 {object} response.Envelope{data=object}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/me/theme [put]
func (userHandler *UserHandler) UpdateTheme(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var themeRequest dto.UpdateThemeRequest
	if err := decodeAndValidate(writer, request, &themeRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := userHandler.service.SaveTheme(request.Context(), userID, themeRequest.Theme); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, dto.ThemeResponse{Theme: themeRequest.Theme})
}

// UpdateLanguage handles PUT /api/v1/me/language: saves the current user's UI language.
func (userHandler *UserHandler) UpdateLanguage(writer http.ResponseWriter, request *http.Request) {
	userID, err := currentUserID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var languageRequest dto.UpdateLanguageRequest
	if err := decodeAndValidate(writer, request, &languageRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := userHandler.service.SaveLanguage(request.Context(), userID, languageRequest.Language); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, dto.LanguageResponse{Language: languageRequest.Language})
}

// GetByID handles GET /api/v1/users/{id}.
//
// @Summary      Get user by id
// @Description  Returns user id with its details, or "User not found" when the user may not see it.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Customers / Users list, Add / edit customer or user, Customer / user detail
// @Tags         Users
// @ID           getUserById
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope{data=dto.UserResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/users/{id} [get]
func (userHandler *UserHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	user, err := userHandler.service.GetByID(request.Context(), actor, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, user)
}

// Create handles POST /api/v1/users.
//
// @Summary      Create user
// @Description  Adds a user with its customer details and machine access. the user must be allowed to give the role (see accessFor).
// @Description  machineFilter (supervisor, operator, qc, shop head) names the company's organization levels userCountryId … userLineId (null = any) that narrow the form's machine list; with allMachines the user may access every company machine on a line under them, machines added later included, and machineIds are not stored.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid password, Please select valid machines, Please select valid machine filter (a level that is not one chain of the company's organization).
// @Description
// @Description  Used by pages: Add / edit customer or user
// @Tags         Users
// @ID           createUser
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateUserRequest true "Request body"
// @Success      201 {object} response.Envelope{data=dto.UserResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/users [post]
func (userHandler *UserHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var userRequest dto.CreateUserRequest
	if err := decodeAndValidate(writer, request, &userRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	user, err := userHandler.service.Create(request.Context(), actor, userRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, user)
}

// Update handles PUT /api/v1/users/{id}.
//
// @Summary      Update user
// @Description  Replaces the user and its customer details. Customer details the user does not have yet are created. password sets a new password (8 to 72 characters); empty keeps the current one. The machine access and machineFilter are replaced (null removes the filter).
// @Description  machineFilter (supervisor, operator, qc, shop head) names the company's organization levels userCountryId … userLineId (null = any) that narrow the form's machine list; with allMachines the user may access every company machine on a line under them, machines added later included, and machineIds are not stored.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid password, Please enter valid status, Please select valid machines, Please select valid machine filter.
// @Description
// @Description  Used by pages: Customers / Users list, Add / edit customer or user
// @Tags         Users
// @ID           updateUser
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateUserRequest true "Request body"
// @Success      200 {object} response.Envelope{data=dto.UserResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/users/{id} [put]
func (userHandler *UserHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var userRequest dto.UpdateUserRequest
	if err := decodeAndValidate(writer, request, &userRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	user, err := userHandler.service.Update(request.Context(), actor, id, userRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, user)
}

// Delete handles DELETE /api/v1/users/{id}.
//
// @Summary      Delete user
// @Description  Soft-deletes the user with its profiles and addresses, and removes its machine access and machine filter. the user must be allowed to manage the user.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Customers / Users list
// @Tags         Users
// @ID           deleteUser
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/users/{id} [delete]
func (userHandler *UserHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := userHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "User deleted")
}
