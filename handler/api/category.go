package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	catId := r.Context().Value("id").(string)
	if catId == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idLogin, _ := strconv.Atoi(catId)
	categories, err := c.categoryService.GetCategories(r.Context(), idLogin)
	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(categories)
}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	userID := r.Context().Value("id").(string)
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}
	if category.Type == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	idlogin, _ := strconv.Atoi(userID)
	category2 := entity.Category{}
	category2.UserID = idlogin
	category2.Type = category.Type
	categories, err := c.categoryService.StoreCategory(r.Context(), &category2)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	data := map[string]interface{}{
		"user_id":     idlogin,
		"category_id": categories.ID,
		"message":     "success create new category",
	}
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(data)
}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	r.Context().Value("id")
	userID := r.URL.Query().Get("category_id")
	deleteuserID, _ := strconv.Atoi(userID)
	deleteuserID2 := entity.Category{}
	// deleteuserID2.UserID = deleteuserID2.UserID
	// deleteuserID2.ID = deleteuserID2.ID
	err := c.categoryService.DeleteCategory(r.Context(), deleteuserID)
	if err != nil {
		w.WriteHeader(500)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	data := map[string]interface{}{
		"user_id":     deleteuserID2.UserID,
		"category_id": deleteuserID2.ID,
		"message":     "success delete category",
	}
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(data)
}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
