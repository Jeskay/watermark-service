package watermark

import (
	"context"
	"errors"
	"image"
	"net/http"
	"os"
	"strings"
	"time"
	pictureproto "watermark-service/api/v1/protos/picture"
	"watermark-service/internal"
	"watermark-service/internal/util"
	"watermark-service/internal/watermark"

	"github.com/go-kit/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var logger log.Logger

type watermarkService struct {
	ORMInstance      *gorm.DB
	DBAvailable      bool
	Dsn              string
	pictureAvailable bool
	pictureClient    pictureproto.PictureClient
	storage          watermark.Storage
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
	}
	if err == nil {
		err = watermark.InitDb(db)
	} else {
		service.DBAvailable = false
		go service.Reconnect()
	}
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(pictureServiceAddr, opts...)
	if err != nil {
		logger.Log("Dialing", "PictureService", "Failed:", err)
		service.pictureAvailable = false
	} else {
		service.pictureClient = pictureproto.NewPictureClient(conn)
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
		logger.Log("Reconnect Failed with error: ", err)
		logger.Log("Attempt to reconnect after ", idleTime, " seconds")
		time.Sleep(idleTime * time.Second)
	}
	logger.Log("Reconnected", "to:", w.Dsn)
}

func (d *watermarkService) Add(ctx context.Context, logo image.Image, image image.Image, text string, fill bool, pos internal.Position) (string, error) {
	claimedUser, ok := ctx.Value("user").(*internal.User)
	if !ok {
		return "", nil
	}
	data := util.ImageToBytes(logo, ".png")
	Logo := &pictureproto.Image{Data: data, Type: ".png"}
	data = util.ImageToBytes(image, ".png")
	Image := &pictureproto.Image{Data: data, Type: ".png"}
	resp, err := d.pictureClient.Create(ctx, &pictureproto.CreateRequest{
		Logo:  Logo,
		Image: Image,
		Text:  text,
		Fill:  fill,
		Pos:   pictureproto.Position(pictureproto.Position_value[string(pos)]),
	})
	if err != nil || resp.Err != "" {
		return "", err
	}
	url, err := d.storage.Upload(ctx, "text.png", resp.Image)
	if err != nil {
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

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
