package memory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitermakov/otusgo-hw/hw12_13_14_15_calendar/internal/model"
)

func TestUserMemoryRepo(t *testing.T) {
	t.Run("complex test", func(t *testing.T) {
		userRepo := UserRepo{}
		ctx := context.Background()
		users := []model.User{
			{
				Name:  "user 1",
				Email: "user1@mail.ru",
			}, {
				Name:  "user 2",
				Email: "user2@yandex.ru",
			}, {
				Name:  "user 3",
				Email: "user3@ya.ru",
			}, {
				Name:  "user 4",
				Email: "user4@inbox.ru",
			}, {
				Name:  "user 5",
				Email: "user5@test.ru",
			},
		}
		for i, user := range users {
			input := model.UserCreate{
				Name:  user.Name,
				Email: user.Email,
			}
			newUser, err := userRepo.Add(ctx, input)
			require.NoError(t, err)

			users[i].ID = newUser.ID
		}

		users[1].Email = "user2_updated@yandex.ru"
		_ = userRepo.Update(ctx, model.UserUpdate{Email: &users[1].Email}, model.UserSearch{ID: &users[1].ID})

		users[2].Name = "user 3 updated"
		_ = userRepo.Update(ctx, model.UserUpdate{Name: &users[2].Name}, model.UserSearch{Email: &users[2].Email})

		_ = userRepo.Delete(ctx, model.UserSearch{ID: &users[0].ID})
		_ = userRepo.Delete(ctx, model.UserSearch{ID: &users[4].ID})

		actual, _ := userRepo.GetList(ctx, model.UserSearch{})

		require.ElementsMatch(t, users[1:4], actual)

		actual, _ = userRepo.GetList(ctx, model.UserSearch{ID: &users[1].ID})
		require.Equal(t, 1, len(actual))

		actual, _ = userRepo.GetList(ctx, model.UserSearch{Email: &users[2].Email})
		require.Equal(t, 1, len(actual))

		actual, _ = userRepo.GetList(ctx, model.UserSearch{Email: &users[0].Email})
		require.Equal(t, 0, len(actual))
	})
}
