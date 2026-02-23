CREATE SCHEMA IF NOT EXISTS workouts;

CREATE TABLE IF NOT EXISTS workouts.workout_types (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS workouts.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    telegram_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    equipment_ids INTEGER[],
    goals JSONB DEFAULT '[]',
    experience INTEGER DEFAULT 1,
    limitations TEXT[],
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS workouts.equipments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL,
    is_home_gym BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS workouts.exercises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(255),
    muscle_groups TEXT[],
    difficulty INTEGER DEFAULT 1,
    duration_minutes INTEGER,
    image_path VARCHAR(255),
    video_path VARCHAR(255),
    equipment_ids INTEGER[],
    instructions TEXT[],
    common_mistakes TEXT[],
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS workouts.workouts (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES workouts.users(id),
    name VARCHAR(255),
    workout_type_id INTEGER NOT NULL REFERENCES workouts.workout_types(id),
    duration_minutes INTEGER,
    status VARCHAR(255) DEFAULT 'planned',
    workout_type VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS workouts.workout_exercises (
    id SERIAL PRIMARY KEY,
    workout_id INTEGER NOT NULL REFERENCES workouts.workouts(id) ON DELETE CASCADE,
    exercise_id INTEGER NOT NULL REFERENCES workouts.exercises(id),
    order_index INTEGER,
    sets_count INTEGER,
    reps_count INTEGER,
    rest_seconds INTEGER,
    weight_kg DOUBLE PRECISION,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS workouts.sets (
    id SERIAL PRIMARY KEY,
    workout_exercise_id INTEGER NOT NULL REFERENCES workouts.workout_exercises(id) ON DELETE CASCADE,
    set_number INTEGER,
    weight_kg DOUBLE PRECISION,
    reps_done INTEGER,
    rest_taken_seconds INTEGER,
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS workouts.weight_histories (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES workouts.users(id),
    exercise_id INTEGER NOT NULL REFERENCES workouts.exercises(id),
    weight_kg DOUBLE PRECISION NOT NULL,
    reps_count INTEGER NOT NULL,
    recorded_at TIMESTAMP WITH TIME ZONE
);
