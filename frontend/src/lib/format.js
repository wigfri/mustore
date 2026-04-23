export function formatPrice(price, currency) {
  return `${new Intl.NumberFormat('ru-RU').format(price)} ${currency}`
}
