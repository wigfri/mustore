export const emptyInstrument = () => ({
  name: '',
  brand: '',
  category: 'guitar',
  type: '',
  description: '',
  price: 0,
  currency: 'RUB',
  stock: 0,
  sku: '',
  image_url: '',
  is_active: true,
})

export function rowToForm(row) {
  if (!row) {
    return emptyInstrument()
  }
  return {
    name: row.name,
    brand: row.brand,
    category: row.category,
    type: row.type,
    description: row.description,
    price: row.price,
    currency: row.currency,
    stock: row.stock,
    sku: row.sku,
    image_url: row.image_url,
    is_active: row.is_active,
  }
}
