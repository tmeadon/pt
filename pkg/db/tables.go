package db

type tableMap struct {
	SchemaVersionTable      string
	MusclesTable            string
	EquipmentTable          string
	ExerciseCategoriesTable string
	ExercisesTable          string
	ExerciseEquipmentTable  string
	ExerciseMusclesTable    string
}

var tables = tableMap{
	SchemaVersionTable:      "schema_version",
	MusclesTable:            "muscles",
	EquipmentTable:          "equipment",
	ExerciseCategoriesTable: "exercise_categories",
	ExercisesTable:          "exercises",
	ExerciseEquipmentTable:  "exercise_equipment",
	ExerciseMusclesTable:    "exercise_muscles",
}

// var tableTemplates []string = []string{
// 	`CREATE TABLE IF NOT EXISTS {{.MusclesTable}} (
// 		id INTEGER NOT NULL PRIMARY KEY,
// 		name TEXT NOT NULL,
// 		simple_name TEXT NOT NULL,
// 		is_front BOOLEAN
// 	);`,
// 	`CREATE TABLE IF NOT EXISTS {{.EquipmentTable}} (
// 		id INTEGER NOT NULL PRIMARY KEY,
// 		name TEXT NOT NULL
// 	);`,
// 	`CREATE TABLE IF NOT EXISTS {{.ExerciseCategoriesTable}} (
// 		id INTEGER NOT NULL PRIMARY KEY,
// 		name TEXT NOT NULL
// 	);`,
// 	`CREATE TABLE IF NOT EXISTS {{.ExercisesTable}} (
// 		id INTEGER NOT NULL PRIMARY KEY,
// 		name TEXT NOT NULL,
// 		description TEXT,
// 		category_id INTEGER NOT NULL,
// 		FOREIGN KEY (category_id) REFERENCES {{.ExerciseCategoriesTable}} (id)
// 	);`,
// 	`CREATE TABLE IF NOT EXISTS {{.ExerciseEquipmentTable}} (
// 		id INTEGER NOT NULL PRIMARY KEY,
// 		exercise_id INTEGER NOT NULL,
// 		equipment_id INTEGER NOT NULL,
// 		FOREIGN KEY (exercise_id) REFERENCES {{.ExercisesTable}} (id),
// 		FOREIGN KEY (equipment_id) REFERENCES {{.EquipmentTable}} (id)
// 	);`,
// 	`CREATE TABLE IF NOT EXISTS {{.ExerciseMusclesTable}} (
// 		id INTEGER NOT NULL PRIMARY KEY,
// 		exercise_id INTEGER NOT NULL,
// 		muscle_id INTEGER NOT NULL,
// 		is_primary BOOLEAN,
// 		FOREIGN KEY (exercise_id) REFERENCES {{.ExercisesTable}} (id),
// 		FOREIGN KEY (muscle_id) REFERENCES {{.MusclesTable}} (id)
// 	);`,
// }

// func CreateTables() error {
// 	for _, t := range tableTemplates {
// 		cmd, err := buildTableCommand(t)
// 		if err != nil {
// 			return err
// 		}

// 		stmt, err := DB.Prepare(cmd)
// 		if err != nil {
// 			return err
// 		}

// 		_, err = stmt.Exec()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
