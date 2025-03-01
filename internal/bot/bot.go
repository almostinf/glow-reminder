package bot

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "time/tzdata"

	"github.com/almostinf/glow-reminder/internal/domain"
	"github.com/almostinf/glow-reminder/internal/usecase"
	"github.com/almostinf/glow-reminder/pkg/clock"
	"github.com/almostinf/glow-reminder/pkg/logger"
	"github.com/google/uuid"
	telebot "gopkg.in/telebot.v4"
)

var (
	// Universal markup builders.
	menu           = &telebot.ReplyMarkup{ResizeKeyboard: true}
	paginationMenu = &telebot.ReplyMarkup{}
	colourMenu     = &telebot.ReplyMarkup{}
	effectMenu     = &telebot.ReplyMarkup{}

	// Reply buttons.
	btnHelp          = menu.Text("ℹ️ Help")
	btnAddReminder   = menu.Text("➕ New Reminder")
	btnListReminders = menu.Text("📂 Reminders List")

	// Pagination buttons.
	btnPrev = paginationMenu.Data("⬅️", "pagination_prev")
	btnNext = paginationMenu.Data("➡️", "pagination_next")

	// Inline buttons.
	btnColourRed      = colourMenu.Data("🔴 Red", "colour_red")
	btnColourBlue     = colourMenu.Data("🔵 Blue", "colour_blue")
	btnColourGreen    = colourMenu.Data("🟢 Green", "colour_green")
	btnEffectStatic   = effectMenu.Data("🗿 Static", "effect_static")
	btnEffectBlinking = effectMenu.Data("✨ Blinking", "effect_blinking")

	startMsg               = "👋 Hello! It's a reminder bot"
	choosingTimeMsg        = "🚀 Please enter the time in format 'YYYY-MM-DD HH:MM'"
	tryAgainAddReminderMsg = "⚠️ Please start by clicking ➕ button"
	tryAgainMsg            = "⚠️ Please try again"
	helpMsg                = "Help:\n" +
		"- Use the 📂 button to view scheduled reminders\n" +
		"- Use the ➕ button to create a new reminder\n" +
		"- Use the 🗑 button to delete existing reminder\n" +
		"- Use the ⬅️ and ➡️ buttons to scroll through the list of reminders"

	timeFormat = "2006-01-02 15:04"
	limit      = int64(5)
)

var _ Bot = (*bot)(nil)

type Bot interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type bot struct {
	*telebot.Bot

	userStates      map[int64]*userState
	m               *sync.RWMutex
	logger          logger.Logger
	reminderUsecase usecase.ReminderUsecase
	clock           clock.Clock
}

func New(cfg Config, logger logger.Logger, reminderUsecase usecase.ReminderUsecase, clock clock.Clock) (*bot, error) {
	tbot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.Token,
		Poller: &telebot.LongPoller{Timeout: cfg.PollerTimeout},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create new telebot: %w", err)
	}

	menu.Reply(
		menu.Row(btnListReminders, btnAddReminder, btnHelp),
	)

	paginationMenu.Inline(
		paginationMenu.Row(btnPrev, btnNext),
	)

	colourMenu.Inline(
		colourMenu.Row(btnColourRed, btnColourGreen, btnColourBlue),
	)

	effectMenu.Inline(
		effectMenu.Row(btnEffectStatic, btnEffectBlinking),
	)

	return &bot{
		Bot: tbot,

		userStates:      make(map[int64]*userState),
		m:               &sync.RWMutex{},
		logger:          logger,
		reminderUsecase: reminderUsecase,
		clock:           clock,
	}, nil
}

func (b *bot) setUserState(userID int64, us *userState) {
	b.m.Lock()
	defer b.m.Unlock()

	b.userStates[userID] = us
}

func (b *bot) getUserState(userID int64) (*userState, bool) {
	b.m.RLock()
	defer b.m.RUnlock()

	us, ok := b.userStates[userID]
	return us, ok
}

func (b *bot) Start(_ context.Context) error {
	b.Handle("/start", b.handleStart())
	b.Handle(&btnHelp, b.handleHelp())
	b.Handle(&btnAddReminder, b.handleAddReminder())
	b.Handle(telebot.OnText, b.handleText())
	b.Handle(telebot.OnCallback, b.handleCallback())
	b.Handle(&btnListReminders, b.handleListReminders())

	go func() {
		b.Bot.Start()
	}()

	return nil
}

func (b *bot) Stop(_ context.Context) error {
	b.Bot.Stop()
	return nil
}

func (b *bot) handleStart() func(c telebot.Context) error {
	return func(c telebot.Context) error {
		b.setUserState(c.Sender().ID, &userState{
			s: menuState,
		})
		return c.Send(startMsg, menu)
	}
}

func (b *bot) handleHelp() func(c telebot.Context) error {
	return func(c telebot.Context) error {
		b.setUserState(c.Sender().ID, &userState{
			s: menuState,
		})
		return c.Send(helpMsg)
	}
}

