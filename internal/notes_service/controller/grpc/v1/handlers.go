package v1

import (
	"context"
	"fmt"
	pb_notes_model "github.com/almalii/grpc-contracts/gen/go/notes_service/model/v1"
	pb_notes_service "github.com/almalii/grpc-contracts/gen/go/notes_service/service/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type NotesServer struct {
	pb_notes_service.UnimplementedNotesServiceServer
}

func NewNotesServer(unimplementedNotesServiceServer pb_notes_service.UnimplementedNotesServiceServer) *NotesServer {
	return &NotesServer{UnimplementedNotesServiceServer: unimplementedNotesServiceServer}
}

func (n *NotesServer) CreateNote(context.Context, *pb_notes_model.CreateNoteRequest) (*pb_notes_model.NoteResponse, error) {
	fmt.Println("CreateNote")
	return &pb_notes_model.NoteResponse{
		Id: "123",
	}, nil
}

func (n *NotesServer) GetNote(context.Context, *pb_notes_model.NoteIDRequest) (*pb_notes_model.NoteResponse, error) {
	fmt.Println("GetNote")
	return &pb_notes_model.NoteResponse{
		Id: "123",
	}, nil
}

//func (n *NotesServer) GetNotes(*pb_notes_model.AuthorIDRequest, pb_notes_model.NotesService_GetNotesServer) error {
//
//}

func (n *NotesServer) UpdateNote(context.Context, *pb_notes_model.UpdateNoteRequest) (*pb_notes_model.NoteResponse, error) {
	fmt.Println("UpdateNote")
	return &pb_notes_model.NoteResponse{
		Id: "123",
	}, nil
}

func (n *NotesServer) DeleteNote(context.Context, *pb_notes_model.NoteIDRequest) (*emptypb.Empty, error) {
	fmt.Println("DeleteNote")
	return &emptypb.Empty{}, nil
}
