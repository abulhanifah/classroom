package models

import (
	"fmt"
	"sort"
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
}

type Seat struct {
	ID        string    `json:"id,omitempty" form:"id,omitempty" query:"id,omitempty" gorm:"type:char(36)"`
	Name      string    `json:"name" gorm:"type:varchar(50)" validate:"required|string"`
	SeatType  string    `json:"seat_type" gorm:"type:varchar(50)" validate:"required|string"` // student,teacher
	ClassID   int       `json:"class.id,omitempty" gorm:"index:classroom_id"`
	ClassName string    `json:"class.name,omitempty" query:"class.name,omitempty" gorm:"-"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type Learning struct {
	ID        string    `json:"id,omitempty" form:"id,omitempty" query:"id,omitempty" gorm:"type:char(36)"`
	ClassID   int       `json:"class.id,omitempty" gorm:"type:char(36);index:learning_class_id"`
	ClassName string    `json:"class.name,omitempty" gorm:"-"`
	SeatID    string    `json:"seat.id,omitempty" gorm:"type:char(36);index:learning_seat_id"`
	SeatName  string    `json:"seat.name,omitempty" gorm:"-"`
	UserID    string    `json:"user.id,omitempty" gorm:"type:char(36);index:learning_user_id"`
	UserName  string    `json:"user.name,omitempty" gorm:"-"`
	ExpiredAt time.Time `json:"expired_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
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
	temp := []Classroom{}
	helpers.GetDB(ctx).Table("classrooms").Where("name = ? and (deleted_at is null or deleted_at ='0000-00-00 00:00:00')", cr.Name).Scan(&temp)
	if len(temp) > 0 {
		return helpers.GeneralErrorMessage(400, "Class name `"+cr.Name+"` has been exist", map[string]interface{}{})
	}
	helpers.GetDB(ctx).Create(cr)
	CreateSeats(cr.Rows, cr.Columns, cr.ID, ctx)
	class := CheckInResponse{ClassID: cr.ID, Rows: cr.Rows, Columns: cr.Columns}
	return class.SetOccupiedClass(ctx)
}

func (cr *CheckInResponse) SetOccupiedClass(ctx helpers.Context) map[string]interface{} {
	cr.SetAvailableSeats(ctx)
	cr.SetOccupiedSeats(ctx)
	return helpers.StructToMap(cr)
}

