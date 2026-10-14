package handler

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"strconv"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// maxDrawingBodyBytes is the largest POST /parts/drawings body: the PDF plus room
// for the multipart headers.
const maxDrawingBodyBytes = dto.PartDrawingMaxBytes + 1<<20

// PartHandler handles HTTP requests for parts. The superadmin sees and changes every
// customer's parts; everyone else reads their company's and managers change them
// (service.PartService decides which).
type PartHandler struct {
	service *service.PartService
}

// NewPartHandler creates a PartHandler that uses partService.
func NewPartHandler(partService *service.PartService) *PartHandler {
	return &PartHandler{service: partService}
}

// GetAll handles GET /api/v1/parts[?userId=<customer id>].
func (partHandler *PartHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	customerID, err := queryID(request, "userId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	parts, err := partHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, parts)
}

// GetByID handles GET /api/v1/parts/{id}[?userId=<customer id>].
func (partHandler *PartHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	customerID, err := queryID(request, "userId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	part, err := partHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, part)
}

// Create handles POST /api/v1/parts.
func (partHandler *PartHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var partRequest dto.CreatePartRequest
	if err := decodeAndValidate(writer, request, &partRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	part, err := partHandler.service.Create(request.Context(), actor, partRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, part)
}

// Update handles PUT /api/v1/parts/{id}.
func (partHandler *PartHandler) Update(writer http.ResponseWriter, request *http.Request) {
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

	var partRequest dto.UpdatePartRequest
	if err := decodeAndValidate(writer, request, &partRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	part, err := partHandler.service.Update(request.Context(), actor, id, partRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, part)
}

// Delete handles DELETE /api/v1/parts/{id}.
func (partHandler *PartHandler) Delete(writer http.ResponseWriter, request *http.Request) {
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

	if err := partHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Part deleted")
}

// UploadDrawing handles POST /api/v1/parts/drawings: a multipart form whose "file"
// field is the part drawing PDF. The returned file name goes in a process's drawingFile.
func (partHandler *PartHandler) UploadDrawing(writer http.ResponseWriter, request *http.Request) {
	request.Body = http.MaxBytesReader(writer, request.Body, maxDrawingBodyBytes)
	file, header, err := request.FormFile("file")
	if err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			response.SendResponse(writer, request, response.NewError(http.StatusBadRequest, "Part drawing must be 10 MB or smaller"), nil)
			return
		}
		response.SendResponse(writer, request, response.NewError(http.StatusBadRequest, "Please upload part drawing (PDF)"), nil)
		return
	}
	defer file.Close()

	content, err := io.ReadAll(io.LimitReader(file, dto.PartDrawingMaxBytes+1))
	if err != nil {
		response.SendResponse(writer, request, response.NewError(http.StatusBadRequest, "Please upload part drawing (PDF)"), nil)
		return
	}
	drawingRequest := dto.UploadPartDrawingRequest{Name: header.Filename, Content: content}
	if err := drawingRequest.Validate(); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	drawing, err := partHandler.service.UploadDrawing(request.Context(), drawingRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, drawing)
}

// Drawing handles GET /api/v1/parts/{id}/processes/{processId}/drawing: the drawing
// PDF of the process, shown inline in the browser.
func (partHandler *PartHandler) Drawing(writer http.ResponseWriter, request *http.Request) {
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
	processID, err := strconv.ParseInt(request.PathValue("processId"), 10, 64)
	if err != nil || processID <= 0 {
		response.SendResponse(writer, request, response.NewError(http.StatusBadRequest, "Please enter valid process id"), nil)
		return
	}

	file, name, err := partHandler.service.OpenDrawing(request.Context(), actor, id, processID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	writer.Header().Set("Content-Type", "application/pdf")
	writer.Header().Set("Content-Disposition", mime.FormatMediaType("inline", map[string]string{"filename": name}))
	writer.Header().Set("X-Content-Type-Options", "nosniff")
	writer.Header().Set("Cache-Control", "private, no-store")
	http.ServeContent(writer, request, "", info.ModTime(), file)
}
