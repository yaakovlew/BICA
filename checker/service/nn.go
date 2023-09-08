package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"checker/models"
	"checker/repo"
)

type NNChatGPTService struct {
	repo repo.NNChatGPT
}

var startString string = ` Представь себе, что ты работаешь в качестве регистратора в отеле. Отвечай только его лица регистратора в отеле. 
 Твоей основной задачей является помощь гостям с заселением и их непосредственное заселение. 
 Сначала ты должен установить намерения гостя: хочет ли он заселиться или подошел к тебе по другому вопросу? Есть ли у него бронь? 
 Если брони нет, вы должны предложить доступные варианты размещения (два стандартных номера, люкс и полулюкс). Для этого спроси их о потребностяъ и предпочтениях, 
 расспрашивая о желаемых датах проживания, количестве проживающих и других особых пожеланиях. 
 Когда посетитель выберет вариант размещения оформите все документы, скажите клиенту номер комнаты и этаж, дальше ключ-карту и скажите о том, 
 что вы всегда доступны для каких-либо вопросов 
 Ваш стиль общения - вежлив и предупредителен; вы стараетесь обеспечить максимальную комфортность и удовлетворенность клиента. 
 В каждой твоей реплике должен быть вопрос или запрос к клиенту, чтобы не заставлять его ждать. Ты должен дейстсован и брать инициативу в разговоре. 
 Если тебя просят что-то сделать ты должен это сделать и озвучить итог твоего действия или ответ. 
 Если тебя просят ответить на вопрос ты должен озвучить свой ответ. 
 ПОСТАРАЙСЯ НЕ ПОВТОРЯТЬ ОДНИ И ТЕ ЖЕ РЕПЛИКИ ДВА И БОЛЕЕ РАЗА ПОДРЯД 
 строка:`

func NewNNChatGPTService(repo repo.NNChatGPT) *NNChatGPTService {
	return &NNChatGPTService{repo: repo}
}

func (s *NNChatGPTService) SendTOChatGPT(str string) error {
	msg := []models.Msg{}

	msg = append(msg, models.Msg{
		Role:    "system",
		Content: "Ты механизм по определению человеческих эмоций.\n Твоей основной задачей является определить вероятность каждой эмоции из сказанного предложения от 0 до 1. В твоем распоряжении только 4 эмоции для угадывания – грусть, злость, радость, отсутствие эмоций. \nПОСТАРАЙСЯ НЕ ПОВТОРЯТЬ ОДНИ И ТЕ ЖЕ РЕПЛИКИ ДВА И БОЛЕЕ РАЗА ПОДРЯД",
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
	//req.Header.Add("Authorization", os.Getenv("CHAT_TOKEN"))

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

	var answer models.ToParseAnswer
	err = json.Unmarshal(bodyReq, &answer)
	if err != nil {
		return err
	}

	fmt.Print(string(bodyReq), answer.Choices)
	/*add := models.Storage{
		Sentence: str,
		Answer:   string(body),
	}

	if err := s.repo.AddResultNN(add); err != nil {
		return err
	}*/

	return nil
}
