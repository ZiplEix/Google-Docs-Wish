package database

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Document struct {
	Title        string
	ID           string
	UserId       string
	LastModified time.Time
	Author       string
}

func NewDocument(data map[string]interface{}, id ...string) *Document {
	doc := &Document{
		Title:        "",
		ID:           "",
		UserId:       "",
		LastModified: time.Time{},
		Author:       "",
	}

	if title, ok := data["title"].(string); ok {
		doc.Title = title
	}

	if len(id) > 0 {
		doc.ID = id[0]
	}

	if userId, ok := data["userId"].(string); ok {
		doc.UserId = userId
	}

	if lastModified, ok := data["last_modified"].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, lastModified)
		if err != nil {
			fmt.Printf("error parsing time: %v\n", err)
			doc.LastModified = time.Now()
		} else {
			doc.LastModified = parsedTime
		}
	}

	if author, ok := data["author"].(string); ok {
		doc.Author = author
	}

	return doc
}

func (doc *Document) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"title":         doc.Title,
		"id":            doc.ID,
		"userId":        doc.UserId,
		"last_modified": doc.LastModified,
		"author":        doc.Author,
	}
}

func CreateNewDocInDb(userId string) (*Document, error) {
	user, err := GetUserFromId(userId)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	doc := NewDocument(map[string]interface{}{
		"title":  "New Document",
		"userId": userId,
		"author": user.Email,
	})

	docRef, wr, err := FirestoreClient.Collection("documents").Add(context.Background(), doc.ToMap())
	if err != nil {
		return nil, fmt.Errorf("error creating document: %v", err)
	}

	fmt.Printf("Document created: %v\n", wr)

	doc.ID = docRef.ID

	_, err = docRef.Set(context.Background(), doc.ToMap())
	if err != nil {
		return nil, fmt.Errorf("error updating document: %v", err)
	}

	return doc, nil
}

func (doc *Document) CreateInDb() (string, error) {
	docRef, wr, err := FirestoreClient.Collection("documents").Add(context.Background(), doc.ToMap())
	if err != nil {
		return "", fmt.Errorf("error creating document: %v", err)
	}

	fmt.Printf("Document created: %v\n", wr)

	doc.ID = docRef.ID

	_, err = docRef.Set(context.Background(), doc.ToMap())
	if err != nil {
		return "", fmt.Errorf("error updating document: %v", err)
	}

	return doc.ID, nil
}

func GetDocumentFromUserId(userId string) ([]*Document, error) {
	docsIter := FirestoreClient.Collection("documents").Where("userId", "==", userId).Documents(context.Background())
	docs, err := docsIter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting documents: %v", err)
	}

	var documents []*Document
	for _, doc := range docs {
		document := NewDocument(doc.Data(), doc.Ref.ID)
		documents = append(documents, document)
	}

	return documents, nil
}

func DeleteDocumentById(docId string) error {
	docRef := FirestoreClient.Collection("documents").Doc(docId)
	_, err := docRef.Delete(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting document: %v", err)
	}

	return nil
}

func SearchDocument(query string, userId string) ([]*Document, error) {
	docsIter := FirestoreClient.Collection("documents").Where("userId", "==", userId).Documents(context.Background())
	docs, err := docsIter.GetAll()
	if err != nil {
		return nil, fmt.Errorf("error getting documents: %v", err)
	}

	var documents []*Document
	query = strings.ToLower(query)
	for _, doc := range docs {
		title := strings.ToLower(doc.Data()["title"].(string))
		if strings.Contains(title, query) {
			document := NewDocument(doc.Data(), doc.Ref.ID)
			documents = append(documents, document)
		}
	}

	return documents, nil
}
