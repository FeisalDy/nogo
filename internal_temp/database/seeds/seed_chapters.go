package seeds

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// Chapter model for seeding
type Chapter struct {
	gorm.Model
	NovelId   uint `json:"novel_id" gorm:"not null;uniqueIndex:idx_novel_chapter_unique"`
	Number    int  `json:"number" gorm:"not null;uniqueIndex:idx_novel_chapter_unique"`
	WordCount *int `json:"word_count"`
}

// ChapterTranslation model for seeding
type ChapterTranslation struct {
	gorm.Model
	ChapterId    uint   `json:"chapter_id" gorm:"not null;uniqueIndex:idx_chapter_lang_unique"`
	Language     string `json:"language" gorm:"not null;uniqueIndex:idx_chapter_lang_unique"`
	Title        string `json:"title" gorm:"not null"`
	Content      string `json:"content" gorm:"type:text"`
	TranslatorId *uint  `json:"translator_id"`
}

// SeedChapters seeds sample chapters for existing novels
func SeedChapters(db *gorm.DB) error {
	log.Println("🌱 Seeding chapters...")

	// Get author user
	var author User
	if err := db.Where("email = ?", "author1@example.com").First(&author).Error; err != nil {
		log.Println("⚠️  Author user not found, skipping chapter seeding")
		return nil
	}

	// Get all novels
	var novels []Novel
	if err := db.Find(&novels).Error; err != nil {
		log.Println("⚠️  No novels found, skipping chapter seeding")
		return nil
	}

	// Seed 5 chapters for each novel
	for _, novel := range novels {
		for chapterNum := 1; chapterNum <= 5; chapterNum++ {
			// Check if chapter exists
			var existingChapter Chapter
			result := db.Where("novel_id = ? AND number = ?", novel.ID, chapterNum).First(&existingChapter)

			if result.Error == gorm.ErrRecordNotFound {
				chapter := Chapter{
					NovelId:   novel.ID,
					Number:    chapterNum,
					WordCount: intPtr(2500 + (chapterNum * 100)),
				}

				if err := db.Create(&chapter).Error; err != nil {
					log.Printf("⚠️  Failed to seed chapter %d for novel %d: %v", chapterNum, novel.ID, err)
					continue
				}

				// Create translations for this chapter
				translations := []ChapterTranslation{
					{
						ChapterId:    chapter.ID,
						Language:     "en-US",
						Title:        fmt.Sprintf("Chapter %d: The Beginning of the Journey", chapterNum),
						Content:      generateChapterContent("en-US", chapterNum),
						TranslatorId: &author.ID,
					},
					{
						ChapterId:    chapter.ID,
						Language:     "id-ID",
						Title:        fmt.Sprintf("Bab %d: Permulaan Perjalanan", chapterNum),
						Content:      generateChapterContent("id-ID", chapterNum),
						TranslatorId: &author.ID,
					},
				}

				for _, trans := range translations {
					if err := db.Create(&trans).Error; err != nil {
						log.Printf("⚠️  Failed to seed chapter translation: %v", err)
					}
				}

				log.Printf("✅ Created chapter %d for novel ID %d", chapterNum, novel.ID)
			} else if result.Error != nil {
				return result.Error
			} else {
				log.Printf("⏭️  Chapter %d for novel %d already exists", chapterNum, novel.ID)
			}
		}
	}

	log.Println("✅ Chapters seeding completed")
	return nil
}

// Helper function to generate sample chapter content
func generateChapterContent(language string, chapterNum int) string {
	if language == "id-ID" {
		return fmt.Sprintf(`Ini adalah konten dari bab %d. Dalam bab ini, protagonis menghadapi tantangan baru yang akan menguji kemampuan dan tekadnya.

Dia berdiri di tepi jurang, menatap ke kegelapan di bawah. Angin malam berhembus kencang, membuat jubahnya berkibar. Dia tahu bahwa satu langkah salah bisa berarti akhir dari perjalanannya.

"Aku harus melanjutkan," gumamnya pada diri sendiri. "Terlalu banyak yang bergantung padaku."

Dengan langkah mantap, dia melangkah maju menuju takdirnya...

[Konten bab selanjutnya akan berlanjut dengan petualangan yang mendebarkan dan pengembangan karakter yang mendalam. Pembaca akan melihat pertumbuhan protagonis saat ia menghadapi berbagai rintangan dan musuh yang semakin kuat.]`, chapterNum)
	}

	return fmt.Sprintf(`This is the content of chapter %d. In this chapter, the protagonist faces new challenges that will test their abilities and resolve.

He stood at the edge of the cliff, staring into the darkness below. The night wind blew fiercely, making his robe flutter. He knew that one wrong step could mean the end of his journey.

"I must continue," he muttered to himself. "Too much depends on me."

With steady steps, he moved forward toward his destiny...

[The rest of the chapter will continue with exciting adventures and deep character development. Readers will witness the protagonist's growth as they face various obstacles and increasingly powerful enemies.]`, chapterNum)
}
