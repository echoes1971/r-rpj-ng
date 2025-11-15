package dblayer

/*
CREATE TABLE `rprj_dbversion` (

	`model_name` varchar(255) NOT NULL,
	`version` int(11) NOT NULL,
	PRIMARY KEY (`model_name`),
	KEY `rprj_dbversion_0` (`model_name`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
type DBVersion struct {
	DBEntity
}

func NewDBVersion() *DBVersion {
	columns := []Column{
		{Name: "model_name", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "version", Type: "int(11)", Constraints: []string{"NOT NULL"}},
	}
	keys := []string{"model_name"}
	return &DBVersion{
		DBEntity: *NewDBEntity(
			"DBVersion",
			"dbversion",
			columns,
			keys,
			[]ForeignKey{},
			make(map[string]any),
		),
	}
}

/*
CREATE TABLE `rprj_users` (

	`id` varchar(16) NOT NULL,
	`login` varchar(255) NOT NULL,
	`pwd` varchar(255) NOT NULL,
	`pwd_salt` varchar(4) DEFAULT '',
	`fullname` text DEFAULT NULL,
	`group_id` varchar(16) NOT NULL,
	PRIMARY KEY (`id`),
	KEY `rprj_users_0` (`id`),
	KEY `rprj_users_1` (`group_id`)

) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
*/
type DBUser struct {
	DBEntity
}

func NewDBUser() *DBUser {
	columns := []Column{
		{Name: "id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
		{Name: "login", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "pwd", Type: "varchar(255)", Constraints: []string{"NOT NULL"}},
		{Name: "pwd_salt", Type: "varchar(4)", Constraints: []string{}},
		{Name: "fullname", Type: "text", Constraints: []string{}},
		{Name: "group_id", Type: "varchar(16)", Constraints: []string{"NOT NULL"}},
	}
	keys := []string{"id"}
	foreignKeys := []ForeignKey{
		{Column: "group_id", RefTable: "groups", RefColumn: "id"},
	}
	return &DBUser{
		DBEntity: *NewDBEntity(
			"DBUser",
			"users",
			columns,
			keys,
			foreignKeys,
			make(map[string]any),
		),
	}
}
