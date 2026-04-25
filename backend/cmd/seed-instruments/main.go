// Command seed-instruments inserts demo musical instruments into Postgres.
// Requires CONFIG_PATH (same as the main app). By default runs only if the catalog is empty; use -force to append another demo batch.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/wigfri/mustore/app"
	"github.com/wigfri/mustore/domain/models"
	"github.com/wigfri/mustore/domain/repositories"
	"github.com/wigfri/mustore/services/config"
)

type seedRow struct {
	name        string
	slug        string
	brand       string
	category    string
	typ         string
	description string
	price       int64
	currency    string
	stock       int
	sku         string
	imageURL    string
}

// Demo catalog: mix of categories and price ranges (RUB by default).
var seeds = []seedRow{
	{
		name: "Fender Player Stratocaster", slug: "fender-player-stratocaster", brand: "Fender",
		category: "guitar", typ: "Электрогитара",
		description: "Классический звук трёх синглов, удобный гриф Modern C. Подходит для рока, блюза и фанка.",
		price: 124900, currency: "RUB", stock: 4,
		sku: "MUSTORE-SEED-001",
		imageURL: "https://images.unsplash.com/photo-1564186763535-ebb21ef5277f?w=800&q=80",
	},
	{
		name: "Gibson Les Paul Standard 60s", slug: "gibson-les-paul-standard-60s", brand: "Gibson",
		category: "guitar", typ: "Электрогитара",
		description: "Махагон корпус и гриф, накладка палисандр, звукосниматели Burstbucker. Тёплый sustain и плотный перегруз.",
		price: 389000, currency: "RUB", stock: 2,
		sku: "MUSTORE-SEED-002",
		imageURL: "https://images.unsplash.com/photo-1550985616-10811853b76d?w=800&q=80",
	},
	{
		name: "Ibanez RG470DX", slug: "ibanez-rg470dx", brand: "Ibanez",
		category: "guitar", typ: "Электрогитара",
		description: "Тремоло Edge Zero II, тонкий гриф Wizard III — для быстрых соло и метала.",
		price: 52990, currency: "RUB", stock: 6,
		sku: "MUSTORE-SEED-003",
		imageURL: "https://images.unsplash.com/photo-1510915361894-db8b60106cb1?w=800&q=80",
	},
	{
		name: "Yamaha Pacifica 112V", slug: "yamaha-pacifica-112v", brand: "Yamaha",
		category: "guitar", typ: "Электрогитара",
		description: "Универсальный стартовый инструмент: HSS, качественный тремоло и стабильный строй.",
		price: 27990, currency: "RUB", stock: 12,
		sku: "MUSTORE-SEED-004",
		imageURL: "https://images.unsplash.com/photo-1525201548942-d8732f661e1e?w=800&q=80",
	},
	{
		name: "Taylor GS Mini Mahogany", slug: "taylor-gs-mini-mahogany", brand: "Taylor",
		category: "guitar", typ: "Акустическая гитара",
		description: "Компактный корпус с богатым низом. Отличный вариант для дома и репетиций.",
		price: 62900, currency: "RUB", stock: 5,
		sku: "MUSTORE-SEED-005",
		imageURL: "https://images.unsplash.com/photo-1510915361894-db8b60106cb1?w=800&q=80",
	},
	{
		name: "Casio CDP-S360", slug: "casio-cdp-s360", brand: "Casio",
		category: "piano", typ: "Цифровое пианино",
		description: "88 взвешенных клавиш, встроенные ритмы и уроки. Лёгкий корпус для перевозки.",
		price: 45990, currency: "RUB", stock: 8,
		sku: "MUSTORE-SEED-006",
		imageURL: "https://images.unsplash.com/photo-1552422535-c45813c61732?w=800&q=80",
	},
	{
		name: "Roland FP-30X", slug: "roland-fp-30x", brand: "Roland",
		category: "piano", typ: "Цифровое пианино",
		description: "Клавиатура PHA-4, сэмплы SuperNATURAL, Bluetooth для приложений Roland Piano App.",
		price: 89900, currency: "RUB", stock: 3,
		sku: "MUSTORE-SEED-007",
		imageURL: "https://images.unsplash.com/photo-1573871669414-0090d7b5179f?w=800&q=80",
	},
	{
		name: "Kawai CA401", slug: "kawai-ca401", brand: "Kawai",
		category: "piano", typ: "Цифровое пианино",
		description: "Деревянные клавиши Grand Feel Compact, звук Shigeru Kawai SK-EX, премиальная мебельная стойка.",
		price: 329000, currency: "RUB", stock: 1,
		sku: "MUSTORE-SEED-008",
		imageURL: "https://images.unsplash.com/photo-1520523839897-bd0b52f945a0?w=800&q=80",
	},
	{
		name: "Pearl Roadshow RS525S", slug: "pearl-roadshow-rs525s", brand: "Pearl",
		category: "drums", typ: "Акустическая ударная установка",
		description: "Полный комплект 22\" + томы + флор + малый, тарелки и стойки в комплекте.",
		price: 54900, currency: "RUB", stock: 4,
		sku: "MUSTORE-SEED-009",
		imageURL: "https://images.unsplash.com/photo-1519892300165-cb5582dfdef2?w=800&q=80",
	},
	{
		name: "Yamaha Stage Custom Birch 20\"", slug: "yamaha-stage-custom-birch-20", brand: "Yamaha",
		category: "drums", typ: "Акустическая ударная установка",
		description: "Корпуса из берёзы, яркий атакующий звук. Конфигурация fusion: 10/12/14/20.",
		price: 89900, currency: "RUB", stock: 2,
		sku: "MUSTORE-SEED-010",
		imageURL: "https://images.unsplash.com/photo-1543443288-1f90aef46548?w=800&q=80",
	},
	{
		name: "Yamaha YAS-280", slug: "yamaha-yas-280", brand: "Yamaha",
		category: "wind", typ: "Альт-саксофон",
		description: "Учебная модель с устойчивым строем и лёгким откликом. Футляр в комплекте.",
		price: 92000, currency: "RUB", stock: 3,
		sku: "MUSTORE-SEED-011",
		imageURL: "https://images.unsplash.com/photo-1573871669414-0090d7b5179f?w=800&q=80",
	},
	{
		name: "Buffet Crampon E12 France", slug: "buffet-crampon-e12-france", brand: "Buffet Crampon",
		category: "wind", typ: "Кларнет B♭",
		description: "Гренадилла или ABS по серии, кольца посеребрённые, кейс Bam опционально.",
		price: 185000, currency: "RUB", stock: 2,
		sku: "MUSTORE-SEED-012",
		imageURL: "https://images.unsplash.com/photo-1519682337058-a94d519337bc?w=800&q=80",
	},
	{
		name: "Thomastik Dominant 4/4", slug: "thomastik-dominant-4-4", brand: "Thomastik",
		category: "string", typ: "Струны для скрипки",
		description: "Синтетическое ядро, тёплый тембр, стабильный строй. Полный комплект.",
		price: 4500, currency: "RUB", stock: 25,
		sku: "MUSTORE-SEED-013",
		imageURL: "https://images.unsplash.com/photo-1465821185615-80fd45f943ca?w=800&q=80",
	},
	{
		name: "Shure SM58", slug: "shure-sm58", brand: "Shure",
		category: "electronic", typ: "Динамический микрофон",
		description: "Стандарт сцены и студии: прочный корпус, кардиоида, вокальный звук с присутствием.",
		price: 10900, currency: "RUB", stock: 30,
		sku: "MUSTORE-SEED-014",
		imageURL: "https://images.unsplash.com/photo-1598653222000-6b7b7a552625?w=800&q=80",
	},
	{
		name: "Boss DS-1", slug: "boss-ds-1", brand: "Boss",
		category: "accessory", typ: "Педаль эффектов",
		description: "Классический дисторшн: от лёгкого кранча до насыщенного соло-саунда.",
		price: 6500, currency: "RUB", stock: 18,
		sku: "MUSTORE-SEED-015",
		imageURL: "https://images.unsplash.com/photo-1556590652-473933b10637?w=800&q=80",
	},
	{
		name: "D'Addario EXL110 (10-46)", slug: "daddario-exl110", brand: "D'Addario",
		category: "accessory", typ: "Струны для электрогитары",
		description: "Никелированная обмотка, сбалансированное натяжение, сет из трёх комплектов в упаковке магазина.",
		price: 850, currency: "RUB", stock: 80,
		sku: "MUSTORE-SEED-016",
		imageURL: "https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=800&q=80",
	},
}

