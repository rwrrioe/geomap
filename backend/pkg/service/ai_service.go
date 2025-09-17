package service

//TODO ADD cached fast responses improve prompts
//TODO ADD not found errors
import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rwrrioe/geomap/backend/pkg/entities"
	"github.com/rwrrioe/geomap/backend/pkg/repository"
	"google.golang.org/genai"
)

type AIPredictService struct {
	problemRepo repository.ProblemRepository
}

func NewAIPredictService(problemRepo repository.ProblemRepository) *AIPredictService {
	return &AIPredictService{
		problemRepo: problemRepo,
	}
}
func InitAI(ctx context.Context) (*genai.Client, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("failed to load .env")
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

func (s *AIPredictService) PredictForDistrict(ctx context.Context, districtID int) error {
	var extendedAIAnswer entities.ExtendedAIResponse
	districtStat, err := s.problemRepo.GetAnalysisByDistrict(ctx, districtID)
	if err != nil {
		return err
	}
	cityStat, err := s.problemRepo.GetAnalysisByCity(ctx)
	if err != nil {
		return err
	}
	client, err := InitAI(ctx)
	if err != nil {
		return err
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"extended_answer": {Type: genai.TypeString},
				"status":          {Type: genai.TypeString},
			},
			Required: []string{"extended_answer", "status"},
		},
	}

	prompt := fmt.Sprint(`
У тебя есть анализ по средним значениям и проблемам в городе Алматы по одному данному району. В первом наборе данных type_id - айди типа проблемы, type_name - тип проблемы, problems_count - число проблем в данном районе, solved_count - число решенных проблем в данном районе, 
imp_avg - среднее по шкале важности проблем в данном районе(от 1 до 10). Во втором наборе данных усредненные данные по всему городу: problem_count - число проблем во всем городе, status_count - число решенных проблем во всем городе, imp_avg - среднее важности проблем по всему городу(от 1 до 10)
Ты должен интерпретировать эти данные, 
сделать анализ обощить статистику, сравнить со значениями по городу.Ты должен сделать будущие конкретные прогнозы для данного района основанные на типах возникаемых проблем и их частотею Сделай 4-5 содержательных предложений.
Строго следуй конфигу и структуре не добавляй лишних комментариев. Статистика ниже. 
`, districtStat, cityStat)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to generate AI response:%w", err)
	}

	fmt.Println("end generating, unmarshalling")
	if err := json.Unmarshal([]byte(result.Text()), &extendedAIAnswer); err != nil {
		return fmt.Errorf("failed to unmarshal AI response:%w", err)
	}

	err = s.problemRepo.CacheAIResponse(ctx, &extendedAIAnswer, districtID)
	if err != nil {
		return err
	}

	return nil

}

