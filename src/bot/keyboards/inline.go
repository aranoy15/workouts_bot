package keyboards

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func CreateWorkoutKeyboard(workoutID uint) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				WorkoutStart,
				fmt.Sprintf("workout:start:%d", workoutID),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				WorkoutEdit,
				fmt.Sprintf("workout:edit:%d", workoutID),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				WorkoutDelete,
				fmt.Sprintf("workout:delete:%d", workoutID),
			),
		),
	)

	return keyboard
}

func CreateExerciseKeyboard(exerciseID uint) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseDetails,
				fmt.Sprintf("exercise:details:%d", exerciseID),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseVideo,
				fmt.Sprintf("exercise:video:%d", exerciseID),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseAdd,
				fmt.Sprintf("exercise:add:%d", exerciseID),
			),
		),
	)

	return keyboard
}

func CreateSetKeyboard(
	workoutExerciseID uint,
	setNumber int,
) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				SetComplete,
				fmt.Sprintf("set:complete:%d:%d", workoutExerciseID, setNumber),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				SetSkip,
				fmt.Sprintf("set:skip:%d:%d", workoutExerciseID, setNumber),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				SetPause,
				fmt.Sprintf("set:pause:%d:%d", workoutExerciseID, setNumber),
			),
		),
	)

	return keyboard
}

func CreateGoalSelectionKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				GoalMuscleGain,
				"goal:muscle_gain",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				GoalStrength,
				"goal:strength",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				GoalEndurance,
				"goal:endurance",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				GoalWeightLoss,
				"goal:weight_loss",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavBack,
				"settings:main",
			),
		),
	)

	return keyboard
}

func CreateEquipmentKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				EquipmentHome,
				"equipment:home",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				EquipmentGym,
				"equipment:gym",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				EquipmentNone,
				"equipment:none",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				EquipmentCustom,
				"equipment:custom",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavBack,
				"settings:main",
			),
		),
	)

	return keyboard
}

func CreateWorkoutTypeKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				WorkoutTypeSplit,
				"workout_type:split",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				WorkoutTypePushPull,
				"workout_type:push_pull",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				WorkoutTypeFullBody,
				"workout_type:fullbody",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				WorkoutTypeCustom,
				"workout_type:custom",
			),
		),
	)

	return keyboard
}

func CreateMyWorkoutsKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				WorkoutRefresh,
				"workouts:refresh",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				WorkoutStats,
				"workouts:stats",
			),
		),
	)

	return keyboard
}

func CreateExercisesKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseCompound,
				"exercises:compound",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseIsolation,
				"exercises:isolation",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseStrength,
				"exercises:strength",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseCardio,
				"exercises:cardio",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseBodyweight,
				"exercises:bodyweight",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseHIIT,
				"exercises:hiit",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExerciseEndurance,
				"exercises:endurance",
			),
		),
	)

	return keyboard
}

func CreateSettingsKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				SettingsGoals,
				"settings:goals",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				SettingsEquipment,
				"settings:equipment",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				SettingsExperience,
				"settings:experience",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				SettingsLimitations,
				"settings:limitations",
			),
		),
	)

	return keyboard
}

func CreateExperienceLevelKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExpBeginner,
				"experience:1",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				ExpIntermediate,
				"experience:3",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				ExpAdvanced,
				"experience:5",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				ExpExpert,
				"experience:7",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavBack,
				"settings:main",
			),
		),
	)

	return keyboard
}

func CreateMuscleGroupKeyboard() tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				MuscleChest,
				"muscle:chest",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				MuscleBack,
				"muscle:back",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				MuscleShoulders,
				"muscle:shoulders",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				MuscleBiceps,
				"muscle:biceps",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				MuscleTriceps,
				"muscle:triceps",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				MuscleLegs,
				"muscle:legs",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				MuscleGlutes,
				"muscle:glutes",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				MuscleAbs,
				"muscle:abs",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavBack,
				"exercises:main",
			),
		),
	)

	return keyboard
}

func CreateWorkoutDurationKeyboard(workoutType string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				Duration30,
				fmt.Sprintf("duration:%s:30", workoutType),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				Duration45,
				fmt.Sprintf("duration:%s:45", workoutType),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				Duration60,
				fmt.Sprintf("duration:%s:60", workoutType),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				Duration90,
				fmt.Sprintf("duration:%s:90", workoutType),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavBack,
				"workout_type:main",
			),
		),
	)

	return keyboard
}

func CreateConfirmationKeyboard(action string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavYes,
				fmt.Sprintf("confirm:%s:yes", action),
			),
			tgbotapi.NewInlineKeyboardButtonData(
				NavNo,
				fmt.Sprintf("confirm:%s:no", action),
			),
		),
	)

	return keyboard
}

func CreateBackKeyboard(callbackData string) tgbotapi.InlineKeyboardMarkup {
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				NavBack,
				callbackData,
			),
		),
	)

	return keyboard
}