func (b *bot) handleAddReminder() func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		b.setUserState(c.Sender().ID, &userState{
			s: timeChoosingState,
			reminder: domain.Reminder{
				UserID: userID,
			},
		})

		return c.Send(choosingTimeMsg)
	}
}

func (b *bot) handleText() func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userID := c.Sender().ID
		us, ok := b.getUserState(userID)
		if !ok {
			us.s = menuState
			b.setUserState(userID, us)
			return c.Send(tryAgainAddReminderMsg)
		}

		switch us.s {
		case timeChoosingState:
			return b.handleChoosingTime(c)
		case textEnteringState:
			return b.handleTextEntering(c)
		default:
			us.s = menuState
			b.setUserState(userID, us)
			return c.Send(tryAgainAddReminderMsg)
		}
	}
}

func (b *bot) handleCallback() func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userID := c.Sender().ID
		us, ok := b.getUserState(userID)
		if !ok {
			us.s = menuState
			b.setUserState(userID, us)
			return c.Send(tryAgainMsg)
		}

		switch us.s {
		case colourChoosingState:
			return b.handleChoosingColour(c)
		case effectChoosingState:
			return b.handleChoosingEffect(c)
		case listRemindersState:
			return b.waitingListReminders(c)
		default:
			us.s = menuState
			b.setUserState(userID, us)
			return c.Send(tryAgainMsg)
		}
	}
}

func (b *bot) handleChoosingTime(c telebot.Context) error {
	userID := c.Sender().ID
	us, ok := b.getUserState(userID)
	if !ok || us.s != timeChoosingState {
		us.s = menuState
		b.setUserState(userID, us)
		return c.Send(tryAgainAddReminderMsg)
	}

	mskLocation, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		us.s = menuState
		b.setUserState(userID, us)
		b.logger.Error("Failed to load Moscow location", map[string]interface{}{
			"error": err.Error(),
		})
		return c.Send("❌ Failed to load Moscow location")
	}

	timeStr := c.Text()
	parsedTime, err := time.ParseInLocation(timeFormat, timeStr, mskLocation)
	if err != nil {
		us.s = menuState
		b.setUserState(userID, us)
		return c.Send("❌ Invalid time format. Please use 'YYYY-MM-DD HH:MM'")
	}

	us.reminder.ScheduledAt = parsedTime
	us.s = textEnteringState

	b.setUserState(userID, us)

	return c.Send("🚀 Please enter a reminder text")
}

func (b *bot) handleTextEntering(c telebot.Context) error {
	userID := c.Sender().ID
	us, ok := b.getUserState(userID)
	if !ok || us.s != textEnteringState {
		us.s = menuState
		b.setUserState(userID, us)
		return c.Send(tryAgainAddReminderMsg)
	}

	msg := c.Text()

	us.reminder.Msg = msg
	us.s = colourChoosingState

	b.setUserState(userID, us)

	return c.Send("🚀 Choose an effect colour", colourMenu)
}

func (b *bot) handleChoosingColour(c telebot.Context) error {
	userID := c.Sender().ID
	us, ok := b.getUserState(userID)
	if !ok || us.s != colourChoosingState {
		us.s = menuState
		b.setUserState(userID, us)
		return c.Send(tryAgainAddReminderMsg)
	}

	var colour domain.Colour
	switch c.Callback().Data {
	case "\fcolour_red":
		colour = domain.Red
	case "\fcolour_green":
		colour = domain.Green
	case "\fcolour_blue":
		colour = domain.Blue
	default:
		us.s = menuState
		b.setUserState(userID, us)
		b.logger.Error("invalid colour choosing", map[string]interface{}{
			"user_id":       userID,
			"callback_date": c.Callback().Unique,
		})
		return c.Send("❌ Invalid colour. Please choose red, green or blue")
	}

	us.reminder.Colour = colour
	us.s = effectChoosingState

	b.setUserState(userID, us)

	return c.Send("🚀 Choose an effect mode", effectMenu)
}

func (b *bot) handleChoosingEffect(c telebot.Context) error {
	userID := c.Sender().ID
	us, ok := b.getUserState(userID)
	if !ok || us.s != effectChoosingState {
		us.s = menuState
		b.setUserState(userID, us)
		return c.Send(tryAgainAddReminderMsg)
	}

	var mode domain.Mode
	switch c.Callback().Data {
	case "\feffect_static":
		mode = domain.Static
	case "\feffect_blinking":
		mode = domain.Blinking
	default:
		us.s = menuState
		b.setUserState(userID, us)
		b.logger.Error("invalid effect choosing", map[string]interface{}{
			"user_id":       userID,
			"callback_date": c.Callback().Data,
		})
		return c.Send("❌ Invalid mode. Please choose blinking or static")
	}

	us.reminder.Mode = mode
	us.reminder.ID = uuid.New()
	us.reminder.CreatedAt = b.clock.NowUTC()
	us.reminder.UpdatedAt = b.clock.NowUTC()
	us.s = menuState

	b.setUserState(userID, us)

	if err := b.reminderUsecase.CreateReminder(context.TODO(), us.reminder); err != nil {
		us.s = menuState
		b.setUserState(userID, us)
		b.logger.Error("failed to CreateReminder", map[string]interface{}{
			"user_id":  userID,
			"reminder": us.reminder,
			"err":      err.Error(),
		})
		return c.Send(tryAgainAddReminderMsg)
	}

	return c.Send("✅ Reminder successfully created")
}

