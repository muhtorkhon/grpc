package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	pb "learn_grpc/proto"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedEducationServiceServer
}

func (s *server) AddStudent(ctx context.Context, req *pb.StudentRequest) (*pb.StudentResponse, error) {
	log.Printf("Yangi oquvchi: %s, Kurs: %s", req.Name, req.CourseId)
	return &pb.StudentResponse{
		StudentId: "STU123", // Bu yerda haqiqiy ID generatsiya qilinishi mumkin
		Message:   fmt.Sprintf("%s muvaffaqiyatli qoshildi!", req.Name),
	}, nil
}

func (s *server) ListCourses(req *pb.CoursesRequest, stream pb.EducationService_ListCoursesServer) error {
	courses := []struct{ id, name, teacher string }{
		{"C1", "Go dasturlash", "Ahmad"},
		{"C2", "Python", "Zarina"},
		{"C3", "Data Science", "Jasur"},
	}
	for _, course := range courses {
		stream.Send(&pb.CourseResponse{
			CourseId: course.id,
			Name:     course.name,
			Teacher:  course.teacher,
		})
	}
	return nil
}

func (s *server) SubmitPayments(stream pb.EducationService_SubmitPaymentsServer) error {
	var total float32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.PaymentResponse{
				Message:     "Tolovlar qabul qilindi",
				TotalAmount: total,
			})
		}
		if err != nil {
			return err
		}
		total += req.Amount
		log.Printf("Tolov: %s uchun %.2f", req.StudentId, req.Amount)
	}
}

func (s *server) Chat(stream pb.EducationService_ChatServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Printf("%s: %s", msg.Sender, msg.Content)
		// Echo qaytarish (haqiqiy chatda logika qo‘shiladi)
		stream.Send(&pb.ChatMessage{
			Sender:  "Server",
			Content: "Javob: " + msg.Content,
		})
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Portni tinglashda xato: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterEducationServiceServer(s, &server{})
	log.Printf("Server %v da ishga tushdi", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Serverni ishga tushirishda xato: %v", err)
	}
}