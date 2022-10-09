package product

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"goes/webServiceWithGoPluralSight/cors"
	"golang.org/x/net/websocket"
)

const productsPath = "products"

func SetupRoute(apiBasePath string) {
	handleProducts := http.HandlerFunc(productsHandler)
	handleProduct := http.HandlerFunc(productHandler)
	http.Handle("/websocket", websocket.Handler(productSocket))
	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsPath), cors.Middleware(handleProducts))
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsPath), cors.Middleware(handleProduct))

}
func productHandler(writer http.ResponseWriter, request *http.Request) {
	urlPathSegments := strings.Split(request.URL.Path, fmt.Sprintf("%s/", productsPath))
	if len(urlPathSegments[1:]) > 1 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}
	product, err := getProductByID(productID)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	if product == nil {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	switch request.Method {
	case http.MethodGet:

		productJson, err := json.Marshal(product)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.Write(productJson)
		return
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := io.ReadAll(request.Body)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedProduct.ProductId != productID {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		err = updateProduct(updatedProduct)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		err := removeProductByid(productID)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodOptions:
		return // cors.Middleware handles the logic for us
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		productList, err := getProductList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		productsJson, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(productsJson)
	case http.MethodPost:
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var newProduct Product
		err = json.Unmarshal(bodyBytes, &newProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if newProduct.ProductId != 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = insertProduct(newProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}
