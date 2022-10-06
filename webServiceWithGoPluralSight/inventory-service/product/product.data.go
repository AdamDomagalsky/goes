package product

import (
	"database/sql"

	"goes/webServiceWithGoPluralSight/database"
)

func getProductByID(id int) (*Product, error) {
	row := database.DbConn.QueryRow(`SELECT productId, 
       	manufacturer,
       	sku,
       	upc,
       	pricePerUnit,
       	quantityOnHand,
       	productName 
		FROM products
		WHERE productId = ?`, id)
	var product Product
	err := row.Scan(&product.ProductId,
		&product.Manufacturer,
		&product.Sku,
		&product.Upc,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &product, nil
}

func removeProductByid(id int) error {
	_, err := database.DbConn.Query(`DELETE FROM products WHERE productId = ?`, id)
	if err != nil {
		return err
	}
	return err
}

func getProductList() ([]Product, error) {
	results, err := database.DbConn.Query(`SELECT productId, 
       	manufacturer,
       	sku,
       	upc,
       	pricePerUnit,
       	quantityOnHand,
       	productName 
		FROM products`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		err := results.Scan(&product.ProductId,
			&product.Manufacturer,
			&product.Sku,
			&product.Upc,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName,
		)
		products = append(products, product)

		if err != nil {
			return nil, err
		}
	}
	return products, nil
}

func updateProduct(product Product) error {
	_, err := database.DbConn.Exec(`
	UPDATE products SET
	manufacturer=?,
	sku=?,
	upc=?,
	pricePerUnit=CAST(? AS DECIMAL (13,2)),
	quantityOnHand=?,
	productName=?
	WHERE productId=?`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductId)
	if err != nil {
		return err
	}
	return nil
}

func insertProduct(product Product) (int, error) {
	result, err := database.DbConn.Exec(`
	INSERT INTO products (manufacturer, sku, upc, pricePerUnit, quantityOnHand, productName) VALUES (?,?,?,?,?,?)`,
		product.Manufacturer,
		product.Sku,
		product.Upc,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName)
	if err != nil {
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(insertId), nil
}
