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
func CreateCity(c *fiber.Ctx) error {
	var input models.CityInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	var newID int
	currentTime := time.Now()
	err := db.DB.QueryRow("INSERT INTO cities (name, create_at, updated_at) VALUES ($1, $2, $3) RETURNING id",
		input.Name, currentTime, currentTime).Scan(&newID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Nama sudah terdaftar"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat kota: " + err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "kota berhasil dibuat", "id": newID})
}

func GetAllCities(c *fiber.Ctx) error {
	rows, err := db.DB.Query("SELECT id, name FROM cities WHERE deleted_at IS NULL")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar kota"+err.Error()})
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var city models.City
		if err := rows.Scan(&city.ID,&city.Name); err != nil {
			log.Printf("Error saat scan cities: %v", err)
			continue
		}
		cities = append(cities, city)
	}

	return c.Status(http.StatusOK).JSON(cities)
}

func GetCityByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID kota tidak valid"})
	}

	var city models.City
	row := db.DB.QueryRow("SELECT id, name FROM cities WHERE id = $1 AND deleted_at IS NULL", id)
	err = row.Scan(&city.ID,&city.Name)

	if err == sql.ErrNoRows {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "kota tidak ditemukan"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil kota"})
	}

	return c.Status(http.StatusOK).JSON(city)
}

func UpdateCity(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID kota tidak valid"})
	}

	var input models.CityInput 
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	var existingID int
	err = db.DB.QueryRow("SELECT id FROM cities WHERE id = $1", id).Scan(&existingID)
	if err == sql.ErrNoRows {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "kota tidak ditemukan"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memeriksa kota"})
	}

	currentTime := time.Now()

	_, err = db.DB.Exec("UPDATE cities SET name = $1, updated_at = $2 WHERE id = $3",
		input.Name, currentTime, id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Username sudah terdaftar oleh kota lain"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui kota: " + err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "kota berhasil diperbarui"})
}

func DeleteCity(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID Kota tidak valid"})
	}
	currentTime := time.Now()
	res, err := db.DB.Exec("update cities SET deleted_at = $1 WHERE id = $2", currentTime, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus Kota"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Kota tidak ditemukan"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Kota berhasil dihapus"})
}