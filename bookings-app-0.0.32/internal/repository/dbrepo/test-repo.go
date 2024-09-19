package dbrepo

import (
	"time"

	"github.com/tsawler/bookings-app/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {

	return 1, nil
}

func (m *testDBRepo) InserRoomRestriction(r models.RoomRestriction) error {
	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) {

	return true, nil
}

func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room

	return rooms, nil
}

func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {

	var room models.Room

	return room, nil
}
func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User
	return u, nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {
	var res []models.Reservation
	return res, nil
}
