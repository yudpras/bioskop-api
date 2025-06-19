package	handler
import (
	"bioskop/db"
	"bioskop/models"
	"net/http"
	"time"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/lib/pq"
	"log"
	"strconv"
	"database/sql"
)
func CreateUser(c *fiber.Ctx) error {
	var input models.UserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengenkripsi password"})
	}

	var newID int
	currentTime := time.Now()
	err = db.DB.QueryRow("INSERT INTO users (name, username, password, create_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id",
		input.Name, input.Username, string(hashedPassword), currentTime, currentTime).Scan(&newID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Username sudah terdaftar"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat pengguna: " + err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "Pengguna berhasil dibuat", "id": newID})
}

func GetAllUsers(c *fiber.Ctx) error {
	rows, err := db.DB.Query("SELECT id, name, username FROM users WHERE deleted_at IS NULL")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar pengguna"+err.Error()})
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Username); err != nil {
			log.Printf("Error saat scan user: %v", err)
			continue
		}
		users = append(users, user)
	}

	return c.Status(http.StatusOK).JSON(users)
}

func GetUserByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}

	var user models.User
	row := db.DB.QueryRow("SELECT id, name, username FROM users WHERE id = $1 AND deleted_at IS NULL", id)
	err = row.Scan(&user.ID, &user.Name, &user.Username)

	if err == sql.ErrNoRows {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil pengguna"})
	}

	return c.Status(http.StatusOK).JSON(user)
}

func UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}

	var input models.UserInput 
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Permintaan tidak valid"})
	}

	var existingID int
	err = db.DB.QueryRow("SELECT id FROM users WHERE id = $1", id).Scan(&existingID)
	if err == sql.ErrNoRows {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
	} else if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memeriksa pengguna"})
	}

	var hashedPassword string
	if input.Password != "" { 
		hashedBytes, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengenkripsi password baru"})
		}
		hashedPassword = string(hashedBytes)
	} else {
		var oldPassword string
		err = db.DB.QueryRow("SELECT password FROM users WHERE id = $1", id).Scan(&oldPassword)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil password lama"})
		}
		hashedPassword = oldPassword
	}

	currentTime := time.Now()

	_, err = db.DB.Exec("UPDATE users SET name = $1, username = $2, password = $3, updated_at = $4 WHERE id = $5",
		input.Name, input.Username, hashedPassword, currentTime, id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Username sudah terdaftar oleh pengguna lain"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui pengguna: " + err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Pengguna berhasil diperbarui"})
}

func DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "ID pengguna tidak valid"})
	}
	currentTime := time.Now()
	res, err := db.DB.Exec("update users SET deleted_at = $1 WHERE id = $2", currentTime, id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghapus pengguna"})
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Pengguna berhasil dihapus"})
}