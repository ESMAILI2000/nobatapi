package uscase

import (
	"NobatApi/model"
	"NobatApi/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppointmentUsecase interface {
	SaveAppointment(a *model.Appointment) error
	GetAppointmentsByDate(date string) ([]model.Appointment, error)
	IsSlotAvailable(date, timeSlot string) (bool, error)
	UpdateStatus(id primitive.ObjectID, status string) error
	GetAllAppointments() ([]model.Appointment, error)
}

var AppointmentUse *appointmentUsecaseImpl

type appointmentUsecaseImpl struct {
	appointmentRepository repository.AppointmentRepository
}

func NewAppointmentUsecase(r repository.AppointmentRepository) *appointmentUsecaseImpl {
	return &appointmentUsecaseImpl{r}
}

func (u *appointmentUsecaseImpl) SaveAppointment(a *model.Appointment) error {
	// تنظیم زمان ایجاد
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	return u.appointmentRepository.Save(a)
}

func (u *appointmentUsecaseImpl) GetAppointmentsByDate(date string) ([]model.Appointment, error) {
	return u.appointmentRepository.FindByDate(date)
}

func (u *appointmentUsecaseImpl) IsSlotAvailable(date, timeSlot string) (bool, error) {
	taken, err := u.appointmentRepository.IsSlotTaken(date, timeSlot)
	if err != nil {
		return false, err
	}
	return !taken, nil
}

func (u *appointmentUsecaseImpl) UpdateStatus(id primitive.ObjectID, status string) error {
	return u.appointmentRepository.UpdateStatus(id, status)
}

func (u *appointmentUsecaseImpl) GetAllAppointments() ([]model.Appointment, error) {
	return u.appointmentRepository.FindAll()
}
