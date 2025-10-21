package seeds

import (
	"log"

	"gorm.io/gorm"
)

// Novel model for seeding
type Novel struct {
	gorm.Model
	OriginalLanguage string  `json:"original_language" gorm:"not null;index"`
	OriginalAuthor   *string `json:"original_author" gorm:"not null;index"`
	Status           *string `json:"status" gorm:"index"`
	Source           *string `json:"source"`
	WordCount        *int    `json:"word_count"`
	CoverMediaId     *uint   `json:"cover_media_id"`
	CreatedBy        *uint   `json:"created_by" gorm:"index"`
}

// NovelTranslation model for seeding
type NovelTranslation struct {
	gorm.Model
	Title        string  `json:"title" gorm:"not null"`
	Synopsis     *string `json:"synopsis" gorm:"type:text"`
	TranslatorId *uint   `json:"translator_id"`
	NovelId      uint    `json:"novel_id" gorm:"not null;uniqueIndex:idx_novel_lang_unique"`
	Language     string  `json:"language" gorm:"not null;uniqueIndex:idx_novel_lang_unique"`
}

// NovelGenre junction table
type NovelGenre struct {
	NovelID uint `gorm:"primaryKey;index"`
	GenreID uint `gorm:"primaryKey;index"`
}

// NovelTag junction table
type NovelTag struct {
	NovelID uint `gorm:"primaryKey;index"`
	TagID   uint `gorm:"primaryKey;index"`
}

// SeedNovels seeds sample novels with translations
func SeedNovels(db *gorm.DB) error {
	log.Println("🌱 Seeding novels...")

	// Get author user
	var author User
	if err := db.Where("email = ?", "author1@example.com").First(&author).Error; err != nil {
		log.Println("⚠️  Author user not found, skipping novel seeding")
		return nil
	}

	novels := []struct {
		Novel        Novel
		Translations []NovelTranslation
		GenreSlugs   []string
		TagSlugs     []string
	}{
		{
			Novel: Novel{
				OriginalLanguage: "zh-CN",
				OriginalAuthor:   strPtr("李明"),
				Status:           strPtr("ongoing"),
				Source:           strPtr("https://example.com/novel1"),
				WordCount:        intPtr(500000),
				CreatedBy:        &author.ID,
			},
			Translations: []NovelTranslation{
				{
					Title:    "The Legendary Cultivator",
					Synopsis: strPtr("A young man embarks on a journey to become the strongest cultivator in the realm. Facing countless trials and powerful enemies, he must master ancient techniques and forge his own path to immortality."),
					Language: "en-US",
				},
				{
					Title:    "Kultivator Legendaris",
					Synopsis: strPtr("Seorang pemuda memulai perjalanan untuk menjadi kultivator terkuat di alam semesta. Menghadapi berbagai ujian dan musuh yang kuat, ia harus menguasai teknik kuno dan membentuk jalannya sendiri menuju keabadian."),
					Language: "id-ID",
				},
			},
			GenreSlugs: []string{"fantasy", "action", "martial-arts"},
			TagSlugs:   []string{"cultivation", "weak-to-strong", "overpowered-mc", "magic"},
		},
		{
			Novel: Novel{
				OriginalLanguage: "ja-JP",
				OriginalAuthor:   strPtr("田中太郎"),
				Status:           strPtr("ongoing"),
				Source:           strPtr("https://example.com/novel2"),
				WordCount:        intPtr(300000),
				CreatedBy:        &author.ID,
			},
			Translations: []NovelTranslation{
				{
					Title:    "Reborn in Another World",
					Synopsis: strPtr("After dying in a tragic accident, a salary worker finds himself reborn in a fantasy world with game-like mechanics. Armed with knowledge from his previous life, he sets out to live freely and enjoy his second chance."),
					Language: "en-US",
				},
				{
					Title:    "Terlahir Kembali di Dunia Lain",
					Synopsis: strPtr("Setelah meninggal dalam kecelakaan tragis, seorang pekerja kantoran mendapati dirinya terlahir kembali di dunia fantasi dengan mekanik seperti game. Dipersenjatai dengan pengetahuan dari kehidupan sebelumnya, ia bertekad untuk hidup bebas dan menikmati kesempatan keduanya."),
					Language: "id-ID",
				},
			},
			GenreSlugs: []string{"fantasy", "comedy", "slice-of-life"},
			TagSlugs:   []string{"isekai", "reincarnation", "system", "adventure"},
		},
		{
			Novel: Novel{
				OriginalLanguage: "ko-KR",
				OriginalAuthor:   strPtr("김철수"),
				Status:           strPtr("completed"),
				Source:           strPtr("https://example.com/novel3"),
				WordCount:        intPtr(800000),
				CreatedBy:        &author.ID,
			},
			Translations: []NovelTranslation{
				{
					Title:    "Shadow Monarch",
					Synopsis: strPtr("In a world where dungeons and monsters have become reality, the weakest hunter receives a mysterious power that allows him to rise from the dead and command an army of shadows. His journey from the weakest to the strongest begins."),
					Language: "en-US",
				},
				{
					Title:    "Raja Bayangan",
					Synopsis: strPtr("Di dunia di mana dungeon dan monster telah menjadi kenyataan, pemburu terlemah menerima kekuatan misterius yang memungkinkannya bangkit dari kematian dan mengendalikan pasukan bayangan. Perjalanannya dari yang terlemah ke yang terkuat dimulai."),
					Language: "id-ID",
				},
			},
			GenreSlugs: []string{"action", "fantasy", "horror"},
			TagSlugs:   []string{"weak-to-strong", "system", "dungeon", "monster", "overpowered-mc"},
		},
	}

	for i, novelData := range novels {
		// Check if novel exists (by author and language combination)
		var existingNovel Novel
		result := db.Where("original_language = ? AND original_author = ?",
			novelData.Novel.OriginalLanguage,
			*novelData.Novel.OriginalAuthor).First(&existingNovel)

		if result.Error == gorm.ErrRecordNotFound {
			// Create novel
			if err := db.Create(&novelData.Novel).Error; err != nil {
				log.Printf("⚠️  Failed to seed novel %d: %v", i+1, err)
				return err
			}

			// Create translations
			for _, trans := range novelData.Translations {
				trans.NovelId = novelData.Novel.ID
				trans.TranslatorId = &author.ID
				if err := db.Create(&trans).Error; err != nil {
					log.Printf("⚠️  Failed to seed translation for novel %d: %v", i+1, err)
				}
			}

			// Add genres
			for _, genreSlug := range novelData.GenreSlugs {
				var genre Genre
				if err := db.Where("slug = ?", genreSlug).First(&genre).Error; err == nil {
					novelGenre := NovelGenre{NovelID: novelData.Novel.ID, GenreID: genre.ID}
					db.Create(&novelGenre)
				}
			}

			// Add tags
			for _, tagSlug := range novelData.TagSlugs {
				var tag Tag
				if err := db.Where("slug = ?", tagSlug).First(&tag).Error; err == nil {
					novelTag := NovelTag{NovelID: novelData.Novel.ID, TagID: tag.ID}
					db.Create(&novelTag)
				}
			}

			log.Printf("✅ Created novel: %s", novelData.Translations[0].Title)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("⏭️  Novel already exists: %s", *novelData.Novel.OriginalAuthor)
		}
	}

	log.Println("✅ Novels seeding completed")
	return nil
}
