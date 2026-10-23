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

// GetAll handles GET /api/v1/users[?filterRole=<role name>][&companyId=<customer id>]:
// the users the current user may see, e.g. ?filterRole=customer for the customer list.
// filterRole narrows the list by the listed users' role; the caller is always the session user.
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
	filter := dto.UserListFilter{Role: strings.ToLower(strings.TrimSpace(request.URL.Query().Get("filterRole"))), CompanyID: companyID}

	users, err := userHandler.service.GetAll(request.Context(), actor, filter)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, users)
}

// Me handles GET /api/v1/me: the current user with its customer details.
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
