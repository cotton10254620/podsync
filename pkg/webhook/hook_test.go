package webhook

import (
	"testing"
	"time"

	"github.com/go-pg/pg"
	"github.com/mxpv/patreon-go"
	"github.com/mxpv/podsync/pkg/models"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	pledge := createPledge()

	hook := createHandler(t)
	err := hook.Handle(pledge, patreon.EventCreatePledge)
	require.NoError(t, err)

	model := &models.Pledge{PledgeID: 12345}
	err = hook.db.Select(model)
	require.NoError(t, err)
	require.Equal(t, pledge.Attributes.AmountCents, model.AmountCents)
}

func TestUpdate(t *testing.T) {
	pledge := createPledge()

	hook := createHandler(t)
	err := hook.Handle(pledge, patreon.EventCreatePledge)
	require.NoError(t, err)

	pledge.Attributes.AmountCents = 999

	err = hook.Handle(pledge, patreon.EventUpdatePledge)
	require.NoError(t, err)

	model := &models.Pledge{PledgeID: 12345}
	err = hook.db.Select(model)
	require.NoError(t, err)
	require.Equal(t, 999, model.AmountCents)
}

func TestDelete(t *testing.T) {
	pledge := createPledge()
	hook := createHandler(t)

	err := hook.Handle(pledge, patreon.EventCreatePledge)
	require.NoError(t, err)

	err = hook.Handle(pledge, patreon.EventDeletePledge)
	require.NoError(t, err)
}

func TestFindPledge(t *testing.T) {
	pledge := createPledge()
	hook := createHandler(t)

	err := hook.Handle(pledge, patreon.EventCreatePledge)
	require.NoError(t, err)

	res, err := hook.FindPledge(pledge.ID)
	require.NoError(t, err)
	require.Equal(t, res.AmountCents, pledge.Attributes.AmountCents)
}

func createHandler(t *testing.T) *Handler {
	opts, err := pg.ParseURL("postgres://postgres:@localhost/podsync?sslmode=disable")
	if err != nil {
		require.NoError(t, err)
	}

	db := pg.Connect(opts)

	_, err = db.Model(&models.Pledge{}).Where("1=1").Delete()
	require.NoError(t, err)

	return NewHookHandler(db)
}

func createPledge() *patreon.Pledge {
	pledge := &patreon.Pledge{
		ID:   "12345",
		Type: "pledge",
	}

	pledge.Attributes.AmountCents = 400
	pledge.Attributes.CreatedAt = patreon.NullTime{Valid: true, Time: time.Now().UTC()}

	pledge.Relationships.Patron = &patreon.PatronRelationship{}
	pledge.Relationships.Patron.Data.ID = "67890"

	return pledge
}
