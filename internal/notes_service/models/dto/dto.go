package dto

import (
	pb_notes_model "github.com/almalii/grpc-contracts/gen/go/notes_service/model/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	var tagsList []*pb_notes_model.TagsList
	for _, tag := range resp.Tags {
		tagsList = append(tagsList, &pb_notes_model.TagsList{Tags: []string{tag}})
	}

	return &pb_notes_model.GetNoteResponse{
		Id:        resp.ID.String(),
		Title:     resp.Title,
		Body:      resp.Body,
		Tags:      tagsList, //TODO: wtf?
		Author:    resp.Author.String(),
		CreatedAt: timestamppb.New(resp.CreatedAt),
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
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
		Title: req.Title,
		Body:  req.Body,
		//Tags:  req.Tags, //TODO: update contract
	}
}

func NewUpdateNoteResponse(resp usecase.UpdateNoteInput) *pb_notes_model.UpdateNoteResponse {
	var tagsList []*pb_notes_model.TagsList
	for _, tag := range *resp.Tags {
		tagsList = append(tagsList, &pb_notes_model.TagsList{Tags: []string{tag}})
	}

	return &pb_notes_model.UpdateNoteResponse{
		Title:     *resp.Title,
		Body:      *resp.Body,
		Tags:      tagsList, //TODO: wtf?
		UpdatedAt: timestamppb.New(resp.UpdatedAt),
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
