package main

import (
	"context"
	"fmt"
	"io"
	"log"
	pb "learn_grpc/proto"
	"time"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ulanishda xato: %v", err)
	}
	defer conn.Close()
	c := pb.NewEducationServiceClient(conn)

	// 1. O‘quvchi qo‘shish
	addStudent(c)

	// 2. Kurslar ro‘yxatini olish
	listCourses(c)

	// 3. To‘lovlarni yuborish
	submitPayments(c)

	// 4. Chat sinovi
	chat(c)
}

func addStudent(c pb.EducationServiceClient) {
	resp, err := c.AddStudent(context.Background(), &pb.StudentRequest{
		Name:     "Ali",
		Age:      20,
		CourseId: "C1",
	})
	if err != nil {
		log.Fatalf("O‘quvchi qo‘shishda xato: %v", err)
	}
	log.Printf("Javob: %s (ID: %s)", resp.Message, resp.StudentId)
}

func listCourses(c pb.EducationServiceClient) {
	stream, err := c.ListCourses(context.Background(), &pb.CoursesRequest{})
	if err != nil {
		log.Fatalf("Kurslarni olishda xato: %v", err)
	}
	for {
		course, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Stream xatosi: %v", err)
		}
		log.Printf("Kurs: %s - %s (O‘qituvchi: %s)", course.CourseId, course.Name, course.Teacher)
	}
}

func submitPayments(c pb.EducationServiceClient) {
	stream, err := c.SubmitPayments(context.Background())
	if err != nil {
		log.Fatalf("To‘lov streamida xato: %v", err)
	}
	payments := []struct{ studentId string; amount float32 }{
		{"STU123", 100.0},
		{"STU123", 150.0},
	}
	for _, p := range payments {
		stream.Send(&pb.PaymentRequest{StudentId: p.studentId, Amount: p.amount})
	}
	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("To‘lov javobida xato: %v", err)
	}
	log.Printf("To‘lov natijasi: %s, Jami: %.2f", resp.Message, resp.TotalAmount)
}

func chat(c pb.EducationServiceClient) {
	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatalf("Chat streamida xato: %v", err)
	}
	go func() {
		for i := 0; i < 3; i++ {
			stream.Send(&pb.ChatMessage{Sender: "Ali", Content: fmt.Sprintf("Xabar %d", i)})
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Chat qabul qilishda xato: %v", err)
		}
		log.Printf("Chat javobi: %s - %s", msg.Sender, msg.Content)
	}
}