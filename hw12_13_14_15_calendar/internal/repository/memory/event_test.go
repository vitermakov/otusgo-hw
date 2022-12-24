package memory

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
)

func TestEventMemoryRepo(t *testing.T) {
	t.Run("complex test", func(t *testing.T) {
		eventRepo := EventRepo{}
		ctx := context.Background()
		baseDate := time.Now()

		userID1, _ := uuid.Parse("ab8e3706-7ad8-11ed-95f7-d00d1b9e4cfe")
		userID2, _ := uuid.Parse("90bdce82-7ad8-11ed-99c1-d00d1b9e4cfe")
		userID3, _ := uuid.Parse("973454b8-7ae0-11ed-97ae-d00d1b9e4cfe")

		events := []model.Event{
			{
				Title:    "title 1",
				Date:     baseDate.Add(10 * time.Hour),
				Duration: time.Minute * 45,
				Owner:    &model.User{ID: userID1},
			}, {
				Title:    "title 2",
				Date:     baseDate.Add(3 * 24 * time.Hour),
				Duration: time.Hour,
				Owner:    &model.User{ID: userID2},
			}, {
				Title:    "title 3",
				Date:     baseDate.Add(4 * 24 * time.Hour),
				Duration: time.Hour,
				Owner:    &model.User{ID: userID3},
			}, {
				Title:    "title 4",
				Date:     baseDate.Add(20 * 24 * time.Hour),
				Duration: time.Hour * 2,
				Owner:    &model.User{ID: userID1},
			}, {
				Title:    "title 5",
				Date:     baseDate.Add(40 * 24 * time.Hour),
				Duration: time.Hour * 3,
				Owner:    &model.User{ID: userID2},
			},
		}
		for i, event := range events {
			input := model.EventCreate{
				Title:    event.Title,
				Date:     event.Date,
				Duration: int(event.Duration.Minutes()),
				OwnerID:  event.Owner.ID,
			}
			newUser, err := eventRepo.Add(ctx, input)
			require.NoError(t, err)

			events[i].ID = newUser.ID
			events[i].CreatedAt = newUser.CreatedAt
			events[i].UpdatedAt = newUser.UpdatedAt
		}

		actual, _ := eventRepo.GetList(ctx, model.EventSearch{})
		require.ElementsMatch(t, events, actual)

		actual, _ = eventRepo.GetList(ctx, model.EventSearch{ID: &events[1].ID})
		require.Equal(t, 1, len(actual))

		actual, _ = eventRepo.GetList(ctx, model.EventSearch{NotID: &events[4].ID})
		require.Equal(t, events[0:4], actual)

		actual, _ = eventRepo.GetList(ctx, model.EventSearch{OwnerID: &userID1})
		require.Equal(t, 2, len(actual))

		actual, _ = eventRepo.GetList(ctx, model.EventSearch{
			DateRange: &model.DateRange{DateStart: baseDate.Add(2 * 24 * time.Hour), Duration: time.Hour * 24 * 3},
		})
		require.Equal(t, 2, len(actual))

		actual, _ = eventRepo.GetList(ctx, model.EventSearch{
			DateRange: &model.DateRange{
				DateStart: baseDate.Add(3*24*time.Hour + time.Minute*30),
				Duration:  time.Hour * 3,
			},
			TacDuration: true,
		})
		require.Equal(t, 1, len(actual))

		actual, _ = eventRepo.GetList(ctx, model.EventSearch{
			DateRange: &model.DateRange{
				DateStart: baseDate.Add(3*24*time.Hour + time.Minute*90),
				Duration:  time.Hour * 3,
			},
			TacDuration: true,
		})
		require.Equal(t, 0, len(actual))
	})
}
