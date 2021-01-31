package grpcserver

import (
	"bytes"
	"context"
	"net"
	"os"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage"
	"github.com/anfilat/otus-go/hw12_13_14_15_calendar/internal/storage/initstorage"
)

type SuiteTest struct {
	suite.Suite
	client   CalendarClient
	conn     *grpc.ClientConn
	grpcSrv  *grpc.Server
	listener *bufconn.Listener
	app      app.App
	logg     logger.Logger
	db       storage.Storage
}

func (s *SuiteTest) SetupTest() {
	ctx := context.Background()

	var buf bytes.Buffer
	s.logg, _ = logger.New("", &buf, "")

	dbConnect := os.Getenv("PQ_TEST")
	s.db, _ = initstorage.New(ctx, dbConnect == "", dbConnect)

	s.app = app.New(s.logg, s.db)

	s.conn, _ = grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(s)))
	s.client = NewCalendarClient(s.conn)

	_ = s.app.DeleteAll(ctx)
}

func dialer(s *SuiteTest) func(context.Context, string) (net.Conn, error) {
	s.listener = bufconn.Listen(1024 * 1024)

	s.grpcSrv = grpc.NewServer()
	RegisterCalendarServer(s.grpcSrv, NewService(s.app))

	go func() {
		_ = s.grpcSrv.Serve(s.listener)
	}()

	return func(context.Context, string) (net.Conn, error) {
		return s.listener.Dial()
	}
}

func (s *SuiteTest) TearDownTest() {
	ctx := context.Background()

	_ = s.conn.Close()
	s.grpcSrv.GracefulStop()
	_ = s.listener.Close()
	_ = s.app.DeleteAll(ctx)
	_ = s.db.Close(ctx)
}

func (s *SuiteTest) NewCommonEvent() *Event {
	var eventStart = time.Now().Add(2 * time.Hour)
	var eventStop = eventStart.Add(time.Hour)
	notification := 4 * time.Hour

	return &Event{
		Id:           0,
		Title:        "some event",
		Start:        timestamppb.New(eventStart),
		Stop:         timestamppb.New(eventStop),
		Description:  "the event",
		UserId:       1,
		Notification: durationpb.New(notification),
	}
}

func (s *SuiteTest) EqualEvents(event1, event2 *Event) {
	s.Require().Equal(event1.Title, event2.Title)
	s.Require().Equal(event1.Description, event2.Description)
	s.Require().Equal(event1.Start.AsTime().Unix(), event2.Start.AsTime().Unix())
	s.Require().Equal(event1.Stop.AsTime().Unix(), event2.Stop.AsTime().Unix())
	s.Require().Equal(event1.UserId, event2.UserId)
	if event1.Notification == nil || event2.Notification == nil {
		s.Require().Equal(event1.Notification, event2.Notification)
	} else {
		s.Require().Equal(event1.Notification.AsDuration(), event2.Notification.AsDuration())
	}
}

func (s *SuiteTest) AddEvent(event *Event) int32 {
	ctx := context.Background()
	createRes, err := s.client.Create(ctx, event)
	s.Require().NoError(err)
	return createRes.Id
}
