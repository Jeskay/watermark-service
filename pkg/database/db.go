package database

import (
	"context"
	"errors"
	"image"
	"net/http"
	"strings"
	watermarkproto "watermark-service/api/v1/protos/watermark"
	"watermark-service/internal"
	"watermark-service/internal/database"
	"watermark-service/internal/util"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type dbService struct {
	orm                *gorm.DB
	watermarkAvailable bool
	watermarkClient    watermarkproto.WatermarkClient
	storage            database.Storage
}

func NewService(dbORM *gorm.DB, watermarkServiceAddr string, cloudName, apiKey, secretKey string) *dbService {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(watermarkServiceAddr, opts...)
	if err != nil {
		return &dbService{
			orm:                dbORM,
			watermarkAvailable: false,
		}
	}
	c := watermarkproto.NewWatermarkClient(conn)
	return &dbService{
		orm:                dbORM,
		watermarkAvailable: true,
		watermarkClient:    c,
		storage:            database.NewCloudinaryStorage(cloudName, apiKey, secretKey),
	}
}

func (d *dbService) Add(ctx context.Context, logo image.Image, image image.Image, text string, fill bool, pos internal.Position) (string, error) {
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return "", nil
	}
	data := util.ImageToBytes(logo, ".png")
	Logo := &watermarkproto.Image{Data: data, Type: ".png"}
	data = util.ImageToBytes(image, ".png")
	Image := &watermarkproto.Image{Data: data, Type: ".png"}
	resp, err := d.watermarkClient.Create(ctx, &watermarkproto.CreateRequest{
		Logo:  Logo,
		Image: Image,
		Text:  text,
		Fill:  fill,
		Pos:   watermarkproto.Position(watermarkproto.Position_value[string(pos)]),
	})
	if err != nil || resp.Err != "" {
		return "", err
	}
	url, err := d.storage.Upload(ctx, "text.png", resp.Image)
	if err != nil {
		return "", nil
	}
	newDoc := database.Document{
		AuthorId: claimedUser.ID,
		Title:    "TestImage",
		ImageUrl: url,
	}
	result := d.orm.Create(&newDoc)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return "", errors.New("Document already exists")
	} else if result.Error != nil {
		return "", errors.New(result.Error.Error())
	}
	return url, nil
}

func (d *dbService) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return nil, nil
	}
	var result []database.Document
	res := d.orm.Find(&result, "author_id = ?", claimedUser.ID)
	if res.Error != nil {
		return nil, res.Error
	}
	docs := make([]internal.Document, len(result))
	for i, doc := range result {
		docs[i] = internal.Document{
			ID:       doc.ID,
			AuthorId: doc.AuthorId,
			Title:    doc.Title,
			ImageUrl: doc.ImageUrl,
		}
	}
	return docs, nil
}

func (d *dbService) Remove(ctx context.Context, ticketId string) (int, error) {
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return http.StatusUnauthorized, nil
	}
	var result []database.Document
	r := d.orm.Model(&database.Document{}).Find(&result, "author_id = ? AND image_url = ?", claimedUser.ID, ticketId)
	if r.Error != nil {
		return http.StatusInternalServerError, r.Error
	}
	if len(result) == 0 {
		return http.StatusNotFound, nil
	}
	r = d.orm.Delete(&database.Document{}, "image_url = ?", ticketId)
	if r.Error != nil {
		return http.StatusInternalServerError, r.Error
	}
	err := d.storage.Delete(ctx, ticketId)
	if err != nil {
		d.orm.Model(&database.Document{}).Where("image_url", ticketId).Update("deleted_at", nil)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (d *dbService) ServiceStatus(_ context.Context) (int, error) {
	return http.StatusOK, nil
}
