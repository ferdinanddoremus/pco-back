package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("24xq4zj5fruqgx6")
		if err != nil {
			return err
		}

		// add
		new_strava_segment := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "2jjmhdkf",
			"name": "strava_segment",
			"type": "url",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"exceptDomains": [],
				"onlyDomains": []
			}
		}`), new_strava_segment); err != nil {
			return err
		}
		collection.Schema.AddField(new_strava_segment)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("24xq4zj5fruqgx6")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("2jjmhdkf")

		return dao.SaveCollection(collection)
	})
}
