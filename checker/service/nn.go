package service

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
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
	body := `{
  "model": "gpt-3.5-turbo",
  "messages": [
    {
      "role": "system",
      "content": "Ты механизм по определению человеческих эмоций.\nТвоей основной задачей является определить вероятность каждой эмоции из сказанного предложения от 0 до 1.\nВ твоем распоряжении только 26 эмоций для угадывания – восхищение, веселье, злость, раздражение, одобрение, забота, непонимание, любопытство, желание, разочарование, неодобрение, отвращение, смущение, страх, признательность, горе, радость, любовь, нервозность, оптимизм, гордость, осознание, облегчение, раскаяние, грусть, удивление.\nИспользуй эмоции только из указанного списка!\nПОСТАРАЙСЯ НЕ ПОВТОРЯТЬ ОДНИ И ТЕ ЖЕ РЕПЛИКИ ДВА И БОЛЕЕ РАЗА ПОДРЯД \nНе пиши ничего в ответ на отправленное предложение, кроме перечисления вероятностей эмоций\n"
    },
    {
      "role": "user",
      "content": "Я купил машину, ура!"
    },
    {
      "role": "assistant",
      "content": "восхищение - 0.9\nвеселье - 0.89\nзлость - 0.005\nраздражение - 0.001\nодобрение - 0.1\nзабота - 0.04\nнепонимание - 0.001\nлюбопытство - 0.01\nжелание - 0\nразочарование - 0\nнеодобрение - 0\nотвращение - 0\nсмущение - 0\nстрах -0.001\nпризнательность - 0.1\nгоре - 0\nрадость - 0.99\nлюбовь - 0.54\nнервозность - 0.35\nоптимизм - 0.78\nгордость - 0.69\nосознание - 0.04\nоблегчение - 0.25\nраскаяние - 0.01\nгрусть - 0\nудивление - 0.12"
    },
    {
      "role": "user",
      "content": "У меня не особо радостный день."
    },
    {
      "role": "assistant",
      "content": "восхищение -0\nвеселье - 0\nзлость - 0.5\nраздражение - 0.6\nодобрение - 0\nзабота - 0.1\nнепонимание - 0.12\nлюбопытство - 0.1\nжелание - 0\nразочарование - 0.7\nнеодобрение - 0.75\nотвращение - 0.85\nсмущение - 0.1\nстрах - 0.01\nпризнательность - 0.01\nгоре - 0.9\nрадость - 0.001\nлюбовь - 0\nнервозность - 0.86\nоптимизм - 0\nгордость - 0\nосознание - 0\nоблегчение - 0\nраскаяние - 0.21\nгрусть - 0.95\nудивление - 0.005"
    },
    {
      "role": "user",
      "content":"` + str + `"
    }
  ],
  "temperature": 1,
  "max_tokens": 405,
  "top_p": 1,
  "frequency_penalty": 0,
  "presence_penalty": 0
}`
	fmt.Println(body)
	/*
		jsonData, err := json.Marshal(body)
		if err != nil {
			return err
		}*/

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer([]byte(body)))
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

	fmt.Println(realAnswer)
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
	fmt.Println(feels)
	for _, feel := range feels {
		emotionRate := strings.Split(feel, "-")

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

	sort.Slice(mass, func(i, j int) bool {
		return mass[i].Rate > mass[j].Rate
	})
	var flag string = "false"

	for i := 0; i < len(mass); i++ {
		s, err := strconv.ParseFloat(mass[i].Rate, 64)
		if err == nil {
			continue
		}
		if mass[i].Feeling == realAnswer && s >= 0.6 {
			flag = "true"
		}
	}

	fmt.Println(maxStr)
	queryAdd := "VALUES(" + "'" + sentence + "'" + ", " + "'" + realAnswer + "'" + ", " + "'" + flag + "'"

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
	feelings["грусть"] = "sadness"
	feelings["удивление"] = "surprise"
	feelings["отвращение"] = "disgust"
	//feelings["стыд"] = "shame"
	feelings["нейтральность"] = "neutrality"

	return feelings[str]
}

func RussianName(str string) string {
	feelings := make(map[string]string)

	feelings["anger"] = "злость"
	feelings["fear"] = "страх"
	feelings["joy"] = "радость"
	feelings["sadness"] = "грусть"
	feelings["surprise"] = "удивление"
	feelings["disgust"] = "отвращение"
	feelings["shame"] = "стыд"
	feelings["neutrality"] = "нейтральность"

	/*
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
	*/
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
