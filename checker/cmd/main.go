package main

import (
	"checker/repo"
	"checker/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("error intializing config: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env varibles: %s", err.Error())
	}
	db, err := repo.NewPostgresDB(repo.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("Fatal to connect to DB, because: %s", err.Error())
	}
	repos := repo.NewRepository(db)

	feelings := make(map[string]string)

	feelings["хорошо"] = "good"
	feelings["плохо"] = "bad"

	/*if err := repos.InitTable(feelings); err != nil {
		panic(err)
	}*/

	servicer := service.NewService(repos)

	servicer.ParseCSVFile("/Users/yaakovlew/Desktop/BICA/checker/example.csv")

}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
