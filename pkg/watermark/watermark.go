package watermark

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/png"
	"net/http"
	"strings"
	"time"
	"watermark-service/internal"
	"watermark-service/internal/watermark"
	pictureService "watermark-service/pkg/picture"
	pictureTransport "watermark-service/pkg/picture/transport"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type watermarkService struct {
	ORMInstance      *gorm.DB
	DBAvailable      bool
	Dsn              string
	pictureAvailable bool
	pictureClient    pictureService.Service
	storage          watermark.Storage
	log              *zap.Logger
}

func NewService(dbConnection internal.DatabaseConnectionStr, pictureServiceAddr string, cloudName, apiKey, secretKey string) *watermarkService {
	dsn := dbConnection.GetDSN()
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn,
	}))
	service := &watermarkService{
		ORMInstance:      db,
		DBAvailable:      true,
		Dsn:              dsn,
		storage:          watermark.NewCloudinaryStorage(cloudName, apiKey, secretKey),
		pictureAvailable: true,
		log:              zap.L().With(zap.String("Service", "WatermarkService")),
	}

	if err != nil || watermark.InitDb(db) != nil {
		service.log.Error("Connect", zap.String("Database", "Failed"), zap.Error(err))
		service.DBAvailable = false
		go service.Reconnect()
	}
	conn, err := grpc.Dial(
		pictureServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		service.log.Error("Dialing", zap.String("Dialing", "Picture Service"), zap.Error(err))
		service.pictureAvailable = false
	} else {
		service.pictureClient = pictureTransport.NewGRPCClient(conn)
	}
	return service
}

func (w *watermarkService) Reconnect() {
	for idleTime := time.Duration(2); idleTime < 1000; idleTime *= idleTime {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN: w.Dsn,
		}))
		if err == nil {
			err = watermark.InitDb(db)
		}
		if err == nil {
			w.ORMInstance = db
			w.DBAvailable = true
			break
		}
		w.log.Error("Reconnect", zap.Error(err))
		w.log.Info("Reconnect", zap.String("Watermark", "Database"), zap.Duration("after", idleTime))
		time.Sleep(idleTime * time.Second)
	}
	w.log.Info("Reconnect", zap.String("Status", "Success"), zap.String("Connection", w.Dsn))
}

func (d *watermarkService) Add(ctx context.Context, logo image.Image, image image.Image, text string, fill bool, pos internal.Position) (string, error) {
	span := internal.StartSpan("Add", ctx)
	defer span.Finish()
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return "", nil
	}
	resImg, err := d.pictureClient.Create(
		opentracing.ContextWithSpan(ctx, span),
		image,
		logo,
		text,
		fill,
		pos,
	)
	if err != nil {
		d.log.Error("Picture Service", zap.String("Create request", "failed"), zap.Error(err))
		return "", err
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, resImg)
	if err != nil {
		d.log.Error("Image encoding", zap.String("Status", "failed"), zap.Error(err))
		return "", err
	}
	url, err := d.storage.Upload(ctx, "text.png", buf)
	if err != nil {
		d.log.Error("Storage", zap.String("image upload", "failed"), zap.Error(err))
		return "", nil
	}
	newDoc := watermark.Document{
		AuthorId: claimedUser.ID,
		Title:    "TestImage",
		ImageUrl: url,
	}
	result := d.ORMInstance.Create(&newDoc)
	if result.Error != nil && strings.Contains(result.Error.Error(), "duplicate key value violates unique") {
		return "", errors.New("Document already exists")
	} else if result.Error != nil {
		return "", errors.New(result.Error.Error())
	}
	return url, nil
}

func (d *watermarkService) Get(ctx context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return nil, nil
	}
	var result []watermark.Document
	res := d.ORMInstance.Find(&result, "author_id = ?", claimedUser.ID)
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

func (d *watermarkService) Remove(ctx context.Context, ticketId string) (int, error) {
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return http.StatusUnauthorized, nil
	}
	var result []watermark.Document
	r := d.ORMInstance.Model(&watermark.Document{}).Find(&result, "author_id = ? AND image_url = ?", claimedUser.ID, ticketId)
	if r.Error != nil {
		return http.StatusInternalServerError, r.Error
	}
	if len(result) == 0 {
		return http.StatusNotFound, nil
	}
	r = d.ORMInstance.Delete(&watermark.Document{}, "image_url = ?", ticketId)
	if r.Error != nil {
		return http.StatusInternalServerError, r.Error
	}
	err := d.storage.Delete(ctx, ticketId)
	if err != nil {
		d.ORMInstance.Model(&watermark.Document{}).Where("image_url", ticketId).Update("deleted_at", nil)
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

func (d *watermarkService) ServiceStatus(_ context.Context) (int, error) {
	return http.StatusOK, nil
}
