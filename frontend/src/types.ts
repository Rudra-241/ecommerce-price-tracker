// Shapes mirror what the Go backend serializes. The GORM models currently have
// no json tags, so fields come back PascalCased (ID, Name, Price, ...). The
// price-history endpoint, however, returns an explicit gin.H with snake_case
// keys. Keep these in sync with src/internal/models and handlers/product.go.

export type Direction = 'Increased' | 'Decreased' | 'Unchanged' | 'BelowStart'

export interface Product {
  ID: number
  Name: string
  Price: number
  ImgLink: string
  Url: string
  Direction: Direction
  ModifiedAt: string
}

export interface PriceStamp {
  ID: number
  ProductID: number
  Price: number
  ChangedAt: string
}

export interface PriceHistoryResponse {
  product_id: number
  product_name: string
  current_price: number
  price_history: PriceStamp[]
}
