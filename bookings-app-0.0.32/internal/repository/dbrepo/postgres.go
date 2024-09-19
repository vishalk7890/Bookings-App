package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/tsawler/bookings-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDbRepo) AllUsers() bool {
	return true
}
func (m *postgresDbRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var newID int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) 
             values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`

	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (m *postgresDbRepo) InserRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at) values ($1,$2,$3,$4,$5,$6,$7)`
	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationId,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDbRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `select count(restriction_id) from  room_restrictions where room_id = $1 and $2 < end_date and $3 > start_date;`
	var numRows int
	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	return true, err
}

// select r.id , r.room_name from rooms r where r.id not in (select rr.room_id from room_restrictions rr where '2022-10-21' < rr.end_date and '2022-11-21' > rr.start_date)
func (m *postgresDbRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := `select r.id , r.room_name from rooms r where r.id not in (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)`
	var rooms []models.Room
	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		return rooms, err
	}
	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		log.Fatal("Error scanning rows")
	}

	return rooms, err
}

func (m *postgresDbRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var room models.Room
	query := `select room_name , id from rooms where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(

		&room.RoomName,
		&room.ID,
	)
	if err != nil {
		log.Fatal("Error scanning rows")
		return room, err
	}
	return room, nil
}

func (m *postgresDbRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password , accesslevel, created_at, updated_at from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Fatal("cant find record", err)
		return user, err
	}
	return user, nil
}

func (m *postgresDbRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, access_level = $4, updated_at = $5`
	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.Email, u.AccessLevel, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *postgresDbRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var id int
	var hashedPassword string
	row := m.DB.QueryRowContext(ctx, "select id , password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}
	return id, hashedPassword, nil
}

func (m *postgresDbRepo) AllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var reservation []models.Reservation

	query := ` select r.id, r.first_name, r.last_name, r.email, r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, rm.id,rm_room_name from reservations r left join rooms rm on(r.room_id = rm.id) order by r.start_date asc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return reservation, err
	}
	defer rows.Close()
	for rows.Next() {
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Room.ID,
			&i.Room.RoomName,
		)
		if err != nil {
			log.Fatal("Error scanning rows")
			return reservation, err
		}
		reservation = append(reservation, i)

	}
	if err = rows.Err(); err != nil {
		return reservation, err
	}
	return reservation, nil
}
