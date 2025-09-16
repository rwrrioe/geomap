package service

import (
	"context"
	"google.golang.org/genai"
    "github.com/joho/godotenv"
)


func ServiceRepo struct {
	Repo *repository.ProblemRepository
}

func InitAI() (*genai.Client, error) {
	if err := godotenv.Load(); err != nil {
		return nil , fmt.Errorf("failed to load .env")
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY is not set")
	}

	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return client, nil
}

func(s *ServiceRepo) PredictForDistrict(ctx context.context,districtID int) (*entities.ExtendedAnswer, error) {
	var extendedAnswer entities.ExtendedAnswer
	districtStat, err := repository.GetAnalysisByDistrict(ctx, districtID)
	err != nil {
		return nil, err
	}
	client := InitAI()


		config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"problems": {
							"extended_answer":  {Type: genai.TypeString},
							"status":        {Type: genai.TypeString},
						},
						Required: []string{"extended_answer", "status"},
					},
				},
				Required: []string{"problems"}					
			},

	prompt := fmt.Sprintf(`
У тебя есть анализ по средним значениям и проблемам в городе Алматы по одному данному району. Ты должен интерпретировать эти данные, сделать анализ обощить статистику и сделать будущие прогнозы и риски на основе текущих данных. Сделай 3-5 содержательных предложений.
Строго следуй конфигу и структуре не добавляй лишних комментариев. Статистика ниже:
`, districtStat)

		result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to generate AI response:%w", err)
	}

	fmt.Println("end generating, unmarshalling")
	if err := json.Unmarshal([]byte(result.Text()), &extendedAnswer); err != nil {
		return fmt.Errorf("failed to unmarshal AI response:%w", err)
	}

	return &extendedAnswer, nil
	}

func (s *ServiceRepo) PredictForType(ctx context.context, districtID int) (*entities.ExtendedAnswer, error){
	var extendedAnswer entities.ExtendedAnswer
	typeStat, err := repository.GetAnalysisByType
	if err != nil {
		return nil, err
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"problems": {
							"extended_answer":  {Type: genai.TypeString},
							"status":        {Type: genai.TypeString},
						},
						Required: []string{"extended_answer", "status"},
					},
				},
				Required: []string{"problems"}					
			},

	prompt := fmt.Sprintf(`
У тебя есть анализ по средним значениям и проблемам в городе Алматы по одному данному типу проблем. Ты должен интерпретировать эти данные, сделать анализ обощить статистику и сделать будущие прогнозы и риски на основе текущих данных. Сделай 3-5 содержательных предложений.
Строго следуй конфигу и структуре не добавляй лишних комментариев. Статистика ниже:
`, districtStat)

		result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to generate AI response:%w", err)
	}

	fmt.Println("end generating, unmarshalling")
	if err := json.Unmarshal([]byte(result.Text()), &extendedAnswer); err != nil {
		return fmt.Errorf("failed to unmarshal AI response:%w", err)
	}

	return &extendedAnswer, nil

}



func (s *ServiceRepo) predictForCity(ctx context.context) (*entities.ExtendedAnswer, error) {
	var extendedAnswer entities.ExtendedAnswer
	cityStat, err := repository.GetAnalysisByCity
	if err != nil {
		return nil, err
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"problems": {
							"extended_answer":  {Type: genai.TypeString},
							"status":        {Type: genai.TypeString},
						},
						Required: []string{"extended_answer", "status"},
					},
				},
				Required: []string{"problems"}					
			},

	prompt := fmt.Sprintf(`
У тебя есть анализ по средним значениям и проблемам по всему городу Алматы. Ты должен интерпретировать эти данные, сделать анализ обощить статистику и сделать будущие прогнозы и риски на основе текущих данных. Сделай 3-5 содержательных предложений.
Строго следуй конфигу и структуре не добавляй лишних комментариев. Статистика ниже:
`, districtStat)

		result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to generate AI response:%w", err)
	}

	fmt.Println("end generating, unmarshalling")
	if err := json.Unmarshal([]byte(result.Text()), &extendedAnswer); err != nil {
		return fmt.Errorf("failed to unmarshal AI response:%w", err)
	}

	return &extendedAnswer, nil
}