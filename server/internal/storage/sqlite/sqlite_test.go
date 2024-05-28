package sqlite

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_SaveThumbnail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	type args struct {
		videoUrl  string
		thumbnail []byte
	}

	storage := &Storage{db: db}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "success",
			mock: func() {
				mock.ExpectExec("INSERT INTO thumbnails").
					WithArgs(
						"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
						string([]byte{255, 216, 255, 224}),
					).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			input: args{
				videoUrl:  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
				thumbnail: []byte{255, 216, 255, 224},
			},
			wantErr: false,
		},
		{
			name: "empty videoURL",
			mock: func() {
				mock.ExpectExec("INSERT INTO thumbnails").
					WithArgs(
						"",
						string([]byte{255, 216, 255, 224}),
					).WillReturnError(errors.New("empty videoURL"))
			},
			input: args{
				videoUrl:  "",
				thumbnail: []byte{255, 216, 255, 224},
			},
			wantErr: true,
		},
		{
			name: "empty thumbnail",
			mock: func() {
				mock.ExpectExec("INSERT INTO thumbnails").
					WithArgs(
						"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
						"",
					).WillReturnError(errors.New("empty thumbnail"))
			},
			input: args{
				videoUrl:  "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
				thumbnail: nil,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := storage.SaveThumbnail(context.Background(), tt.input.videoUrl, tt.input.thumbnail)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
