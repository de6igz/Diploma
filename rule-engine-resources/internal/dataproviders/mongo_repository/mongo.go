package mongo

import (
	"context"
	"rule-engine-resources/internal/domain"
	"rule-engine-resources/internal/usecases"
	"strconv"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRuleRepository struct {
	coll   *mongo.Collection
	logger *zerolog.Logger
}

func NewMongoRuleRepository(db *mongo.Database, collName string, logger zerolog.Logger) *MongoRuleRepository {
	return &MongoRuleRepository{
		coll:   db.Collection(collName),
		logger: &logger,
	}
}

func (mr *MongoRuleRepository) GetAllRules(ctx context.Context) ([]domain.Rule, error) {
	mr.logger.Debug().Msg("MongoRuleRepository.GetAllRules: fetching all rules from DB")
	cur, err := mr.coll.Find(ctx, bson.D{})
	if err != nil {
		mr.logger.Error().Err(err).Msg("Failed to fetch rules from Mongo")
		return nil, err
	}
	defer cur.Close(ctx)

	var rules []domain.Rule
	for cur.Next(ctx) {
		var r domain.Rule
		if err := cur.Decode(&r); err != nil {
			mr.logger.Warn().Err(err).Msg("Failed to decode rule doc")
			continue
		}
		rules = append(rules, r)
	}
	mr.logger.Debug().Msgf("Fetched %d rules from Mongo", len(rules))
	return rules, nil
}

// GetRulesByUserAndService – основной метод, возвращает только те правила,
// которые принадлежат userID и serviceName.
func (mr *MongoRuleRepository) GetRulesByUserAndServiceAndProjectId(
	ctx context.Context,
	userID string,
	serviceName string,
	projectId string,
) ([]domain.Rule, error) {

	mr.logger.Debug().Msgf("MongoRuleRepository.GetRulesByUserAndService user=%s service=%s", userID, serviceName)

	userIdInt, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		mr.logger.Error().Err(err).Msg("Failed to parse userID")
		return nil, err
	}

	filter := bson.M{
		"user_id":      userIdInt,
		"service_name": serviceName,
		"project_id":   projectId,
	}
	cur, err := mr.coll.Find(ctx, filter)
	if err != nil {
		mr.logger.Error().Err(err).Msg("Failed to fetch filtered rules from Mongo")
		return nil, err
	}
	defer cur.Close(ctx)

	var rules []domain.Rule
	for cur.Next(ctx) {
		var r domain.Rule
		if err := cur.Decode(&r); err != nil {
			mr.logger.Warn().Err(err).Msg("Failed to decode rule doc")
			continue
		}
		rules = append(rules, r)
	}
	mr.logger.Debug().Msgf("Fetched %d rules for user=%s / service=%s", len(rules), userID, serviceName)
	return rules, nil
}

var _ usecases.RuleRepository = (*MongoRuleRepository)(nil)
