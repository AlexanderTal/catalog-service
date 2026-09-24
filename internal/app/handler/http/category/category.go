package hcategory

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"

	"github.com/AlexanderTal/catalog-service/internal/app/entity"
	rhandler "github.com/AlexanderTal/catalog-service/internal/app/handler/http"
	"github.com/AlexanderTal/catalog-service/internal/app/service"
	"github.com/AlexanderTal/catalog-service/internal/pkg/http/httph"
)

type handler struct {
	srv service.Category
}

func NewHandler(srv service.Category) rhandler.Category {
	return &handler{srv: srv}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req entity.RequestCategoryCreate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httph.HandleError(w, entity.ErrIncorrectParameters)
		return
	}

	if err := req.Validate(); err != nil {
		httph.HandleError(w, err)
		return
	}

	category, err := h.srv.Create(r.Context(), req)
	if err != nil {
		httph.HandleError(w, err)
		return
	}

	resp := entity.ResponseCategoryCreate{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
	}

	httph.SendJSON(w, http.StatusCreated, resp)
}

func (h *handler) Update(w http.ResponseWriter, r *http.Request) {
	guidStr := mux.Vars(r)["guid"]
	guid, err := uuid.FromString(guidStr)
	if err != nil {
		httph.HandleError(w, entity.ErrIncorrectParameters)
		return
	}

	var req entity.RequestCategoryUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httph.HandleError(w, entity.ErrIncorrectParameters)
		return
	}

	if err := req.Validate(); err != nil {
		httph.HandleError(w, err)
		return
	}

	category, err := h.srv.Update(r.Context(), guid, req)
	if err != nil {
		httph.HandleError(w, err)
		return
	}

	resp := entity.ResponseCategoryUpdate{
		GUID:      category.GUID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}

	httph.SendJSON(w, http.StatusOK, resp)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	guidStr := mux.Vars(r)["guid"]
	guid, err := uuid.FromString(guidStr)
	if err != nil {
		httph.HandleError(w, entity.ErrIncorrectParameters)
		return
	}

	if err := h.srv.Delete(r.Context(), guid); err != nil {
		httph.HandleError(w, err)
		return
	}

	httph.SendEmpty(w, http.StatusOK)
}

func (h *handler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.srv.List(r.Context())
	if err != nil {
		httph.HandleError(w, err)
		return
	}

	items := make([]entity.ResponseCategoryListItem, 0, len(categories))
	for _, c := range categories {
		items = append(items, entity.ResponseCategoryListItem{
			GUID:      c.GUID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		})
	}

	resp := entity.ResponseCategoryList{Data: items}

	httph.SendJSON(w, http.StatusOK, resp)
}
