package fs_pkg

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

func Update(ctx context.Context, client *firestore.Client, docPath string, updateStruct []firestore.Update) {
	iter := client.Collection(docPath).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
			fmt.Println("Failed Ref: ", doc.Ref)
		}

		_, err = doc.Ref.Update(ctx, updateStruct)
		if err != nil {
			log.Fatalf("Failed to update: %v", err)
			fmt.Println("Failed Ref: ", doc.Ref)
		}
	}
}
