package handlers

import (
	"NobatApi/model"
	"NobatApi/uscase"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ساختار ساده برای نمایش در تمپلیت (ID به صورت رشته)
type ApptView struct {
	ID           string
	CustomerName string
	Time         string
	Status       string
	Source       string
}

// GET /dashboard
func Dashboard(c *fiber.Ctx) error {
	today := time.Now().Format("2006-01-02")

	appts, err := uscase.AppointmentUse.GetAppointmentsByDate(today)
	if err != nil {
		return c.Render("dashboard", fiber.Map{
			"TodayCount":   0,
			"Appointments": []ApptView{},
			"Today":        today,
		})
	}

	// تبدیل به ساختار نمایشی
	var views []ApptView
	for _, a := range appts {
		views = append(views, ApptView{
			ID:           a.ID.Hex(),
			CustomerName: a.CustomerName,
			Time:         a.Time,
			Status:       a.Status,
			Source:       a.Source,
		})
	}

	return c.Render("dashboard", fiber.Map{
		"TodayCount":   len(views),
		"Appointments": views,
		"Today":        today,
	})
}

// POST /appointment/create
func CreateAppointment(c *fiber.Ctx) error {
	name := c.FormValue("customer_name")
	date := c.FormValue("date")
	timeSlot := c.FormValue("time")
	source := c.FormValue("source")

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}
	if source == "" {
		source = "walkin"
	}

	// چک تداخل
	available, err := uscase.AppointmentUse.IsSlotAvailable(date, timeSlot)
	if err != nil {
		return c.Render("dashboard", fiber.Map{
			"Error": "خطا در بررسی تداخل زمانی",
		})
	}
	if !available {
		return c.Render("dashboard", fiber.Map{
			"Error": "این ساعت قبلاً رزرو شده است",
		})
	}

	appt := &model.Appointment{
		CustomerName: name,
		Date:         date,
		Time:         timeSlot,
		Status:       "confirmed",
		Source:       source,
		CreatedAt:    time.Now(),
	}

	err = uscase.AppointmentUse.SaveAppointment(appt)
	if err != nil {
		return c.Render("dashboard", fiber.Map{
			"Error": "خطا در ثبت نوبت",
		})
	}

	return c.Redirect("/dashboard")
}

// POST /appointment/update-status
func UpdateAppointmentStatus(c *fiber.Ctx) error {
	idStr := c.FormValue("id")
	status := c.FormValue("status")

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(400).SendString("ID نامعتبر")
	}

	err = uscase.AppointmentUse.UpdateStatus(id, status)
	if err != nil {
		return c.Status(500).SendString("خطا در تغییر وضعیت")
	}

	return c.Redirect("/dashboard")
}
