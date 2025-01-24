package repositories

import (
	"database/sql"
	"github.com/Fyefhqdishka/LocFinder/internal/models"
	"log/slog"
)

type LocRepository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewLocRepository(db *sql.DB, log *slog.Logger) *LocRepository {
	return &LocRepository{
		db:  db,
		log: log,
	}
}

func (r *LocRepository) GetByIP(ip string) (string, string, error) {
	query := `SELECT country, city FROM locations WHERE ip_address = $1`
	row := r.db.QueryRow(query, ip)

	var country, city string
	err := row.Scan(&country, &city)
	return country, city, err
}

func (r *LocRepository) Save(ip, country, city string) error {
	query := `INSERT INTO locations (ip_address, country, city, created_at) VALUES ($1, $2, $3, NOW())`
	_, err := r.db.Exec(query, ip, country, city)
	return err
}

func (r *LocRepository) Update(ip, country, city string) error {
	query := `UPDATE locations SET country = $2, city = $3 WHERE ip_address = $1`
	_, err := r.db.Exec(query, ip, country, city)
	return err
}

func (r *LocRepository) Delete(ip string) error {
	query := `DELETE FROM locations WHERE ip_address = $1`
	_, err := r.db.Exec(query, ip)
	return err
}

func (r *LocRepository) GetAll() ([]models.IPLocation, error) {
	query := `SELECT ip_address, country, city FROM locations`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locations []models.IPLocation
	for rows.Next() {
		var location models.IPLocation
		if err := rows.Scan(&location.IP, &location.Country, &location.City); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}
	return locations, nil
}
