package	handler
import (
	"bioskop/db"
	"bioskop/models"
	"net/http"
	"time"
	"github.com/gofiber/fiber/v2"

	"github.com/lib/pq"
	"log"
	"strconv"
	"database/sql"
)
func CreateCinema(c *fiber.Ctx) error {
	var input models.CinemaInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	var newID int
	currentTime := time.Now()
	err := db.DB.QueryRow("INSERT INTO cinemas (cities_id, name, address, phone, create_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
	input.CityId, input.Name, input.Address, input.Phone, currentTime, currentTime).Scan(&newID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Nama sudah terdaftar"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat bioskop: " + err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "bioskop berhasil dibuat", "id": newID})
}

func GetAllCinemas(c *fiber.Ctx) error {
	rows, err := db.DB.Query("SELECT id, cities_id, name, address, phone FROM cinemas WHERE deleted_at IS NULL")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar bioskop"+err.Error()})
	}
	defer rows.Close()

	var cinemas []models.Cinema
	for rows.Next() {
		var cinema models.Cinema
		if err := rows.Scan(&cinema.ID,&cinema.CityId,&cinema.Name,&cinema.Address,&cinema.Phone); err != nil {
			log.Printf("Error saat scan cinemas: %v", err)
			continue
		}
		cinemas = append(cinemas, cinema)
	}

	return c.Status(http.StatusOK).JSON(cinemas)
}

func GetCinemaByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID bioskop tidak valid"})
	}

	var cinema models.Cinema
	row := db.DB.QueryRow("SELECT id, cities_id, name, address, phone FROM cinemas WHERE id = $1 AND deleted_at IS NULL", id)
	err = row.Scan(&cinema.ID,&cinema.CityId,&cinema.Name,&cinema.Address,&cinema.Phone)

	if err == sql.ErrNoRows {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "bioskop tidak ditemukan"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil bioskop"})
	}

	return c.Status(http.StatusOK).JSON(cinema)
}

func UpdateCinema(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID bioskop tidak valid"})
	}

	var input models.CinemaInput 
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	var existingIDcities int
	err = db.DB.QueryRow("SELECT id FROM cities WHERE id = $1", input.CityId).Scan(&existingIDcities)
	if err == sql.ErrNoRows {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "kota tidak ditemukan"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memeriksa kota"})
	}

	var existingIDcinema int
	err = db.DB.QueryRow("SELECT id FROM cinemas WHERE id = $1", id).Scan(&existingIDcinema)
	if err == sql.ErrNoRows {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "bioskop tidak ditemukan"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memeriksa bioskop"})
	}

	currentTime := time.Now()

	_, err = db.DB.Exec("UPDATE cinemas SET cities_id = $1, name = $2, address = $3, phone = $4, updated_at = $5 WHERE id = $6",
	input.CityId, input.Name, input.Address, input.Phone, currentTime, id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "nama sudah terdaftar oleh bioskop lain"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui bioskop: " + err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "bioskop berhasil diperbarui"})
}

func DeleteCinema(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID bioskop tidak valid"})
	}
	currentTime := time.Now()
	res, err := db.DB.Exec("update cinemas SET deleted_at = $1 WHERE id = $2", currentTime, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus bioskop"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "bioskop tidak ditemukan"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "bioskop berhasil dihapus"})
}