export const categories = [
  { value: 'guitar', label: 'Гитары' },
  { value: 'piano', label: 'Клавишные' },
  { value: 'drums', label: 'Ударные' },
  { value: 'wind', label: 'Духовые' },
  { value: 'string', label: 'Смычковые' },
  { value: 'accessory', label: 'Аксессуары' },
  { value: 'electronic', label: 'Электроника' },
]

export const currencies = ['RUB', 'USD', 'EUR']

export function categoryLabel(value) {
  return categories.find((c) => c.value === value)?.label ?? value
}
