package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateWorkoutKeyboard(workoutID uint) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("▶️ Начать", fmt.Sprintf("workout:start:%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("✏️ Редактировать", fmt.Sprintf("workout:edit:%d", workoutID)),
			tgbotapi.NewInlineKeyboardButtonData("🗑️ Удалить", fmt.Sprintf("workout:delete:%d", workoutID)),
		),
	)

	return keyboard
}

func CreateExerciseKeyboard(exerciseID uint) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📖 Подробнее", fmt.Sprintf("exercise:details:%d", exerciseID)),
			tgbotapi.NewInlineKeyboardButtonData("🎥 Видео", fmt.Sprintf("exercise:video:%d", exerciseID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("➕ Добавить в тренировку", fmt.Sprintf("exercise:add:%d", exerciseID)),
		),
	)

	return keyboard
}

func CreateSetKeyboard(workoutExerciseID uint, setNumber int) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅ Завершить подход", fmt.Sprintf("set:complete:%d:%d", workoutExerciseID, setNumber)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⏭️ Пропустить", fmt.Sprintf("set:skip:%d:%d", workoutExerciseID, setNumber)),
			tgbotapi.NewInlineKeyboardButtonData("⏸️ Пауза", fmt.Sprintf("set:pause:%d:%d", workoutExerciseID, setNumber)),
		),
	)

	return keyboard
}

func CreateGoalSelectionKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("💪 Набор массы", "goal:muscle_gain"),
			tgbotapi.NewInlineKeyboardButtonData("🏋️ Сила", "goal:strength"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏃 Выносливость", "goal:endurance"),
			tgbotapi.NewInlineKeyboardButtonData("🔥 Похудение", "goal:weight_loss"),
		),
	)

	return keyboard
}

func CreateEquipmentKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 Домашний зал", "equipment:home"),
			tgbotapi.NewInlineKeyboardButtonData("🏋️ Тренажерный зал", "equipment:gym"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚫 Без оборудования", "equipment:none"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Настроить", "equipment:custom"),
		),
	)

	return keyboard
}
