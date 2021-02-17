package models

import (
	"fmt"
	"time"

	"github.com/abulhanifah/classroom/helpers"
)

type Classroom struct {
	ID        int       `json:"id,omitempty" form:"id,omitempty" query:"id,omitempty"`
	Name      string    `json:"name" gorm:"type:varchar(50)" validate:"required"`
	Rows      int       `json:"rows,omitempty" form:"rows,omitempty" query:"rows,omitempty" validate:"required"`
	Columns   int       `json:"columns,omitempty" form:"columns,omitempty" query:"columns,omitempty" validate:"required"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"-"`
}

type Seat struct {
	ID        string    `json:"id,omitempty" form:"id,omitempty" query:"id,omitempty" gorm:"type:char(36)"`
	Name      string    `json:"name" gorm:"type:varchar(50)" validate:"required|string"`
	SeatType  string    `json:"seat_type" gorm:"type:varchar(50)" validate:"required|string"` // student,teacher
	ClassID   int       `json:"class.id,omitempty" gorm:"index:classroom_id"`
	ClassName string    `json:"class.name,omitempty" query:"class.name,omitempty" gorm:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"-"`
}

type Learning struct {
	ID        int       `json:"id,omitempty" form:"id,omitempty" query:"id,omitempty"`
	SeatID    int       `json:"seat.id,omitempty" gorm:"type:char(36);index:learning_seat_id"`
	UserID    int       `json:"user.id,omitempty" gorm:"type:char(36);index:learning_user_id"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"-"`
}

type OccupiedSeat struct {
	Seat        string
	StudentName string
}

type CheckInResponse struct {
	ClassID       int            `json:"class_id,omitempty" name:"class_id"`
	Rows          int            `json:"rows,omitempty"`
	Columns       int            `json:"columns,omitempty"`
	Teacher       string         `json:"teacher,omitempty"`
	AvailableSeat []string       `json:"available_seats"`
	OccupiedSeat  []OccupiedSeat `json:"occupied_seats"`
	Message       string         `json:"message,omitempty"`
}

var Abjad string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func (cr *Classroom) CreateClass(ctx helpers.Context) map[string]interface{} {
	isValid, msg := helpers.Validate(ctx, cr)
	if !isValid {
		return msg
	}
	c := Classroom{}
	helpers.GetDB(ctx).Where(Classroom{Name: cr.Name}).First(&c)
	if c.ID != 0 {
		return helpers.GeneralErrorMessage(400, "Class name `"+c.Name+"` has been exist", map[string]interface{}{})
	}
	helpers.GetDB(ctx).Create(cr)
	class := CheckInResponse{ClassID: cr.ID, Rows: cr.Rows, Columns: cr.Columns}
	return class.SetOccupiedClass(ctx)
}

func (cr *CheckInResponse) SetOccupiedClass(ctx helpers.Context) map[string]interface{} {
	CreateSeats(cr.Rows, cr.Columns, ctx)
	cr.SetAvailableSeats(ctx)
	cr.SetOccupiedSeats(ctx)
	return helpers.StructToMap(cr)
}

func (cr *CheckInResponse) SetAvailableSeats(ctx helpers.Context) {
	seats := []Seat{}
	cr.AvailableSeat = []string{}
	helpers.GetDB(ctx).
		Raw("Select s.name from seats as s where s.seat_type = 'student' and not exists(select 1 from learnings as l where l.seat_id = s.id)").
		Scan(&seats)
	for _, s := range seats {
		cr.AvailableSeat = append(cr.AvailableSeat, s.Name)
	}
}

func (cr *CheckInResponse) SetOccupiedSeats(ctx helpers.Context) {
	seats := []OccupiedSeat{}
	teacher := Seat{}
	cr.OccupiedSeat = []OccupiedSeat{}
	helpers.GetDB(ctx).
		Raw("Select s.name,u.name as student_name from seats as s inner join learnings as l on l.seat_id = s.id inner join users as u on u.id = l.user_id where s.seat_type = 'student'").
		Scan(&seats)
	helpers.GetDB(ctx).
		Raw("Select s.name from seats as s inner join learnings as l on l.seat_id = s.id where s.seat_type = 'teacher'").
		Scan(&teacher)
	for _, s := range seats {
		cr.OccupiedSeat = append(cr.OccupiedSeat, s)
	}
	if teacher.Name != "" {
		cr.Teacher = "in"
	} else {
		cr.Teacher = "out"
	}
}

func CreateSeats(rows, cols int, ctx helpers.Context) {
	var seats []string = make([]string, 0)
	for i := 0; i < cols; i++ {
		for j := 1; j <= rows; j++ {
			seats = append(seats, fmt.Sprintf("%d%s", j, string(Abjad[i])))
		}
	}
	for _, seat := range seats {
		s := Seat{ID: helpers.NewUUID(), Name: seat, SeatType: "student"}
		helpers.GetDB(ctx).Create(s)
	}
	s := Seat{ID: helpers.NewUUID(), Name: "Teacher", SeatType: "teacher"}
	helpers.GetDB(ctx).Create(s)
}

func BookSeat(ctx helpers.Context, clasId int, userId, inOut string) map[string]interface{} {
	class := CheckInResponse{ClassID: clasId}
	return class.SetOccupiedClass(ctx)
}
