package dto

import (
	pb_notes_model "github.com/almalii/grpc-contracts/gen/go/notes_service/model/v1"
	"github.com/google/uuid"
	"notes-rew/internal/notes_service/models"
	"notes-rew/internal/notes_service/usecase"
)

func NewCreateNoteInput(currentUser uuid.UUID, req *pb_notes_model.CreateNoteRequest) usecase.CreateNoteInput {
	return usecase.CreateNoteInput{
		Title: req.Title,
		Body:  req.Body,
		//Tags:   req.Tags, TODO: update contract
		Author: currentUser,
	}
}

func NewCreateNoteResponse(idNote uuid.UUID, resp usecase.CreateNoteInput) *pb_notes_model.NoteResponse {
	return &pb_notes_model.NoteResponse{
		Id:    idNote.String(),
		Title: resp.Title,
		Body:  resp.Body,
		//Tags:  resp.Tags, TODO: update contract
		Author: resp.Author.String(),
	}
}

func NewGetNoteInput(req *pb_notes_model.NoteIDRequest) uuid.UUID {
	noteID, err := uuid.Parse(req.Id)
	if err != nil {
		return uuid.Nil
	}
	return noteID
}

func NewGetNoteResponse(resp models.NoteOutput) *pb_notes_model.NoteResponse {
	return &pb_notes_model.NoteResponse{
		Id:    resp.ID.String(),
		Title: resp.Title,
		Body:  resp.Body,
		//Tags:  resp.Tags, TODO: update contract
		Author:    resp.Author.String(),
		CreatedAt: resp.CreatedAt.String(), //TODO: изменить на timestamp
		UpdatedAt: resp.UpdatedAt.String(), //TODO: изменить на timestamp
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
	var notes []*pb_notes_model.NoteResponse
	for _, note := range resp {
		notes = append(notes, NewGetNoteResponse(note))
	}

	return &pb_notes_model.NoteResponseList{
		Notes: notes,
	}
}

func NewUpdateNoteInput(req *pb_notes_model.UpdateNoteRequest) usecase.UpdateNoteInput {
	return usecase.UpdateNoteInput{
		Title: req.Title,
		Body:  req.Body,
		//Tags:   req.Tags, TODO: update contract
	}
}

func NewUpdateNoteResponse(noteId, authorId uuid.UUID, resp usecase.UpdateNoteInput) *pb_notes_model.NoteResponse {
	return &pb_notes_model.NoteResponse{
		Id:    noteId.String(),
		Title: *resp.Title,
		Body:  *resp.Body,
		//Tags:  *resp.Tags, TODO: update contract
		Author: authorId.String(),
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
