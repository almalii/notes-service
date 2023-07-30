package v1

import (
	"context"
	pb_notes_model "github.com/almalii/grpc-contracts/gen/go/notes_service/model/v1"
	pb_notes_service "github.com/almalii/grpc-contracts/gen/go/notes_service/service/v1"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
	"notes-rew/internal/notes_service/models"
	"notes-rew/internal/notes_service/models/dto"
	"notes-rew/internal/notes_service/usecase"
)

type NoteUsecase interface {
	CreateNote(ctx context.Context, req usecase.CreateNoteInput) (uuid.UUID, error)
	ReadNote(ctx context.Context, noteID, currentUserID uuid.UUID) (models.NoteOutput, error)
	ReadAllNotes(ctx context.Context, currentUserID uuid.UUID) ([]models.NoteOutput, error)
	UpdateNote(ctx context.Context, id uuid.UUID, req usecase.UpdateNoteInput) error
	DeleteNote(ctx context.Context, id uuid.UUID) error
}

type NotesServer struct {
	usecase   NoteUsecase
	validator *validator.Validate
	pb_notes_service.UnimplementedNotesServiceServer
}

func (n *NotesServer) CreateNote(
	ctx context.Context,
	req *pb_notes_model.CreateNoteRequest,
) (*pb_notes_model.NoteIDResponse, error) {
	currentUserID := ctx.Value("userID").(uuid.UUID)

	input := dto.NewCreateNoteInput(currentUserID, req)

	noteID, err := n.usecase.CreateNote(ctx, input)
	if err != nil {
		logrus.Error("error creating note: ", err)
		return nil, err
	}

	resp := dto.NewCreateNoteResponse(noteID)

	return resp, nil
}

func (n *NotesServer) GetNote(
	ctx context.Context,
	req *pb_notes_model.NoteIDRequest,
) (*pb_notes_model.GetNoteResponse, error) {
	noteID := dto.NewGetNoteInput(req)

	currentUserID := ctx.Value("userID").(uuid.UUID)

	note, err := n.usecase.ReadNote(ctx, noteID, currentUserID)
	if err != nil {
		logrus.Error("error getting note: ", err)
		return nil, err
	}

	resp := dto.NewGetNoteResponse(note)

	return resp, nil
}

func (n *NotesServer) GetNotes(
	ctx context.Context,
	req *pb_notes_model.AuthorIDRequest,
) (*pb_notes_model.NoteResponseList, error) {
	currentUserID := ctx.Value("userID").(uuid.UUID)

	notes, err := n.usecase.ReadAllNotes(ctx, currentUserID)
	if err != nil {
		logrus.Error("error getting notes: ", err)
		return nil, err
	}

	resp := dto.NewGetNotesResponse(notes)

	return resp, nil
}

func (n *NotesServer) UpdateNote(
	ctx context.Context,
	req *pb_notes_model.UpdateNoteRequest,
) (*pb_notes_model.UpdateNoteResponse, error) {
	input := dto.NewUpdateNoteInput(req)
	noteID := dto.NewCurrentNoteID(req)

	currentUserID := ctx.Value("userID").(uuid.UUID)

	_, err := n.usecase.ReadNote(ctx, noteID, currentUserID)
	if err != nil {
		logrus.Error("error getting note: ", err)
		return nil, err
	}

	err = n.usecase.UpdateNote(ctx, noteID, input)
	if err != nil {
		logrus.Error("error updating note: ", err)
		return nil, err
	}

	resp := dto.NewUpdateNoteResponse(input)

	return resp, nil
}

func (n *NotesServer) DeleteNote(
	ctx context.Context,
	req *pb_notes_model.NoteIDRequest,
) (*emptypb.Empty, error) {
	noteID := dto.NewDeleteNoteInput(req)
	currentUserID := ctx.Value("userID").(uuid.UUID)

	_, err := n.usecase.ReadNote(ctx, noteID, currentUserID)
	if err != nil {
		logrus.Error("error getting note: ", err)
		return nil, err
	}

	err = n.usecase.DeleteNote(ctx, noteID)
	if err != nil {
		logrus.Error("error deleting note: ", err)
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func NewNotesServer(
	usecase NoteUsecase,
	unimplementedNotesServiceServer pb_notes_service.UnimplementedNotesServiceServer,
) *NotesServer {
	return &NotesServer{
		usecase:                         usecase,
		UnimplementedNotesServiceServer: unimplementedNotesServiceServer,
	}
}
