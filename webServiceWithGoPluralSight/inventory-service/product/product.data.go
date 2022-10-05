package product

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("loading products form the file...")
	prodMap, err := loadProductMap()
	productMap.m = prodMap
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d products loaded...\n", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}
	file, _ := os.ReadFile(fileName)
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		log.Fatal(err)
	}
	prodMap := make(map[int]Product)
	for i := 0; i < len(productList); i++ {
		prodMap[productList[i].ProductId] = productList[i]
	}

	return prodMap, nil
}

func getProductByID(id int) *Product {
	productMap.RLock()
	defer productMap.RUnlock()
	if product, ok := productMap.m[id]; ok {
		return &product
	}
	return nil
}

func removeProductByid(id int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, id)
}

func getProductList() []Product {
	productMap.RLock()
	products := make([]Product, 0, len(productMap.m))
	for _, value := range productMap.m {
		products = append(products, value)
	}
	productMap.RUnlock()
	return products
}

func getProductIds() []int {
	productMap.RLock()
	var productIds []int
	for key := range productMap.m {
		productIds = append(productIds, key)
	}
	productMap.RUnlock()
	sort.Ints(productIds)
	return productIds
}

func getNextProductID() int {
	productIDs := getProductIds()
	return productIDs[len(productIDs)-1] + 1
}

func addOrUpdateProduct(product Product) (int, error) {
	addOrUpdateID := -1
	if product.ProductId > 0 {
		oldProduct := getProductByID(product.ProductId)
		if oldProduct == nil {
			return 0, fmt.Errorf("product [%d] doesn't exist", product.ProductId)
		}
		addOrUpdateID = product.ProductId
	} else {
		addOrUpdateID = getNextProductID()
		product.ProductId = addOrUpdateID
	}
	productMap.Lock()
	productMap.m[addOrUpdateID] = product
	productMap.Unlock()
	return addOrUpdateID, nil
}