func (s *AIPredictService) PredictForType(ctx context.Context, typeID int) error {
	var extendedAIAnswer entities.ExtendedAIResponse
	typeStat, err := s.problemRepo.GetAnalysisByType(ctx, typeID)
	if err != nil {
		return err
	}
	cityStat, err := s.problemRepo.GetAnalysisByCity(ctx)
	if err != nil {
		return err
	}

	client, err := InitAI(ctx)
	if err != nil {
		return err
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"extended_answer": {Type: genai.TypeString},
				"status":          {Type: genai.TypeString},
			},
			Required: []string{"extended_answer", "status"},
		},
	}

	prompt := fmt.Sprint(`
У тебя есть анализ по средним значениям и проблемам в городе Алматы по одному данному типу проблем. В первом наборе данных district_id - айди районе, district_name - имя района, problems_count - число проблем в данном районе, solved_count - число решенных проблем в данном районе, 
imp_avg - среднее по шкале важности проблем в данном районе(от 1 до 10). Во втором наборе данных усредненные данные по всему городу: problem_count - число проблем во всем городе, status_count - число решенных проблем во всем городе, imp_avg - среднее важности проблем по всему городу(от 1 до 10)
Ты должен интерпретировать эти данные, 
сделать анализ обощить статистику, указав критические районы с данным типом проблем.Ты должен делать будущие конкретные прогнозы на основе типа проблемы и сравнения с данными по городу. Сделай 4-5 содержательных предложений.
Строго следуй конфигу и структуре не добавляй лишних комментариев. Статистика ниже. 
`, typeStat, cityStat)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to generate AI response:%w", err)
	}

	fmt.Println("end generating, unmarshalling")
	if err := json.Unmarshal([]byte(result.Text()), &extendedAIAnswer); err != nil {
		return fmt.Errorf("failed to unmarshal AI response:%w", err)
	}

	err = s.problemRepo.CacheAIResponse(ctx, &extendedAIAnswer, typeID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AIPredictService) PredictForCity(ctx context.Context) error {
	var extendedAIAnswer entities.ExtendedAIResponse
	cityStat, err := s.problemRepo.GetAnalysisByCity(ctx)
	if err != nil {
		return err
	}

	client, err := InitAI(ctx)
	if err != nil {
		return err
	}

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type: genai.TypeObject,
			Properties: map[string]*genai.Schema{
				"extended_answer": {Type: genai.TypeString},
				"status":          {Type: genai.TypeString},
			},
			Required: []string{"extended_answer", "status"},
		},
	}

	prompt := fmt.Sprint(`
У тебя есть анализ по средним значениям и проблемам в городе Алматы. В наборе данных усредненные данные по всему городу: problem_count - число проблем во всем городе, status_count - число решенных проблем во всем городе, imp_avg - среднее важности проблем по всему городу(от 1 до 10)
Ты должен:
1) Сделать анализ и выделить статистику 
2) Выделить ключевые районы и их проблемы.
3)Сделать будущие конкретные прогнозы и риски на основе текущих данных
4) найти решения.
Выделить конкретный план действий и районы.Сделай 4-5 содержательных предложений.
Строго следуй конфигу и структуре не добавляй лишних комментариев. Статистика ниже. :
`, cityStat)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		config,
	)
	if err != nil {
		return fmt.Errorf("failed to generate AI response:%w", err)
	}

	fmt.Println("end generating, unmarshalling")
	if err := json.Unmarshal([]byte(result.Text()), &extendedAIAnswer); err != nil {
		return fmt.Errorf("failed to unmarshal AI response:%w", err)
	}

	err = s.problemRepo.CacheAIResponse(ctx, &extendedAIAnswer, 1)
	if err != nil {
		return err
	}

	fmt.Println(extendedAIAnswer.AnswerText)
	return nil
}

func (s *AIPredictService) GetAnalysisByCity(ctx context.Context) (*entities.ExtendedAIResponse, error) {
	analysis, err := s.problemRepo.GetAIResponseById(ctx, -1)
	if err != nil {

		err := s.PredictForCity(ctx)
		if err != nil {
			log.Fatal(err)
		}

		analysis, err := s.problemRepo.GetAIResponseById(ctx, -1)
		if err != nil {
			log.Fatal(err)
		}

		analysisDTO := entities.ExtendedAIResponse{
			AnswerText: analysis.ResponseText,
			Status:     analysis.Status,
		}
		return &analysisDTO, nil
	}

	analysisDTO := entities.ExtendedAIResponse{
		AnswerText: analysis.ResponseText,
		Status:     analysis.Status,
	}

	return &analysisDTO, nil
}

func (s *AIPredictService) GetAnalysisByDistrict(ctx context.Context, districtID int) (*entities.ExtendedAIResponse, error) {
	analysis, err := s.problemRepo.GetAIResponseById(ctx, districtID)
	if err != nil {

		err := s.PredictForCity(ctx)
		if err != nil {
			log.Fatal(err)
		}

		analysis, err := s.problemRepo.GetAIResponseById(ctx, districtID)
		if err != nil {
			log.Fatal(err)
		}

		analysisDTO := entities.ExtendedAIResponse{
			AnswerText: analysis.ResponseText,
			Status:     analysis.Status,
		}
		return &analysisDTO, nil
	}

	analysisDTO := entities.ExtendedAIResponse{
		AnswerText: analysis.ResponseText,
		Status:     analysis.Status,
	}

	return &analysisDTO, nil
}

func (s *AIPredictService) GetAnalysisByType(ctx context.Context, typeID int) (*entities.ExtendedAIResponse, error) {
	analysis, err := s.problemRepo.GetAIResponseById(ctx, typeID)
	if err != nil {

		err := s.PredictForCity(ctx)
		if err != nil {
			log.Fatal(err)
		}

		analysis, err := s.problemRepo.GetAIResponseById(ctx, typeID)
		if err != nil {
			log.Fatal(err)
		}

		analysisDTO := entities.ExtendedAIResponse{
			AnswerText: analysis.ResponseText,
			Status:     analysis.Status,
		}
		return &analysisDTO, nil
	}

	analysisDTO := entities.ExtendedAIResponse{
		AnswerText: analysis.ResponseText,
		Status:     analysis.Status,
	}

	return &analysisDTO, nil
}
