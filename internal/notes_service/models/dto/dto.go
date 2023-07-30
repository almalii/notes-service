package dto

import (
	pb_notes_model "github.com/almalii/grpc-contracts/gen/go/notes_service/model/v1"
	"github.com/google/uuid"
	"notes-rew/internal/notes_service/models"
	"notes-rew/internal/notes_service/usecase"
)

func NewCreateNoteInput(currentUser uuid.UUID, req *pb_notes_model.CreateNoteRequest) usecase.CreateNoteInput {
	return usecase.CreateNoteInput{
		Title:  req.Title,
		Body:   req.Body,
		Tags:   req.Tags,
		Author: currentUser,
	}
}

func NewCreateNoteResponse(idNote uuid.UUID) *pb_notes_model.NoteIDResponse {
	return &pb_notes_model.NoteIDResponse{
		Id: idNote.String(),
	}
}

func NewGetNoteInput(req *pb_notes_model.NoteIDRequest) uuid.UUID {
	noteID, err := uuid.Parse(req.Id)
	if err != nil {
		return uuid.Nil
	}
	return noteID
}

func NewGetNoteResponse(resp models.NoteOutput) *pb_notes_model.GetNoteResponse {
	return &pb_notes_model.GetNoteResponse{
		Id:        resp.ID.String(),
		Title:     resp.Title,
		Body:      resp.Body,
		Tags:      resp.Tags,
		Author:    resp.Author.String(),
		CreatedAt: resp.CreatedAt.String(),
		UpdatedAt: resp.UpdatedAt.String(),
	}
}

func NewGetNotesInput(req *pb_notes_model.AuthorIDRequest) uuid.UUID {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return uuid.Nil
	}
	return userID
}

func NewGetNotesResponse(resp []models.NoteOutput) *pb_notes_model.NoteResponseList {
	var notes []*pb_notes_model.GetNoteResponse
	for _, note := range resp {
		notes = append(notes, NewGetNoteResponse(note))
	}

	return &pb_notes_model.NoteResponseList{
		Notes: notes,
	}
}

func NewUpdateNoteInput(req *pb_notes_model.UpdateNoteRequest) usecase.UpdateNoteInput {
	return usecase.UpdateNoteInput{
		Title: &req.Title,
		Body:  &req.Body,
		Tags:  &req.Tags,
	}
}

func NewUpdateNoteResponse(resp usecase.UpdateNoteInput) *pb_notes_model.UpdateNoteResponse {
	return &pb_notes_model.UpdateNoteResponse{
		Title:     *resp.Title,
		Body:      *resp.Body,
		Tags:      *resp.Tags,
		UpdatedAt: resp.UpdatedAt.String(),
	}
}

func NewCurrentNoteID(req *pb_notes_model.UpdateNoteRequest) uuid.UUID {
	noteID, err := uuid.Parse(req.Id)
	if err != nil {
		return uuid.Nil
	}
	return noteID
}

func NewDeleteNoteInput(req *pb_notes_model.NoteIDRequest) uuid.UUID {
	noteID, err := uuid.Parse(req.Id)
	if err != nil {
		return uuid.Nil
	}
	return noteID
}