func (b *bot) handleListReminders() func(c telebot.Context) error {
	return func(c telebot.Context) error {
		userID := c.Sender().ID

		b.setUserState(userID, &userState{
			s: listRemindersState,
		})

		return b.listReminders(context.TODO(), c, domain.GetRemindersParams{
			UserID: userID,
			Limit:  uint64(limit),
		})
	}
}

func (b *bot) listReminders(ctx context.Context, c telebot.Context, params domain.GetRemindersParams) error {
	reminders, err := b.reminderUsecase.GetReminders(ctx, domain.GetRemindersParams{
		UserID: params.UserID,
		Offset: params.Offset,
		Limit:  params.Limit,
	})
	if err != nil {
		b.setUserState(params.UserID, &userState{
			s: menuState,
		})
		b.logger.Error("failed to GetReminders", map[string]interface{}{
			"user_id":   params.UserID,
			"reminders": reminders,
			"err":       err.Error(),
		})
		return c.Send(tryAgainAddReminderMsg)
	}

	for _, reminder := range reminders {
		var colour string
		switch reminder.Colour {
		case domain.Red:
			colour = "🔴 Red"
		case domain.Green:
			colour = "🟢 Green"
		case domain.Blue:
			colour = "🔵 Blue"
		}

		var mode string
		switch reminder.Mode {
		case domain.Blinking:
			mode = "✨ Blinking"
		case domain.Static:
			mode = "🗿 Static"
		}

		reminderMsg := "🗓 Reminder\n" +
			"Message: " + reminder.Msg + "\n" +
			"Scheduled At: " + reminder.ScheduledAt.Format(timeFormat) + "\n" +
			"Colour: " + colour + "\n" +
			"Mode: " + mode

		deleteMenu := &telebot.ReplyMarkup{}
		btnDelete := deleteMenu.Data("🗑 Delete", fmt.Sprintf("delete_reminder:%s", reminder.ID))
		deleteMenu.Inline(deleteMenu.Row(btnDelete))

		if err = c.Send(reminderMsg, deleteMenu); err != nil {
			b.setUserState(params.UserID, &userState{
				s: menuState,
			})
			b.logger.Error("failed to Send reminderMsg", map[string]interface{}{
				"user_id":      params.UserID,
				"reminder_msg": reminderMsg,
				"err":          err.Error(),
			})
			return c.Send(tryAgainAddReminderMsg)
		}
	}

	return c.Send("⚙️ Control menu", paginationMenu)
}

func (b *bot) waitingListReminders(c telebot.Context) error {
	userID := c.Sender().ID
	us, ok := b.getUserState(userID)
	if !ok || us.s != listRemindersState {
		us.s = menuState
		b.setUserState(userID, us)
		return c.Send(tryAgainAddReminderMsg)
	}

	callbackSplitted := strings.Split(c.Callback().Data, ":")

	if len(callbackSplitted) > 1 {
		return b.deleteReminder(c, userID, us, callbackSplitted[1])
	}

	switch callbackSplitted[0] {
	case "\fpagination_prev":
		if us.offset-limit >= 0 {
			us.offset -= limit
		}
	case "\fpagination_next":
		us.offset += limit
	}

	b.setUserState(userID, us)

	return b.listReminders(context.TODO(), c, domain.GetRemindersParams{
		UserID: userID,
		Offset: uint64(us.offset),
		Limit:  uint64(limit),
	})
}

func (b *bot) deleteReminder(c telebot.Context, userID int64, us *userState, reminderID string) error {
	id, err := uuid.Parse(reminderID)
	if err != nil {
		us.s = menuState
		b.setUserState(userID, us)
		b.logger.Error("failed to parse uuid", map[string]interface{}{
			"user_id":     userID,
			"reminder_id": reminderID,
		})
		return c.Send(tryAgainAddReminderMsg)
	}

	if err = b.reminderUsecase.DeleteReminder(context.TODO(), id); err != nil {
		us.s = menuState
		b.setUserState(userID, us)
		b.logger.Error("failed to DeleteReminder", map[string]interface{}{
			"user_id":     userID,
			"reminder_id": id.String(),
			"err":         err.Error(),
		})
		return c.Send(tryAgainAddReminderMsg)
	}

	return c.Send("✅ Reminder successfully deleted")
}
