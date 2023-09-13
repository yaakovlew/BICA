package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"checker/models"
	"checker/repo"
)

type NNChatGPTService struct {
	repo repo.NNChatGPT
}

func NewNNChatGPTService(repo repo.NNChatGPT) *NNChatGPTService {
	return &NNChatGPTService{repo: repo}
}

func (s *NNChatGPTService) SendTOChatGPT(str, realAnswer string) error {
	msg := []models.Msg{}

	msg = append(msg, models.Msg{
		Role:    "system",
		Content: "Ты механизм по определению человеческих эмоций.\\nТвоей основной задачей является определить вероятность каждой эмоции из сказанного предложения от 0 до 1.\\nВ твоем распоряжении только 8 эмоций для угадывания – радость, удивление, страх, отвращение, злость, грусть, стыд, нейтральность.\\nИспользуй эмоции только из указанного списка!\\nПОСТАРАЙСЯ НЕ ПОВТОРЯТЬ ОДНИ И ТЕ ЖЕ РЕПЛИКИ ДВА И БОЛЕЕ РАЗА ПОДРЯД \\nНе пиши ничего в ответ на отправленное предложение, кроме перечисления вероятностей эмоций\\n",
	})
	msg = append(msg, models.Msg{
		Role:    "assistant",
		Content: str,
	})

	body := models.BodyToRequest{
		Model:            "gpt-3.5-turbo",
		Messages:         msg,
		Temperature:      1,
		MaxTokens:        256,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CHAT_TOKEN"))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyReq, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(string(bodyReq))
	var answer models.ToParseAnswer
	err = json.Unmarshal(bodyReq, &answer)
	if err != nil {
		return err
	}

	if len(answer.Choices) == 0 {
		return nil
	}

	if answer.Choices[0] == (models.MsgAnswer{}) {
		return nil
	}

	return s.produce(str, realAnswer, answer.Choices[0].MSG.Content)
}

func ParseString(str string) []models.ToAddInDB {
	feelings := []models.ToAddInDB{}
	feels := strings.Split(str, "\n")

	for _, feel := range feels {
		emotionRate := strings.Split(feel, ":")

		emotionRate[1] = strings.ReplaceAll(emotionRate[1], " ", "")
		feelings = append(feelings, models.ToAddInDB{
			Feeling: strings.ReplaceAll(emotionRate[0], " ", ""),
			Rate:    emotionRate[1],
		})
	}
	fmt.Println(feelings)
	return feelings
}

func (s *NNChatGPTService) produce(sentence, realAnswer, request string) error {
	query := fmt.Sprintf("INSERT INTO nn_answer(sentence, real_answer, current_answer")
	mass := ParseString(request)
	fmt.Println(mass)
	if len(mass) == 0 {
		return nil
	}
	if len(mass) != 0 {
		if mass[0].Feeling == "Вероятностиэмоций" {
			mass = mass[1:]
		}
	}
	var max float64 = 0
	var maxStr string = ""
	for _, mas := range mass {
		f, _ := strconv.ParseFloat(mas.Rate, 64)
		fmt.Println(max, f)
		if f >= max {
			max = f
			maxStr = mas.Feeling
		}
	}
	fmt.Println(maxStr)
	queryAdd := "VALUES(" + "'" + sentence + "'" + ", " + "'" + RussianName(realAnswer) + "'" + ", " + "'" + maxStr + "'"

	for _, mas := range mass {
		query = query + ", " + EnglishName(mas.Feeling)
		queryAdd = queryAdd + ", " + string(mas.Rate)
	}
	query = query + ") " + queryAdd + ")"
	fmt.Println(query)

	return s.repo.Query(query)
}

func EnglishName(str string) string {
	feelings := make(map[string]string)

	feelings["злость"] = "anger"
	feelings["страх"] = "fear"
	feelings["радость"] = "joy"
	feelings["любовь"] = "love"
	feelings["грусть"] = "sadness"
	feelings["удивление"] = "surprise"
	feelings["отвращение"] = "disgust"
	feelings["стыд"] = "shame"
	feelings["нейтральность"] = "neutrality"
	feelings["Злость"] = "anger"
	feelings["Страх"] = "fear"
	feelings["Радость"] = "joy"
	feelings["Любовь"] = "love"
	feelings["Грусть"] = "sadness"
	feelings["Удивление"] = "surprise"
	feelings["Отвращение"] = "disgust"
	feelings["Стыд"] = "shame"
	feelings["Нейтральность"] = "neutrality"
	feelings["восхищение"] = "admiration"
	feelings["веселье"] = "amusement"
	feelings["раздражение"] = "annoyance"
	feelings["одобрение"] = "approval"
	feelings["забота"] = "caring"
	feelings["непонимание"] = "confusion"
	feelings["любопытство"] = "curiosity"
	feelings["желание"] = "desire"
	feelings["разочарование"] = "disappointment"
	feelings["неодобрение"] = "disapproval"
	feelings["смущение"] = "embarrassment"
	feelings["возбуждение"] = "excitement"
	feelings["признательность"] = "gratitude"
	feelings["горе"] = "grief"
	feelings["нервозность"] = "nervousness"
	feelings["оптимизм"] = "optimism"
	feelings["гордость"] = "pride"
	feelings["осознание"] = "realization"
	feelings["облегчение"] = "relief"
	feelings["раскаяние"] = "remorse"

	feelings["Восхищение"] = "admiration"
	feelings["Веселье"] = "amusement"
	feelings["Раздражение"] = "annoyance"
	feelings["Одобрение"] = "approval"
	feelings["Забота"] = "caring"
	feelings["Непонимание"] = "confusion"
	feelings["Любопытство"] = "curiosity"
	feelings["Желание"] = "desire"
	feelings["Разочарование"] = "disappointment"
	feelings["Неодобрение"] = "disapproval"
	feelings["Смущение"] = "embarrassment"
	feelings["Возбуждение"] = "excitement"
	feelings["Признательность"] = "gratitude"
	feelings["Горе"] = "grief"
	feelings["Нервозность"] = "nervousness"
	feelings["Оптимизм"] = "optimism"
	feelings["Гордость"] = "pride"
	feelings["Осознание"] = "realization"
	feelings["Облегчение"] = "relief"
	feelings["Раскаяние"] = "remorse"
	return feelings[str]
}

func RussianName(str string) string {
	feelings := make(map[string]string)

	feelings["anger"] = "злость"
	feelings["fear"] = "страх"
	feelings["joy"] = "радость"
	feelings["love"] = "любовь"
	feelings["sadness"] = "грусть"
	feelings["surprise"] = "удивление"
	feelings["disgust"] = "отвращение"
	feelings["shame"] = "стыд"
	feelings["neutrality"] = "нейтральность"
	feelings["admiration"] = "восхищение"
	feelings["amusement"] = "веселье"
	feelings["annoyance"] = "раздражение"
	feelings["approval"] = "одобрение"
	feelings["caring"] = "забота"
	feelings["confusion"] = "непонимание"
	feelings["curiosity"] = "любопытство"
	feelings["desire"] = "желание"
	feelings["disappointment"] = "разочарование"
	feelings["disapproval"] = "неодобрение"
	feelings["embarrassment"] = "смущение"
	feelings["excitement"] = "возбуждение"
	feelings["gratitude"] = "признательность"
	feelings["grief"] = "горе"
	feelings["nervousness"] = "нервозность"
	feelings["optimism"] = "оптимизм"
	feelings["pride"] = "гордость"
	feelings["realization"] = "осознание"
	feelings["relief"] = "облегчение"
	feelings["remorse"] = "раскаяние"

	return feelings[str]
}

func (s *NNChatGPTService) ParseCSVFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			panic(err)
		}
		fmt.Println(record[0], record[1])
		if err := s.SendTOChatGPT(record[0], record[1]); err != nil {
			fmt.Println(err)
		}
		time.Sleep(20 * time.Second)
	}
}
