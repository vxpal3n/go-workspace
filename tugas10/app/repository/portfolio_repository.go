package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"tugas10/app/model"
)

var (
	ErrPortfolioNotFound = errors.New("portofolio belum dibuat")
	ErrProjectLimit      = errors.New("batas jumlah proyek tercapai")
)

const MaxProjectsPerPortfolio = 20

type PortfolioRepository interface {
	EnsureIndexes(ctx context.Context) error
	FindByStudent(ctx context.Context, studentID int) (model.Portfolio, error)
	Upsert(ctx context.Context, p model.Portfolio) (model.Portfolio, error)
	AddProject(ctx context.Context, studentID int, project model.Project) (model.Portfolio, error)
	DeleteProject(ctx context.Context, studentID int, projectID string) (model.Portfolio, error)
	StatsTechnologies(ctx context.Context, limit int) ([]model.TechCount, error)
}

type portfolioDocument struct {
	ID           bson.ObjectID         `bson:"_id,omitempty"`
	StudentID    int                   `bson:"student_id"`
	Headline     string                `bson:"headline"`
	Bio          string                `bson:"bio"`
	Projects     []projectDocument     `bson:"projects"`
	Certificates []certificateDocument `bson:"certificates"`
	CreatedAt    time.Time             `bson:"created_at"`
	UpdatedAt    time.Time             `bson:"updated_at"`
}

type projectDocument struct {
	ID           bson.ObjectID `bson:"_id"`
	Title        string        `bson:"title"`
	Technologies []string      `bson:"technologies"`
	Link         string        `bson:"link,omitempty"`
	Year         int           `bson:"year"`
}

type certificateDocument struct {
	ID     bson.ObjectID `bson:"_id"`
	Name   string        `bson:"name"`
	Issuer string        `bson:"issuer"`
	Date   string        `bson:"date"`
}

func (d portfolioDocument) toModel() model.Portfolio {
	projects := make([]model.Project, 0, len(d.Projects))
	for _, p := range d.Projects {
		techs := p.Technologies
		if techs == nil {
			techs = []string{}
		}
		projects = append(projects, model.Project{
			ID: p.ID.Hex(), Title: p.Title, Technologies: techs,
			Link: p.Link, Year: p.Year,
		})
	}
	certs := make([]model.Certificate, 0, len(d.Certificates))
	for _, c := range d.Certificates {
		certs = append(certs, model.Certificate{
			ID: c.ID.Hex(), Name: c.Name, Issuer: c.Issuer, Date: c.Date,
		})
	}
	return model.Portfolio{
		ID:           d.ID.Hex(),
		StudentID:    d.StudentID,
		Headline:     d.Headline,
		Bio:          d.Bio,
		Projects:     projects,
		Certificates: certs,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

type portfolioMongoRepository struct {
	coll *mongo.Collection
}

func NewPortfolioRepository(db *mongo.Database) PortfolioRepository {
	return &portfolioMongoRepository{coll: db.Collection("student_portfolios")}
}

func (r *portfolioMongoRepository) EnsureIndexes(ctx context.Context) error {
	_, err := r.coll.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "student_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("student_unique"),
		},
		{
			Keys:    bson.D{{Key: "projects.technologies", Value: 1}},
			Options: options.Index().SetName("tech_multikey"),
		},
	})
	if err != nil {
		return fmt.Errorf("membuat index student_portfolios: %w", err)
	}
	return nil
}

func (r *portfolioMongoRepository) FindByStudent(ctx context.Context, studentID int) (model.Portfolio, error) {
	var doc portfolioDocument
	if err := r.coll.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&doc); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Portfolio{}, ErrPortfolioNotFound
		}
		return model.Portfolio{}, fmt.Errorf("mengambil portofolio: %w", err)
	}
	return doc.toModel(), nil
}