func main() {
	force := flag.Bool("force", false, "insert demo rows even if the catalog is not empty (new UUIDs for id/slug/sku)")
	flag.Parse()

	if os.Getenv("CONFIG_PATH") == "" {
		log.Fatal("CONFIG_PATH is required (e.g. ./config/dev.yaml)")
	}

	cfg := config.Make()
	conn, err := app.InitDB(cfg)
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	existing, err := conn.Instrument().All(repositories.InstrumentListFilter{Limit: 1, Offset: 0})
	if err != nil {
		log.Fatalf("list instruments: %v", err)
	}
	if len(existing) > 0 && !*force {
		log.Println("catalog is not empty; skipping seed (run with -force to add another demo batch)")
		return
	}

	now := time.Now()
	inserted := 0
	for _, s := range seeds {
		id := uuid.New()
		slug := s.slug
		sku := s.sku
		if *force {
			slug = fmt.Sprintf("%s-%s", s.slug, id.String()[:8])
			sku = fmt.Sprintf("MUSTORE-%s", id.String())
		}

		inst := &models.Instrument{
			Id:          id,
			Name:        s.name,
			Slug:        slug,
			Brand:       s.brand,
			Category:    s.category,
			Type:        s.typ,
			Description: s.description,
			Price:       s.price,
			Currency:    s.currency,
			Stock:       s.stock,
			SKU:         sku,
			ImageURL:    s.imageURL,
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := inst.Validate(); err != nil {
			log.Fatalf("validate %q: %v", s.name, err)
		}
		if _, err := conn.Instrument().Insert(inst); err != nil {
			log.Fatalf("insert %q: %v", s.name, err)
		}
		inserted++
	}

	log.Printf("seeded %d instruments", inserted)
}
