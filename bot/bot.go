package bot

import (
	"NobatApi/config"
	"NobatApi/model"
	"NobatApi/uscase"
	"log"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotService struct {
	Bot    *tgbotapi.BotAPI
	Config *config.Config
}

func NewBotService(cfg *config.Config) *BotService {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Println("⚠️ ربات تلگرام متصل نشد:", err)
		log.Println("ℹ️ پنل وب بدون ربات ادامه میدهد")
		return &BotService{
			Bot:    nil,
			Config: cfg,
		}
	}

	return &BotService{
		Bot:    bot,
		Config: cfg,
	}
}

func (b *BotService) Start() {
	if b.Bot == nil {
		log.Println("ℹ️ ربات غیرفعال است (اتصال به تلگرام ممکن نشد)")
		return
	}

	log.Printf("✅ ربات تلگرام روشن شد: %s", b.Bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		go b.handleMessage(update.Message)
	}
}

// وضعیت کاربر در حین گفتگو
type UserSession struct {
	State string
	Name  string
	Date  string
}

var userSessions = make(map[int64]*UserSession)

func (b *BotService) handleMessage(msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	// اگر سشن وجود نداشت، بساز
	if userSessions[chatID] == nil {
		userSessions[chatID] = &UserSession{}
	}

	switch msg.Text {
	case "/start":
		userSessions[chatID].State = "waiting_name"
		b.sendMessage(chatID, "سلام! 👋\nبه ربات نوبت‌دهی خوش آمدید.\nلطفاً نام و نام خانوادگی خود را وارد کنید:")
	default:
		b.handleState(chatID, msg.Text)
	}
}

func (b *BotService) handleState(chatID int64, text string) {
	session := userSessions[chatID]
	if session == nil {
		return
	}

	switch session.State {
	case "waiting_name":
		session.Name = text
		session.State = "waiting_date"
		b.sendMessage(chatID, "عالی!\nحالا تاریخ نوبت را وارد کنید (مثلاً 1405/06/12):")

	case "waiting_date":
		session.Date = text
		session.State = "waiting_time"
		b.sendMessage(chatID, "حالا ساعت نوبت را وارد کنید (مثلاً 14:30):")

	case "waiting_time":
		name := session.Name
		date := session.Date

		// ساخت نوبت
		appt := &model.Appointment{
			CustomerName:   name,
			CustomerChatID: chatID,
			Date:           date,
			Time:           text,
			Status:         "pending",
			CreatedAt:      time.Now(),
		}

		// ذخیره نوبت
		err := uscase.AppointmentUse.SaveAppointment(appt)
		if err != nil {
			b.sendMessage(chatID, "❌ خطا در ثبت نوبت. لطفاً دوباره تلاش کنید.")
			delete(userSessions, chatID)
			return
		}

		b.sendMessage(chatID, "✅ نوبت شما با موفقیت ثبت شد!\n\n📅 تاریخ: "+date+"\n🕐 ساعت: "+text+"\n\nمنتظر تأیید نهایی از طرف ما باشید.")

		// حذف سشن بعد از اتمام
		delete(userSessions, chatID)
	}
}

func (b *BotService) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	b.Bot.Send(msg)
}

// برای تبدیل chatID به string اگر جایی لازم شد
func chatIDToString(chatID int64) string {
	return strconv.FormatInt(chatID, 10)
}