func (r *portfolioMongoRepository) Upsert(ctx context.Context, p model.Portfolio) (model.Portfolio, error) {
	t := now()

	projDocs := make([]projectDocument, 0, len(p.Projects))
	for _, proj := range p.Projects {
		techs := proj.Technologies
		if techs == nil {
			techs = []string{}
		}
		projDocs = append(projDocs, projectDocument{
			ID: bson.NewObjectID(), Title: proj.Title, Technologies: techs,
			Link: proj.Link, Year: proj.Year,
		})
	}
	certDocs := make([]certificateDocument, 0, len(p.Certificates))
	for _, c := range p.Certificates {
		certDocs = append(certDocs, certificateDocument{
			ID: bson.NewObjectID(), Name: c.Name, Issuer: c.Issuer, Date: c.Date,
		})
	}

	update := bson.M{
		"$set": bson.M{
			"headline":     p.Headline,
			"bio":          p.Bio,
			"projects":     projDocs,
			"certificates": certDocs,
			"updated_at":   t,
		},
		"$setOnInsert": bson.M{"created_at": t},
	}

	opts := options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After)

	var result portfolioDocument
	if err := r.coll.FindOneAndUpdate(ctx,
		bson.M{"student_id": p.StudentID}, update, opts,
	).Decode(&result); err != nil {
		return model.Portfolio{}, fmt.Errorf("menyimpan portofolio: %w", err)
	}
	return result.toModel(), nil
}

func (r *portfolioMongoRepository) AddProject(ctx context.Context, studentID int, project model.Project) (model.Portfolio, error) {
	techs := project.Technologies
	if techs == nil {
		techs = []string{}
	}
	projDoc := projectDocument{
		ID: bson.NewObjectID(), Title: project.Title, Technologies: techs,
		Link: project.Link, Year: project.Year,
	}

	limitField := "projects." + strconv.Itoa(MaxProjectsPerPortfolio-1)
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var updated portfolioDocument
	err := r.coll.FindOneAndUpdate(ctx,
		bson.M{
			"student_id":    studentID,
			limitField:      bson.M{"$exists": false},
		},
		bson.M{
			"$push": bson.M{"projects": projDoc},
			"$set":  bson.M{"updated_at": now()},
		},
		opts,
	).Decode(&updated)
	if err == nil {
		return updated.toModel(), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.Portfolio{}, fmt.Errorf("menambah proyek: %w", err)
	}

	count, err := r.coll.CountDocuments(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("memeriksa portofolio: %w", err)
	}
	if count == 0 {
		return model.Portfolio{}, ErrPortfolioNotFound
	}
	return model.Portfolio{}, ErrProjectLimit
}

func (r *portfolioMongoRepository) DeleteProject(ctx context.Context, studentID int, projectID string) (model.Portfolio, error) {
	oid, err := bson.ObjectIDFromHex(projectID)
	if err != nil {
		return model.Portfolio{}, ErrInvalidID
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated portfolioDocument
	err = r.coll.FindOneAndUpdate(ctx,
		bson.M{
			"student_id":     studentID,
			"projects._id":   oid,
		},
		bson.M{
			"$pull": bson.M{"projects": bson.M{"_id": oid}},
			"$set":  bson.M{"updated_at": now()},
		},
		opts,
	).Decode(&updated)
	if err == nil {
		return updated.toModel(), nil
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return model.Portfolio{}, fmt.Errorf("menghapus proyek: %w", err)
	}

	count, err := r.coll.CountDocuments(ctx, bson.M{"student_id": studentID})
	if err != nil {
		return model.Portfolio{}, fmt.Errorf("memeriksa portofolio: %w", err)
	}
	if count == 0 {
		return model.Portfolio{}, ErrPortfolioNotFound
	}
	return model.Portfolio{}, ErrNotFound
}

func (r *portfolioMongoRepository) StatsTechnologies(ctx context.Context, limit int) ([]model.TechCount, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	pipeline := mongo.Pipeline{
		{{Key: "$unwind", Value: "$projects"}},
		{{Key: "$unwind", Value: "$projects.technologies"}},
		{{Key: "$group", Value: bson.M{
			"_id":   "$projects.technologies",
			"total": bson.M{"$sum": 1},
		}}},
		{{Key: "$sort", Value: bson.D{
			{Key: "total", Value: -1},
			{Key: "_id", Value: 1},
		}}},
		{{Key: "$limit", Value: int64(limit)}},
	}

	cursor, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("menghitung teknologi: %w", err)
	}
	defer cursor.Close(ctx)

	var out []struct {
		Technology string `bson:"_id"`
		Total      int    `bson:"total"`
	}
	if err := cursor.All(ctx, &out); err != nil {
		return nil, fmt.Errorf("membaca statistik teknologi: %w", err)
	}

	result := make([]model.TechCount, 0, len(out))
	for _, o := range out {
		if o.Technology == "" {
			continue
		}
		result = append(result, model.TechCount{Technology: o.Technology, Total: o.Total})
	}
	return result, nil
}