func (cr *CheckInResponse) SetAvailableSeats(ctx helpers.Context) {
	seats := []Seat{}
	cr.AvailableSeat = []string{}
	helpers.GetDB(ctx).
		Raw("Select s.name from seats as s where s.seat_type = 'student' and s.class_id = ? and not exists(select 1 from learnings as l where l.seat_id = s.id) order by s.name", cr.ClassID).
		Scan(&seats)
	for _, s := range seats {
		cr.AvailableSeat = append(cr.AvailableSeat, s.Name)
	}
	sort.Strings(cr.AvailableSeat)
}
func (cr *CheckInResponse) SetOccupiedSeats(ctx helpers.Context) {
	seats := []OccupiedSeat{}
	teacher := Seat{}
	cr.OccupiedSeat = []OccupiedSeat{}
	helpers.GetDB(ctx).
		Raw("Select s.name as seat,u.name as student_name from seats as s inner join learnings as l on l.seat_id = s.id inner join users as u on u.id = l.user_id where s.seat_type = 'student' and l.class_id = ?", cr.ClassID).
		Scan(&seats)
	helpers.GetDB(ctx).
		Raw("Select s.name from seats as s inner join learnings as l on l.seat_id = s.id where s.seat_type = 'teacher' and s.class_id = ?", cr.ClassID).
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

func CreateSeats(rows, cols, classId int, ctx helpers.Context) {
	var seats []string = make([]string, 0)
	for i := 0; i < cols; i++ {
		for j := 1; j <= rows; j++ {
			seats = append(seats, fmt.Sprintf("%d%s", j, string(Abjad[i])))
		}
	}
	for _, seat := range seats {
		s := Seat{ID: helpers.NewUUID(), Name: seat, SeatType: "student", ClassID: classId}
		helpers.GetDB(ctx).Create(s)
	}
	s := Seat{ID: helpers.NewUUID(), Name: "Teacher", SeatType: "teacher", ClassID: classId}
	helpers.GetDB(ctx).Create(s)
}

func BookSeat(ctx helpers.Context, classId int, userId, inOut string) map[string]interface{} {
	class := CheckInResponse{ClassID: classId}
	cr := Classroom{}
	u := User{}
	helpers.GetDB(ctx).Table("classrooms").Where(Classroom{ID: classId}).Scan(&cr)
	helpers.GetDB(ctx).Table("users").Where(User{ID: userId}).Scan(&u)
	if u.ID == "" {
		return helpers.GeneralErrorMessage(400, "User id `"+userId+"` was not found", map[string]interface{}{})
	} else if cr.ID == 0 {
		return helpers.GeneralErrorMessage(400, fmt.Sprintf("Class id `%d` was not found", classId), map[string]interface{}{})
	} else {
		class.Rows = cr.Rows
		class.Columns = cr.Columns
		if u.RoleID == 3 { //3 = student
			class.SetAvailableSeats(ctx)
			if len(class.AvailableSeat) == 0 {
				class.Message = fmt.Sprintf("Hi %s, the class is fully seated", u.Name)
			} else {
				seat := Seat{}
				learning := Learning{}
				if inOut == "in" {
					helpers.GetDB(ctx).Table("seats").Where(Seat{Name: class.AvailableSeat[0]}).Scan(&seat)
					helpers.GetDB(ctx).Table("learnings").Where(Learning{UserID: u.ID, ClassID: classId, SeatID: seat.ID}).
						Joins("inner join seats on learnings.seat_id = seats.id").
						Joins("inner join users on learnings.user_id = users.id").
						Select("learnings.*,seats.name as seat_name,users.name as user_name").
						Scan(&learning)
					fmt.Println("masuk sini id", learning.ID)
					if learning.ID == "" {
						learning = Learning{SeatID: seat.ID, UserID: u.ID, ID: helpers.NewUUID(), ClassID: classId}
						helpers.GetDB(ctx).Create(learning)
						class.Message = fmt.Sprintf("Hi %s, your seat is %s", u.Name, seat.Name)
					} else {
						class.Message = fmt.Sprintf("Hi %s, your seat is %s", u.Name, learning.SeatName)
					}
				} else {
					helpers.GetDB(ctx).Table("learnings").Where(Learning{UserID: u.ID, ClassID: classId}).
						Joins("inner join seats on learnings.seat_id = seats.id").
						Joins("inner join users on learnings.user_id = users.id").
						Select("learnings.*,seats.name as seat_name,users.name as user_name").
						Scan(&learning)
					if learning.ID == "" {
						class.Message = fmt.Sprintf("Hi %s, you are not register in this class", u.Name)
					} else {
						helpers.GetDB(ctx).Where(Learning{UserID: u.ID, ClassID: classId}).Delete(learning)
						class.Message = fmt.Sprintf("Hi %s, %s is now available for other students", u.Name, learning.SeatName)
					}

				}
			}
		} else if u.RoleID == 2 {
			seat := Seat{}
			helpers.GetDB(ctx).Table("seats").Where(Seat{ClassID: classId, SeatType: "teacher"}).Scan(&seat)
			learning := Learning{}
			helpers.GetDB(ctx).Table("learnings").Where(Learning{UserID: u.ID, ClassID: classId}).
				Joins("inner join seats on learnings.seat_id = seats.id").
				Joins("inner join users on learnings.user_id = users.id").
				Select("learnings.*,seats.name as seat_name,users.name as user_name").
				Scan(&learning)
			if inOut == "in" {
				if learning.ID == "" {
					learning = Learning{SeatID: seat.ID, UserID: u.ID, ID: helpers.NewUUID(), ClassID: classId}
					helpers.GetDB(ctx).Create(learning)
				}
				class.Message = fmt.Sprintf("Hi %s, you has been set as teacher", u.Name)
			} else {
				if learning.ID != "" {
					helpers.GetDB(ctx).Where(Learning{UserID: u.ID, ClassID: classId}).Delete(learning)
				}
				class.Message = fmt.Sprintf("Hi %s, you has been unset as teacher", u.Name)
			}
		}
	}
	return cr.GetClassById(ctx, classId)
}

func (cr *Classroom) GetClassById(ctx helpers.Context, classId int) map[string]interface{} {
	helpers.GetDB(ctx).Table("classrooms").Where(Classroom{ID: classId}).Scan(&cr)
	class := CheckInResponse{ClassID: classId, Rows: cr.Rows, Columns: cr.Columns}
	return class.SetOccupiedClass(ctx)
}
